package rest

import (
	"context"
	"errors"
	discovery "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v2"
)

type Server interface {
	Fetch(context.Context, *discovery.DiscoveryRequest) (*discovery.DiscoveryResponse, error)
}

type Callbacks interface {
	// OnFetchRequest is called for each Fetch request. Returning an error will end processing of the
	// request and respond with an error.
	OnFetchRequest(context.Context, *discovery.DiscoveryRequest) error
	// OnFetchResponse is called immediately prior to sending a response.
	OnFetchResponse(*discovery.DiscoveryRequest, *discovery.DiscoveryResponse)
}

func NewServer(config cache.ConfigFetcher, callbacks Callbacks) Server {
	return &server{cache: config, callbacks: callbacks}
}

type server struct {
	cache     cache.ConfigFetcher
	callbacks Callbacks
}

func (s *server) Fetch(ctx context.Context, req *discovery.DiscoveryRequest) (*discovery.DiscoveryResponse, error) {
	if s.callbacks != nil {
		if err := s.callbacks.OnFetchRequest(ctx, req); err != nil {
			return nil, err
		}
	}
	resp, err := s.cache.Fetch(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("missing response")
	}
	out, err := resp.GetDiscoveryResponse()
	if s.callbacks != nil {
		s.callbacks.OnFetchResponse(req, out)
	}
	return out, err
}
