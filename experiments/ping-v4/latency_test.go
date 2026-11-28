package latency

import (
	"context"
	"net/netip"
	"testing"
)

func TestInvalidConfiguration(t *testing.T) {
	elapsed, err := MeasureOffline(
		context.Background(),
		nil,
		nil,
		0,
		"",
	)

	if err == nil || elapsed != 0 {
		t.Fatal(
			"invalid configuration accepted",
		)
	}
}

func TestCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(
		context.Background(),
	)
	cancel()

	addresses := []netip.Addr{
		netip.MustParseAddr("10.0.0.2"),
	}

	dns := []netip.Addr{
		netip.MustParseAddr("1.1.1.1"),
	}

	elapsed, err := MeasureOffline(
		ctx,
		addresses,
		dns,
		1420,
		"test",
	)

	if err == nil || elapsed != 0 {
		t.Fatal(
			"cancelled context accepted",
		)
	}
}
