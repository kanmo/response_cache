# response_cache

This Go library provides a simple and efficient way to cache gRPC response data for idempotent RPC calls. It is designed to help you cache responses and handle repeated requests efficiently.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Initialize Cache Module](#initialize-cache-module)
  - [Registering Types](#registering-types)
  - [Using GetOrSetCache](#using-getorsetcache)
  - [Redis Get and Set Methods](#redis-get-and-set-methods)
  - [Sample gRPC Interceptor](#sample-grpc-interceptor)
- [Example](#example)
- [Testing](#testing)
- [License](#license)

## Installation

To use this library in your project, add it as a dependency:

```bash
go get github.com/kanmo/response_cache
```

## Usage

### Initialize Cache Module

Before using the GetOrSetCache method, you must initialize the ResponseCache module with an implementation of the RedisRepository interface. This interface defines methods for interacting with Redis, including Get, Set, and Delete.

```go
import (
	"github.com/kanmo/response_cache"
	"github.com/yourusername/response_cache/pb" // Assuming your proto-generated package
)

func init() {
	// Register the types that will be cached
	response_cache.RegisterType(&pb.UserResponse{})
}

// Example RedisRepository implementation
type MyRedisRepository struct {
	// Your Redis client or connection pool here
}

func (r *MyRedisRepository) Get(ctx context.Context, key string, value any) error {
	// Implement your Redis GET logic here
	// Remember to Unmarshal the value before returning
	return nil
}

func (r *MyRedisRepository) Set(ctx context.Context, key string, value any) error {
	// Implement your Redis SET logic here
	// Remember to Marshal the value before setting it in Redis
	return nil
}

func (r *MyRedisRepository) Delete(ctx context.Context, key string) error {
	// Implement your Redis DELETE logic here
	return nil
}

func main() {
	redisRepo := &MyRedisRepository{}
	cache := response_cache.NewResponseCache(redisRepo)
}
```

### Registering Types

Before using the GetOrSetCache method, you must register the RPC response types that will be cached. This ensures that the library knows how to handle the specific types being cached.

```go
import (
	"github.com/kanmo/response_cache"
	"github.com/kanmo/response_cache/pb" // Assuming your proto-generated package
)

func init() {
	// Register the types that will be cached
	response_cache.RegisterType(&pb.UserResponse{})
}
```

### Using GetOrSetCache

To use the GetOrSetCache method, follow these steps:

1. Create a CacheData instance: Use the NewResponseCache method to create a cache object.
2. 	Call the GetOrSetCache method: Use the GetOrSetCache method to attempt to retrieve a cached response or execute the handler function and cache the response.

```go
ctx := context.Background()
cacheKey := response_cache.GenerateCacheKey(userId, methodName, idempotencyKey)

cache := response_cache.NewResponseCache(redisRepo)

response, err := cache.GetOrSetCache(ctx, req, func(ctx context.Context, req interface{}) (interface{}, error) {
// Execute the actual RPC handler here
return &pb.UserResponse{UserId: "1234", UserName: "Test User"}, nil
})

if err != nil {
// Handle error
}

if response != nil {
// Use the cached or newly generated response
}
```

### Redis Get and Set Methods

When implementing the Get and Set methods in your RedisRepository interface, it’s important to handle the marshalling and unmarshalling of the cache data:

- Get: Retrieve the data from Redis and unmarshal it into the appropriate type.
- Set: Marshal the data before storing it in Redis.

Example implementation:

```go
import (
	"context"
	"encoding/json"
	"github.com/kanmo/response_cache"
)

func (r *MyRedisRepository) Get(ctx context.Context, key string, value any) error {
    // Retrieve the data from Redis (assuming JSON stored)
    data, err := r.redisClient.Get(ctx, key).Result()
    if err != nil {
        return err
    }

    // Cast value to CacheData and unmarshal
    cacheData, ok := value.(*response_cache.CacheData)
    if !ok {
        return errors.New("failed to cast value to CacheData")
    }

    return cacheData.Unmarshal([]byte(data))
}

func (r *MyRedisRepository) Set(ctx context.Context, key string, value any) error {
    // Cast value to CacheData and marshal
    cacheData, ok := value.(*response_cache.CacheData)
    if !ok {
        return errors.New("failed to cast value to CacheData")
    }

    data, err := cacheData.Marshal()
    if err != nil {
        return err
    }

    // Store the data in Redis
    return r.redisClient.Set(ctx, key, data, 0).Err()
}
```

### Sample gRPC Interceptor

Below is a sample gRPC interceptor that uses the GetOrSetCache method to handle idempotent requests:

```go
import (
    "context"
    "google.golang.org/grpc"
    "github.com/kanmo/response_cache"
    "github.com/kanmo/response_cache/pb"
)

func (i *Interceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    idempotencyKey, existIdempotencyKey := getIdempotencyKey(req)
    if !existIdempotencyKey {
        return handler(ctx, req)
    }

    cacheKey := response_cache.GenerateCacheKey(userId, info.FullMethod, idempotencyKey)
    cache := response_cache.NewResponseCache(redisRepo)

    response, err := cache.GetOrSetCache(ctx, req, handler)
    if err != nil {
        return nil, err
    }

    return response, nil
}

func getIdempotencyKey(req interface{}) (string, bool) {
    if r, ok := req.(idempotencyKeyGetter); ok {
        return r.GetIdempotencyKey(), true
    }
    return "", false
}
```

## Example

Here’s an example of how to integrate this library into your project:

```go
package main

import (
	"context"
	"fmt"
	"github.com/kanmo/response_cache"
	"github.com/kanmo/response_cache/pb"
)

func main() {
	ctx := context.Background()
	cacheKey := response_cache.GenerateCacheKey("user123", "/example/method", "key123")
	cache := response_cache.NewResponseCache(redisClient)

	response, err := cache.GetOrSetCache(ctx, &pb.UserRequest{IdempotencyKey: "key123"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return &pb.UserResponse{UserId: "1234", UserName: "Test User"}, nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response:", response)
}
```

## Testing

To run tests for this library, use the following command:

```bash
go test ./...
```
Ensure that you have the necessary dependencies installed and properly configured, including protoc for compiling Protocol Buffers.

## License

This project is licensed under the MIT License - see the [LICENSE.txt](LICENSE.txt) for details.




