package model

type batchRecord[V any] struct {
	key    string
	setter func(v V)
}

type loader[V any] struct {
	loaders []batchRecord[V]
	keyMaps map[string]bool
	batch   func(keys []string) (map[string]V, error)
}

func newBatchLoader[V any]() *loader[V] {
	return &loader[V]{keyMaps: map[string]bool{}}
}

func (l *loader[V]) SetBatch(f func(keys []string) (map[string]V, error)) {
	l.batch = f
}

func (l *loader[V]) Add(key string, setter func(v V)) {
	l.keyMaps[key] = true
	l.loaders = append(l.loaders, batchRecord[V]{key: key, setter: setter})
}

func (l loader[V]) Keys() []string {
	keys := []string{}
	for key := range l.keyMaps {
		keys = append(keys, key)
	}

	return keys
}

func (l loader[V]) Load() error {
	keys := []string{}
	for _, key := range l.Keys() {
		keys = append(keys, key)
	}

	maps, err := l.batch(keys)
	if err != nil {
		return err
	}

	for _, loader := range l.loaders {
		loader.setter(maps[loader.key])
	}

	return nil
}
