package psrpc

import (
	"math/rand"
	"sync"
	"time"
	"unicode"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type lockedRandSource struct {
	mu  sync.Mutex
	src rand.Source
}

func (s *lockedRandSource) Int63() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.src.Int63()
}

func (s *lockedRandSource) Seed(seed int64) {
	s.mu.Lock()
	s.src.Seed(seed)
	s.mu.Unlock()
}

var idRNG = rand.New(&lockedRandSource{src: rand.NewSource(time.Now().UnixNano())})

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func readIDChars(b []byte) {
	var n int
	for {
		r := idRNG.Int63()
		for i := 0; i < 10; i++ {
			if int(r&0x3f) < len(alphabet) {
				b[n] = alphabet[r&0x3f]
				n++
				if n == len(b) {
					return
				}
			}
			r >>= 6
		}
	}
}

func formatID(prefix string) string {
	b := make([]byte, 16)
	copy(b, prefix)
	readIDChars(b[4:])
	return string(b)
}

func newRequestID() string {
	return formatID("REQ_")
}

func newStreamID() string {
	return formatID("STR_")
}

const lowerHex = "0123456789abcdef"

var channelChar = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0030, 0x0039, 1}, // 0-9
		{0x0041, 0x005a, 1}, // A-Z
		{0x005f, 0x005f, 1}, // _
		{0x0061, 0x007a, 1}, // a-z
	},
	LatinOffset: 4,
}

func appendSanitizedChannelPart(buf []byte, s string) []byte {
	for _, r := range s {
		if unicode.Is(channelChar, r) {
			buf = append(buf, byte(r))
		} else if r < 0x10000 {
			buf = append(buf, `u+`...)
			for s := 12; s >= 0; s -= 4 {
				buf = append(buf, lowerHex[r>>uint(s)&0xF])
			}
		} else {
			buf = append(buf, `U+`...)
			for s := 28; s >= 0; s -= 4 {
				buf = append(buf, lowerHex[r>>uint(s)&0xF])
			}
		}
	}
	return buf
}

func appendChannelParts[T any](buf []byte, parts ...T) []byte {
	var prefix bool
	for _, t := range parts {
		if prefix {
			buf = append(buf, '|')
		}
		l := len(buf)
		switch v := any(t).(type) {
		case string:
			buf = appendSanitizedChannelPart(buf, v)
		case []string:
			buf = appendChannelParts(buf, v...)
		}
		prefix = len(buf) > l
	}
	return buf
}

func channelPartsLen[T any](parts ...T) int {
	var n int
	for _, t := range parts {
		switch v := any(t).(type) {
		case string:
			n += len(v) + 1
		case []string:
			n += channelPartsLen(v...)
		}
	}
	return n
}

func formatChannel(parts ...any) string {
	buf := make([]byte, 0, 4*channelPartsLen(parts...)/3)
	return string(appendChannelParts(buf, parts...))
}

func getRPCChannel(serviceName, rpc string, topic []string) string {
	return formatChannel(serviceName, rpc, topic, "REQ")
}

func getHandlerKey(rpc string, topic []string) string {
	return formatChannel(rpc, topic)
}

func getResponseChannel(serviceName, clientID string) string {
	return formatChannel(serviceName, clientID, "RES")
}

func getClaimRequestChannel(serviceName, clientID string) string {
	return formatChannel(serviceName, clientID, "CLAIM")
}

func getClaimResponseChannel(serviceName, rpc string, topic []string) string {
	return formatChannel(serviceName, rpc, topic, "RCLAIM")
}

func getStreamChannel(serviceName, nodeID string) string {
	return formatChannel(serviceName, nodeID, "STR")
}

func getStreamServerChannel(serviceName, rpc string, topic []string) string {
	return formatChannel(serviceName, rpc, topic, "STR")
}

func serialize(msg proto.Message) ([]byte, error) {
	a, err := anypb.New(msg)
	if err != nil {
		return nil, err
	}

	b, err := proto.Marshal(a)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func deserialize(b []byte) (proto.Message, error) {
	a := &anypb.Any{}
	err := proto.Unmarshal(b, a)
	if err != nil {
		return nil, err
	}

	return a.UnmarshalNew()
}

func serializePayload(m proto.Message) ([]byte, *anypb.Any, error) {
	a, err := anypb.New(m)
	return a.GetValue(), a, err
}

func deserializePayload[T proto.Message](buf []byte, any *anypb.Any) (T, error) {
	if any != nil {
		buf = any.Value
	}
	var p T
	v := p.ProtoReflect().New().Interface().(T)
	return v, proto.Unmarshal(buf, v)
}
