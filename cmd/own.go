package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ownCmd = &cobra.Command{
	Use:     "own id [id...]",
	Aliases: []string{"pwn"},
	Short:   "Toggle owned status of one or more entries",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s, f := loadLootFile()

		for _, arg := range args {
			id, err := s.FindID(arg)
			if err != nil {
				bail(err)
			}
			e, err := s.Get(id)
			if err != nil {
				bail(err)
			}
			e.Owned = !e.Owned
			s.Data[id] = *e

			status := "unowned"
			if e.Owned {
				status = "owned"
			}
			fmt.Println(status, id, "->", truncate(e.Value))
		}

		s.Save(f)
	},
	ValidArgsFunction: completeID,
}

func init() {
	rootCmd.AddCommand(ownCmd)
}
