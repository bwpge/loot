package loot

import (
	"loot/internal/config"
	"loot/internal/entry"
	"loot/internal/state"
)

type (
	Entry  = entry.Entry
	Filter = entry.Filter
	State  = state.State
)

func Config() *config.Config {
	return config.Get()
}

func LoadConfig(path string) error {
	return config.Load(path)
}

func NewState() *State {
	return state.New()
}

func LoadState(path string) (*State, error) {
	return state.LoadState(path)
}
