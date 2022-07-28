package lift_test

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/AndrewHarrisSPU/lift"
)

func ExampleT() {
	r := lift.T[rune]()
	b := lift.T[byte]()

	fmt.Println(r == b)
	// Output:
	// false
}

func ExampleT_alias() {
	i := lift.T[int32]()
	r := lift.T[rune]()

	// because rune is an alias of int32, the output is "true"
	fmt.Println(i == r)

	// Output:
	// true
}

func ExampleTypeOf() {
	t := lift.T[string]()

	fmt.Println(t == lift.TypeOf(""))
	fmt.Println(t == lift.TypeOf("twine"))
	fmt.Println(t == lift.TypeOf([]rune{}))
	// Output:
	// true
	// true
	// false
}

func ExampleEnumIs() {
	a := lift.Wrap(false)
	fmt.Println(lift.EnumIs[bool](a))
	// Output:
	// true
}

// WRAP, UNWRAP

func ExampleWrap() {
	π := lift.Wrap(math.Pi)
	pi, _ := lift.Unwrap[float64](π)
	fmt.Println(pi == math.Pi)
	// Output:
	// true
}

func ExampleWrap_equivalence() {
	sym := lift.Wrap('?')

	// the wrapped symbol is not == to rune's type enumeration symbol
	fmt.Println(lift.T[rune]() == sym)

	// EnumIs[rune] evaluates true, however
	fmt.Println(lift.EnumIs[rune](sym))
	// Output:
	// false
	// true
}

func ExampleUnwrap() {
	bits := byte(0)
	sym := lift.Wrap(bits)

	if _, ok := lift.Unwrap[byte](sym); ok {
		fmt.Print("got a byte")
	}
	// Output:
	// got a byte
}

// This works. A pointer to a [strings.Builder] is an [io.Writer].
func ExampleUnwrapAs_interface() {
	var b strings.Builder
	sym := lift.Wrap(&b)

	if _, ok := lift.UnwrapAs[io.Writer](sym); ok {
		fmt.Printf("unwrapped a Writer!")
	}
	// Output:
	// unwrapped a Writer!
}

// This won't work. [UnwrapAs] performs type assertion, not conversion.
func ExampleUnwrapAs_noConversion() {
	type bit bool
	sym := lift.Wrap(bit(true))

	if bit, ok := lift.UnwrapAs[bool](sym); ok {
		fmt.Printf("%v\n", bit)
	} else {
		fmt.Println("?")
	}
	// Output:
	// ?
}

// MAP

func ExampleMap() {
	// Map creation
	m := lift.NewMap[string](
		lift.Def[int]("red"),
		lift.Def[float64]("blue"),
	)

	type vector2d [2]float64

	// Storing to a Map
	m.Store(
		lift.Def[complex128]("C"),
		lift.Def[vector2d]("R2"),
	)

	// Loading from a Map
	intColor, _ := lift.LoadTypeOf(m, 1)
	floatColor, _ := lift.LoadSym(m, lift.TypeOf(1.0))
	complexDomain, _ := lift.Load[complex128](m)
	vector2dDomain, _ := lift.Load[vector2d](m)

	fmt.Printf("ints are %s, float64s are %s\n", intColor, floatColor)
	fmt.Printf("complex128s live in %s (not %s)", complexDomain, vector2dDomain)
	// Output:
	// ints are red, float64s are blue
	// complex128s live in C (not R2)
}

func ExampleDef() {
	type text struct{}
	type mauve struct{}

	crayons := lift.NewMap[int](
		lift.Def[text](0x222222),
		lift.Def[mauve](0xa17188),
	)

	headline, _ := lift.Load[mauve](crayons)
	paragraph, _ := lift.Load[text](crayons)

	fmt.Printf("#%06x, #%06x", headline, paragraph)
	// Output:
	// #a17188, #222222
}

func ExampleDefSym() {
	items := lift.NewMap[string](
		lift.DefSym(lift.Any, "could be anything"),
	)

	fish := lift.Wrap(any("it's a fish!"))
	anything, _ := lift.LoadSym(items, fish)

	fmt.Printf(anything)
	// Output:
	// could be anything
}

func ExampleMap_Store() {
	type Queen struct{}

	black := lift.NewMap[rune]()
	black.Store(
		lift.Def[Queen]('♛'),
	)

	piece, _ := lift.Load[Queen](black)
	fmt.Printf("%c", piece)
	// Output:
	// ♛
}

func ExampleMap_Delete() {
	m := lift.NewMap[bool](
		lift.Def[rune](true),
		lift.Def[string](true),
	)

	m.Delete(lift.T[rune]())
	if _, ok := lift.Load[rune](m); !ok {
		fmt.Println("No entry found")
	}
	// Output:
	// No entry found
}

func ExampleLoad() {
	masks := lift.NewMap[int](
		lift.Def[uint8](0xff),
	)

	mask, _ := lift.Load[uint8](masks)
	fmt.Println(mask&0x1234_5678 == 0x78)

	if _, ok := lift.Load[string](masks); !ok {
		fmt.Println("oops")
	}
	// Output:
	// true
	// oops
}

func ExampleLoadKey() {
	names := lift.NewMap[string](
		lift.Def[bool]("Boolean"),
		lift.Def[string]("string"),
	)

	name1, _ := lift.LoadTypeOf(names, false)
	name2, _ := lift.LoadTypeOf(names, "false")

	fmt.Print(name1, ", ", name2)
	// Output:
	// Boolean, string
}

func ExampleLoadSym() {
	type hammer struct{}
	type screwdriver struct{}

	tools := lift.NewMap[int](
		lift.Def[hammer](2),
		lift.Def[screwdriver](17),
	)

	var count int
	for _, tool := range tools.Keys() {
		n, _ := lift.LoadSym(tools, tool)
		count += n
	}

	fmt.Printf("I have %d tools", count)
	// Output:
	// I have 19 tools
}

func ExampleLoadTypeOf() {
	fmtrs := lift.NewMap[string](
		lift.Def[uint8]("%08x"),
	)

	x := uint8(127)
	fmtr, _ := lift.LoadTypeOf(fmtrs, x)
	fmt.Printf(fmtr, x)
	// Output:
	// 0000007f
}

// Corner cases of [Sym] around empty-ish or any-ish values are reasonable.
// [Sym] is an interface type, so the zero value of a [Sym] is nil-ish, and will cause panic.
func Example_e_emptyAnyNil() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("oops!", r)
		}
	}()

	items := lift.NewMap[string](
		lift.Def[struct{}]("the empty struct"),
		lift.Def[lift.Empty]("the lift.Empty struct"),
		lift.Def[any]("anything"),
		lift.Def[lift.Sym]("the type enumeration of lift.Sym"),
	)

	report := func(tag string, item string) {
		fmt.Printf("%-12s %s,\n", tag, item)
	}

	// Three kinds of the empty struct{}
	empty := struct{}{}
	item, _ := lift.LoadTypeOf(items, empty)
	report("empty i", item)

	item, _ = lift.Load[lift.Empty](items)
	report("empty ii", item)

	type local struct{}
	item, _ = lift.LoadTypeOf(items, local{})
	report("empty iii", item)

	// Three kinds of any
	item, _ = lift.Load[any](items)
	report("any i", item)

	item, _ = lift.LoadTypeOf(items, any("other thing"))
	report("any ii", item)

	item, _ = lift.LoadSym(items, lift.Any)
	report("any iii", item)

	// LoadTypeOf infers the type enumeration of lift.Sym from anything wrapped
	item, _ = lift.LoadTypeOf(items, lift.Wrap(any(nil)))
	report("sym i", item)

	// Doesn't crash yet, as we're passing the type enumeration of lift.Sym
	var NilSym lift.Sym
	item, _ = lift.LoadTypeOf(items, NilSym)
	report("sym ii", item)

	// PANIC ENSUES - we've passed a raw nil
	item, _ = lift.LoadSym(items, NilSym)
	report("nil i", item)

	// (unreached) PANIC ENSUES:
	items.Store(
		lift.DefSym(NilSym, "panic"),
	)

	// Output:
	// empty i      the empty struct,
	// empty ii     the lift.Empty struct,
	// empty iii    ,
	// any i        anything,
	// any ii       anything,
	// any iii      anything,
	// sym i        the type enumeration of lift.Sym,
	// sym ii       the type enumeration of lift.Sym,
	// oops! runtime error: invalid memory address or nil pointer dereference
}
