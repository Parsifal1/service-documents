package services

var cache = make(map[string]interface{})

// Set устанавливает значение в кеш
func Set(key string, value interface{}) {
	cache[key] = value
}

// Get возвращает значение из кеша
func Get(key string) (interface{}, bool) {
	value, found := cache[key]
	return value, found
}

// Invalidate удаляет значение из кеша
func Invalidate(key string) {
	delete(cache, key)
}
