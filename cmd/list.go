package cmd

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"loot/internal/entry"
	"loot/internal/state"

	"github.com/spf13/cobra"
)

var (
	listAll   bool
	listFlags bool
	listTags  []string
	listHosts []string
)

var listCmd = &cobra.Command{
	Use:     "list [filter]",
	Aliases: []string{"ls"},
	Short:   "List entries in the loot file",
	Run: func(cmd *cobra.Command, args []string) {
		s, _ := loadLootFile()
		if !listFlags {
			printEntries(s, entry.Filter{
				ID:    args,
				Tags:  listTags,
				Hosts: listHosts,
			})
		}
		if listAll {
			fmt.Println()
		}
		if listFlags || listAll {
			printFlags(s)
		}
	},
	ValidArgsFunction: idCompletion,
}

func printEntries(s *state.State, filter entry.Filter) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	header := []string{"ID", "VALUE", "TAGS", "HOSTS", "COMMENT"}
	fmt.Fprintln(w, strings.Join(header, "\t"))

	data := s.Filter(filter)
	for _, k := range slices.Sorted(maps.Keys(data)) {
		v := s.Data[k]
		fields := []string{
			k,
			truncate(v.Value),
			truncate(strings.Join(v.Tags, ", ")),
			truncate(strings.Join(v.Hosts, ", ")),
			truncate(v.Comment),
		}
		fmt.Fprintln(w, strings.Join(fields, "\t"))
	}
	w.Flush()
}

func printFlags(s *state.State) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	header := []string{"FLAG", "TYPE", "HOST"}
	fmt.Fprintln(w, strings.Join(header, "\t"))

	for id, flag := range s.Flags {
		fmt.Fprintln(w, strings.Join([]string{id, flag.Owner, flag.Host}, "\t"))
	}
	w.Flush()
}

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "List both entries and flags")
	listCmd.Flags().BoolVarP(&listFlags, "flags", "f", false, "Only list captured flags")
	listCmd.Flags().
		StringSliceVarP(&listTags, "tag", "t", []string{}, "Only display entries matching given tags")
	listCmd.Flags().
		StringSliceVarP(&listHosts, "host", "H", []string{}, "Only display entries matching given hosts")
	rootCmd.AddCommand(listCmd)
}
