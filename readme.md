# Gecs

An Entity-Component-System storage implement with golang. It's a sparse sets based inspired by `EnTT`.
* Fast iterator.
* Reuse entity id.
* Combine include/exclude Component query supported.

 All the technique/features/idea are from [EnTT](https://github.com/skypjack/entt). (use the amazing ecs implement [Entt](https://github.com/skypjack/entt) if you code cpp.)

As go does not support generic now. So we should add extra info to identify different component type.(It's called component id in the project)

[![Build Status](https://travis-ci.com/tutumagi/gecs.svg?branch=master)](https://travis-ci.com/tutumagi/gecs) 
[![codecov](https://codecov.io/gh/tutumagi/gecs/branch/master/graph/badge.svg)](https://codecov.io/gh/tutumagi/gecs)
[![codebeat badge](https://codebeat.co/badges/d8005100-a652-456e-a95e-cf11f40c90d6)](https://codebeat.co/projects/github-com-tutumagi-gecs-master)
![license](https://img.shields.io/github/license/tutumagi/gesc) 

## Getting Started

```go
// pre define your component data
type Name string
type Age int
type Position struct {
    X int
    Y int
    Z int
}

func main() {
    // create a registry
    	registry := NewRegistry()

    // register all components what you will use later
    // It's to identify component type
	NameID := registry.RegisterComponent("name")
	AgeID := registry.RegisterComponent("age")
	PositionID := registry.RegisterComponent("position")

	// create an entity
	entity := registry.Create()

	// assign component data
	registry.Assign(entity, NameID, Name("tufei"))
	registry.Assign(entity, AgeID, Age(12))

	// check whether an entity has given components data or not
	registry.Has(entity, NameID)        // true
    registry.Has(entity, NameID, AgeID) // true
    registry.Has(entity, PositionID)    // false

	// iterator entities by givin component
	registry.View(NameID, AgeID).Each(func(entity EntityID, datas map[ComponentID]interface{}) {
		name := datas[NameID].(Name) // the component type bind componentID must be consistent with you assign before
		age := datas[AgeID].(Age)

		fmt.Printf("name is %v\n", name)
		fmt.Printf("age is %v\n", age)
	})
}
```

## Prerequisites

GO 1.14 and above


## TODO

* [x] exclude test
* [ ] avoid gc as much as possible
* [ ] notify when update component data(create/update/remove) 
* [ ] avoid iterate empty component data
* [ ] snapshot
* [ ] group
* [ ] database backend
* [ ] benchmark

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details