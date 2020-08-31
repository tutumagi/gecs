package gecs

import "fmt"

// checkComponentFamily the all components name -> componentID
// Registry will register all components of the checkComponentFamily when init.
var checkComponentFamily _KVStrInt = make(_KVStrInt, 0)

var checkSingletonFamily _KVStrInt = make(_KVStrInt, 0)

// RegisterComponent register component with name
func RegisterComponent(name string) ComponentID {
	if _, ok := checkComponentFamily[name]; ok {
		panic(fmt.Sprintf("register same name component %v", name))
	}

	cid := len(checkComponentFamily)
	checkComponentFamily[name] = cid

	return ComponentID(cid)
}

// RegisterSingleton register singleton with name
func RegisterSingleton(name string) SingletonID {
	if _, ok := checkSingletonFamily[name]; ok {
		panic(fmt.Sprintf("register same name component %v", name))
	}

	cid := len(checkSingletonFamily)
	checkSingletonFamily[name] = cid

	return SingletonID(cid)
}
