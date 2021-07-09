// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package scheduler

import (
	"github.com/cjdelisle/matterfoss-server/v5/app"
	tjobs "github.com/cjdelisle/matterfoss-server/v5/jobs/interfaces"
)

type PluginsJobInterfaceImpl struct {
	App *app.App
}

func init() {
	app.RegisterJobsPluginsJobInterface(func(s *app.Server) tjobs.PluginsJobInterface {
		a := app.New(app.ServerConnector(s))
		return &PluginsJobInterfaceImpl{a}
	})
}
