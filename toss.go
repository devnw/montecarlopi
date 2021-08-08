// Copyright © 2019 Developer Network, LLC
//
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.

package montecarlopi

import (
	"context"
	"encoding/json"
	"math/rand"

	"go.atomizer.io/engine"
)

// Ensure compile time adherence to interface
var _ engine.Atom = (*Toss)(nil)

// Toss randomly tosses darts
type Toss struct {
	Value int `json:"value"`
}

// Process test method
func (t *Toss) Process(
	ctx context.Context,
	conductor engine.Conductor,
	electron *engine.Electron,
) (result []byte, err error) {
	// Step 1: Generate my Random X/Y Coordinates

	/* #nosec G404 */
	x := rand.Float64()

	/* #nosec G404 */
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
