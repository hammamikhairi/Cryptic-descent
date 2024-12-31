package helpers

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/constraints"
)

func Distance(v1, v2 rl.Vector2) float32 {
	dx := v1.X - v2.X
	dy := v1.Y - v2.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// Min returns the smaller of two comparable values
func Min[T constraints.Ordered](val1, val2 T) T {
	if val1 < val2 {
		return val1
	}
	return val2
}
