# lift
Lift is a generics-inflected Go package for lifting Go types

# Why `lift`?

Loosely (and imprecisely) inspired by notions of function lifting, `lift` intends to provide some glue for lifting in Go.
The central idea of `lift` is type enumeration symbols - runtime values which predicate branching and matching logic over types.

This overlaps with a number of existing mechanisms in Go, from type assertions to `reflect`. `lift` is distinguished by being
generics-inflected, and narrow in scope. Eventually, `lift` is trying to be a small but complete abstraction for runtime dispatch gadgetry.
The package is at the moment entirely experimental in nature. I hope you find the results interesting or informative or useful :)

# What is in `lift`?

## `Sym`

`Sym` are type enumeration symbols. Additionally, `Sym` may carry wrapped values in conjunction with type enumertaion symbols.
This puts a generics-inflected take on familiar ideas.

## `lift.Map`

Type enumeration symbols are a natural fit as map keys.
The `Map` type uses `Sym` keys, and may be parameterized to any `V` value. For example:

```
	piece, _ := lift.Load[ Queen ]( black )
```

loads a piece associated with some type `Queen`, from some map `black`.

## Runtime dispatch gadgets

- `lift/conv` has a `Converter` type, for indexing conversion functions. For example:

```
	rgb, _ := conv.To[RGB](cv, hex)
```

converting to some type `RGB`, from some `hex` value, with some converter `cv`.

- `lift` package examples explore other runtime dispatch gadgetry.

# How does `lift` work?

There are two types of type enumeration symbols at work in `lift`. First, the internal `enum`:

```
	type enum[T any] struct{}
```

Sensibly enough:
 - Any type `T` implies a corresponding type `enum[T]`
 - An `enum[T]` is `==` to all similarly flavored `enum[T]`
 - An `enum[U]` is not `==` to any differently flavored `enum[V]`, for any `V`

Second, the exported `Sym`:

```
	type Sym interface { ... }
```

An `enum[T]` satisfies the `Sym` interface. When an `enum[T]` is boxed as a `Sym`, the comparison rules apply still apply, modulo the boxing.
Concretely, if two `Sym` share a dynamic type `enum[T]`, their dynamic values will never differ.

Not all `Sym` box an `enum[T]`. But, as an invariant promoted by the lift package, all (non-nil) `Sym` derive `enum[T]` via an internal method:

```
	enum() Sym
```

So, while a `Sym` is not type-parametric, it can propagate a type flavor for the purposes of comparison.
In turn, this permits branching or matching on type flavors (even, or especially, where a type parameter associated with `Sym` flavor isn't present).

# Is it efficient?

`lift` is not a performance hack. So far, performance hasn't been a priority; writing this package has been more about exploring the concepts. Examining performance more thoroughly is a future goal.

Glaringly, `lift.Map` is sluggish. In its defense, `lift.Map` is very general, and easy to write for. For some usages, a more performant alternative (e.g., a `switch` for fixed sets) is tractable. Something other than a built-in `map` for the underlying implementation of `lift.Map` could be interesting to explore.

That said, some preliminary benchmarking suggests that using `lift` as an alternative to equivalently elaborate schemes based on e.g. `interface{}` boxing is not a performance loss. Contrasting `lift` with `reflect`, `unsafe` etc. also seems worthwhile to look at.