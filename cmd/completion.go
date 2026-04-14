package cmd

import (
	"maps"
	"slices"

	"github.com/spf13/cobra"
)

func completeID(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	values := []string{}
	s := loadLootFileNoErr()
	if s == nil {
		return values, cobra.ShellCompDirectiveNoFileComp
	}

	for k, v := range s.Data {
		desc := truncate(v.Value)
		if v.Comment != "" {
			desc += " (" + v.Comment + ")"
		}
		values = append(values, k+"\t"+desc)
	}

	return values, cobra.ShellCompDirectiveNoFileComp
}

func completeFlag(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	values := []string{}
	s := loadLootFileNoErr()
	if s == nil {
		return values, cobra.ShellCompDirectiveNoFileComp
	}

	for k, v := range s.Flags {
		desc := truncate(v.Type)
		if v.Host != "" {
			desc += " (type: " + v.Host + ")"
		}
		values = append(values, k+"\t"+desc)
	}

	return values, cobra.ShellCompDirectiveNoFileComp
}

func completeTag(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	values := []string{
		"username",
		"password",
		"credential",
		"domain",
	}
	s := loadLootFileNoErr()
	if s == nil {
		return values, cobra.ShellCompDirectiveNoFileComp
	}

	set := make(map[string]struct{})
	for _, e := range s.Data {
		for _, t := range e.Tags {
			set[t] = struct{}{}
		}
	}

	return slices.Collect(maps.Keys(set)), cobra.ShellCompDirectiveNoFileComp
}

func completeHost(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	s := loadLootFileNoErr()
	if s == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	set := make(map[string]struct{})
	for _, e := range s.Data {
		for _, h := range e.Hosts {
			set[h] = struct{}{}
		}
	}

	return slices.Collect(maps.Keys(set)), cobra.ShellCompDirectiveNoFileComp
}
