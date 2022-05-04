// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"github.com/cjdelisle/matterfoss-server/v6/app"
	"github.com/cjdelisle/matterfoss-server/v6/app/request"
	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/cjdelisle/matterfoss-server/v6/shared/i18n"
)

type ShortcutsProvider struct {
}

const (
	CmdShortcuts = "shortcuts"
)

func init() {
	app.RegisterCommandProvider(&ShortcutsProvider{})
}

func (*ShortcutsProvider) GetTrigger() string {
	return CmdShortcuts
}

func (*ShortcutsProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdShortcuts,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_shortcuts.desc"),
		AutoCompleteHint: "",
		DisplayName:      T("api.command_shortcuts.name"),
	}
}

func (*ShortcutsProvider) DoCommand(a *app.App, c *request.Context, args *model.CommandArgs, message string) *model.CommandResponse {
	// This command is handled client-side and shouldn't hit the server.
	return &model.CommandResponse{
		Text:         args.T("api.command_shortcuts.unsupported.app_error"),
		ResponseType: model.CommandResponseTypeEphemeral,
	}
}
