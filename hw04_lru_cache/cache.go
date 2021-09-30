package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	// возвращаемое значение - флаг, присутствовал ли элемент в кэше

	// если элемент присутствует в словаре, то обновить его значение и переместить элемент в начало очереди
	if item, ok := lc.items[key]; ok {
		item.Value = value
		lc.queue.PushFront(item)
		return true
	}

	// если элемента нет в словаре, то добавить в словарь и в начало очереди
	lc.queue.PushFront(value)
	lc.items[key] = lc.queue.Front()

	// если размер очереди больше ёмкости кэша,
	// то необходимо удалить последний элемент из очереди и его значение из словаря
	if lc.capacity < lc.queue.Len() {
		tail := lc.queue.Back()

		// нужно найти и удалить этот элемент из мапы
		for itemKey, itemValue := range lc.items {
			if itemValue == tail {
				delete(lc.items, itemKey)
			}
		}

		// и удалить из очереди
		lc.queue.Remove(tail)
	}

	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	// если элемент присутствует в словаре, то переместить элемент в начало очереди и вернуть его значение и true
	if item, ok := lc.items[key]; ok {
		lc.queue.PushFront(item)
		return item.Value, true
	}

	// если элемента нет в словаре, то вернуть nil и false
	return nil, false
}

func (lc *lruCache) Clear() {
	// очистка емкости
	lc.capacity = 0

	// очистка очереди
	lc.queue = NewList()

	// очистка мапы
	lc.items = make(map[Key]*ListItem, lc.capacity)
}

// type cacheItem struct {
//	 key   string
//	 value interface{}
// }

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
