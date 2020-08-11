package entt

import (
	"gonut/engine/algo"
	"runtime/debug"
	"testing"
)

func Test_SparseSetEntity(t *testing.T) {
	s := NewSparseSet()
	s.Reserve(100)
	s.Construct(EntityID(5))
	s.Construct(EntityID(1))
	s.Construct(EntityID(4))

	algo.Each(s.Iterator(), func(data interface{}) {
		t.Logf("entity %v", data)
	})

	t.Logf("SparseSet %v", s)
	AssertBool(t, s.Has(EntityID(5)), true)
	AssertBool(t, s.Has(EntityID(1)), true)
	AssertBool(t, s.Has(EntityID(4)), true)
	AssertBool(t, s.Has(EntityID(10)), false)

	AssertInt(t, s.Get(EntityID(5)), 0)
	AssertInt(t, s.Get(EntityID(1)), 1)
	AssertInt(t, s.Get(EntityID(4)), 2)

	s.Destroy(EntityID(5))
	t.Logf("s %v", s)
	AssertBool(t, s.Has(EntityID(1)), true)
	AssertBool(t, s.Has(EntityID(4)), true)
}

type _Data1 struct {
	name string
}

func Test_SparseSet2(t *testing.T) {
	s := NewSparseSet2(ComponentID(1))
	s.Reserve(100)
	s.Construct(EntityID(5), &_Data1{"tufei5"})
	s.Construct(EntityID(1), &_Data1{"tufei1"})
	s.Construct(EntityID(4), &_Data1{"tufei4"})
	s.Construct(EntityID(10), &_Data1{"tufei10"})

	t.Logf("s %v", s)
	AssertBool(t, s.Has(EntityID(5)), true)
	AssertBool(t, s.Has(EntityID(1)), true)
	AssertBool(t, s.Has(EntityID(4)), true)
	AssertBool(t, s.Has(EntityID(10)), true)

	AssertString(t, s.Get(EntityID(5)).(*_Data1).name, "tufei5")
	AssertString(t, s.Get(EntityID(1)).(*_Data1).name, "tufei1")
	AssertString(t, s.Get(EntityID(4)).(*_Data1).name, "tufei4")
	AssertString(t, s.Get(EntityID(10)).(*_Data1).name, "tufei10")

	s.Destroy(EntityID(5))
	AssertBool(t, s.Has(EntityID(5)), false)
	s.Construct(EntityID(5), &_Data1{"tufei5"})
	// s.Destroy(EntityID(5))
	t.Logf("s %v", s)
	AssertBool(t, s.Has(EntityID(1)), true)
	AssertBool(t, s.Has(EntityID(4)), true)
	AssertBool(t, s.Has(EntityID(5)), true)

}

func AssertBool(tb testing.TB, expression bool, expect bool) {
	if expression != expect {
		tb.Errorf("should be %v %v", expect, string(debug.Stack()))
	}
}

func AssertInt(tb testing.TB, expression int, expect int) {
	if expression != expect {
		tb.Errorf("should be %v %v", expect, string(debug.Stack()))
	}
}

func AssertString(tb testing.TB, expression string, expect string) {
	if expression != expect {
		tb.Errorf("should be %v %v", expect, string(debug.Stack()))
	}
}
