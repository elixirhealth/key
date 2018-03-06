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
	SamplePublicKeysRequest
	SamplePublicKeysResponse
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
	KeyType    KeyType  `protobuf:"varint,2,opt,name=key_type,json=keyType,enum=keyapi.KeyType" json:"key_type,omitempty"`
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

func (m *AddPublicKeysRequest) GetKeyType() KeyType {
	if m != nil {
		return m.KeyType
	}
	return KeyType_AUTHOR
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

type SamplePublicKeysRequest struct {
	OfEntityId        string `protobuf:"bytes,1,opt,name=of_entity_id,json=ofEntityId" json:"of_entity_id,omitempty"`
	NPublicKeys       uint32 `protobuf:"varint,2,opt,name=n_public_keys,json=nPublicKeys" json:"n_public_keys,omitempty"`
	RequesterEntityId string `protobuf:"bytes,3,opt,name=requester_entity_id,json=requesterEntityId" json:"requester_entity_id,omitempty"`
}

func (m *SamplePublicKeysRequest) Reset()                    { *m = SamplePublicKeysRequest{} }
func (m *SamplePublicKeysRequest) String() string            { return proto.CompactTextString(m) }
func (*SamplePublicKeysRequest) ProtoMessage()               {}
func (*SamplePublicKeysRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *SamplePublicKeysRequest) GetOfEntityId() string {
	if m != nil {
		return m.OfEntityId
	}
	return ""
}

func (m *SamplePublicKeysRequest) GetNPublicKeys() uint32 {
	if m != nil {
		return m.NPublicKeys
	}
	return 0
}

func (m *SamplePublicKeysRequest) GetRequesterEntityId() string {
	if m != nil {
		return m.RequesterEntityId
	}
	return ""
}

type SamplePublicKeysResponse struct {
	PublicKeyDetails []*PublicKeyDetail `protobuf:"bytes,1,rep,name=public_key_details,json=publicKeyDetails" json:"public_key_details,omitempty"`
}

func (m *SamplePublicKeysResponse) Reset()                    { *m = SamplePublicKeysResponse{} }
func (m *SamplePublicKeysResponse) String() string            { return proto.CompactTextString(m) }
func (*SamplePublicKeysResponse) ProtoMessage()               {}
func (*SamplePublicKeysResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *SamplePublicKeysResponse) GetPublicKeyDetails() []*PublicKeyDetail {
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
func (*PublicKeyDetail) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

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
	proto.RegisterType((*SamplePublicKeysRequest)(nil), "keyapi.SamplePublicKeysRequest")
	proto.RegisterType((*SamplePublicKeysResponse)(nil), "keyapi.SamplePublicKeysResponse")
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
	SamplePublicKeys(ctx context.Context, in *SamplePublicKeysRequest, opts ...grpc.CallOption) (*SamplePublicKeysResponse, error)
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

func (c *keyClient) SamplePublicKeys(ctx context.Context, in *SamplePublicKeysRequest, opts ...grpc.CallOption) (*SamplePublicKeysResponse, error) {
	out := new(SamplePublicKeysResponse)
	err := grpc.Invoke(ctx, "/keyapi.Key/SamplePublicKeys", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Key service

type KeyServer interface {
	AddPublicKeys(context.Context, *AddPublicKeysRequest) (*AddPublicKeysResponse, error)
	GetPublicKeys(context.Context, *GetPublicKeysRequest) (*GetPublicKeysResponse, error)
	SamplePublicKeys(context.Context, *SamplePublicKeysRequest) (*SamplePublicKeysResponse, error)
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

func _Key_SamplePublicKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SamplePublicKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyServer).SamplePublicKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/keyapi.Key/SamplePublicKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyServer).SamplePublicKeys(ctx, req.(*SamplePublicKeysRequest))
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
		{
			MethodName: "SamplePublicKeys",
			Handler:    _Key_SamplePublicKeys_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/keyapi/key.proto",
}

func init() { proto.RegisterFile("pkg/keyapi/key.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 404 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0x4f, 0x6f, 0xaa, 0x40,
	0x10, 0xc0, 0x5d, 0x49, 0xfc, 0x33, 0xea, 0x93, 0xb7, 0x4f, 0x23, 0xf1, 0x3d, 0x23, 0x8f, 0x13,
	0xf1, 0x40, 0x13, 0x7b, 0xe8, 0xd9, 0x44, 0xd2, 0x36, 0x26, 0x6d, 0xb3, 0xb5, 0xe9, 0xad, 0x04,
	0xcb, 0xda, 0x10, 0x2c, 0x6c, 0x05, 0x0f, 0x7b, 0xeb, 0x37, 0xe8, 0x37, 0x6e, 0x1a, 0x40, 0x40,
	0xf0, 0x4f, 0x2f, 0x3d, 0xb9, 0x99, 0xd9, 0xf9, 0xed, 0xcc, 0xf8, 0x03, 0x3a, 0xcc, 0x79, 0x39,
	0x73, 0x28, 0x37, 0x99, 0x1d, 0xfe, 0x68, 0x6c, 0xed, 0x05, 0x1e, 0xae, 0xc4, 0x11, 0xe5, 0x1d,
	0x41, 0x67, 0x62, 0x59, 0x77, 0x9b, 0xc5, 0xca, 0x7e, 0x9e, 0x51, 0xee, 0x13, 0xfa, 0xb6, 0xa1,
	0x7e, 0x80, 0xff, 0x42, 0x9d, 0xba, 0x81, 0x1d, 0x70, 0xc3, 0xb6, 0x24, 0x24, 0x23, 0xb5, 0x4e,
	0x6a, 0x71, 0xe0, 0xda, 0xc2, 0x23, 0xa8, 0x39, 0x94, 0x1b, 0x01, 0x67, 0x54, 0x2a, 0xcb, 0x48,
	0xfd, 0x35, 0x6e, 0x6b, 0x31, 0x50, 0x9b, 0x51, 0x3e, 0xe7, 0x8c, 0x92, 0xaa, 0x13, 0x1f, 0xf0,
	0x10, 0x1a, 0x2c, 0xa2, 0x1b, 0x0e, 0xe5, 0xbe, 0x24, 0xc8, 0x82, 0xda, 0x24, 0xc0, 0xd2, 0x07,
	0x95, 0x1e, 0x74, 0x0b, 0x1d, 0xf8, 0xcc, 0x73, 0x7d, 0xaa, 0x5c, 0x40, 0xe7, 0x92, 0x06, 0xfb,
	0xad, 0x7d, 0x4b, 0x7c, 0x82, 0x6e, 0xa1, 0x30, 0x26, 0x62, 0x1d, 0x70, 0x56, 0x69, 0x58, 0x34,
	0x30, 0xed, 0x95, 0x2f, 0x21, 0x59, 0x50, 0x1b, 0xe3, 0x5e, 0x32, 0x41, 0x5a, 0x37, 0x8d, 0xf2,
	0x44, 0x64, 0xf9, 0x80, 0xaf, 0x7c, 0x20, 0xe8, 0xdd, 0x9b, 0xaf, 0x6c, 0x45, 0xf7, 0x9b, 0x93,
	0xa1, 0xe9, 0x2d, 0x8d, 0xe2, 0xea, 0xc0, 0x5b, 0xea, 0xc9, 0xf2, 0x14, 0x68, 0xb9, 0xc6, 0xee,
	0x00, 0xe1, 0x06, 0x5b, 0xa4, 0xe1, 0x66, 0x30, 0xac, 0xc1, 0x9f, 0x75, 0x0c, 0xa4, 0xeb, 0x1d,
	0x98, 0x10, 0xc1, 0x7e, 0xa7, 0xa9, 0x84, 0xa9, 0x98, 0x20, 0xed, 0x37, 0xf4, 0xb3, 0x43, 0x73,
	0x68, 0x17, 0x2e, 0xe1, 0x01, 0x40, 0x46, 0x8e, 0x26, 0x6d, 0x92, 0x7a, 0x5a, 0x98, 0x57, 0xa8,
	0x7c, 0x42, 0x21, 0xe1, 0xb4, 0x42, 0xa3, 0xff, 0x50, 0xdd, 0xc6, 0x30, 0x40, 0x65, 0xf2, 0x30,
	0xbf, 0xba, 0x25, 0x62, 0x29, 0x3c, 0x13, 0x7d, 0x32, 0xd5, 0x89, 0x88, 0xc6, 0x9f, 0x08, 0x84,
	0xf0, 0xcd, 0x1b, 0x68, 0xe5, 0x64, 0xc2, 0xff, 0x12, 0xea, 0x21, 0xcb, 0xfb, 0x83, 0x23, 0xd9,
	0xad, 0x81, 0xa5, 0x90, 0x97, 0x53, 0x29, 0xe3, 0x1d, 0x52, 0x33, 0xe3, 0x1d, 0xf4, 0x4f, 0x29,
	0xe1, 0x47, 0x10, 0x8b, 0x7f, 0x14, 0x1e, 0x26, 0x45, 0x47, 0x9c, 0xea, 0xcb, 0xc7, 0x2f, 0x24,
	0xe0, 0x45, 0x25, 0xfa, 0xae, 0xcf, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x90, 0xe3, 0xb4, 0xdb,
	0xef, 0x03, 0x00, 0x00,
}
