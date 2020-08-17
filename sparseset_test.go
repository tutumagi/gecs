package entt

import (
	"testing"

	. "github.com/go-playground/assert/v2"
)

func Test_SparseSetEntity(t *testing.T) {
	s := NewSparseSet()
	Equal(t, s.Empty(), true)
	s.Reserve(100)
	Equal(t, s.Empty(), true)
	s.Emplace(EntityID(5))
	Equal(t, s.Empty(), false)
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

	s.Destroy(EntityID(1))
	s.Destroy(EntityID(4))
	Equal(t, s.Empty(), true)
}

type _Data1 struct {
	name string
}

func Test_SparseSwap(t *testing.T) {
	s := NewSparseSet()

	entity5 := EntityID(5)
	entity1 := EntityID(1)
	entity4 := EntityID(4)
	entity10 := EntityID(10)

	s.Emplace(entity5)
	s.Emplace(entity1)
	s.Emplace(entity4)
	s.Emplace(entity10)

	Equal(t, s.Has(entity5), true)
	Equal(t, s.Has(entity1), true)
	Equal(t, s.Has(entity4), true)
	Equal(t, s.Has(entity10), true)

	s.swap(entity1, entity5)

	Equal(t, s.Has(entity5), true)
	Equal(t, s.Has(entity1), true)
	Equal(t, s.Has(entity4), true)
	Equal(t, s.Has(entity10), true)
}
