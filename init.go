package montecarlopi

import (
	"context"

	"github.com/benjivesterby/alog"
	"github.com/benjivesterby/atomizer"
)

func init() {
	ctx := context.Background()

	// Register the monte carlo atoms
	if err := atomizer.Register(ctx, "montecarlo", &MonteCarlo{}); err == nil {
		if err = atomizer.Register(ctx, "toss", &Toss{}); err != nil {
			alog.Error(err)
		}
	} else {
		alog.Error(err)
	}
}
