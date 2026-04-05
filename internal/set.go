package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type HashSet map[string]struct{}

func (s *HashSet) Add(value string) {
	(*s)[hashStr(value)] = struct{}{}
}

func (s *HashSet) Remove(value string) {
	delete(*s, hashStr(value))
}

func (s *HashSet) Contains(value string) bool {
	_, found := (*s)[hashStr(value)]
	return found
}

func (s HashSet) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	return json.Marshal(keys)
}

func (s *HashSet) UnmarshalJSON(data []byte) error {
	var items []string
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	*s = make(HashSet)
	for _, item := range items {
		(*s)[item] = struct{}{}
	}
	return nil
}

func hashStr(value string) string {
	h := sha256.New()
	h.Write([]byte(value))
	sha256Hash := h.Sum(nil)
	return hex.EncodeToString(sha256Hash)
}
