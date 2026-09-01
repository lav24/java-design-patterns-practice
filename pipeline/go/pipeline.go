package main

// Handler is a single stage: takes one I, returns one O. Go has no interfaces
// to implement here — a plain function already satisfies this shape.
type Handler[I any, O any] func(I) O

// Pipeline wraps one handler, which may itself be several stages nested
// together (see AddHandler).
type Pipeline[I any, O any] struct {
	run Handler[I, O]
}

// NewPipeline starts a pipeline with its first stage.
func NewPipeline[I any, O any](h Handler[I, O]) Pipeline[I, O] {
	return Pipeline[I, O]{run: h}
}

// AddHandler returns a NEW pipeline whose single handler runs p's chain
// first, then feeds the result into next. Go methods can't introduce their
// own type parameter beyond the receiver's, so this has to be a standalone
// function rather than a Pipeline method (unlike Java's Pipeline.addHandler).
func AddHandler[I any, O any, K any](p Pipeline[I, O], next Handler[O, K]) Pipeline[I, K] {
	return Pipeline[I, K]{run: func(input I) K {
		return next(p.run(input))
	}}
}

// Execute runs the whole chain against input.
func (p Pipeline[I, O]) Execute(input I) O {
	return p.run(input)
}
