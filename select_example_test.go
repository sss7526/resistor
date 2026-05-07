package resistor_test

import (
	"fmt"
	"github.com/sss7526/resistor"
)

func ExampleSelectStandardResistor() {
	req := resistor.SelectionRequest{
		Resistance: 487,
	}

	res, _ := resistor.SelectStandardResistor(req)

	fmt.Println(res.SelectedResistance)
	fmt.Println(res.Bands)

	// Output:
	// 470
	// [yellow violet brown gold]
}
