// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/keyapi/key.proto

/*
Package keyapi is a generated protocol buffer package.

It is generated from these files:
	pkg/keyapi/key.proto

It has these top-level messages:
	AddPublicKeysRequest
	AddPublicKeysResponse
	GetPublicKeysRequest
	GetPublicKeysResponse
	PublicKeyDetail
*/
package keyapi

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type KeyType int32

const (
	KeyType_AUTHOR KeyType = 0
	KeyType_READER KeyType = 1
)

var KeyType_name = map[int32]string{
	0: "AUTHOR",
	1: "READER",
}
var KeyType_value = map[string]int32{
	"AUTHOR": 0,
	"READER": 1,
}

func (x KeyType) String() string {
	return proto.EnumName(KeyType_name, int32(x))
}
func (KeyType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type AddPublicKeysRequest struct {
	EntityId   string   `protobuf:"bytes,1,opt,name=entity_id,json=entityId" json:"entity_id,omitempty"`
	KeyType    string   `protobuf:"bytes,2,opt,name=key_type,json=keyType" json:"key_type,omitempty"`
	PublicKeys [][]byte `protobuf:"bytes,3,rep,name=public_keys,json=publicKeys,proto3" json:"public_keys,omitempty"`
}

func (m *AddPublicKeysRequest) Reset()                    { *m = AddPublicKeysRequest{} }
func (m *AddPublicKeysRequest) String() string            { return proto.CompactTextString(m) }
func (*AddPublicKeysRequest) ProtoMessage()               {}
func (*AddPublicKeysRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *AddPublicKeysRequest) GetEntityId() string {
	if m != nil {
		return m.EntityId
	}
	return ""
}

func (m *AddPublicKeysRequest) GetKeyType() string {
	if m != nil {
		return m.KeyType
	}
	return ""
}

func (m *AddPublicKeysRequest) GetPublicKeys() [][]byte {
	if m != nil {
		return m.PublicKeys
	}
	return nil
}

type AddPublicKeysResponse struct {
}

func (m *AddPublicKeysResponse) Reset()                    { *m = AddPublicKeysResponse{} }
func (m *AddPublicKeysResponse) String() string            { return proto.CompactTextString(m) }
func (*AddPublicKeysResponse) ProtoMessage()               {}
func (*AddPublicKeysResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type GetPublicKeysRequest struct {
	PublicKeys [][]byte `protobuf:"bytes,3,rep,name=public_keys,json=publicKeys,proto3" json:"public_keys,omitempty"`
}

func (m *GetPublicKeysRequest) Reset()                    { *m = GetPublicKeysRequest{} }
func (m *GetPublicKeysRequest) String() string            { return proto.CompactTextString(m) }
func (*GetPublicKeysRequest) ProtoMessage()               {}
func (*GetPublicKeysRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetPublicKeysRequest) GetPublicKeys() [][]byte {
	if m != nil {
		return m.PublicKeys
	}
	return nil
}

type GetPublicKeysResponse struct {
	PublicKeyDetails []*PublicKeyDetail `protobuf:"bytes,1,rep,name=public_key_details,json=publicKeyDetails" json:"public_key_details,omitempty"`
}

func (m *GetPublicKeysResponse) Reset()                    { *m = GetPublicKeysResponse{} }
func (m *GetPublicKeysResponse) String() string            { return proto.CompactTextString(m) }
func (*GetPublicKeysResponse) ProtoMessage()               {}
func (*GetPublicKeysResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GetPublicKeysResponse) GetPublicKeyDetails() []*PublicKeyDetail {
	if m != nil {
		return m.PublicKeyDetails
	}
	return nil
}

type PublicKeyDetail struct {
	PublicKey []byte  `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	EntityId  string  `protobuf:"bytes,2,opt,name=entity_id,json=entityId" json:"entity_id,omitempty"`
	KeyType   KeyType `protobuf:"varint,3,opt,name=key_type,json=keyType,enum=keyapi.KeyType" json:"key_type,omitempty"`
}

func (m *PublicKeyDetail) Reset()                    { *m = PublicKeyDetail{} }
func (m *PublicKeyDetail) String() string            { return proto.CompactTextString(m) }
func (*PublicKeyDetail) ProtoMessage()               {}
func (*PublicKeyDetail) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *PublicKeyDetail) GetPublicKey() []byte {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

func (m *PublicKeyDetail) GetEntityId() string {
	if m != nil {
		return m.EntityId
	}
	return ""
}

func (m *PublicKeyDetail) GetKeyType() KeyType {
	if m != nil {
		return m.KeyType
	}
	return KeyType_AUTHOR
}

func init() {
	proto.RegisterType((*AddPublicKeysRequest)(nil), "keyapi.AddPublicKeysRequest")
	proto.RegisterType((*AddPublicKeysResponse)(nil), "keyapi.AddPublicKeysResponse")
	proto.RegisterType((*GetPublicKeysRequest)(nil), "keyapi.GetPublicKeysRequest")
	proto.RegisterType((*GetPublicKeysResponse)(nil), "keyapi.GetPublicKeysResponse")
	proto.RegisterType((*PublicKeyDetail)(nil), "keyapi.PublicKeyDetail")
	proto.RegisterEnum("keyapi.KeyType", KeyType_name, KeyType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Key service

type KeyClient interface {
	AddPublicKeys(ctx context.Context, in *AddPublicKeysRequest, opts ...grpc.CallOption) (*AddPublicKeysResponse, error)
	GetPublicKeys(ctx context.Context, in *GetPublicKeysRequest, opts ...grpc.CallOption) (*GetPublicKeysResponse, error)
}

type keyClient struct {
	cc *grpc.ClientConn
}

func NewKeyClient(cc *grpc.ClientConn) KeyClient {
	return &keyClient{cc}
}

func (c *keyClient) AddPublicKeys(ctx context.Context, in *AddPublicKeysRequest, opts ...grpc.CallOption) (*AddPublicKeysResponse, error) {
	out := new(AddPublicKeysResponse)
	err := grpc.Invoke(ctx, "/keyapi.Key/AddPublicKeys", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyClient) GetPublicKeys(ctx context.Context, in *GetPublicKeysRequest, opts ...grpc.CallOption) (*GetPublicKeysResponse, error) {
	out := new(GetPublicKeysResponse)
	err := grpc.Invoke(ctx, "/keyapi.Key/GetPublicKeys", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Key service

type KeyServer interface {
	AddPublicKeys(context.Context, *AddPublicKeysRequest) (*AddPublicKeysResponse, error)
	GetPublicKeys(context.Context, *GetPublicKeysRequest) (*GetPublicKeysResponse, error)
}

func RegisterKeyServer(s *grpc.Server, srv KeyServer) {
	s.RegisterService(&_Key_serviceDesc, srv)
}

func _Key_AddPublicKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddPublicKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyServer).AddPublicKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/keyapi.Key/AddPublicKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyServer).AddPublicKeys(ctx, req.(*AddPublicKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Key_GetPublicKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPublicKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyServer).GetPublicKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/keyapi.Key/GetPublicKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyServer).GetPublicKeys(ctx, req.(*GetPublicKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Key_serviceDesc = grpc.ServiceDesc{
	ServiceName: "keyapi.Key",
	HandlerType: (*KeyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddPublicKeys",
			Handler:    _Key_AddPublicKeys_Handler,
		},
		{
			MethodName: "GetPublicKeys",
			Handler:    _Key_GetPublicKeys_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/keyapi/key.proto",
}

func init() { proto.RegisterFile("pkg/keyapi/key.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xcd, 0x6a, 0xb3, 0x40,
	0x14, 0x86, 0x33, 0x11, 0xf2, 0x73, 0x92, 0xef, 0x4b, 0x18, 0x0c, 0xb1, 0x3f, 0xa1, 0xd6, 0x95,
	0x64, 0x61, 0xc1, 0x2e, 0xba, 0x16, 0x22, 0x6d, 0x09, 0xb4, 0x65, 0x48, 0xb7, 0x95, 0xa4, 0x1e,
	0x8a, 0x18, 0xe2, 0xb4, 0x33, 0x59, 0xcc, 0x05, 0xf5, 0x3e, 0x8b, 0x9a, 0x68, 0xb5, 0x96, 0xae,
	0x94, 0x73, 0x7c, 0x9f, 0x19, 0x9e, 0x57, 0xd0, 0x79, 0xfc, 0x76, 0x15, 0xa3, 0x5a, 0xf3, 0x28,
	0x7d, 0x38, 0xfc, 0x23, 0x91, 0x09, 0xed, 0xe4, 0x13, 0x2b, 0x01, 0xdd, 0x0b, 0xc3, 0xa7, 0xfd,
	0x66, 0x1b, 0xbd, 0x2e, 0x51, 0x09, 0x86, 0xef, 0x7b, 0x14, 0x92, 0x9e, 0x41, 0x1f, 0x77, 0x32,
	0x92, 0x2a, 0x88, 0x42, 0x83, 0x98, 0xc4, 0xee, 0xb3, 0x5e, 0x3e, 0xb8, 0x0f, 0xe9, 0x09, 0xf4,
	0x62, 0x54, 0x81, 0x54, 0x1c, 0x8d, 0x76, 0xb6, 0xeb, 0xc6, 0xa8, 0x56, 0x8a, 0x23, 0xbd, 0x80,
	0x01, 0xcf, 0x60, 0x41, 0x8c, 0x4a, 0x18, 0x9a, 0xa9, 0xd9, 0x43, 0x06, 0xbc, 0xe0, 0x5b, 0x53,
	0x98, 0xd4, 0x0e, 0x14, 0x3c, 0xd9, 0x09, 0xb4, 0x6e, 0x40, 0xbf, 0x45, 0xf9, 0xf3, 0x26, 0x7f,
	0x12, 0x5f, 0x60, 0x52, 0x0b, 0xe6, 0x44, 0xea, 0x03, 0x2d, 0x93, 0x41, 0x88, 0x72, 0x1d, 0x6d,
	0x85, 0x41, 0x4c, 0xcd, 0x1e, 0xb8, 0x53, 0x27, 0x17, 0xe0, 0x14, 0xb9, 0x45, 0xb6, 0x67, 0x63,
	0x5e, 0x1d, 0x08, 0x4b, 0xc1, 0xa8, 0xf6, 0x11, 0x9d, 0x01, 0x94, 0xe4, 0x4c, 0xcf, 0x90, 0xf5,
	0x8b, 0x60, 0x55, 0x5e, 0xbb, 0x26, 0x6f, 0xfe, 0x4d, 0x9e, 0x66, 0x12, 0xfb, 0xbf, 0x3b, 0x3a,
	0xde, 0x65, 0x99, 0x4b, 0x2c, 0x6c, 0xce, 0x2f, 0xa1, 0x7b, 0x98, 0x51, 0x80, 0x8e, 0xf7, 0xbc,
	0xba, 0x7b, 0x64, 0xe3, 0x56, 0xfa, 0xce, 0x7c, 0x6f, 0xe1, 0xb3, 0x31, 0x71, 0x3f, 0x09, 0x68,
	0xe9, 0x99, 0x0f, 0xf0, 0xaf, 0xe2, 0x95, 0x9e, 0x1f, 0xa9, 0x4d, 0xfd, 0x9e, 0xce, 0x7e, 0xd9,
	0x1e, 0xca, 0x68, 0xa5, 0xbc, 0x8a, 0xd5, 0x92, 0xd7, 0xd4, 0x52, 0xc9, 0x6b, 0xac, 0xc2, 0x6a,
	0x6d, 0x3a, 0xd9, 0x7f, 0x77, 0xfd, 0x15, 0x00, 0x00, 0xff, 0xff, 0x19, 0x59, 0x61, 0xd4, 0x8f,
	0x02, 0x00, 0x00,
}
