// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/state.proto

package state

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_0493a22f06b3cb67, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type Command struct {
	Kind                 string   `protobuf:"bytes,1,opt,name=Kind,proto3" json:"Kind,omitempty"`
	ID                   int64    `protobuf:"varint,2,opt,name=ID,proto3" json:"ID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_0493a22f06b3cb67, []int{1}
}

func (m *Command) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Command.Unmarshal(m, b)
}
func (m *Command) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Command.Marshal(b, m, deterministic)
}
func (m *Command) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command.Merge(m, src)
}
func (m *Command) XXX_Size() int {
	return xxx_messageInfo_Command.Size(m)
}
func (m *Command) XXX_DiscardUnknown() {
	xxx_messageInfo_Command.DiscardUnknown(m)
}

var xxx_messageInfo_Command proto.InternalMessageInfo

func (m *Command) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *Command) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

type Data struct {
	Content              []byte   `protobuf:"bytes,1,opt,name=Content,proto3" json:"Content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Data) Reset()         { *m = Data{} }
func (m *Data) String() string { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()    {}
func (*Data) Descriptor() ([]byte, []int) {
	return fileDescriptor_0493a22f06b3cb67, []int{2}
}

func (m *Data) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Data.Unmarshal(m, b)
}
func (m *Data) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Data.Marshal(b, m, deterministic)
}
func (m *Data) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Data.Merge(m, src)
}
func (m *Data) XXX_Size() int {
	return xxx_messageInfo_Data.Size(m)
}
func (m *Data) XXX_DiscardUnknown() {
	xxx_messageInfo_Data.DiscardUnknown(m)
}

var xxx_messageInfo_Data proto.InternalMessageInfo

func (m *Data) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type Kind struct {
	Name                 string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Kind) Reset()         { *m = Kind{} }
func (m *Kind) String() string { return proto.CompactTextString(m) }
func (*Kind) ProtoMessage()    {}
func (*Kind) Descriptor() ([]byte, []int) {
	return fileDescriptor_0493a22f06b3cb67, []int{3}
}

func (m *Kind) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Kind.Unmarshal(m, b)
}
func (m *Kind) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Kind.Marshal(b, m, deterministic)
}
func (m *Kind) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Kind.Merge(m, src)
}
func (m *Kind) XXX_Size() int {
	return xxx_messageInfo_Kind.Size(m)
}
func (m *Kind) XXX_DiscardUnknown() {
	xxx_messageInfo_Kind.DiscardUnknown(m)
}

var xxx_messageInfo_Kind proto.InternalMessageInfo

func (m *Kind) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Domain struct {
	Name                 string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Domain) Reset()         { *m = Domain{} }
func (m *Domain) String() string { return proto.CompactTextString(m) }
func (*Domain) ProtoMessage()    {}
func (*Domain) Descriptor() ([]byte, []int) {
	return fileDescriptor_0493a22f06b3cb67, []int{4}
}

func (m *Domain) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Domain.Unmarshal(m, b)
}
func (m *Domain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Domain.Marshal(b, m, deterministic)
}
func (m *Domain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Domain.Merge(m, src)
}
func (m *Domain) XXX_Size() int {
	return xxx_messageInfo_Domain.Size(m)
}
func (m *Domain) XXX_DiscardUnknown() {
	xxx_messageInfo_Domain.DiscardUnknown(m)
}

var xxx_messageInfo_Domain proto.InternalMessageInfo

func (m *Domain) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Empty)(nil), "state.Empty")
	proto.RegisterType((*Command)(nil), "state.Command")
	proto.RegisterType((*Data)(nil), "state.Data")
	proto.RegisterType((*Kind)(nil), "state.Kind")
	proto.RegisterType((*Domain)(nil), "state.Domain")
}

func init() { proto.RegisterFile("proto/state.proto", fileDescriptor_0493a22f06b3cb67) }

var fileDescriptor_0493a22f06b3cb67 = []byte{
	// 224 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x2f, 0x2e, 0x49, 0x2c, 0x49, 0xd5, 0x03, 0xb3, 0x85, 0x58, 0xc1, 0x1c, 0x25, 0x76,
	0x2e, 0x56, 0xd7, 0xdc, 0x82, 0x92, 0x4a, 0x25, 0x5d, 0x2e, 0x76, 0xe7, 0xfc, 0xdc, 0xdc, 0xc4,
	0xbc, 0x14, 0x21, 0x21, 0x2e, 0x16, 0xef, 0xcc, 0xbc, 0x14, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce,
	0x20, 0x30, 0x5b, 0x88, 0x8f, 0x8b, 0xc9, 0xd3, 0x45, 0x82, 0x49, 0x81, 0x51, 0x83, 0x39, 0x88,
	0xc9, 0xd3, 0x45, 0x49, 0x81, 0x8b, 0xc5, 0x25, 0xb1, 0x24, 0x51, 0x48, 0x02, 0xa4, 0x2d, 0xaf,
	0x24, 0x35, 0xaf, 0x04, 0xac, 0x9c, 0x27, 0x08, 0xc6, 0x55, 0x92, 0x82, 0x98, 0x02, 0x32, 0xcd,
	0x2f, 0x31, 0x37, 0x15, 0x66, 0x1a, 0x88, 0xad, 0x24, 0xc3, 0xc5, 0xe6, 0x92, 0x9f, 0x9b, 0x98,
	0x99, 0x87, 0x4d, 0xd6, 0xa8, 0x93, 0x91, 0x8b, 0x35, 0x18, 0xe4, 0x3a, 0x21, 0x2d, 0x2e, 0x4e,
	0xf7, 0xd4, 0x12, 0xa8, 0x52, 0x1e, 0x3d, 0x88, 0xfb, 0xc1, 0xee, 0x95, 0xe2, 0x85, 0xf2, 0x20,
	0x92, 0x4a, 0x0c, 0x42, 0xaa, 0x5c, 0xac, 0x2e, 0x45, 0x20, 0x75, 0xdc, 0x50, 0x19, 0x90, 0xed,
	0x52, 0x30, 0x0e, 0xc8, 0xb1, 0x4a, 0x0c, 0x06, 0x8c, 0x42, 0xba, 0x5c, 0xdc, 0xc1, 0xa9, 0x79,
	0x29, 0x30, 0xbf, 0xf2, 0x41, 0xe5, 0xa1, 0x7c, 0x29, 0x14, 0x4b, 0x94, 0x18, 0x92, 0xd8, 0xc0,
	0xa1, 0x65, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x93, 0xc0, 0x53, 0xcc, 0x42, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// StateClient is the client API for State service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StateClient interface {
	GetDomain(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Domain, error)
	Drain(ctx context.Context, in *Kind, opts ...grpc.CallOption) (State_DrainClient, error)
	SendCommand(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Empty, error)
}

type stateClient struct {
	cc *grpc.ClientConn
}

func NewStateClient(cc *grpc.ClientConn) StateClient {
	return &stateClient{cc}
}

func (c *stateClient) GetDomain(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Domain, error) {
	out := new(Domain)
	err := c.cc.Invoke(ctx, "/state.State/GetDomain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateClient) Drain(ctx context.Context, in *Kind, opts ...grpc.CallOption) (State_DrainClient, error) {
	stream, err := c.cc.NewStream(ctx, &_State_serviceDesc.Streams[0], "/state.State/Drain", opts...)
	if err != nil {
		return nil, err
	}
	x := &stateDrainClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type State_DrainClient interface {
	Recv() (*Data, error)
	grpc.ClientStream
}

type stateDrainClient struct {
	grpc.ClientStream
}

func (x *stateDrainClient) Recv() (*Data, error) {
	m := new(Data)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *stateClient) SendCommand(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/state.State/SendCommand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StateServer is the server API for State service.
type StateServer interface {
	GetDomain(context.Context, *Empty) (*Domain, error)
	Drain(*Kind, State_DrainServer) error
	SendCommand(context.Context, *Command) (*Empty, error)
}

// UnimplementedStateServer can be embedded to have forward compatible implementations.
type UnimplementedStateServer struct {
}

func (*UnimplementedStateServer) GetDomain(ctx context.Context, req *Empty) (*Domain, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDomain not implemented")
}
func (*UnimplementedStateServer) Drain(req *Kind, srv State_DrainServer) error {
	return status.Errorf(codes.Unimplemented, "method Drain not implemented")
}
func (*UnimplementedStateServer) SendCommand(ctx context.Context, req *Command) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCommand not implemented")
}

func RegisterStateServer(s *grpc.Server, srv StateServer) {
	s.RegisterService(&_State_serviceDesc, srv)
}

func _State_GetDomain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateServer).GetDomain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/state.State/GetDomain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateServer).GetDomain(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _State_Drain_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Kind)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StateServer).Drain(m, &stateDrainServer{stream})
}

type State_DrainServer interface {
	Send(*Data) error
	grpc.ServerStream
}

type stateDrainServer struct {
	grpc.ServerStream
}

func (x *stateDrainServer) Send(m *Data) error {
	return x.ServerStream.SendMsg(m)
}

func _State_SendCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Command)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateServer).SendCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/state.State/SendCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateServer).SendCommand(ctx, req.(*Command))
	}
	return interceptor(ctx, in, info, handler)
}

var _State_serviceDesc = grpc.ServiceDesc{
	ServiceName: "state.State",
	HandlerType: (*StateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDomain",
			Handler:    _State_GetDomain_Handler,
		},
		{
			MethodName: "SendCommand",
			Handler:    _State_SendCommand_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Drain",
			Handler:       _State_Drain_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/state.proto",
}
