package cmd

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"loot/internal/ui"
	"loot/loot"

	"github.com/spf13/cobra"
)

var listLong bool

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List entries in the loot file",
	Run: func(cmd *cobra.Command, args []string) {
		s, _ := loadLootFile()
		if len(s.Data) == 0 {
			fmt.Fprintln(os.Stderr, "no entries to show")
			return
		}

		if listLong {
			listPrintLong(s.Data)
		} else {
			listPrintShort(s.Data)
		}
	},
	ValidArgsFunction: idCompletion,
}

func init() {
	listCmd.Flags().
		BoolVarP(&listLong, "long", "l", false, "Display entries in detailed format (exclusive with -r)")
	rootCmd.AddCommand(listCmd)
}

func listPrintShort(data map[string]loot.Entry) {
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

func listPrintLong(data map[string]loot.Entry) {
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
