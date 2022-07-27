// Package lift uses generic gadgets to lift types to values,
// referred to in documentation as type enumerations or type enumeration symbols.
// When lifted as values, type enumerations can predicate branching or matching logic
// over the universe of types.
//
// The package intends to make production and consumption of type enumeration symbols narrow and ergonomic.
// The most fundamental component of lift is the lift.Sym type, an exported type enumeraion symbol.
// See the README.md file for further exposition on internal mechanics.
//
// The lift.Map type is provided as an immediate and highly general application of type enumeration symbols,
// using them as map keys.
//
// Function-flavored type enumerations (e.g. from func(T), distinct from T) can predicate
// runtime dispatch gadgetry. For example, the lift/conv package leverages lift to index
// conversion functions from given source and destination types. Package examples explore
// other gadgetry, roughly in order of increasing elaboration.
package lift

import (
	"fmt"
)

// GADGETS

// A Sym is an interface protecting internal methods for producing type enumeration symbols.
type Sym interface {
	enum() Sym
	exfiltrate() any
}

// enum is a T-flavored, internal type enumeration symbol.
type enum[T any] struct{}

func (e enum[T]) enum() Sym {
	return e
}

func (e enum[T]) exfiltrate() any {
	return any(e)
}

// EnumIs determines equivalence of two type enumerations:
// one derived from T, the other from the argument
func EnumIs[T any](sym Sym) bool {
	return enum[T]{} == sym.enum()
}

// Any is the type enumeration symbol of interface{}
var Any = enum[any]{}

// Empty is an empty structure
type Empty struct{}

// SYM

// The (function) T returns a (type) T-flavored type enumeration symbol.
func T[T any]() Sym {
	return enum[T]{}
}

// TypeOf returns a T-flavored type enumeration symbol.
func TypeOf[T any](T) Sym {
	return enum[T]{}
}

type wrapped[T any] struct {
	t T
}

func (w wrapped[T]) enum() Sym {
	return enum[T]{}
}

func (w wrapped[T]) exfiltrate() any {
	return any(w.t)
}

// Wrap produces a Sym, like T or TypeOf. Unlike T or TypeOf,
// the result may be unwrapped to recover a wrapped value.
func Wrap[T any](t T) Sym {
	return wrapped[T]{t: t}
}

// Unwrap recovers a value of type T. It is successful when
// the Sym to unwrap was produced by Wrap, and T precisely matches
// the flavor of T inferred by Wrap.
func Unwrap[T any](sym Sym) (t T, ok bool) {
	if w, ok := sym.(wrapped[T]); ok {
		return w.t, true
	}
	return t, false
}

// UnwrapAs resembles Unwrap, but is successful when the wrapped
// value satisfies an interface type T.
func UnwrapAs[T any](sym Sym) (t T, ok bool) {
	t, ok = sym.exfiltrate().(T)
	return
}

// MustUnwrap is a fail-fast version of Unwrap.
// Unwrapping is expected to succed, and failure to unwrap panics.
func MustUnwrap[T any](sym Sym) T {
	if w, ok := sym.(wrapped[T]); ok {
		return w.t
	}
	panic(fmt.Errorf("MustUnwrap: want %T, got %T", enum[T]{}, sym.enum()))
}

// MustUnwrapAs is a fail-fast version of UnwrapAs.
// Unwrapping is expected to succed, and failure to unwrap panics.
func MustUnwrapAs[T any](sym Sym) T {
	if t, ok := sym.exfiltrate().(T); ok {
		return t
	}
	panic(fmt.Errorf("MustUnwrapAs: want %T, got %T", enum[T]{}, sym.enum()))
}

// MAP

// Map defines associations between type enumerations and values of type V.
type Map[V any] struct {
	defs map[Sym]V
}

// Entry encapsulates a definition of a single Map association.
type Entry[V any] struct {
	k Sym
	v V
}

// NewMap returns an initialized Map, with any provided definitions stored.
func NewMap[V any](defs ...Entry[V]) Map[V] {
	m := Map[V]{
		defs: make(map[Sym]V, len(defs)),
	}
	m.Store(defs...)
	return m
}

// Def constructs Map entries.
func Def[K any, V any](v V) Entry[V] {
	return Entry[V]{enum[K]{}, v}
}

// DefSym constructs Map entries. Unlike Def, the key flavor is already lifted in the Sym.
func DefSym[V any](sym Sym, v V) Entry[V] {
	return Entry[V]{sym.enum(), v}
}

// Store stores a variadic list of entries in a Map.
func (m Map[V]) Store(defs ...Entry[V]) {
	for _, def := range defs {
		m.defs[def.k] = def.v
	}
}

// Delete removes a variadic list of type enumerations from a Map.
func (m Map[V]) Delete(keys ...Sym) {
	for _, key := range keys {
		delete(m.defs, key.enum())
	}
}

// Len returns the number of definitions present in a Map.
func (m Map[V]) Len() int {
	return len(m.defs)
}

// Keys collects the type enumeration keys defined for a Map.
func (m Map[V]) Keys() []Sym {
	keys := make([]Sym, len(m.defs))
	i := 0
	for k := range m.defs {
		keys[i] = k
		i++
	}
	return keys
}

// Entries collects definitions present in a Map.
func (m Map[V]) Entries() []Entry[V] {
	entries := make([]Entry[V], len(m.defs))
	i := 0
	for k, v := range m.defs {
		entries[i] = Entry[V]{k: k, v: v}
		i++
	}
	return entries
}

// Load returns a value from a Map, if found.
func Load[K any, V any](m Map[V]) (v V, ok bool) {
	v, ok = m.defs[Sym(enum[K]{})]
	return
}

// LoadTypeOf resembles Load, with type parameter K inferred from a second argument.
func LoadTypeOf[K any, V any](m Map[V], _ K) (v V, ok bool) {
	v, ok = Load[K](m)
	return
}

// LoadSym resembles Load, where the type enumeration key is lifted in the second argument.
func LoadSym[V any](m Map[V], sym Sym) (v V, ok bool) {
	v, ok = m.defs[sym.enum()]
	return
}
