package state

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand/v2"
	"os"
	"path"
	"strings"

	"github.com/bwpge/loot/internal"
	"github.com/bwpge/loot/internal/entry"
)

var (
	errEntryNotFound   = errors.New("entry not found")
	errAmbiguousPrefix = errors.New("value matches multiple entries")
)

type Entry = entry.Entry

type State struct {
	Hashes internal.HashSet      `json:"hashes"`
	Data   map[string]Entry      `json:"data"`
	Flags  map[string]entry.Flag `json:"flags"`
}

func New() *State {
	return &State{
		Data:  make(map[string]Entry),
		Flags: make(map[string]entry.Flag),
	}
}

func Load(path string) (*State, error) {
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

	s.Data[id] = e
	s.Hashes.Add(e.Value)

	return id
}

func (s *State) Remove(id string) error {
	if _, found := s.Data[id]; !found {
		return errEntryNotFound
	}

	delete(s.Data, id)
	s.UpdateHashes()

	return nil
}

func (s *State) Get(id string) (*Entry, error) {
	e, found := s.Data[id]
	if !found {
		return nil, errEntryNotFound
	}
	return &e, nil
}

func (s *State) FindID(prefix string) (string, error) {
	return find(s.Data, prefix)
}

func (s *State) FindFlag(prefix string) (string, error) {
	return find(s.Flags, prefix)
}

func find[T any](data map[string]T, prefix string) (string, error) {
	result := ""
	for k := range data {
		if strings.HasPrefix(k, prefix) {
			if result == "" {
				result = k
			} else {
				return "", errAmbiguousPrefix
			}
		}
	}

	if result == "" {
		return "", errEntryNotFound
	}

	return result, nil
}

func (s *State) Filter(f entry.Filter) map[string]Entry {
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
			// sane default behavior for terms not using special patterns
			if !strings.ContainsAny(t, "*?[]-\\") {
				t += "*"
			}

			for _, vt := range v.Tags {
				if match, _ := path.Match(t, vt); match {
					return true
				}
			}
		}
		for _, h := range f.Hosts {
			// sane default behavior for terms not using special patterns
			if !strings.ContainsAny(h, "*?[]-\\") {
				h += "*"
			}

			for _, vh := range v.Hosts {
				if match, _ := path.Match(h, vh); match {
					return true
				}
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

func (s *State) Capture(flag string, owner string, host string) {
	s.Flags[flag] = entry.Flag{Owner: owner, Host: host}
}

func (s *State) Clear() {
	clear(s.Data)
	clear(s.Flags)
	clear(s.Hashes)
}

func (s *State) UpdateHashes() {
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
