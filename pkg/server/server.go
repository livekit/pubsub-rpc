package server

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/frostbyte73/core"
	"go.uber.org/atomic"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/proto"

	"github.com/livekit/psrpc"
	"github.com/livekit/psrpc/internal/bus"
	"github.com/livekit/psrpc/pkg/info"
)

type rpcHandler interface {
	close(force bool)
}

type RPCServer struct {
	*info.ServiceDefinition
	psrpc.ServerOpts

	bus bus.MessageBus

	mu       sync.RWMutex
	handlers map[string]rpcHandler
	active   atomic.Int32
	shutdown core.Fuse
}

func NewRPCServer(sd *info.ServiceDefinition, b bus.MessageBus, opts ...psrpc.ServerOption) *RPCServer {
	s := &RPCServer{
		ServiceDefinition: sd,
		ServerOpts:        getServerOpts(opts...),
		bus:               b,
		handlers:          make(map[string]rpcHandler),
		shutdown:          core.NewFuse(),
	}

	return s
}

func RegisterHandler[RequestType proto.Message, ResponseType proto.Message](
	s *RPCServer,
	rpc string,
	topic []string,
	svcImpl func(context.Context, RequestType) (ResponseType, error),
	affinityFunc AffinityFunc[RequestType],
) error {
	if s.shutdown.IsBroken() {
		return psrpc.ErrServerClosed
	}

	i := s.GetInfo(rpc, topic)

	key := i.GetHandlerKey()
	s.mu.RLock()
	_, ok := s.handlers[key]
	s.mu.RUnlock()
	if ok {
		return errors.New("handler already exists")
	}

	// create handler
	h, err := newRPCHandler(s, i, svcImpl, s.ChainedInterceptor, affinityFunc)
	if err != nil {
		return err
	}

	s.active.Inc()
	h.onCompleted = func() {
		s.active.Dec()
		s.mu.Lock()
		delete(s.handlers, key)
		s.mu.Unlock()
	}

	s.mu.Lock()
	s.handlers[key] = h
	s.mu.Unlock()

	h.run(s)
	return nil
}

func RegisterStreamHandler[RequestType proto.Message, ResponseType proto.Message](
	s *RPCServer,
	rpc string,
	topic []string,
	svcImpl func(psrpc.ServerStream[ResponseType, RequestType]) error,
	affinityFunc StreamAffinityFunc,
) error {
	if s.shutdown.IsBroken() {
		return psrpc.ErrServerClosed
	}

	i := s.GetInfo(rpc, topic)

	key := i.GetHandlerKey()
	s.mu.RLock()
	_, ok := s.handlers[key]
	s.mu.RUnlock()
	if ok {
		return errors.New("handler already exists")
	}

	// create handler
	h, err := newStreamRPCHandler(s, i, svcImpl, affinityFunc)
	if err != nil {
		return err
	}

	s.active.Inc()
	h.onCompleted = func() {
		s.active.Dec()
		s.mu.Lock()
		delete(s.handlers, key)
		s.mu.Unlock()
	}

	s.mu.Lock()
	s.handlers[key] = h
	s.mu.Unlock()

	h.run(s)
	return nil
}

func (s *RPCServer) DeregisterHandler(rpc string, topic []string) {
	i := s.GetInfo(rpc, topic)
	key := i.GetHandlerKey()
	s.mu.RLock()
	h, ok := s.handlers[key]
	s.mu.RUnlock()
	if ok {
		h.close(true)
	}
}

func (s *RPCServer) Publish(ctx context.Context, rpc string, topic []string, msg proto.Message) error {
	i := s.GetInfo(rpc, topic)
	return s.bus.Publish(ctx, i.GetRPCChannel(), msg)
}

func (s *RPCServer) Close(force bool) {
	s.shutdown.Once(func() {
		s.mu.RLock()
		handlers := maps.Values(s.handlers)
		s.mu.RUnlock()

		var wg sync.WaitGroup
		for _, h := range handlers {
			wg.Add(1)
			h := h
			go func() {
				h.close(force)
				wg.Done()
			}()
		}
		wg.Wait()
	})

	if !force {
		for s.active.Load() > 0 {
			time.Sleep(time.Millisecond * 100)
		}
	}
}
