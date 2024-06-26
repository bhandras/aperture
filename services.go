package aperture

import (
	"context"
	"time"

	"github.com/lightninglabs/aperture/l402"
	"github.com/lightninglabs/aperture/mint"
	"github.com/lightninglabs/aperture/proxy"
)

// staticServiceLimiter provides static restrictions for services.
//
// TODO(wilmer): use etcd instead.
type staticServiceLimiter struct {
	capabilities map[l402.Service]l402.Caveat
	constraints  map[l402.Service][]l402.Caveat
	timeouts     map[l402.Service]l402.Caveat
}

// A compile-time constraint to ensure staticServiceLimiter implements
// mint.ServiceLimiter.
var _ mint.ServiceLimiter = (*staticServiceLimiter)(nil)

// newStaticServiceLimiter instantiates a new static service limiter backed by
// the given restrictions.
func newStaticServiceLimiter(
	proxyServices []*proxy.Service) *staticServiceLimiter {

	capabilities := make(map[l402.Service]l402.Caveat)
	constraints := make(map[l402.Service][]l402.Caveat)
	timeouts := make(map[l402.Service]l402.Caveat)

	for _, proxyService := range proxyServices {
		s := l402.Service{
			Name:  proxyService.Name,
			Tier:  l402.BaseTier,
			Price: proxyService.Price,
		}

		if proxyService.Timeout > 0 {
			timeouts[s] = l402.NewTimeoutCaveat(
				proxyService.Name,
				proxyService.Timeout,
				time.Now,
			)
		}

		capabilities[s] = l402.NewCapabilitiesCaveat(
			proxyService.Name, proxyService.Capabilities,
		)
		for cond, value := range proxyService.Constraints {
			caveat := l402.Caveat{Condition: cond, Value: value}
			constraints[s] = append(constraints[s], caveat)
		}
	}

	return &staticServiceLimiter{
		capabilities: capabilities,
		constraints:  constraints,
		timeouts:     timeouts,
	}
}

// ServiceCapabilities returns the capabilities caveats for each service. This
// determines which capabilities of each service can be accessed.
func (l *staticServiceLimiter) ServiceCapabilities(ctx context.Context,
	services ...l402.Service) ([]l402.Caveat, error) {

	res := make([]l402.Caveat, 0, len(services))
	for _, service := range services {
		capabilities, ok := l.capabilities[service]
		if !ok {
			continue
		}
		res = append(res, capabilities)
	}

	return res, nil
}

// ServiceConstraints returns the constraints for each service. This enforces
// additional constraints on a particular service/service capability.
func (l *staticServiceLimiter) ServiceConstraints(ctx context.Context,
	services ...l402.Service) ([]l402.Caveat, error) {

	res := make([]l402.Caveat, 0, len(services))
	for _, service := range services {
		constraints, ok := l.constraints[service]
		if !ok {
			continue
		}
		res = append(res, constraints...)
	}

	return res, nil
}

// ServiceTimeouts returns the timeout caveat for each service. This enforces
// an expiration time for service access if enabled.
func (l *staticServiceLimiter) ServiceTimeouts(ctx context.Context,
	services ...l402.Service) ([]l402.Caveat, error) {

	res := make([]l402.Caveat, 0, len(services))
	for _, service := range services {
		timeout, ok := l.timeouts[service]
		if !ok {
			continue
		}
		res = append(res, timeout)
	}

	return res, nil
}
