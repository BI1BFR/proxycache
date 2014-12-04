package proxycache

type Proxy interface {
	Load(key string) (value []byte, ok bool)
	Save(key string, value []byte) (ok bool)
}
