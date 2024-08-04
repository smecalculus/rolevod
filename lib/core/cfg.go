package core

// port
type Keeper interface {
	Load(key string, v any) error
}
