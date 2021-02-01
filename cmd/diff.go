package cmd

import (
	"context"
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

		errCh := make(chan error)
		defer close(errCh)

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		sourceChan := diff.ExecuteChecksum(res.Source(), errCh)
		targetChan := diff.ExecuteChecksum(res.Target(), errCh)

		fileOps := diff.Compare(sourceChan, targetChan)
		fileOps = fileOps.Transform()
		for _, x := range fileOps {
			diff.PrintOperation(x.Path, x.Operation)
		}

		cancel()

		select {
		case <-ctx.Done():
			break
		case err = <-errCh:
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
