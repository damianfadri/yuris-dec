package dsa

type Iterator[T any] struct {
	Index		int
	Items		[]T
}

func NewIterator[T any](items []T) *Iterator[T] {
	return &Iterator[T]{0, items}
}

func (it *Iterator[T]) HasNext() bool {
	return it.Index < len(it.Items)
}

func (it *Iterator[T]) Next() *T {
	var item *T
	if it.HasNext() {
		item = &it.Items[it.Index]
		it.Index = it.Index + 1
	}
	
	return item
}