package resistor_test

import (
	"fmt"
	"github.com/sss7526/resistor"
)

func ExampleNearestStandard() {

	v, _ := resistor.NearestStandard(500, resistor.E24, resistor.RoundNearest)
	fmt.Println(v)

	v, _ = resistor.NearestStandard(500, resistor.E12, resistor.RoundUp)
	fmt.Println(v)

	// Output:
	// 510
	// 560
}
