package linkedlist

// Item of List
type Item struct {
	value      interface{}
	prev, next *Item
}

// List of Items
type List struct {
	len         int
	first, last *Item
}

// Value of item
func (item *Item) Value() interface{} {
	if item == nil {
		return nil
	}
	return item.value
}

// Next item
func (item *Item) Next() *Item {
	if item == nil {
		return nil
	}
	return item.next
}

// Prev item
func (item *Item) Prev() *Item {
	if item == nil {
		return nil
	}
	return item.prev
}

// Len of list
func (list *List) Len() int {
	return list.len
}

// First element of list
func (list *List) First() *Item {
	return list.first
}

// Last element of list
func (list *List) Last() *Item {
	return list.last
}

// PushBack insert element after last element
func (list *List) PushBack(v interface{}) {
	item := Item{value: v}
	if list.Len() == 0 {
		list.first = &item
		list.last = &item
	} else {
		item.prev = list.last
		list.last.next = &item
		list.last = &item
	}
	list.len++
}

// PushFront insert element before first element
func (list *List) PushFront(v interface{}) {
	item := Item{value: v}
	if list.Len() == 0 {
		list.first = &item
		list.last = &item
	} else {
		item.next = list.first
		list.first.prev = &item
		list.first = &item
	}
	list.len++
}

// Remove element from list
func (list *List) Remove(item Item) {
	switch {
	case list.len == 0:
		return
	case list.len == 1 && *list.first == item:
		list.first = nil
		list.last = nil
		list.len--
	case item.prev == nil && *list.first == item:
		list.first = list.first.next
		if list.first != nil {
			list.first.prev = nil
		}
		list.len--
	case item.next == nil && *list.last == item:
		list.last = list.last.prev
		if list.last != nil {
			list.last.next = nil
		}
		list.len--
	case item.prev != nil && item.next != nil && *item.prev.next == item && *item.next.prev == item:
		item.prev.next = item.next
		item.next.prev = item.prev
		list.len--
	}
}
