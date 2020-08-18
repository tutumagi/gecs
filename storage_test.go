package gecs

import (
	"testing"

	. "github.com/go-playground/assert/v2"
)

func Test_StorageEmplaceDestroy(t *testing.T) {
	comID := ComponentID(10)
	stringStorage := NewStorage(comID)
	Equal(t, stringStorage.Size(), 0)
	Equal(t, stringStorage.Capacity(), 0)

	entityFoo := EntityID(0)
	stringStorage.Emplace(entityFoo, "foo")
	stringStorage.Has(entityFoo)
	Equal(t, stringStorage.Has(entityFoo), true)
	Equal(t, stringStorage.Get(entityFoo).(string), "foo")
	Equal(t, stringStorage.Size(), 1)

	entityBar := EntityID(1)
	stringStorage.Emplace(entityBar, "bar")
	Equal(t, stringStorage.Has(entityBar), true)
	Equal(t, stringStorage.Get(entityBar).(string), "bar")

	Equal(t, stringStorage.Size(), 2)

	stringStorage.Destroy(entityFoo)
	Equal(t, stringStorage.Has(entityFoo), false)
	Equal(t, stringStorage.Size(), 1)

	Equal(t, stringStorage.TryGet(entityFoo), nil)
	Equal(t, stringStorage.Has(entityBar), true)
	Equal(t, stringStorage.Get(entityBar).(string), "bar")
}

func Test_StorageCapacity(t *testing.T) {
	comID := ComponentID(10)
	stringStorage := NewStorage(comID)
	Equal(t, stringStorage.Size(), 0)
	Equal(t, stringStorage.Capacity(), 0)

	stringStorage.Reserve(1000)
	Equal(t, stringStorage.Size(), 0)
	Equal(t, stringStorage.Capacity(), 1000)
}
