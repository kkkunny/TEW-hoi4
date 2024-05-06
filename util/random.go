package util

import (
	"math/rand"
	"time"

	"golang.org/x/exp/constraints"
)

func RandomInt[Int constraints.Integer](b, e Int) Int {
	return b + Int(rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(int64(e-b)))
}
