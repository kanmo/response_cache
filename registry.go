package response_cache

import (
	"reflect"
	"sync"
)

var Registry = NewResponseTypeRegistry()

func NewResponseTypeRegistry() *ResponseTypeRegistry {
	return &ResponseTypeRegistry{
		types: map[string]reflect.Type{},
	}
}

type ResponseTypeRegistry struct {
	types map[string]reflect.Type
	mux   sync.Mutex
}

func (r *ResponseTypeRegistry) Register(name string, t reflect.Type) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.types[name] = t
}

func (r *ResponseTypeRegistry) Get(name string) (reflect.Type, bool) {

	t, ok := r.types[name]
	return t, ok
}
