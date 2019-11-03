package montecarlopi

import "github.com/benjivesterby/atomizer"

func init() {
	// Register the monte carlo atoms
	atomizer.Register(nil, "montecarlo", MonteCarlo{})
	atomizer.Register(nil, "toss", Toss{})
}
