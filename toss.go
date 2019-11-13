package montecarlopi

import (
	"context"
	"encoding/json"
	"math/rand"

	"github.com/benjivesterby/atomizer"
)

// Toss randomly tosses darts
type Toss struct {
	Value int `json:"value"`
}

// ID test method
func (t *Toss) ID() string {
	return "toss"
}

// Process test method
func (t *Toss) Process(ctx context.Context, conductor atomizer.Conductor, electron *atomizer.Electron) (result []byte, err error) {
	// Step 1: Generate my Random X/Y Coordinates
	x := rand.Float64()
	y := rand.Float64()

	// Step 2: Return my Random X/Y Coordinates
	t.Value = t.dsquared(x, y)

	// Step 3: Marshal the job struct passing the value of the toss back to the requestor
	result, err = json.Marshal(t)

	return result, err
}

func (t *Toss) dsquared(x, y float64) int {
	sq := (x * x) + (y * y)

	in := 0
	if sq <= 1 {
		in = 1
	}

	return in
}

// TODO: Implement for project
// Toss takes in the randomly generated x and y coordinates
// determines if they are less than 1 (or inside the circle)
// and returns a 1 if they are and 0 if they aren't so that
// the results can be added without needing explicit logic
// for incrementation
// int Toss(double x, double y) {
// 	double dsquared = 0;

// 	dsquared = (x*x) + (y*y);

// 	return dsquared <= 1 ? 1 : 0;
// }
