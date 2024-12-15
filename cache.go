package memoria

import "fmt"

// implement here the different cache types here
type CachePolicy interface {
	Eject(m *Memoria, requriedSpace uint64) error
	Insert(m *Memoria, key string, val []byte) error
}

type defaultCachePolicy struct{}

func (dc *defaultCachePolicy) Eject(m *Memoria, requriedSpace uint64) error {
	spaceFreed := uint64(0)
	for key, val := range m.cache {
		if spaceFreed >= requriedSpace {
			break
		}
		valSize := uint64(len(val))
		m.cacheSize -= valSize
		delete(m.cache, key)
		spaceFreed += valSize
	}
	return nil
}

func (dc *defaultCachePolicy) Insert(m *Memoria, key string, val []byte) error {
	valueSize := uint64(len(val))
	if m.cacheSize+valueSize > m.MaxCacheSize {
		return fmt.Errorf("defaultCachePolicy: Failded to make room for value (%d/%d)", valueSize, m.MaxCacheSize)
	}
	m.cache[key] = val
	m.cacheSize += valueSize
	return nil
}
