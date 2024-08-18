package response_cache

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	pb "response_cache/pb/test/proto"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockRedisRepository struct {
	mock.Mock
}

func (m *MockRedisRepository) Get(ctx context.Context, key string, value any) error {
	args := m.Called(ctx, key, value)
	//if cacheData, ok := args.Get(1).(*CacheData); ok && cacheData != nil {
	//	*value.(*CacheData) = *cacheData
	//}
	//return args.Error(0)
	if args.Get(0) == nil {
		if cacheData, ok := value.(*CacheData); ok {
			cacheData.Value = nil
		}
	}
	return args.Error(0)
}

func (m *MockRedisRepository) Set(ctx context.Context, key string, value any) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockRedisRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func TestGetOrCache(t *testing.T) {
	ctx := context.Background()
	method := "/test/method"
	userId := "test_user_id"
	idempotencyKey := "key123"

	mockRedis := new(MockRedisRepository)
	rc := NewResponseCache(mockRedis)
	key := GenerateCacheKey(userId, method, idempotencyKey)
	handlerFunc := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &pb.UserResponse{
			UserId:   userId,
			UserName: "Test User",
		}, nil
	}

	t.Run("cache does not exist", func(t *testing.T) {
		mockRedis.On("Get", ctx, key, mock.Anything).
			Return(nil).Run(func(args mock.Arguments) {
			if dest, ok := args.Get(2).(*CacheData); ok {
				dest.Value = nil
			}
		}).Once()
		mockRedis.On("Set", ctx, key, mock.Anything).
			Return(nil).Twice()
		mockRedis.On("Delete", ctx, key).
			Return(nil).Once()

		resp, err := rc.GetOrSetCache(ctx, key, handlerFunc)
		assert.NoError(t, err)
		assert.Equal(t, "test_user_id", resp.(*pb.UserResponse).UserId)
		assert.Equal(t, "Test User", resp.(*pb.UserResponse).UserName)

		mockRedis.AssertCalled(
			t,
			"Set",
			ctx,
			key,
			mock.MatchedBy(func(value interface{}) bool {
				cacheData, ok := value.(*CacheData)
				if !ok {
					return false
				}
				assert.Equal(t, fmt.Sprintf("%s:%s:%s", userId, method, idempotencyKey), cacheData.Key)

				userResp, ok := cacheData.Value.(*pb.UserResponse)
				if !ok {
					return false
				}
				return userResp.UserId == "test_user_id" && userResp.UserName == "Test User"
			}))

	})

}
