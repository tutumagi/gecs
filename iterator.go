package gecs

type IIterator interface {
	Prev() IIterator
	Next() IIterator
	Equal(other IIterator) bool
	Begin() IIterator
	End() IIterator
	Data() interface{}
}

func Each(items IIterator, f func(data interface{})) {
	for first := items.Begin(); !first.Equal(items.End()); {
		f(first.Data())
		first = first.Next()
	}
}

func ReverseEach(items IIterator, f func(data interface{})) {
	for first := items.End(); !first.Equal(items.Begin()); {
		first = first.Prev()
		f(first.Data())
	}
}

func AllOf(items IIterator, f func(data interface{}) bool) bool {
	for first := items.Begin(); !first.Equal(items.End()); {
		if f(first.Data()) == false {
			return false
		}
		first = first.Next()
	}
	return true
}
