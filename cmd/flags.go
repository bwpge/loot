package cmd

import (
	"github.com/spf13/cobra"
)

var flagsCmd = &cobra.Command{
	Use:   "flags",
	Short: "List flags in the loot file (same as list --flags)",
	Run: func(cmd *cobra.Command, args []string) {
		s, _ := loadLootFile()
		printFlags(s)
	},
	ValidArgsFunction: cobra.NoFileCompletions,
}

func init() {
	rootCmd.AddCommand(flagsCmd)
}
