// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/cjdelisle/matterfoss-server/v5/utils/fileutils"
)

func main() {
	// Print angry message to use matterfoss command directly
	fmt.Println(`
------------------------------------ ERROR ------------------------------------------------
The platform binary has been deprecated, please switch to using the new matterfoss binary.
The platform binary will be removed in a future version.
-------------------------------------------------------------------------------------------
	`)

	// Execve the real MM binary
	args := os.Args
	args[0] = "matterfoss"
	args = append(args, "--platform")

	realMatterfoss := fileutils.FindFile("matterfoss")
	if realMatterfoss == "" {
		realMatterfoss = fileutils.FindFile("bin/matterfoss")
	}

	if realMatterfoss == "" {
		fmt.Println("Could not start Matterfoss, use the matterfoss command directly: failed to find matterfoss")
	} else if err := syscall.Exec(realMatterfoss, args, nil); err != nil {
		fmt.Printf("Could not start Matterfoss, use the matterfoss command directly: %s\n", err.Error())
	}
}
