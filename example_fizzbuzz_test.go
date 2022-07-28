package lift_test

import (
	"fmt"

	"github.com/AndrewHarrisSPU/lift"
)

// This example demonstrates lifting a function.
// If the lifted function is passed a Sym that doesn't match the desired type,
// the return is zero valued.
func Example_1FizzBuzz() {
	sayFizz := liftFizzBuzz( func(_ fizz) string {
		return "fizz"
	})
	sayBuzz := liftFizzBuzz( func(_ buzz) string {
		return "buzz"
	})

	for i := 1; i < 31; i++ {
		f, b := parseFizzBuzz( i )
		if res := sayFizz( f ) + sayBuzz( b ); res != "" {
			fmt.Println(i, res)
		}
	}
	// Output:
	// 3 fizz
	// 5 buzz
	// 6 fizz
	// 9 fizz
	// 10 buzz
	// 12 fizz
	// 15 fizzbuzz
	// 18 fizz
	// 20 buzz
	// 21 fizz
	// 24 fizz
	// 25 buzz
	// 27 fizz
	// 30 fizzbuzz
}

type fizz struct{}
type buzz struct{}

func liftFizzBuzz[T any]( fn func(T) string ) func(lift.Sym) string {
	return func( sym lift.Sym ) string {
		if t, ok := lift.Unwrap[T](sym); ok {
			return fn(t)
		}
		return ""
	}
}

func parseFizzBuzz( i int ) (f, b lift.Sym ){
	if i % 3 == 0 {
		f = lift.Wrap( fizz{} )
	}
	if i % 5 == 0 {
		b = lift.Wrap( buzz{} )
	}
	return
}