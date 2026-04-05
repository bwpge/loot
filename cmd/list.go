package cmd

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"loot/internal"
	"loot/internal/ui"

	"github.com/spf13/cobra"
)

var (
	listLong bool
	listTag  string
	listHost string
)

var listCmd = &cobra.Command{
	Use:     "list [filter]",
	Aliases: []string{"ls"},
	Short:   "List entries in the loot file",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s, _ := loadLootFile()
		if len(s.Data) == 0 {
			fmt.Fprintln(os.Stderr, "no entries to show")
			return
		}

		filter := ""
		if len(args) == 1 {
			filter = args[0]
		}

		filtered := make(map[string]internal.Entry)
		for k, v := range s.Data {
			if (filter == "" || strings.HasPrefix(k, filter)) &&
				(listTag == "" || slices.Contains(v.Tags, listTag)) &&
				(listHost == "" || slices.Contains(v.Hosts, listHost)) {
				filtered[k] = v
			}
		}
		skipped := len(s.Data) - len(filtered)

		if listLong {
			listPrintLong(filtered)
		} else {
			listPrintShort(filtered)
		}

		if skipped > 0 {
			fmt.Printf("\nentries not shown: %d", skipped)
		}
	},
	ValidArgsFunction: idCompletion,
}

func init() {
	listCmd.Flags().BoolVarP(&listLong, "long", "l", false, "Display entries in detailed format")
	listCmd.Flags().StringVarP(&listTag, "tag", "t", "", "Only display entries with this tag")
	listCmd.Flags().StringVarP(&listHost, "host", "H", "", "Only display entries with this host")
	rootCmd.AddCommand(listCmd)
}

func listPrintShort(data map[string]internal.Entry) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	header := []string{"ID", "VALUE", "COMMENT", "TAGS", "HOSTS"}
	fmt.Fprintln(w, strings.Join(header, "\t"))

	for _, k := range slices.Sorted(maps.Keys(data)) {
		v := data[k]
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
}

func listPrintLong(data map[string]internal.Entry) {
	for _, k := range slices.Sorted(maps.Keys(data)) {
		v := data[k]
		value := v.Value
		if strings.Contains(value, "\n") {
			value = "\n" + value
		}

		fmt.Println(ui.ID(k))
		fmt.Printf("  %s %s\n", ui.Header("Value:"), ui.Value(value))
		fmt.Printf("  %s %s\n", ui.Header("Comment:"), v.Comment)
		fmt.Printf("  %s %s\n", ui.Header("Tags:"), strings.Join(v.Tags, ", "))
		fmt.Printf("  %s %s\n", ui.Header("Hosts:"), strings.Join(v.Hosts, ", "))
		fmt.Println()
	}
}
