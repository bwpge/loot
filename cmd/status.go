package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the loot file status",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		s, f := loadLootFile()
		fmt.Println("loot file: ", f)
		fmt.Println("entries:   ", len(s.Data))
		fmt.Println("flags:     ", len(s.Flags))
	},
	ValidArgsFunction: cobra.NoFileCompletions,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
