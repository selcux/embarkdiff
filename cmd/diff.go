package cmd

import (
	"errors"

	"github.com/selcux/embarkdiff/diff"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compares given directories",
	Long:  `Compares given contents of the given directories`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res := diff.NewResource()
		err := res.Load()
		if err != nil {
			return err
		}

		if !res.Validate() {
			return errors.New("`source` and `target` are required")
		}

		sourceChan, err := diff.ExecuteChecksum(res.Source())
		if err != nil {
			return err
		}

		targetChan, err := diff.ExecuteChecksum(res.Target())
		if err != nil {
			return err
		}

		diff.Compare(&diff.DirWithChannel{sourceChan, res.Source()},
			&diff.DirWithChannel{targetChan, res.Target()})

		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
