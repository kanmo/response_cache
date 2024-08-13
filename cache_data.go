package response_cache

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

// CacheData is a struct that holds the cache data
// Key is the cache key, it is consisted of userId, methodName, and idempotencyKey
// Value is the gRPC response
// TypeName is the type of the response
// UnmarshalValue is the unmarshaled response
// MarshalValue is the marshaled response
type CacheData struct {
	Key            string
	Value          interface{}
	TypeName       string
	UnmarshalValue interface{}
	MarshalValue   []byte
}

const LockValue string = "locked"

func Create(methodName string, userId string, idempotencyKey string) *CacheData {
	return &CacheData{
		Key: GenerateCacheKey(userId, methodName, idempotencyKey),
	}
}

func GenerateCacheKey(keys ...string) string {
	return strings.TrimSpace(strings.Join(keys, ":"))
}

// Marshal marshals the cache data
// It needs to specify the type to restore, so json.Marshal is applied twice only for the response
func (c *CacheData) Marshal() ([]byte, error) {
	if c.Value == LockValue {
		c.TypeName = "string"
	} else {
		c.TypeName = reflect.TypeOf(c.Value).Elem().Name()
	}
	bytes, err := json.Marshal(c.Value)
	if err != nil {
		return nil, err
	}
	c.MarshalValue = bytes
	return json.Marshal(c)
}

func (c *CacheData) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}

	t, ok := Registry.Get(c.TypeName)
	if !ok {
		return errors.New("unknown type: " + c.TypeName)
	}

	instance := reflect.New(t).Interface()
	if err := json.Unmarshal(c.MarshalValue, instance); err != nil {
		return err
	}
	c.UnmarshalValue = instance

	return nil
}
