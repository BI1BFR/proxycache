package volume

type Volume struct {
	*Cache
	*Buffer
}

func NewVolume(maxEntry int) *Volume {
	return &Volume{
		Cache:  NewCache(maxEntry),
		Buffer: NewBuffer(),
	}
}

func (v *Volume) Set(key string, value []byte, ttw int64) {
	entry := &Entry{key, value}
	v.Cache.Set(entry)
	v.Buffer.Set(entry, ttw)
}

func (v *Volume) Get(key string) []byte {
	entry := v.Cache.Get(key)
	if entry != nil {
		return entry.Value
	}

	entry = v.Buffer.Get(key)
	if entry != nil {
		return entry.Value
	} else {
		return nil
	}
}
