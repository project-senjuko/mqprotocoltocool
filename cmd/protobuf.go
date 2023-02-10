////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2023-present qianjunakasumi <i@qianjunakasumi.ren>                                                    /
//                            project-senjuko/mqprotocoltocool Contributors                                            /
//                                                                                                                     /
//           https://github.com/qianjunakasumi                                                                         /
//           https://github.com/project-senjuko/mqprotocoltocool/graphs/contributors                                   /
//                                                                                                                     /
//   This Source Code Form is subject to the terms of the Mozilla Public                                               /
//   License, v. 2.0. If a copy of the MPL was not distributed with this                                               /
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.                                                          /
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// protobufCmd represents the protobuf command
var protobufCmd = &cobra.Command{
	Use:   "protobuf",
	Short: "Convert Java file to proto file",
	Long:  `Convert Java file to proto file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("protobuf called")
	},
}

func init() {
	rootCmd.AddCommand(protobufCmd)
	protobufCmd.PersistentFlags().String("input", "", "Input directory")
	protobufCmd.PersistentFlags().String("output", "", "Output directory")
}
