package response_cache

import (
	"context"
	"errors"
)

type RedisRepository interface {
	Get(ctx context.Context, key string, value any) error
	Set(ctx context.Context, key string, value any) error
	Delete(ctx context.Context, key string) error
}

type ResponseCache struct {
	redis RedisRepository
}

func NewResponseCache(redis RedisRepository) *ResponseCache {
	return &ResponseCache{
		redis: redis,
	}
}

func (rc *ResponseCache) GetOrSetCache(
	ctx context.Context,
	key string,
	handlerFunc func(context.Context, interface{}) (interface{}, error),
) (interface{}, error) {
	var cache CacheData
	if err := rc.redis.Get(ctx, key, &cache); err != nil {
		return nil, err
	}

	if cache.Value != nil {
		if cache.Value == LockValue {
			return nil, errors.New("locked")
		}
		return cache.Value, nil
	}

	resp, err := func() (_ interface{}, err error) {
		err = rc.Lock(ctx, key)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err2 := rc.Unlock(ctx, key); err2 != nil {
				err = err2
			}
		}()
		return handlerFunc(ctx, key)
	}()
	if err != nil {
		return nil, err
	}
	if err = rc.redis.Set(ctx, key, &CacheData{Key: key, Value: resp}); err != nil {
		return nil, err
	}

	return resp, nil
}

func (rc *ResponseCache) Lock(ctx context.Context, key string) error {
	lockCache := &CacheData{Key: key, Value: LockValue}
	if err := rc.redis.Set(ctx, key, lockCache); err != nil {
		return err
	}
	return nil
}

func (rc *ResponseCache) Unlock(ctx context.Context, key string) error {
	if err := rc.redis.Delete(ctx, key); err != nil {
		return err
	}
	return nil
}
