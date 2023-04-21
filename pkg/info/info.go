package info

import (
	"sync"

	"github.com/livekit/psrpc"
)

type ServiceDefinition struct {
	Name    string
	ID      string
	Methods sync.Map
}

type MethodInfo struct {
	AffinityEnabled bool
	Multi           bool
	RequireClaim    bool
}

type RequestInfo struct {
	psrpc.RPCInfo
	AffinityEnabled bool
	RequireClaim    bool
}

func (s *ServiceDefinition) RegisterMethod(name string, affinityEnabled, multi, requireClaim bool) {
	s.Methods.Store(name, &MethodInfo{
		AffinityEnabled: affinityEnabled,
		Multi:           multi,
		RequireClaim:    requireClaim,
	})
}

func (s *ServiceDefinition) GetInfo(rpc string, topic []string) *RequestInfo {
	v, _ := s.Methods.Load(rpc)
	m := v.(*MethodInfo)

	return &RequestInfo{
		RPCInfo: psrpc.RPCInfo{
			Service: s.Name,
			Method:  rpc,
			Topic:   topic,
			Multi:   m.Multi,
		},
		AffinityEnabled: m.AffinityEnabled,
		RequireClaim:    m.RequireClaim,
	}
}
