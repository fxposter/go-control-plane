package types

import (
	"time"

	"github.com/golang/protobuf/proto"
)

// Resource is the base interface for the xDS payload.
type Resource interface {
	proto.Message
}

// ResourceWithTtl is a Resource with an optional TTL.
type ResourceWithTtl struct {
	Resource Resource

	Ttl *time.Duration
}

// MarshaledResource is an alias for the serialized binary array.
type MarshaledResource = []byte

// SkipFetchError is the error returned when the cache fetch is short
// circuited due to the client's version already being up-to-date.
type SkipFetchError struct{}

// Error satisfies the error interface
func (e SkipFetchError) Error() string {
	return "skip fetch: version up to date"
}

// ResponseType enumeration of supported response types
type ResponseType int

const (
	Cluster ResponseType = iota
	Endpoint
	Listener
	Route
	Secret
	Runtime
	ExtensionConfig
	UnknownType // token to count the total number of supported types
)
