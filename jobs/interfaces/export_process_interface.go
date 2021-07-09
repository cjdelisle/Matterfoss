// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package interfaces

import (
	"github.com/cjdelisle/matterfoss-server/v5/model"
)

type ExportProcessInterface interface {
	MakeWorker() model.Worker
}
