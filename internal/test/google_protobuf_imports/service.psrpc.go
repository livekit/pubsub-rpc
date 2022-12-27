// Code generated by protoc-gen-psrpc v0.2.1, DO NOT EDIT.
// source: service.proto

package google_protobuf_imports

import context "context"

import psrpc "github.com/livekit/psrpc"
import google_protobuf "google.golang.org/protobuf/types/known/emptypb"
import google_protobuf1 "google.golang.org/protobuf/types/known/wrapperspb"

// ====================
// Svc Client Interface
// ====================

type SvcClient interface {
	Send(context.Context, *google_protobuf1.StringValue, ...psrpc.RequestOption) (*google_protobuf.Empty, error)
}

// ========================
// Svc ServerImpl Interface
// ========================

type SvcServerImpl interface {
	Send(context.Context, *google_protobuf1.StringValue) (*google_protobuf.Empty, error)
}

// ====================
// Svc Server Interface
// ====================

type SvcServer interface {
	// Close and wait for pending RPCs to complete
	Shutdown()

	// Close immediately, without waiting for pending RPCs
	Kill()
}

// ==========
// Svc Client
// ==========

type svcClient struct {
	client *psrpc.RPCClient
}

// NewSvcClient creates a psrpc client that implements the SvcClient interface.
func NewSvcClient(clientID string, bus psrpc.MessageBus, opts ...psrpc.ClientOption) (SvcClient, error) {
	rpcClient, err := psrpc.NewRPCClient("Svc", clientID, bus, opts...)
	if err != nil {
		return nil, err
	}

	return &svcClient{
		client: rpcClient,
	}, nil
}

func (c *svcClient) Send(ctx context.Context, req *google_protobuf1.StringValue, opts ...psrpc.RequestOption) (*google_protobuf.Empty, error) {
	return psrpc.RequestSingle[*google_protobuf.Empty](ctx, c.client, "Send", "", req, opts...)
}

// ==========
// Svc Server
// ==========

type svcServer struct {
	svc SvcServerImpl
	rpc *psrpc.RPCServer
}

// NewSvcServer builds a RPCServer that will route requests
// to the corresponding method in the provided svc implementation.
func NewSvcServer(serverID string, svc SvcServerImpl, bus psrpc.MessageBus, opts ...psrpc.ServerOption) (SvcServer, error) {
	s := psrpc.NewRPCServer("Svc", serverID, bus, opts...)

	var err error
	err = psrpc.RegisterHandler(s, "Send", "", svc.Send, nil)
	if err != nil {
		s.Close(false)
		return nil, err
	}

	return &svcServer{
		svc: svc,
		rpc: s,
	}, nil
}

func (s *svcServer) Shutdown() {
	s.rpc.Close(false)
}

func (s *svcServer) Kill() {
	s.rpc.Close(true)
}

var psrpcFileDescriptor0 = []byte{
	// 166 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0xcd, 0xbf, 0x0a, 0xc2, 0x40,
	0x0c, 0xc7, 0xf1, 0x41, 0x71, 0x28, 0xb8, 0x74, 0x10, 0x39, 0xff, 0x3c, 0x42, 0x0a, 0xba, 0x3a,
	0x29, 0x3e, 0x41, 0xc1, 0xc1, 0xa5, 0xb4, 0x35, 0x96, 0x83, 0xf6, 0x2e, 0x24, 0xb9, 0x8a, 0x6f,
	0x2f, 0xde, 0xd1, 0x45, 0xd7, 0xdf, 0xe7, 0x1b, 0x92, 0x2d, 0x05, 0x79, 0xb4, 0x2d, 0x02, 0xb1,
	0x57, 0x9f, 0xef, 0x48, 0x98, 0x5a, 0xb0, 0x4e, 0x91, 0x5d, 0xdd, 0x83, 0xa2, 0x28, 0x04, 0xc1,
	0x0a, 0x07, 0xd2, 0xb7, 0xd9, 0x74, 0xde, 0x77, 0x3d, 0x16, 0x31, 0x6e, 0xc2, 0xb3, 0x88, 0x73,
	0xba, 0x35, 0xfb, 0x5f, 0x7c, 0x71, 0x4d, 0x84, 0x2c, 0xc9, 0x0f, 0x97, 0x6c, 0x56, 0x8e, 0x6d,
	0x7e, 0xca, 0xe6, 0x25, 0xba, 0x47, 0xbe, 0x85, 0xd4, 0xc3, 0xd4, 0x43, 0xa9, 0x6c, 0x5d, 0x77,
	0xab, 0xfb, 0x80, 0x66, 0xf5, 0xa7, 0xd7, 0xef, 0xab, 0xb3, 0xb9, 0xaf, 0x8b, 0x24, 0xd5, 0x24,
	0x95, 0x1d, 0xc8, 0xb3, 0x4a, 0xb3, 0x88, 0xcb, 0xf1, 0x13, 0x00, 0x00, 0xff, 0xff, 0x0a, 0xfa,
	0x30, 0x79, 0xd4, 0x00, 0x00, 0x00,
}
