package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/bwpge/loot/internal/state"
	"github.com/bwpge/loot/internal/ui"
	"github.com/spf13/cobra"
)

const lootFileName = "loot.json"

var (
	Version  = "0.1.0"
	Commit   = ""
	lootFile string
)

func buildVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return Version
	}

	if Commit == "" {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				Commit = setting.Value
				if len(Commit) > 7 {
					Commit = Commit[:7]
				}
				break
			}
		}
	}

	if Commit != "" {
		return Version + " (" + Commit + ")"
	}

	return Version
}

var rootCmd = &cobra.Command{
	Use:     "loot",
	Version: buildVersion(),
	Short:   "Tool for storing and organizing loot during offensive security operations",
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
	rootCmd.SetErrPrefix(ui.ColorErr.Sprint("error:"))
	rootCmd.PersistentFlags().
		StringVarP(&lootFile, "loot-file", "L", "", "explicit loot file path to use")
	rootCmd.SetVersionTemplate(
		`{{printf "%s %s" .Name .Version}}`,
	)
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

func loadLootFile() (*state.State, string) {
	f, err := findLootFile()
	if err != nil {
		bail(err)
	}
	if _, err = os.Stat(f); errors.Is(err, os.ErrNotExist) {
		bail("no loot file found (run", ui.Cli("loot init"), "to create)")
	}
	s, err := state.Load(f)
	if err != nil {
		bail(err)
	}
	return s, f
}

func loadLootFileNoErr() *state.State {
	f, err := findLootFile()
	if err != nil {
		return nil
	}
	if _, err = os.Stat(f); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	s, err := state.Load(f)
	if err != nil {
		return nil
	}
	return s
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
