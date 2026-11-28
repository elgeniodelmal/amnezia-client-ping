package latency

import (
	"context"
	"errors"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"github.com/amnezia-vpn/amneziawg-go/v3/conn"
	"github.com/amnezia-vpn/amneziawg-go/v3/device"
	"github.com/amnezia-vpn/amneziawg-go/v3/tun/netstack"
)

// MeasureOffline requires an inactive peer identity
// and correctly routed outer UDP sockets.
// Never log private keys or raw UAPI configuration.
func MeasureOffline(
	ctx context.Context,
	addresses []netip.Addr,
	dns []netip.Addr,
	mtu int,
	uapi string,
) (time.Duration, error) {

	if len(addresses) == 0 ||
		len(dns) == 0 ||
		mtu < 576 ||
		mtu > 9000 ||
		strings.TrimSpace(uapi) == "" {
		return 0, errors.New(
			"invalid probe configuration",
		)
	}

	if err := ctx.Err(); err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(
		ctx, 12*time.Second,
	)
	defer cancel()

	tun, network, err := netstack.CreateNetTUN(
		addresses, dns, mtu,
	)
	if err != nil {
		return 0, errors.New(
			"netstack initialization failed",
		)
	}

	dev := device.NewDevice(
		tun,
		conn.NewDefaultBind(),
		device.NewLogger(
			device.LogLevelSilent, "",
		),
	)
	defer dev.Close()

	if err = dev.IpcSet(uapi); err != nil {
		return 0, errors.New(
			"peer configuration failed",
		)
	}

	if err = dev.Up(); err != nil {
		return 0, errors.New(
			"tunnel activation failed",
		)
	}

	transport := &http.Transport{
		DialContext:           network.DialContext,
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   8 * time.Second,
		ResponseHeaderTimeout: 8 * time.Second,
	}
	defer transport.CloseIdleConnections()

	client := &http.Client{
		Transport: transport,
		Timeout:   12 * time.Second,
		CheckRedirect: func(
			req *http.Request,
			via []*http.Request,
		) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://www.gstatic.com/generate_204",
		nil,
	)
	if err != nil {
		return 0, errors.New(
			"request initialization failed",
		)
	}

	req.Header.Set("Cache-Control", "no-cache")

	start := time.Now()
	response, err := client.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		return 0, errors.New(
			"HTTPS probe failed",
		)
	}

	defer response.Body.Close()

	if response.StatusCode != 204 {
		return 0, errors.New(
			"unexpected HTTP status",
		)
	}

	return elapsed, nil
}
