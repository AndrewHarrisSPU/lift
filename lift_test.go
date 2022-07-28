package lift_test

import (
	"testing"

	"github.com/AndrewHarrisSPU/lift"
)

// Ensure that MustUnwrap, MustUnwrapAs generate panics
func TestMustPanic(t *testing.T) {
	testMustPanic(t, func() {
		lift.MustUnwrap[struct{}](lift.Wrap(any(nil)))
	})
	testMustPanic(t, func() {
		lift.MustUnwrapAs[struct{}](lift.Wrap(any(nil)))
	})
}

func testMustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("missing panic")
		}
	}()

	fn()
}

// Quick testing for Map methods with no examples
func TestMapExtra(t *testing.T) {
	m := lift.NewMap[struct{}]()

	m.Store(
		lift.Def[int](struct{}{}),
		lift.Def[uint](struct{}{}),
		lift.Def[byte](struct{}{}),
		lift.Def[rune](struct{}{}),
	)

	keys := m.Keys()
	entries := m.Entries()

	if len(keys) != len(entries) || len(keys) != m.Len() {
		t.Errorf("Map method failure")
	}

	m.Delete(m.Keys()...)
	if m.Len() != 0 {
		t.Errorf("Map method failure")
	}
}
