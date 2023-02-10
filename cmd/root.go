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
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mqprotocoltocool",
	Short: "A tool that makes QQ protobuf to standard specification",
	Long: `mqprotocoltocool is an automatic tool that transforms or deobfuscates
QQ Java files into standard specifications that could be used by conch project`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
