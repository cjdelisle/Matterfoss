// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"strings"

	"github.com/cjdelisle/matterfoss-server/v6/app"
	"github.com/cjdelisle/matterfoss-server/v6/app/request"
	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/cjdelisle/matterfoss-server/v6/shared/i18n"
)

type JoinProvider struct {
}

const (
	CmdJoin = "join"
)

func init() {
	app.RegisterCommandProvider(&JoinProvider{})
}

func (*JoinProvider) GetTrigger() string {
	return CmdJoin
}

func (*JoinProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdJoin,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_join.desc"),
		AutoCompleteHint: T("api.command_join.hint"),
		DisplayName:      T("api.command_join.name"),
	}
}

func (*JoinProvider) DoCommand(a *app.App, c *request.Context, args *model.CommandArgs, message string) *model.CommandResponse {
	channelName := strings.ToLower(message)

	if strings.HasPrefix(message, "~") {
		channelName = message[1:]
	}

	channel, err := a.Srv().Store.Channel().GetByName(args.TeamId, channelName, true)
	if err != nil {
		return &model.CommandResponse{Text: args.T("api.command_join.list.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}

	if channel.Name != channelName {
		return &model.CommandResponse{ResponseType: model.CommandResponseTypeEphemeral, Text: args.T("api.command_join.missing.app_error")}
	}

	switch channel.Type {
	case model.ChannelTypeOpen:
		if !a.HasPermissionToChannel(args.UserId, channel.Id, model.PermissionJoinPublicChannels) {
			return &model.CommandResponse{Text: args.T("api.command_join.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
		}
	case model.ChannelTypePrivate:
		if !a.HasPermissionToChannel(args.UserId, channel.Id, model.PermissionReadChannel) {
			return &model.CommandResponse{Text: args.T("api.command_join.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
		}
	default:
		return &model.CommandResponse{Text: args.T("api.command_join.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}

	if appErr := a.JoinChannel(c, channel, args.UserId); appErr != nil {
		return &model.CommandResponse{Text: args.T("api.command_join.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}

	team, appErr := a.GetTeam(channel.TeamId)
	if appErr != nil {
		return &model.CommandResponse{Text: args.T("api.command_join.fail.app_error"), ResponseType: model.CommandResponseTypeEphemeral}
	}

	return &model.CommandResponse{GotoLocation: args.SiteURL + "/" + team.Name + "/channels/" + channel.Name}
}
