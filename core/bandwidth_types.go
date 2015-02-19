package core

import (
	"fmt"
	"log"
	"time"

	// . "github.com/ErebusBat/mikrotik_util/core"
)

// Single point in time sample of the interface state
type InterfaceBandwidthSample struct {
	Interface RbInterface
	IsRx      bool
	Date      time.Time
	ByteCount int64
}

// Calculates the delta between two samples
// The older sample should be the receiver and
// the newer sample should be the parameter
func (first InterfaceBandwidthSample) Delta(second InterfaceBandwidthSample) (sample InterfaceBandwidthSampleDelta) {
	if !first.Interface.Equals(second.Interface) || first.IsRx != second.IsRx {
		return sample
	}
	sample.Interface = first.Interface
	sample.IsRx = first.IsRx
	sample.Duration = second.Date.Sub(first.Date)
	sample.Delta = second.ByteCount - first.ByteCount
	sample.Bits = sample.Delta * 8

	return sample
}

// InterfaceBandwidthSampleDelta is a comparison of a rx OR tx sample over time
type InterfaceBandwidthSampleDelta struct {
	Interface RbInterface
	IsRx      bool
	Duration  time.Duration
	Delta     int64
	Bits      int64
}

// Returns a humanized string representing the bits, i.e. 123.4 Mbps
func (d InterfaceBandwidthSampleDelta) BitsString() string {
	var bits float64 = float64(d.Bits) / d.Duration.Seconds()
	template := "%.1f"

	if bits >= 1000000000 {
		bits = bits / float64(1000000000)
		template += " Gbps"
	} else if bits >= 1000000 {
		bits = bits / float64(1000000)
		template += " Mbps"
	} else if bits >= 1000 {
		bits = bits / float64(1000)
		template += " Kbps"
	}
	return fmt.Sprintf(template, bits)
}

// Returns a string representing the sample, i.e. RB450G-ether1 tx/rx 541.8 Kbps/595.3 Kbps
func (d InterfaceBandwidthSampleDelta) String() string {
	var rxtx string = "tx"
	if d.IsRx {
		rxtx = "rx"
	}
	log.Printf("%#v\n", d)
	return fmt.Sprintf("%s-%s %10s",
		d.Interface.FullName(),
		rxtx,
		d.BitsString(),
	)
}

// InterfaceRxTxSample is a 'complete' (both rx and tx) picture of interface bandwidth
type InterfaceRxTxSample [2]InterfaceBandwidthSampleDelta

// Returns the RX Sample
func (s InterfaceRxTxSample) RX() InterfaceBandwidthSampleDelta {
	if !s[0].IsRx {
		panic(fmt.Sprintf("Non RX interface in RX Slot\n:%#v", s))
	}
	return s[0]
}

// Returns the TX Sample
func (s InterfaceRxTxSample) TX() InterfaceBandwidthSampleDelta {
	if s[1].IsRx {
		panic("Non TX interface in TX Slot")
	}
	return s[1]
}
