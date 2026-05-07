package resistor

import (
	"fmt"
)

func ExampleNearestStandard() {
	
	v, _ := NearestStandard(500, E24, RoundNearest)
	fmt.Println(v)

	v, _ = NearestStandard(500, E12, RoundUp)
	fmt.Println(v)

	// Output:
	// 510
	// 560
}