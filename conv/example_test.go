package conv_test

import (
	"fmt"

	"github.com/AndrewHarrisSPU/lift/conv"
)

// CONVERTER

func Example_converter() {
	// from int to string
	i2s := func(n int) (string, error) {
		return fmt.Sprintf("(int) %d", n), nil
	}
	// from rune to string
	r2s := func(r rune) (string, error) {
		return fmt.Sprintf("(rune) %c", r), nil
	}

	// constructing a new converter
	cv := conv.NewConverter(
		conv.Def(i2s),
	)

	// Storing a conversion
	cv.Store(
		conv.Def(r2s),
	)

	// with destination type as type parameter, and source type inferred:
	a, _ := conv.To[string](cv, 1)
	b, _ := conv.To[string](cv, '1')
	fmt.Println(a)
	fmt.Println(b)

	// converting a type to itself is always defined
	c, _ := conv.To[string](cv, "one")
	fmt.Println(c)
	// Output:
	// (int) 1
	// (rune) 1
	// one
}

// Conversion success or failure can depend on value.
// This is why conv insists on conversion functions with an error return.
func ExampleTo_failure() {
	i2r := func(i int) (rune, error) {
		if 0 <= i && i <= 9 {
			return rune(i + '0'), nil
		}
		return ' ', fmt.Errorf("oops")
	}

	cv := conv.NewConverter(
		conv.Def(i2r),
	)

	if _, err := conv.To[rune](cv, 10); err != nil {
		fmt.Print(err.Error())
	}
	// Output:
	// oops
}

func ExampleConverter_Delete() {
	bit2n := func(bit bool) (n int, err error) {
		if bit {
			n = 1
		}
		return
	}

	cv := conv.NewConverter(
		conv.Def(bit2n),
	)

	cv.Delete(
		conv.Def[bool, int](nil),
	)

	if _, err := conv.To[int](cv, true); err != nil {
		fmt.Println(err.Error())
	}
	// Output:
	// Conversion not found: bool->int
}

func ExampleLookup() {
	type rgb struct{ r, g, b uint8 }

	rgb2hex := func(color rgb) (string, error) {
		return fmt.Sprintf("#%2x%2x%2x", color.r, color.g, color.b), nil
	}

	cv := conv.NewConverter(
		conv.Def(rgb2hex),
	)

	palette := []rgb{
		rgb{0x29, 0xbe, 0xb0},
		rgb{0xe0, 0xb0, 0xff},
	}

	// pretending we're in a scope far, far away ...
	convFunc, _ := conv.Lookup[rgb, string](cv)
	for _, color := range palette {
		hex, _ := convFunc(color)
		fmt.Println(hex)
	}
	// Output:
	// #29beb0
	// #e0b0ff
}