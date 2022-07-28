// The conv package provides a [Coverter], which indexes conversion functions.
package conv

import (
	"fmt"

	"github.com/AndrewHarrisSPU/lift"
)

// A Converter indexes conversion functions
type Converter struct {
	defs			lift.Map[Entry]	
}

// An entry is a wrapped conversion function
type Entry lift.Sym

func NewConverter(defs ...Entry) Converter {
	cv := Converter{
		defs:		lift.NewMap[Entry](),
	}
	cv.Store(defs...)
	return cv
}

// Def wraps a conversion function, yielding an [Entry]
func Def[SRC any, DST any]( convFunc func( SRC )( DST, error )) Entry {
	return Entry( lift.Wrap(convFunc))
}

// Lookup returns a conversion function from source to destination type, if found.
func Lookup[SRC any, DST any]( cv Converter) (func(SRC)(DST,error), bool ){
	sig := lift.T[func(SRC)(DST, error)]()
	if convFunc, ok := lift.LoadSym( cv.defs, sig ); ok {
		return lift.MustUnwrap[func(SRC)(DST,error)](convFunc), true
	}
	return nil, false
}

// To converts the src argument to the provided DST type.
// Conversion may fail if a conversion function isn't found,
// or it may fail if a particular value fails to convert.
func To[DST any, SRC any](cv Converter, src SRC) (dst DST, err error ){
	// If SRC == DST, return src
	if lift.T[SRC]() == lift.T[DST]() {
		return any( src ).(DST), nil
	}

	if convFunc, ok := Lookup[SRC, DST](cv); ok {
		return convFunc( src )
	}
	return dst, fmt.Errorf( "Conversion not found: %T->%T", src, dst )
}

// Store defines a conversion in the [Converter].
func (cv Converter) Store(defs ...Entry) {
	for _, def := range defs {
		cv.defs.Store(
			lift.DefSym( lift.Sym( def ), def ),
		)
	}
}

// Delete removes a defined conversion from the [Converter].
func (cv Converter) Delete(keys ...Entry) {
	for _, key := range keys {
		cv.defs.Delete( key )
	}
}