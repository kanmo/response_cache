package response_cache

import (
	"reflect"
	pb "response_cache/pb/test/proto"
	"testing"
)

func init() {
	Registry.Register("UserResponse", reflect.TypeOf(pb.UserResponse{}))
}

func TestCacheData_MarshalUnmarshal(t *testing.T) {
	// テスト用のUserResponseインスタンスを作成
	originalValue := &pb.UserResponse{
		UserId:   "1234",
		UserName: "Test User",
	}

	// CacheData構造体を作成し、Marshalをテスト
	methodName := reflect.TypeOf(&pb.UserResponse{}).Elem().Name()
	cacheData := Create(methodName, "1234", "idempotencyKey")
	cacheData.Value = originalValue

	marshaledData, err := cacheData.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	newCacheData := &CacheData{}
	err = newCacheData.Unmarshal(marshaledData)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// assert that the unmarshal value is the same as the original value
	unmarshalValue, ok := newCacheData.UnmarshalValue.(*pb.UserResponse)
	if !ok {
		t.Fatalf("Unmarshaled value is not a UserResponse")
	}
	if !reflect.DeepEqual(unmarshalValue, originalValue) {
		t.Errorf("Unmarshaled value does not match original. Got %v, want %v", unmarshalValue, originalValue)
	}

	// verify that the TypeName is correct
	if newCacheData.TypeName != "UserResponse" {
		t.Errorf("TypeName is incorrect. Got %s, want %s", newCacheData.TypeName, "UserResponse")
	}
}
