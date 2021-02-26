// Code generated by protoc-gen-go. DO NOT EDIT.
// source: users.proto

package users

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type GetTrainingBalanceRequest struct {
	UserId               string   `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTrainingBalanceRequest) Reset()         { *m = GetTrainingBalanceRequest{} }
func (m *GetTrainingBalanceRequest) String() string { return proto.CompactTextString(m) }
func (*GetTrainingBalanceRequest) ProtoMessage()    {}
func (*GetTrainingBalanceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{0}
}

func (m *GetTrainingBalanceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTrainingBalanceRequest.Unmarshal(m, b)
}
func (m *GetTrainingBalanceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTrainingBalanceRequest.Marshal(b, m, deterministic)
}
func (m *GetTrainingBalanceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTrainingBalanceRequest.Merge(m, src)
}
func (m *GetTrainingBalanceRequest) XXX_Size() int {
	return xxx_messageInfo_GetTrainingBalanceRequest.Size(m)
}
func (m *GetTrainingBalanceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTrainingBalanceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTrainingBalanceRequest proto.InternalMessageInfo

func (m *GetTrainingBalanceRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

type GetTrainingBalanceResponse struct {
	Amount               int64    `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTrainingBalanceResponse) Reset()         { *m = GetTrainingBalanceResponse{} }
func (m *GetTrainingBalanceResponse) String() string { return proto.CompactTextString(m) }
func (*GetTrainingBalanceResponse) ProtoMessage()    {}
func (*GetTrainingBalanceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{1}
}

func (m *GetTrainingBalanceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTrainingBalanceResponse.Unmarshal(m, b)
}
func (m *GetTrainingBalanceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTrainingBalanceResponse.Marshal(b, m, deterministic)
}
func (m *GetTrainingBalanceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTrainingBalanceResponse.Merge(m, src)
}
func (m *GetTrainingBalanceResponse) XXX_Size() int {
	return xxx_messageInfo_GetTrainingBalanceResponse.Size(m)
}
func (m *GetTrainingBalanceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTrainingBalanceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetTrainingBalanceResponse proto.InternalMessageInfo

func (m *GetTrainingBalanceResponse) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type UpdateTrainingBalanceRequest struct {
	UserId               string   `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	AmountChange         int64    `protobuf:"varint,2,opt,name=amount_change,json=amountChange,proto3" json:"amount_change,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateTrainingBalanceRequest) Reset()         { *m = UpdateTrainingBalanceRequest{} }
func (m *UpdateTrainingBalanceRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateTrainingBalanceRequest) ProtoMessage()    {}
func (*UpdateTrainingBalanceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{2}
}

func (m *UpdateTrainingBalanceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateTrainingBalanceRequest.Unmarshal(m, b)
}
func (m *UpdateTrainingBalanceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateTrainingBalanceRequest.Marshal(b, m, deterministic)
}
func (m *UpdateTrainingBalanceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateTrainingBalanceRequest.Merge(m, src)
}
func (m *UpdateTrainingBalanceRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateTrainingBalanceRequest.Size(m)
}
func (m *UpdateTrainingBalanceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateTrainingBalanceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateTrainingBalanceRequest proto.InternalMessageInfo

func (m *UpdateTrainingBalanceRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *UpdateTrainingBalanceRequest) GetAmountChange() int64 {
	if m != nil {
		return m.AmountChange
	}
	return 0
}

type EmptyResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyResponse) Reset()         { *m = EmptyResponse{} }
func (m *EmptyResponse) String() string { return proto.CompactTextString(m) }
func (*EmptyResponse) ProtoMessage()    {}
func (*EmptyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{3}
}

func (m *EmptyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyResponse.Unmarshal(m, b)
}
func (m *EmptyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyResponse.Marshal(b, m, deterministic)
}
func (m *EmptyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyResponse.Merge(m, src)
}
func (m *EmptyResponse) XXX_Size() int {
	return xxx_messageInfo_EmptyResponse.Size(m)
}
func (m *EmptyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GetTrainingBalanceRequest)(nil), "users.GetTrainingBalanceRequest")
	proto.RegisterType((*GetTrainingBalanceResponse)(nil), "users.GetTrainingBalanceResponse")
	proto.RegisterType((*UpdateTrainingBalanceRequest)(nil), "users.UpdateTrainingBalanceRequest")
	proto.RegisterType((*EmptyResponse)(nil), "users.EmptyResponse")
}

func init() { proto.RegisterFile("users.proto", fileDescriptor_030765f334c86cea) }

var fileDescriptor_030765f334c86cea = []byte{
	// 230 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0x2d, 0x4e, 0x2d,
	0x2a, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0x4c, 0xb8, 0x24, 0xdd,
	0x53, 0x4b, 0x42, 0x8a, 0x12, 0x33, 0xf3, 0x32, 0xf3, 0xd2, 0x9d, 0x12, 0x73, 0x12, 0xf3, 0x92,
	0x53, 0x83, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0xc4, 0xb9, 0xd8, 0x41, 0xaa, 0xe2, 0x33,
	0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xd8, 0x40, 0x5c, 0xcf, 0x14, 0x25, 0x13, 0x2e,
	0x29, 0x6c, 0xba, 0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x85, 0xc4, 0xb8, 0xd8, 0x12, 0x73, 0xf3,
	0x4b, 0xf3, 0x4a, 0xc0, 0xba, 0x98, 0x83, 0xa0, 0x3c, 0xa5, 0x18, 0x2e, 0x99, 0xd0, 0x82, 0x94,
	0xc4, 0x92, 0x54, 0x12, 0xad, 0x13, 0x52, 0xe6, 0xe2, 0x85, 0x18, 0x11, 0x9f, 0x9c, 0x91, 0x98,
	0x97, 0x9e, 0x2a, 0xc1, 0x04, 0x36, 0x97, 0x07, 0x22, 0xe8, 0x0c, 0x16, 0x53, 0xe2, 0xe7, 0xe2,
	0x75, 0xcd, 0x2d, 0x28, 0xa9, 0x84, 0x39, 0xc3, 0xe8, 0x20, 0x23, 0x17, 0x4f, 0x28, 0xc8, 0x93,
	0xc1, 0xa9, 0x45, 0x65, 0x99, 0xc9, 0xa9, 0x42, 0xd1, 0x5c, 0x42, 0x98, 0xae, 0x16, 0x52, 0xd0,
	0x83, 0x04, 0x0b, 0xce, 0x60, 0x90, 0x52, 0xc4, 0xa3, 0x02, 0x62, 0x97, 0x12, 0x83, 0x50, 0x08,
	0x97, 0x28, 0x56, 0xcf, 0x09, 0x29, 0x43, 0x75, 0xe3, 0xf3, 0xba, 0x94, 0x08, 0x54, 0x11, 0x8a,
	0x0f, 0x94, 0x18, 0x92, 0xd8, 0xc0, 0x91, 0x65, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x3a, 0x3e,
	0x5b, 0x8a, 0xbb, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// UsersServiceClient is the client API for UsersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UsersServiceClient interface {
	GetTrainingBalance(ctx context.Context, in *GetTrainingBalanceRequest, opts ...grpc.CallOption) (*GetTrainingBalanceResponse, error)
	UpdateTrainingBalance(ctx context.Context, in *UpdateTrainingBalanceRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
}

type usersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersServiceClient(cc grpc.ClientConnInterface) UsersServiceClient {
	return &usersServiceClient{cc}
}

func (c *usersServiceClient) GetTrainingBalance(ctx context.Context, in *GetTrainingBalanceRequest, opts ...grpc.CallOption) (*GetTrainingBalanceResponse, error) {
	out := new(GetTrainingBalanceResponse)
	err := c.cc.Invoke(ctx, "/users.UsersService/GetTrainingBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) UpdateTrainingBalance(ctx context.Context, in *UpdateTrainingBalanceRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/users.UsersService/UpdateTrainingBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersServiceServer is the server API for UsersService service.
type UsersServiceServer interface {
	GetTrainingBalance(context.Context, *GetTrainingBalanceRequest) (*GetTrainingBalanceResponse, error)
	UpdateTrainingBalance(context.Context, *UpdateTrainingBalanceRequest) (*EmptyResponse, error)
}

// UnimplementedUsersServiceServer can be embedded to have forward compatible implementations.
type UnimplementedUsersServiceServer struct {
}

func (*UnimplementedUsersServiceServer) GetTrainingBalance(ctx context.Context, req *GetTrainingBalanceRequest) (*GetTrainingBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTrainingBalance not implemented")
}
func (*UnimplementedUsersServiceServer) UpdateTrainingBalance(ctx context.Context, req *UpdateTrainingBalanceRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTrainingBalance not implemented")
}

func RegisterUsersServiceServer(s *grpc.Server, srv UsersServiceServer) {
	s.RegisterService(&_UsersService_serviceDesc, srv)
}

func _UsersService_GetTrainingBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTrainingBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).GetTrainingBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UsersService/GetTrainingBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).GetTrainingBalance(ctx, req.(*GetTrainingBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_UpdateTrainingBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTrainingBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).UpdateTrainingBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UsersService/UpdateTrainingBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).UpdateTrainingBalance(ctx, req.(*UpdateTrainingBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UsersService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "users.UsersService",
	HandlerType: (*UsersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTrainingBalance",
			Handler:    _UsersService_GetTrainingBalance_Handler,
		},
		{
			MethodName: "UpdateTrainingBalance",
			Handler:    _UsersService_UpdateTrainingBalance_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "users.proto",
}
