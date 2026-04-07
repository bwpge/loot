package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var showNoNewline bool

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Display an entry value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s, _ := loadLootFile()
		id, err := s.FindID(args[0])
		if err != nil {
			bail(err)
		}
		e, err := s.Get(id)
		if err != nil {
			bail(err)
		}

		fmt.Print(e.Value)
		if !showNoNewline {
			fmt.Println()
		}
	},
	ValidArgsFunction: idCompletion,
}

func init() {
	showCmd.Flags().
		BoolVarP(&showNoNewline, "no-newline", "n", false, "Do not display a trailing newline")
	rootCmd.AddCommand(showCmd)
}
