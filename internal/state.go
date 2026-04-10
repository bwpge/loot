package internal

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand/v2"
	"os"
	"slices"
	"strings"

	"loot/internal/config"
)

type Entry struct {
	Value   string   `json:"value"`
	Comment string   `json:"comment"`
	Tags    []string `json:"tags"`
	Hosts   []string `json:"hosts"`
}

var (
	errEntryNotFound = errors.New("entry id not found")
	errAmbiguousID   = errors.New("id matches multiple entries")
)

type State struct {
	Hashes HashSet          `json:"hashes"`
	Data   map[string]Entry `json:"data"`
}

func NewState() *State {
	return &State{
		Data: make(map[string]Entry),
	}
}

func LoadState(path string) (*State, error) {
	s := State{}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *State) Save(path string) error {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0o644)
	return err
}

func (s *State) ContainsValue(value string) bool {
	return s.Hashes.Contains(value)
}

func (s *State) Add(e Entry) string {
	id := generateID()
	for {
		if _, found := s.Data[id]; !found {
			break
		}
		id = generateID()
	}

	if len(e.Hosts) == 0 {
		e.Hosts = config.Get().DefaultHosts
	}

	s.Data[id] = e
	s.Hashes.Add(e.Value)

	return id
}

func (s *State) Remove(id string) error {
	if _, found := s.Data[id]; !found {
		return errEntryNotFound
	}

	delete(s.Data, id)
	s.updateHashes()

	return nil
}

func (s *State) Get(id string) (*Entry, error) {
	e, found := s.Data[id]
	if !found {
		return nil, errEntryNotFound
	}
	return &e, nil
}

func (s *State) UpdateValue(id string, value string) error {
	if _, found := s.Data[id]; !found {
		return errEntryNotFound
	}

	e := s.Data[id]
	e.Value = value
	s.Data[id] = e
	s.updateHashes()

	return nil
}

func (s *State) FindID(prefix string) (string, error) {
	result := ""
	for k := range s.Data {
		if strings.HasPrefix(k, prefix) {
			if result == "" {
				result = k
			} else {
				return "", errAmbiguousID
			}
		}
	}

	if result == "" {
		return "", errEntryNotFound
	}

	return result, nil
}

type EntryFilter struct {
	ID    []string
	Tags  []string
	Hosts []string
}

func (s *State) Filter(f EntryFilter) map[string]Entry {
	result := make(map[string]Entry)
	if len(f.ID)+len(f.Tags)+len(f.Hosts) == 0 {
		return s.Data
	}

	filterFunc := func(id string) bool {
		for _, i := range f.ID {
			if strings.HasPrefix(id, i) {
				return true
			}
		}
		v := s.Data[id]
		for _, t := range f.Tags {
			if slices.Contains(v.Tags, t) {
				return true
			}
		}
		for _, h := range f.Hosts {
			if slices.Contains(v.Hosts, h) {
				return true
			}
		}
		return false
	}

	for k, v := range s.Data {
		if filterFunc(k) {
			result[k] = v
		}
	}

	return result
}

func (s *State) Clear() {
	clear(s.Data)
	clear(s.Hashes)
}

func (s *State) updateHashes() {
	// PERF: this is really expensive but we're not managing that many values
	// so not worried about it at the moment. will fix if it becomes an issue.
	clear(s.Hashes)
	for _, v := range s.Data {
		s.Hashes.Add(v.Value)
	}
}

const idLen = 6

func generateID() string {
	b := make([]byte, idLen)
	for i := range b {
		b[i] = byte(rand.Uint32())
	}
	return hex.EncodeToString(b)
}
