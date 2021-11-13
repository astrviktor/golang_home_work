package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length int
	head   *ListItem
	tail   *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	// добавление в начало списка, head нужно сдвинуть вправо, элемент - новый head

	// создаем элемент
	item := ListItem{v, nil, nil}
	l.length++

	// если это самый первый элемент, то этот элемент и head и tail
	if l.length == 1 {
		l.head = &item
		l.tail = &item
		return &item
	}

	// если есть элементы, у предыдущего head надо поменять связь слева, у текущего элемента - связь справа
	item.Next = l.head
	l.head.Prev = &item

	// у нас новый head
	l.head = &item

	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	// добавление в конец списка, tail нужно сдвинуть влево, элемент - новый tail

	// создаем элемент
	item := ListItem{v, nil, nil}
	l.length++

	// если это самый первый элемент, то этот элемент и head и tail
	if l.length == 1 {
		l.head = &item
		l.tail = &item
		return &item
	}

	// у предыдущего tail надо поменять связь справа, у текущего элемента - связь слева
	item.Prev = l.tail
	l.tail.Next = &item

	// у нас новый tail
	l.tail = &item

	return &item
}

func (l *list) Remove(i *ListItem) {
	// если удаляем элемент, можем удалить или head или tail или из середины,
	// и еще это может быть удаление самого последнего элемента в списке

	// если длинна 0 то выходим
	if l.length == 0 {
		return
	}

	// если элемент один, удаляем и очищаем head и tail
	if i.Prev == nil && i.Next == nil {
		l.head = nil
		l.tail = nil
		l.length--
		return
	}

	// если элемент head, то у него слева nil, двигаем head вправо
	if i.Prev == nil && i.Next != nil {
		l.head = i.Next
		l.length--
		return
	}

	// если элемент tail, то у него справа nil, двигаем tail влево
	if i.Prev != nil && i.Next == nil {
		l.tail = i.Prev
		l.length--
		return
	}

	// если элемент из центра, у левого надо исправить правую ссылку, а у правого - левую
	if i.Prev != nil && i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		l.length--
		return
	}
}

func (l *list) MoveToFront(i *ListItem) {
	// элемент уже head
	if i.Prev == nil && i.Next != nil {
		return
	}

	// элемент - tail
	if i.Prev != nil && i.Next == nil {
		l.tail = i.Prev
		l.tail.Next = nil
		i.Prev = nil
		i.Next = l.head
		l.head.Prev = i
		l.head = i
		return
	}

	// элемент из центра
	if i.Prev != nil && i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		i.Next = l.head
		i.Prev = nil
		l.head.Prev = i
		l.head = i
		return
	}
}

func NewList() List {
	return new(list)
}
