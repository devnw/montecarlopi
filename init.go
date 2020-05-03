// Copyright © 2019 Developer Network, LLC
//
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.

package montecarlopi

import (
	"github.com/devnw/alog"
	"github.com/devnw/atomizer"
)

func init() {
	// Register the monte carlo atoms
	if err := atomizer.Register(&MonteCarlo{}); err == nil {
		if err = atomizer.Register(&Toss{}); err != nil {
			alog.Error(err)
		}
	} else {
		alog.Error(err)
	}
}
