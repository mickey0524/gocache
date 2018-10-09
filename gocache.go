package gocache

type Cache interface {
	Get(string) (interface{}, error)
	Set(string, interface{}) error
}
