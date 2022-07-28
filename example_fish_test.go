package lift_test

import (
	"fmt"

	"github.com/AndrewHarrisSPU/lift"
)

// This example demonstrates a [Map] with [Sym] values.
// The [Map] is a dispatch table; a [Sym] may be looked up, and a related function returned.
func Example_b_fish() {
	// "recievers"
	type fish string

	type octopus struct {
		arms int
	}

	// "methods"
	fishSlap := func(sym lift.Sym) string {
		f := lift.MustUnwrap[fish](sym)
		return fmt.Sprintf("%sslap", f)
	}

	octopusSlap := func(sym lift.Sym) (slap string) {
		o := lift.MustUnwrap[octopus](sym)
		for i := 0; i < o.arms; i++ {
			slap += "octoslap"
		}
		return
	}

	// a dispatch map
	marineLifeSlap := lift.NewMap[func(lift.Sym) string](
		lift.Def[fish](fishSlap),
		lift.Def[octopus](octopusSlap),
	)

	// the demonstration:
	symbols := []lift.Sym{
		lift.Wrap(fish("trout")),
		lift.Wrap(fish("salmon")),
		lift.Wrap(octopus{8}),
	}

	for _, sym := range symbols {
		slap, _ := lift.LoadSym(marineLifeSlap, sym)
		fmt.Printf("boom! %s!\n", slap(sym))
	}
	// Output:
	// boom! troutslap!
	// boom! salmonslap!
	// boom! octoslapoctoslapoctoslapoctoslapoctoslapoctoslapoctoslapoctoslap!
}
