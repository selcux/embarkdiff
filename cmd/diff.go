/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/selcux/embarkdiff/diff"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compares given directories",
	Long:  `Compares given contents of the given directories`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := diff.Read()
		if err != nil {
			return err
		}

		if !res.Validate() {
			return errors.New("`source` and `target` are required")
		}

		dirInfo, err := diff.NewDirInfo(res.Source)
		if err != nil {
			return err
		}

		for k, v := range dirInfo.Files {
			fmt.Println(k, hex.EncodeToString(v))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
