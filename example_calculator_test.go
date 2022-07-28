package lift_test

import (
	"fmt"

	"github.com/AndrewHarrisSPU/lift"
)

// This example emulates a pocket calculator, modeled as a finite state machine.
// Current state is maintained by a [Map] of transition functions. Inputs are parsed to [Sym].
// The evaluaton loop takes one [Sym], finds the associated edge in the calculator state, and
// dispatches that function.
func Example_d_calculator() {
	calculate("1+2*3=-4=C/=-5C-56=7+8=9=")
	// 1+
	// >        1
	// 2*
	// >        3
	// 3=
	// >        9
	// -4=
	// >        5
	// C/
	// >        0
	// =
	// > DIVZERO!
	// -5C-
	// >        0
	// 56=
	// >      -56
	// 7+
	// >        7
	// 8=
	// >       15
	// 9=
	// >        9
}

func calculate(input string) {
	c := newCalculator()

	for _, r := range input {
		sym := parseCalculatorKey(r)
		edge, _ := lift.LoadSym(c.state, sym)
		edge(c, sym)
	}
}

// CALCULATOR

type calc struct {
	state			lift.Map[edgeFunc]
	acc, res		int
	op
}

type edgeFunc = func(*calc, lift.Sym )
type op = func(*calc) error

func newCalculator() *calc {
	c := new( calc )
	c.state = lift.NewMap[edgeFunc](
		lift.Def[keyC](clear),
		lift.Def[keyEq](eq),
	)
	c.reset()
	c.enterStart()
	return c
}

// STATES

func (c *calc) enterStart() {
	c.state.Store(
		lift.Def[keyOp](eval),
		lift.Def[keyNum](beginAcc),
	)
}

func (c *calc) enterAccumulate() {
	c.state.Store(
		lift.Def[keyOp](eval),
		lift.Def[keyNum](acc),
	)
}

func (c *calc) enterEvaluated() {
	c.state.Store(
		lift.Def[keyOp](store),
		lift.Def[keyNum](resetAcc),
	)
}

func (c *calc) enterErr() {
	c.state.Store(
		lift.Def[keyOp](nop),
		lift.Def[keyNum](nop),
	)
}

// EDGES

func clear( c *calc, _ lift.Sym ){
	c.reset()
	c.enterStart()
}

func eq(c* calc, _ lift.Sym ){
	if err := c.evaluate(); err != nil {
		c.enterErr()
		return
	}
	c.enterEvaluated()
}

func eval( c *calc, sym lift.Sym ){
	if err := c.evaluate(); err != nil {
		c.enterErr()
		return
	}
	store(c, sym)
}

func store( c *calc, sym lift.Sym){
	c.op = lift.MustUnwrap[keyOp](sym)
	c.enterStart()
}

func acc( c *calc, sym lift.Sym ){
	digit := lift.MustUnwrap[keyNum](sym)
	c.acc *= 10
	c.acc += digit
	c.enterAccumulate()
}

func beginAcc( c *calc, sym lift.Sym ){
	digit := lift.MustUnwrap[keyNum](sym)
	if digit == 0 {
		return
	}
	c.acc = 0
	acc(c, sym)
}

func resetAcc( c *calc, sym lift.Sym ){
	c.reset()
	beginAcc(c, sym)
}

func nop( c *calc, _ lift.Sym ) {}

// METHODS

func (c *calc) evaluate() error {
	fmt.Print( "\n> " )
	if err := c.op( c ); err != nil {
		fmt.Println( err.Error() )
		return err
	}
	fmt.Printf("%8d\n", c.res)
	return nil
}

func (c *calc) reset() {
	c.acc, c.res = 0, 0
	c.op = (*calc).add	
}

func (c *calc) add() error {
	c.res += c.acc
	return nil
}

func (c *calc) sub() error {
	c.res -= c.acc
	return nil
}

func (c *calc) mul() error {
	c.res *= c.acc
	return nil
}

func (c *calc) div() error {
	if c.acc == 0 {
		return fmt.Errorf("DIVZERO!")
	}
	c.res /= c.acc
	return nil
}

// PARSING

type keyC struct{}
type keyEq struct{}
type keyOp = func(*calc) error
type keyNum = int

func parseCalculatorKey(r rune) lift.Sym {
	fmt.Printf("%c", r)

	switch r {
	case 'C':
		return lift.Wrap(keyC{})
	case '=':
		return lift.Wrap(keyEq{})
	case '+':
		return lift.Wrap((*calc).add)
	case '-':
		return lift.Wrap((*calc).sub)
	case '*':
		return lift.Wrap((*calc).mul)
	case '/':
		return lift.Wrap((*calc).div)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return lift.Wrap(keyNum(r - '0'))
	default:
		return lift.Wrap(keyC{})
	}
}