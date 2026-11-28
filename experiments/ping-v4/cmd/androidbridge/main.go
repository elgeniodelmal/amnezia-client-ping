package main

/*
#include <stdint.h>
*/
import "C"

import (
	"context"
	"net/netip"
	"strings"
	"sync"

	latency "amnezia-ping-v4"
)

var probeLock sync.Mutex

func parseAddresses(
	input string,
) ([]netip.Addr, error) {

	var result []netip.Addr

	for _, value := range strings.Split(
		input, ",",
	) {
		value = strings.TrimSpace(value)

		if value == "" {
			continue
		}

		if strings.Contains(value, "/") {
			prefix, err := netip.ParsePrefix(value)

			if err != nil {
				return nil, err
			}

			result = append(
				result, prefix.Addr(),
			)
		} else {
			address, err := netip.ParseAddr(value)

			if err != nil {
				return nil, err
			}

			result = append(result, address)
		}
	}

	return result, nil
}

//export AmneziaPingVersion
func AmneziaPingVersion() C.int {
	return 401
}

//export AmneziaPingMeasureOffline
func AmneziaPingMeasureOffline(
	addresses *C.char,
	dns *C.char,
	mtu C.int,
	uapi *C.char,
	vpnActive C.int,
) (result C.int64_t) {

	result = -5

	defer func() {
		if recover() != nil {
			result = -5
		}
	}()

	if vpnActive != 0 {
		return -4
	}

	if addresses == nil ||
		dns == nil ||
		uapi == nil {
		return -2
	}

	if !probeLock.TryLock() {
		return -3
	}

	defer probeLock.Unlock()

	localAddresses, err := parseAddresses(
		C.GoString(addresses),
	)
	if err != nil {
		return -2
	}

	dnsAddresses, err := parseAddresses(
		C.GoString(dns),
	)
	if err != nil {
		return -2
	}

	config := C.GoString(uapi)

	if len(localAddresses) == 0 ||
		len(dnsAddresses) == 0 ||
		int(mtu) < 576 ||
		int(mtu) > 9000 ||
		strings.TrimSpace(config) == "" {
		return -2
	}

	elapsed, err := latency.MeasureOffline(
		context.Background(),
		localAddresses,
		dnsAddresses,
		int(mtu),
		config,
	)

	if err != nil {
		return -1
	}

	ms := elapsed.Milliseconds()

	if ms < 1 {
		ms = 1
	}

	return C.int64_t(ms)
}

func main() {}
