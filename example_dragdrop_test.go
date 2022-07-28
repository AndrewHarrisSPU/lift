package lift_test

import (
	"fmt"

	"github.com/AndrewHarrisSPU/lift"
)

// This example shows a switch-based dispatch table, yielding different logic than a [Map] would.
func Example_c_eventhandlers() {
	// brewing some dummy objects...
	photo := &file{path: "tableflip.jif", data: "(╯°□°)╯︵ ┻━┻"}
	home := &folder{path: "home", locked: false}
	sys := &folder{path: "system", locked: true}

	// brewing some dummy events...
	events := []event{
		newEvent(mouseClick{}, sys, lift.Empty{}),
		newEvent(mouseClick{}, photo, lift.Empty{}),
		newEvent(mouseDrop{}, photo, home),
		newEvent(mouseDrop{}, photo, sys),
		newEvent(mouseClick{}, home, lift.Empty{}),
	}

	for _, ev := range events {
		eventDispatch(ev)
	}
	// Output:
	// tableflip.jif:
	// 	(╯°□°)╯︵ ┻━┻
	// home:
	// 	 tableflip.jif
}

// GADGETS

type opFunc[OP any, SRC any, DST any] func(OP, SRC, DST) bool
type evFunc = func(event) bool

type event struct {
	op, src, dst lift.Sym
	signature    lift.Sym
}

func newEvent[OP any, SRC any, DST any](op OP, src SRC, dst DST) event {
	return event{
		op:        lift.Wrap(op),
		src:       lift.Wrap(src),
		dst:       lift.Wrap(dst),
		signature: lift.T[opFunc[OP, SRC, DST]](),
	}
}

func liftHandler[OP any, SRC any, DST any](fn opFunc[OP, SRC, DST]) evFunc {
	return func(ev event) bool {
		op, opOk := lift.Unwrap[OP](ev.op)
		src, srcOk := lift.Unwrap[SRC](ev.src)
		dst, dstOk := lift.Unwrap[DST](ev.dst)

		if !opOk || !srcOk || !dstOk {
			return false
		}
		return fn(op, src, dst)
	}
}

// DUMMY TYPES

// dummy mouse event types
type mouseClick struct{}
type mouseDrop struct{}

// dummy file types
type file struct {
	path, data string
}

type folder struct {
	path   string
	locked bool
	files  []*file
}

// HANDLERS & DISPATCH

func rejectLockedFolder(ev event) (rejected bool) {
	if dir, isDir := lift.Unwrap[*folder](ev.src); isDir && dir.locked {
		return true
	}
	if dir, isDir := lift.Unwrap[*folder](ev.dst); isDir && dir.locked {
		return true
	}
	return false
}

func openFile(op mouseClick, f *file, _ lift.Empty) bool {
	fmt.Printf("%s:\n\t%s\n", f.path, f.data)
	return true
}

func listFiles(op mouseClick, dir *folder, _ lift.Empty) bool {
	fmt.Printf("%s:\n", dir.path)
	for _, f := range dir.files {
		fmt.Println("\t", f.path)
	}
	return true
}

func moveFile(op mouseDrop, f *file, dir *folder) bool {
	dir.files = append(dir.files, f)
	return true
}

func eventDispatch(ev event) {
	switch {
	case rejectLockedFolder(ev):
	case liftHandler(openFile)(ev):
	case liftHandler(listFiles)(ev):
	case liftHandler(moveFile)(ev):
	}
}
