package linkedlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	item := Item{}
	assert.Nil(t, item.Value())
	item = Item{value: 1}
	assert.Equal(t, item.Value(), 1)
}

func TestNext(t *testing.T) {
	item := Item{}
	assert.Nil(t, item.Next())
	item2 := Item{next: &item}
	assert.Equal(t, item2.Next(), &item)
}

func TestPrev(t *testing.T) {
	item := Item{}
	assert.Nil(t, item.Prev())
	item2 := Item{prev: &item}
	assert.Equal(t, item2.Prev(), &item)
}

func TestEmptyList(t *testing.T) {
	list := List{}
	assert.Nil(t, list.First())
	assert.Nil(t, list.Last())
	assert.Equal(t, list.Len(), 0)
}

func helperPushFirstElement(t *testing.T, list *List, v interface{}) {
	t.Helper()
	assert.NotNil(t, list.First())
	assert.NotNil(t, list.Last())
	assert.Equal(t, list.First(), list.Last())
	assert.Equal(t, list.Last().Value(), v)
	assert.Nil(t, list.First().Prev())
	assert.Nil(t, list.Last().Next())
	assert.Equal(t, list.Len(), 1)
}

func TestPushBack(t *testing.T) {
	list := List{}
	list.PushBack(1)
	helperPushFirstElement(t, &list, 1)
	list.PushBack(2)
	assert.NotNil(t, list.Last())
	assert.Equal(t, list.Last().Value(), 2)
	assert.Equal(t, list.Last().Prev(), list.First())
	assert.Equal(t, list.Last().Prev().Next(), list.Last())
	assert.Equal(t, list.Len(), 2)
	list.PushBack(nil)
	assert.Equal(t, list.Last().Value(), nil)
	assert.Equal(t, list.Len(), 3)
}

func TestPushFront(t *testing.T) {
	list := List{}
	list.PushFront(1)
	helperPushFirstElement(t, &list, 1)
	list.PushFront(2)
	assert.NotNil(t, list.First())
	assert.Equal(t, list.First().Value(), 2)
	assert.Equal(t, list.First().Next(), list.Last())
	assert.Equal(t, list.First().Next().Prev(), list.First())
	assert.Equal(t, list.Len(), 2)
	list.PushFront(nil)
	assert.Equal(t, list.First().Value(), nil)
	assert.Equal(t, list.Len(), 3)
}

func TestRemove(t *testing.T) {
	list := List{}
	fakeItem := Item{value: 1}
	// remove unexisting element from empty list
	list.Remove(fakeItem)
	assert.Nil(t, list.First())
	assert.Nil(t, list.Last())
	assert.Equal(t, list.Len(), 0)
	// fill test data
	list.PushFront(2)
	// remove unexisting element from one-item list(
	list.Remove(fakeItem)
	assert.Equal(t, list.First(), list.Last())
	assert.Equal(t, list.First().Value(), 2)
	assert.Equal(t, list.Len(), 1)
	// remove existing element from one-item list
	list.Remove(*list.First())
	assert.Equal(t, list.Len(), 0)
	assert.Nil(t, list.First())
	assert.Nil(t, list.Last())
	// fill test data
	var middleItem *Item
	for i := 0; i < 10; i++ {
		list.PushBack(i)
		if i == 5 {
			middleItem = list.Last()
		}
	}
	// remove unexisting element from X-item list
	list.Remove(fakeItem)
	assert.Equal(t, list.Len(), 10)
	assert.Equal(t, list.Last().Value(), 9)
	// remove last element from list
	list.Remove(*list.Last())
	assert.Nil(t, list.Last().Next())
	assert.Equal(t, list.Len(), 9)
	assert.Equal(t, list.Last().Value(), 8)
	// remove first element from list
	list.Remove(*list.First())
	assert.Nil(t, list.First().Prev())
	assert.Equal(t, list.Len(), 8)
	assert.Equal(t, list.First().Value(), 1)
	// remove middle element from list
	list.Remove(*middleItem)
	assert.Equal(t, list.Len(), 7)
	assert.Equal(t, middleItem.Next().Prev(), middleItem.Prev())
	assert.Equal(t, middleItem.Prev().Next(), middleItem.Next())

}
