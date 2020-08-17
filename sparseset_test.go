package entt

import (
	"testing"

	. "github.com/go-playground/assert/v2"
)

func Test_SparseSetEntity(t *testing.T) {
	s := NewSparseSet()
	s.Reserve(100)
	s.Emplace(EntityID(5))
	s.Emplace(EntityID(1))
	s.Emplace(EntityID(4))

	Equal(t, s.Has(EntityID(5)), true)
	Equal(t, s.Has(EntityID(1)), true)
	Equal(t, s.Has(EntityID(4)), true)
	Equal(t, s.Has(EntityID(10)), false)

	Equal(t, s.Index(EntityID(5)), 0)
	Equal(t, s.Index(EntityID(1)), 1)
	Equal(t, s.Index(EntityID(4)), 2)

	s.Destroy(EntityID(5))
	t.Logf("s %v", s)
	Equal(t, s.Has(EntityID(5)), false)
	Equal(t, s.Has(EntityID(1)), true)
	Equal(t, s.Has(EntityID(4)), true)
}

type _Data1 struct {
	name string
}

func Test_SparseSet2(t *testing.T) {
	s := NewStorage(ComponentID(1))
	s.Reserve(100)
	s.Emplace(EntityID(5), &_Data1{"tufei5"})
	s.Emplace(EntityID(1), &_Data1{"tufei1"})
	s.Emplace(EntityID(4), &_Data1{"tufei4"})
	s.Emplace(EntityID(10), &_Data1{"tufei10"})

	t.Logf("s %v", s)
	Equal(t, s.Has(EntityID(5)), true)
	Equal(t, s.Has(EntityID(1)), true)
	Equal(t, s.Has(EntityID(4)), true)
	Equal(t, s.Has(EntityID(10)), true)

	Equal(t, s.Get(EntityID(5)).(*_Data1).name, "tufei5")
	Equal(t, s.Get(EntityID(1)).(*_Data1).name, "tufei1")
	Equal(t, s.Get(EntityID(4)).(*_Data1).name, "tufei4")
	Equal(t, s.Get(EntityID(10)).(*_Data1).name, "tufei10")

	s.Destroy(EntityID(5))
	Equal(t, s.Has(EntityID(5)), false)
	s.Emplace(EntityID(5), &_Data1{"tufei5"})
	// s.Destroy(EntityID(5))
	t.Logf("s %v", s)
	Equal(t, s.Has(EntityID(1)), true)
	Equal(t, s.Has(EntityID(4)), true)
	Equal(t, s.Has(EntityID(5)), true)

}
