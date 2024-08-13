package response_cache

// IdempotencyKeyGetter is an interface that defines a method to get idempotency key
type IdempotencyKeyGetter interface {
	GetIdempotencyKey() string
}

// GetIdempotencyKey fetch idempotency key from request
func GetIdempotencyKey(req interface{}) (string, bool) {
	val, ok := req.(IdempotencyKeyGetter)
	if !ok {
		return "", false
	}

	return val.GetIdempotencyKey(), true
}
