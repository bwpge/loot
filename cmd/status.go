package cmd

import (
	"encoding/json"
	"fmt"

	"loot/loot"

	"github.com/spf13/cobra"
)

var statusDump bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the loot file status",
	Run: func(cmd *cobra.Command, args []string) {
		err := cobra.NoArgs(cmd, args)
		if err != nil {
			bail(err)
		}
		s, f := loadLootFile()

		if statusDump {
			printJSON(&s)
			return
		}

		fmt.Println("loot file:     ", f)
		fmt.Println("entries:       ", len(s.Data))
		fmt.Println("unique values: ", len(s.Hashes))
		fmt.Println("config:")
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
