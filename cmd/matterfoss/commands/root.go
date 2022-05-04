// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package commands

import (
	"github.com/spf13/cobra"
)

type Command = cobra.Command

func Run(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}

var RootCmd = &cobra.Command{
	Use:   "matterfoss",
	Short: "Open source, self-hosted Slack-alternative",
	Long:  `Matterfoss offers workplace messaging across web, PC and phones with archiving, search and integration with your existing systems. Documentation available at https://docs.matterfoss.org`,
}

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "", "Configuration file to use.")
}
