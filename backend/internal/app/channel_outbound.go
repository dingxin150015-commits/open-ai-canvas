package app

import (
	"net/http"
	"net/url"
	"time"
)

// The optional flags are only supplied by an already authenticated channel
// path. Ordinary calls remain on the public outbound policy.
func (s *Service) ValidateChannelOutboundURL(rawURL string, flags ...bool) (*url.URL, error) {
	requestedLocal, customRelay := false, false
	if len(flags) > 0 {
		requestedLocal = flags[0]
	}
	if len(flags) > 1 {
		customRelay = flags[1]
	}
	return s.validateChannelOutboundURL(rawURL, requestedLocal, customRelay)
}

// OutboundHTTPClientForChannel keeps the explicit desktop loopback policy
// separate from administrator-configured channel proxies.
func (s *Service) OutboundHTTPClientForChannel(timeout time.Duration, target *url.URL, allowLocal bool, proxyURL string) *http.Client {
	if s.effectiveAllowLocalChannel(allowLocal) && target != nil && isExactDesktopLoopbackHost(target.Hostname()) {
		return outboundHTTPClientWithPolicy(timeout, desktopLoopbackOutboundPolicy(nil))
	}
	return OutboundHTTPClientWithProxy(timeout, proxyURL)
}
