package di

import (
	"fmt"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
)

type Dependency struct {
	deps map[interface{}]interface{}
}

func NewDependency(cfg *config.Config) *Dependency {
	deps := &Dependency{
		deps: make(map[interface{}]interface{}),
	}
	Provide(deps, cfg)
	return deps
}

func Provide[T any](d *Dependency, impl T) {
	key := fmt.Sprintf("%T", new(T))
	d.deps[key] = impl
}

func Get[T any](d *Dependency) T {
	_t := new(T)
	key := fmt.Sprintf("%T", _t)
	val, ok := d.deps[key]
	if !ok {
		panic(key + " was not provided")
	}
	return val.(T)
}
