// Copyright 2024 The Erigon Authors
// This file is part of Erigon.
//
// Erigon is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Erigon is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Erigon. If not, see <http://www.gnu.org/licenses/>.

//go:build windows

package debug

import (
	"io"
	"os"
	"os/signal"

	_debug "github.com/erigontech/erigon-lib/common/debug"
	"github.com/erigontech/erigon-lib/log/v3"
)

func ListenSignals(stack io.Closer, logger log.Logger) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	_debug.GetSigC(&sigc)
	defer signal.Stop(sigc)

	<-sigc
	logger.Info("Got interrupt, shutting down...")
	if stack != nil {
		go stack.Close()
	}
	for i := 10; i > 0; i-- {
		<-sigc
		if i > 1 {
			logger.Warn("Already shutting down, interrupt more to panic.", "times", i-1)
		}
	}
	Exit() // ensure trace and CPU profile data is flushed.
	LoudPanic("boom")
}
