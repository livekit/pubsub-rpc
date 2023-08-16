// Copyright 2023 LiveKit, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"github.com/livekit/psrpc"
	"github.com/livekit/psrpc/internal/bus"
	"github.com/livekit/psrpc/pkg/info"
)

func withStreams() psrpc.ClientOption {
	return func(o *psrpc.ClientOpts) {
		o.EnableStreams = true
	}
}

func getClientOpts(opts ...psrpc.ClientOption) psrpc.ClientOpts {
	o := &psrpc.ClientOpts{
		Timeout:     psrpc.DefaultClientTimeout,
		ChannelSize: bus.DefaultChannelSize,
	}
	for _, opt := range opts {
		opt(o)
	}
	return *o
}

func getRequestOpts(i *info.RequestInfo, options psrpc.ClientOpts, opts ...psrpc.RequestOption) psrpc.RequestOpts {
	o := &psrpc.RequestOpts{
		Timeout: options.Timeout,
	}
	if i.AffinityEnabled {
		o.SelectionOpts = psrpc.SelectionOpts{
			AffinityTimeout:     psrpc.DefaultAffinityTimeout,
			ShortCircuitTimeout: psrpc.DefaultAffinityShortCircuit,
		}
	} else {
		o.SelectionOpts = psrpc.SelectionOpts{
			AcceptFirstAvailable: true,
		}
	}

	for _, opt := range opts {
		opt(o)
	}

	return *o
}
