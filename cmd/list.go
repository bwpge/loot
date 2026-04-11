package cmd

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List entries in the loot file",
	Run: func(cmd *cobra.Command, args []string) {
		err := cobra.NoArgs(cmd, args)
		if err != nil {
			bail(err)
		}

		s, _ := loadLootFile()
		if len(s.Data) == 0 {
			fmt.Fprintln(os.Stderr, "no entries to show")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
		header := []string{"ID", "VALUE", "COMMENT", "TAGS", "HOSTS"}
		fmt.Fprintln(w, strings.Join(header, "\t"))

		for _, k := range slices.Sorted(maps.Keys(s.Data)) {
			v := s.Data[k]
			fields := []string{
				k,
				truncate(v.Value),
				truncate(v.Comment),
				truncate(strings.Join(v.Tags, ", ")),
				truncate(strings.Join(v.Hosts, ", ")),
			}
			fmt.Fprintln(w, strings.Join(fields, "\t"))
		}

		w.Flush()
	},
	ValidArgsFunction: idCompletion,
}

func init() {
	rootCmd.AddCommand(listCmd)
}
