package state

import (
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path"
	"reflect"
	"slices"
	"strings"

	"github.com/bwpge/loot/internal/entry"
)

var (
	errEntryNotFound   = errors.New("entry not found")
	errAmbiguousPrefix = errors.New("value matches multiple entries")
)

type Entry = entry.Entry

type State struct {
	Data  map[string]Entry      `json:"data"`
	Flags map[string]entry.Flag `json:"flags"`
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

func (s *State) Add(e Entry) string {
	id := getID(e.Value)
	s.Data[id] = e

	return id
}

func (s *State) Merge(e Entry) (string, bool) {
	changed := false
	old, found := s.Find(e.Value)
	id := getID(e.Value)
	if !found || reflect.DeepEqual(old, e) {
		return id, changed
	}

	if e.Comment == "" {
		e.Comment = old.Comment
	}
	e.Tags = mergeSlices(old.Tags, e.Tags)
	changed = changed || !slices.Equal(old.Tags, e.Tags)
	e.Hosts = mergeSlices(old.Hosts, e.Hosts)
	changed = changed || !slices.Equal(old.Hosts, e.Hosts)
	e.Owned = e.Owned || old.Owned
	changed = changed || e.Owned != old.Owned
	s.Data[id] = e

	return id, changed
}

func (s *State) Remove(id string) error {
	if _, found := s.Data[id]; !found {
		return errEntryNotFound
	}
	delete(s.Data, id)

	return nil
}

func (s *State) Get(id string) (*Entry, error) {
	e, found := s.Data[id]
	if !found {
		return nil, errEntryNotFound
	}
	return &e, nil
}

func (s *State) Find(value string) (Entry, bool) {
	id := getID(value)
	e, found := s.Data[id]
	return e, found
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

func (s *State) Capture(flag string, f entry.Flag) {
	s.Flags[flag] = f
}

func (s *State) Clear() {
	clear(s.Data)
	clear(s.Flags)
}

func getID(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])[:12]
}

func mergeSlices[T cmp.Ordered](a []T, b []T) []T {
	c := append(a, b...)
	slices.Sort(c)
	return slices.Compact(c)
}
