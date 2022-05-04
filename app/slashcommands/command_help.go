// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"github.com/cjdelisle/matterfoss-server/v6/app"
	"github.com/cjdelisle/matterfoss-server/v6/app/request"
	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/cjdelisle/matterfoss-server/v6/shared/i18n"
)

type HelpProvider struct {
}

const (
	CmdHelp = "help"
)

func init() {
	app.RegisterCommandProvider(&HelpProvider{})
}

func (h *HelpProvider) GetTrigger() string {
	return CmdHelp
}

func (h *HelpProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdHelp,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_help.desc"),
		DisplayName:      T("api.command_help.name"),
	}
}

func (h *HelpProvider) DoCommand(a *app.App, c *request.Context, args *model.CommandArgs, message string) *model.CommandResponse {
	helpLink := *a.Config().SupportSettings.HelpLink

	if helpLink == "" {
		helpLink = model.SupportSettingsDefaultHelpLink
	}

	return &model.CommandResponse{GotoLocation: helpLink}
}
