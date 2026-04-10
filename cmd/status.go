package cmd

import (
	"encoding/json"
	"fmt"

	"loot/internal/ui"
	"loot/loot"

	"github.com/spf13/cobra"
)

var statusDump bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the loot file status",
	Run: func(cmd *cobra.Command, args []string) {
		s, f := loadLootFile()

		if statusDump {
			printJSON(&s)
			return
		}

		fmt.Println(ui.Header("Loot file:"), f)
		fmt.Println(ui.Header("Entries:"), len(s.Data))
		fmt.Println(ui.Header("Unique values:"), len(s.Hashes))
		fmt.Println(ui.Header("Config:"))
		printJSON(loot.Config())
	},
}

func printJSON(v any) {
	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		bail(err)
	}
	fmt.Println(string(j))
}

func init() {
	statusCmd.Flags().
		BoolVarP(&statusDump, "dump", "d", false, "Display the entire loot file as serialized JSON")
	rootCmd.AddCommand(statusCmd)
}
