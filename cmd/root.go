package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"loot/internal"
	"loot/internal/ui"

	"github.com/spf13/cobra"
)

const lootFileName = "loot.json"

var lootFile string

var rootCmd = &cobra.Command{
	Use:   "loot",
	Short: "Tool for storing and organizing loot during offensive security operations",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if lootFile == "" {
			lootFile, err = findLootFile()
		} else {
			lootFile, err = filepath.Abs(lootFile)
		}
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().
		StringVarP(&lootFile, "loot-file", "L", "", "explicit loot file path to use")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func findLootFile() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	curr := cwd

	for {
		path := filepath.Join(curr, lootFileName)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		parentDir := filepath.Dir(curr)
		if parentDir == curr {
			break
		}
		curr = parentDir
	}

	return filepath.Join(cwd, lootFileName), nil
}

func loadLootFile() (*internal.State, string) {
	f, err := findLootFile()
	if err != nil {
		bail(err)
	}
	if _, err = os.Stat(f); errors.Is(err, os.ErrNotExist) {
		bail("no loot file found (run", ui.Cli("loot init"), "to create)")
	}
	s, err := internal.LoadState(f)
	if err != nil {
		bail(err)
	}
	return s, f
}

func loadLootFileNoErr() *internal.State {
	f, err := findLootFile()
	if err != nil {
		return nil
	}
	if _, err = os.Stat(f); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	s, err := internal.LoadState(f)
	if err != nil {
		return nil
	}
	return s
}

func idCompletion(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	values := []string{}
	s := loadLootFileNoErr()
	if s == nil {
		return values, cobra.ShellCompDirectiveNoFileComp
	}

	for k, v := range s.Data {
		desc := truncate(v.Value)
		if v.Comment != "" {
			desc += " (" + v.Comment + ")"
		}
		values = append(values, k+"\t"+desc)
	}

	return values, cobra.ShellCompDirectiveNoFileComp
}

func truncate(v string) string {
	s, _, _ := strings.Cut(v, "\n")
	if len(s) > 50 {
		s = s[:47] + "..."
	}
	return s
}

func listString(v []string) string {
	return "[" + strings.Join(v, ", ") + "]"
}

func bail(a ...any) {
	ui.ColorErr.Fprint(os.Stderr, "error: ")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func warn(a ...any) {
	ui.ColorWarn.Fprint(os.Stderr, "warning: ")
	fmt.Fprintln(os.Stderr, a...)
}
