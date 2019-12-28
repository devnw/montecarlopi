package montecarlopi

import (
	"context"

	"github.com/benjivesterby/alog"
	"github.com/benjivesterby/atomizer"
)

func init() {
	ctx := context.Background()

	// Register the monte carlo atoms
	if err := atomizer.Register(ctx, &MonteCarlo{}); err == nil {
		if err = atomizer.Register(ctx, &Toss{}); err != nil {
			alog.Error(err)
		}
	} else {
		alog.Error(err)
	}
}
