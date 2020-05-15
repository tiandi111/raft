// Code generated by protoc-gen-go. DO NOT EDIT.
// source: raft.proto

package raft

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

type AppendEntriesRequest struct {
	Term                 int32    `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId             int32    `protobuf:"varint,2,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	PrevLogIndex         int32    `protobuf:"varint,3,opt,name=prev_log_index,json=prevLogIndex,proto3" json:"prev_log_index,omitempty"`
	PrevLogTerm          int32    `protobuf:"varint,4,opt,name=prev_log_term,json=prevLogTerm,proto3" json:"prev_log_term,omitempty"`
	LeaderCommit         int32    `protobuf:"varint,5,opt,name=leader_commit,json=leaderCommit,proto3" json:"leader_commit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AppendEntriesRequest) Reset()         { *m = AppendEntriesRequest{} }
func (m *AppendEntriesRequest) String() string { return proto.CompactTextString(m) }
func (*AppendEntriesRequest) ProtoMessage()    {}
func (*AppendEntriesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{0}
}

func (m *AppendEntriesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AppendEntriesRequest.Unmarshal(m, b)
}
func (m *AppendEntriesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AppendEntriesRequest.Marshal(b, m, deterministic)
}
func (m *AppendEntriesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AppendEntriesRequest.Merge(m, src)
}
func (m *AppendEntriesRequest) XXX_Size() int {
	return xxx_messageInfo_AppendEntriesRequest.Size(m)
}
func (m *AppendEntriesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AppendEntriesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AppendEntriesRequest proto.InternalMessageInfo

func (m *AppendEntriesRequest) GetTerm() int32 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesRequest) GetLeaderId() int32 {
	if m != nil {
		return m.LeaderId
	}
	return 0
}

func (m *AppendEntriesRequest) GetPrevLogIndex() int32 {
	if m != nil {
		return m.PrevLogIndex
	}
	return 0
}

func (m *AppendEntriesRequest) GetPrevLogTerm() int32 {
	if m != nil {
		return m.PrevLogTerm
	}
	return 0
}

func (m *AppendEntriesRequest) GetLeaderCommit() int32 {
	if m != nil {
		return m.LeaderCommit
	}
	return 0
}

type AppendEntriesResponse struct {
	Term                 int32    `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Success              bool     `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AppendEntriesResponse) Reset()         { *m = AppendEntriesResponse{} }
func (m *AppendEntriesResponse) String() string { return proto.CompactTextString(m) }
func (*AppendEntriesResponse) ProtoMessage()    {}
func (*AppendEntriesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{1}
}

func (m *AppendEntriesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AppendEntriesResponse.Unmarshal(m, b)
}
func (m *AppendEntriesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AppendEntriesResponse.Marshal(b, m, deterministic)
}
func (m *AppendEntriesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AppendEntriesResponse.Merge(m, src)
}
func (m *AppendEntriesResponse) XXX_Size() int {
	return xxx_messageInfo_AppendEntriesResponse.Size(m)
}
func (m *AppendEntriesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AppendEntriesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AppendEntriesResponse proto.InternalMessageInfo

func (m *AppendEntriesResponse) GetTerm() int32 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type RequestVoteRequest struct {
	Term                 int32    `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	CandidateId          int32    `protobuf:"varint,2,opt,name=candidate_id,json=candidateId,proto3" json:"candidate_id,omitempty"`
	LastLogIndex         int32    `protobuf:"varint,3,opt,name=last_log_index,json=lastLogIndex,proto3" json:"last_log_index,omitempty"`
	LastLogTerm          int32    `protobuf:"varint,4,opt,name=last_log_term,json=lastLogTerm,proto3" json:"last_log_term,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestVoteRequest) Reset()         { *m = RequestVoteRequest{} }
func (m *RequestVoteRequest) String() string { return proto.CompactTextString(m) }
func (*RequestVoteRequest) ProtoMessage()    {}
func (*RequestVoteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{2}
}

func (m *RequestVoteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestVoteRequest.Unmarshal(m, b)
}
func (m *RequestVoteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestVoteRequest.Marshal(b, m, deterministic)
}
func (m *RequestVoteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestVoteRequest.Merge(m, src)
}
func (m *RequestVoteRequest) XXX_Size() int {
	return xxx_messageInfo_RequestVoteRequest.Size(m)
}
func (m *RequestVoteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestVoteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RequestVoteRequest proto.InternalMessageInfo

func (m *RequestVoteRequest) GetTerm() int32 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteRequest) GetCandidateId() int32 {
	if m != nil {
		return m.CandidateId
	}
	return 0
}

func (m *RequestVoteRequest) GetLastLogIndex() int32 {
	if m != nil {
		return m.LastLogIndex
	}
	return 0
}

func (m *RequestVoteRequest) GetLastLogTerm() int32 {
	if m != nil {
		return m.LastLogTerm
	}
	return 0
}

type RequestVoteResponse struct {
	Term                 int32    `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	VoteGranted          bool     `protobuf:"varint,2,opt,name=vote_granted,json=voteGranted,proto3" json:"vote_granted,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestVoteResponse) Reset()         { *m = RequestVoteResponse{} }
func (m *RequestVoteResponse) String() string { return proto.CompactTextString(m) }
func (*RequestVoteResponse) ProtoMessage()    {}
func (*RequestVoteResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{3}
}

func (m *RequestVoteResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestVoteResponse.Unmarshal(m, b)
}
func (m *RequestVoteResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestVoteResponse.Marshal(b, m, deterministic)
}
func (m *RequestVoteResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestVoteResponse.Merge(m, src)
}
func (m *RequestVoteResponse) XXX_Size() int {
	return xxx_messageInfo_RequestVoteResponse.Size(m)
}
func (m *RequestVoteResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestVoteResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RequestVoteResponse proto.InternalMessageInfo

func (m *RequestVoteResponse) GetTerm() int32 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteResponse) GetVoteGranted() bool {
	if m != nil {
		return m.VoteGranted
	}
	return false
}

func init() {
	proto.RegisterType((*AppendEntriesRequest)(nil), "raft.AppendEntriesRequest")
	proto.RegisterType((*AppendEntriesResponse)(nil), "raft.AppendEntriesResponse")
	proto.RegisterType((*RequestVoteRequest)(nil), "raft.RequestVoteRequest")
	proto.RegisterType((*RequestVoteResponse)(nil), "raft.RequestVoteResponse")
}

func init() { proto.RegisterFile("raft.proto", fileDescriptor_b042552c306ae59b) }

var fileDescriptor_b042552c306ae59b = []byte{
	// 328 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xc1, 0x4e, 0xfa, 0x40,
	0x10, 0xc6, 0xff, 0xfd, 0x5b, 0x14, 0xa7, 0xe0, 0x61, 0xd5, 0xa4, 0xc2, 0x45, 0xaa, 0x07, 0x4f,
	0x1c, 0xf4, 0x09, 0x8c, 0x12, 0x83, 0xe1, 0xd4, 0x18, 0xaf, 0xcd, 0xda, 0x1d, 0x48, 0x13, 0xd8,
	0xad, 0xbb, 0x03, 0xf1, 0x45, 0x4c, 0x7c, 0x12, 0x9f, 0xcf, 0xec, 0x6e, 0xc1, 0x22, 0x0d, 0xb7,
	0xdd, 0xdf, 0x7c, 0x99, 0xf9, 0xe6, 0xcb, 0x00, 0x68, 0x3e, 0xa5, 0x61, 0xa9, 0x15, 0x29, 0x16,
	0xda, 0x77, 0xf2, 0x1d, 0xc0, 0xd9, 0x7d, 0x59, 0xa2, 0x14, 0x23, 0x49, 0xba, 0x40, 0x93, 0xe2,
	0xfb, 0x12, 0x0d, 0x31, 0x06, 0x21, 0xa1, 0x5e, 0xc4, 0xc1, 0x65, 0x70, 0xd3, 0x4a, 0xdd, 0x9b,
	0xf5, 0xe1, 0x78, 0x8e, 0x5c, 0xa0, 0xce, 0x0a, 0x11, 0xff, 0x77, 0x85, 0xb6, 0x07, 0x63, 0xc1,
	0xae, 0xe1, 0xa4, 0xd4, 0xb8, 0xca, 0xe6, 0x6a, 0x96, 0x15, 0x52, 0xe0, 0x47, 0x7c, 0xe0, 0x14,
	0x1d, 0x4b, 0x27, 0x6a, 0x36, 0xb6, 0x8c, 0x25, 0xd0, 0xdd, 0xa8, 0x5c, 0xff, 0xd0, 0x89, 0xa2,
	0x4a, 0xf4, 0x62, 0xc7, 0x5c, 0x41, 0xb7, 0x1a, 0x93, 0xab, 0xc5, 0xa2, 0xa0, 0xb8, 0xe5, 0x1b,
	0x79, 0xf8, 0xe0, 0x58, 0x32, 0x82, 0xf3, 0x3f, 0xbe, 0x4d, 0xa9, 0xa4, 0xc1, 0x46, 0xe3, 0x31,
	0x1c, 0x99, 0x65, 0x9e, 0xa3, 0x31, 0xce, 0x76, 0x3b, 0x5d, 0x7f, 0x93, 0xcf, 0x00, 0x58, 0xb5,
	0xf2, 0xab, 0x22, 0xdc, 0xb7, 0xfd, 0x00, 0x3a, 0x39, 0x97, 0xa2, 0x10, 0x9c, 0xf0, 0x37, 0x80,
	0x68, 0xc3, 0x7c, 0x06, 0x73, 0x6e, 0x68, 0x37, 0x03, 0x4b, 0xeb, 0x19, 0x6c, 0x54, 0xf5, 0x0c,
	0x2a, 0x91, 0xcd, 0x20, 0x99, 0xc0, 0xe9, 0x96, 0xad, 0x3d, 0xcb, 0x0d, 0xa0, 0xb3, 0x52, 0x84,
	0xd9, 0x4c, 0x73, 0x49, 0x28, 0xaa, 0x0d, 0x23, 0xcb, 0x9e, 0x3c, 0xba, 0xfd, 0x0a, 0x20, 0x4c,
	0xf9, 0x94, 0xd8, 0x33, 0x74, 0xb7, 0x52, 0x63, 0xbd, 0xa1, 0x3b, 0x89, 0xa6, 0x13, 0xe8, 0xf5,
	0x1b, 0x6b, 0xde, 0x49, 0xf2, 0x8f, 0x3d, 0x42, 0x54, 0xb3, 0xc8, 0x62, 0xaf, 0xde, 0x0d, 0xb3,
	0x77, 0xd1, 0x50, 0x59, 0x77, 0x79, 0x3b, 0x74, 0xd7, 0x78, 0xf7, 0x13, 0x00, 0x00, 0xff, 0xff,
	0xb7, 0xa2, 0x37, 0xf6, 0x9b, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RaftClient is the client API for Raft service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RaftClient interface {
	AppendEntries(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesResponse, error)
	RequestVote(ctx context.Context, in *RequestVoteRequest, opts ...grpc.CallOption) (*RequestVoteResponse, error)
}

type raftClient struct {
	cc *grpc.ClientConn
}

func NewRaftClient(cc *grpc.ClientConn) RaftClient {
	return &raftClient{cc}
}

func (c *raftClient) AppendEntries(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesResponse, error) {
	out := new(AppendEntriesResponse)
	err := c.cc.Invoke(ctx, "/raft.Raft/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftClient) RequestVote(ctx context.Context, in *RequestVoteRequest, opts ...grpc.CallOption) (*RequestVoteResponse, error) {
	out := new(RequestVoteResponse)
	err := c.cc.Invoke(ctx, "/raft.Raft/RequestVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RaftServer is the raft API for Raft service.
type RaftServer interface {
	AppendEntries(context.Context, *AppendEntriesRequest) (*AppendEntriesResponse, error)
	RequestVote(context.Context, *RequestVoteRequest) (*RequestVoteResponse, error)
}

// UnimplementedRaftServer can be embedded to have forward compatible implementations.
type UnimplementedRaftServer struct {
}

func (*UnimplementedRaftServer) AppendEntries(ctx context.Context, req *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (*UnimplementedRaftServer) RequestVote(ctx context.Context, req *RequestVoteRequest) (*RequestVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}

func RegisterRaftServer(s *grpc.Server, srv RaftServer) {
	s.RegisterService(&_Raft_serviceDesc, srv)
}

func _Raft_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft.Raft/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServer).AppendEntries(ctx, req.(*AppendEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Raft_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft.Raft/RequestVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServer).RequestVote(ctx, req.(*RequestVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Raft_serviceDesc = grpc.ServiceDesc{
	ServiceName: "raft.Raft",
	HandlerType: (*RaftServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AppendEntries",
			Handler:    _Raft_AppendEntries_Handler,
		},
		{
			MethodName: "RequestVote",
			Handler:    _Raft_RequestVote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raft.proto",
}
