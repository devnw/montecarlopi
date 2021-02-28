// Copyright © 2019 Developer Network, LLC
//
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.

package montecarlopi

import (
	"atomizer.io/engine"
	"devnw.com/alog"
)

func init() {
	// Register the monte carlo atoms
	err := engine.Register(&MonteCarlo{})
	if err != nil {
		alog.Error(err)
	}

	err = engine.Register(&Toss{})
	if err != nil {
		alog.Error(err)
	}
}
