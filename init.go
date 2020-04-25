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
