package tools

import "math"

func tokenzi(str string) int {
	a := len([]byte(str))
	b := float64(a) * 1.3
	return int(math.Ceil(b))

}
