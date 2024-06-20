package structs

import (
	"encoding/json"
	"google.golang.org/protobuf/runtime/protoimpl"
	"testing"
)

func TestStruct_StructToMap(t *testing.T) {
	type Test struct {
		A string `json:"a"`
		B int    `json:"b"`
	}

	type UserModel struct {
		state         protoimpl.MessageState
		sizeCache     protoimpl.SizeCache
		unknownFields protoimpl.UnknownFields

		Id        int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
		Nickname  string `protobuf:"bytes,2,opt,name=nickname,proto3" json:"nickname,omitempty"`
		AvatorUrl string `protobuf:"bytes,3,opt,name=avatorUrl,proto3" json:"avatorUrl,omitempty"`
		Mobile    string `protobuf:"bytes,4,opt,name=mobile,proto3" json:"mobile,omitempty"`
		Gender    int32  `protobuf:"varint,5,opt,name=gender,proto3" json:"gender,omitempty"`
		CreatedAt string `protobuf:"bytes,6,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
		UpdatedAt string `protobuf:"bytes,7,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
		Status    int32  `protobuf:"varint,8,opt,name=status,proto3" json:"status,omitempty"`
		Test      *Test  `protobuf:"varint,9,opt,name=test,proto3" json:"test,omitempty"`
	}
	user := &UserModel{
		Test: &Test{},
	}

	st, err := New(JSON, user)
	if err != nil {
		t.Fatal(err)
	}

	m := st.ProtoStructToGormMap(WithIgnoreKeys(append(DefaultIgnoreKeys, "id")...))
	b, _ := json.Marshal(m)
	t.Log(string(b))
}
