package mikrotik_util

import (
	"fmt"
	"time"
)

// Single point in time sample of the interface state
type InterfaceBandwidthSample struct {
	Interface MtInterface
	IsRx      bool
	Date      time.Time
	ByteCount int64
}

// Calculates the delta between two samples
// The older sample should be the receiver and
// the newer sample should be the parameter
func (first InterfaceBandwidthSample) Delta(second InterfaceBandwidthSample) (sample InterfaceBandwidthSampleDelta) {
	if first.Interface.Name != second.Interface.Name || first.IsRx != second.IsRx {
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
	Interface MtInterface
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
		panic("Non RX interface in RX Slot")
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

// Samples bandwidth for the given duration and returns the delta.
func (iface MtInterface) SampleBandwidth(isRx bool, sampleSize time.Duration) (InterfaceBandwidthSampleDelta, error) {
	var err error

	// Slice to store samples in
	samples := make([]InterfaceBandwidthSample, 2)
	for x := 0; x < 2; x++ {
		samples[x], err = iface.SampleBytes(isRx)
		if err != nil {
			return InterfaceBandwidthSampleDelta{}, err
		}

		// Don't sleep on the last iteration
		if x < 1 {
			time.Sleep(sampleSize)
		}
	}

	delta := samples[0].Delta(samples[1])
	return delta, nil
}

// Monitors both tx and rx bandwidth at the given period, post results to the returned channel
func (iface MtInterface) MonitorBandwidth(sampleSize time.Duration) (resultsChan chan InterfaceRxTxSample) {
	rxSample := make(chan InterfaceBandwidthSampleDelta)
	txSample := make(chan InterfaceBandwidthSampleDelta)

	clock := time.NewTicker(sampleSize)
	rxClock := make(chan time.Time) // We use Time so one could use a Ticker if need be
	txClock := make(chan time.Time) // We use Time so one could use a Ticker if need be
	resultsChan = make(chan InterfaceRxTxSample)

	// Start the monitor functions
	go iface.MonitorBandwidthRxTx(true, rxClock, rxSample)
	go iface.MonitorBandwidthRxTx(false, txClock, txSample)

	// Clock Loop
	go func() {
		// We do it this way (as opposed to using the clock.C directly) so that they
		// do not block on each other
		for tickTime := range clock.C {
			// Post to the rx and tx clocks as fast as possible
			go func() { rxClock <- tickTime }()
			go func() { txClock <- tickTime }()
		}
	}()

	// Loop forevar!
	go func() {
		for {
			var samples InterfaceRxTxSample
			samples[0] = <-rxSample
			samples[1] = <-txSample
			resultsChan <- samples
		}
	}()

	// Return channel to caller so they can do something useful with it.
	return resultsChan
}

// Monitors either the RX or the TX bandwidth, posting deltas to the given chanel
func (iface MtInterface) MonitorBandwidthRxTx(isRx bool, freqClock chan time.Time, results chan InterfaceBandwidthSampleDelta) {
	samples := make([]InterfaceBandwidthSample, 2)
	// next_idx := 0
	sampleCount := 0
	// for _ = range freqClock {
	for {
		sample, err := iface.SampleBytes(isRx)
		if err != nil {
			results <- InterfaceBandwidthSampleDelta{Delta: -1}
		}

		// Find our storage index
		this_idx := sampleCount % 2
		last_idx := 1 - this_idx // Toogle 0/1
		samples[this_idx] = sample
		sampleCount += 1

		// If we have enough samples then calculate the delta
		// and post it to the channel
		if sampleCount >= 2 {
			prevSample := samples[last_idx]
			results <- prevSample.Delta(sample)
		}

		// Block/Check that the channel is still open, if not then exit
		if _, stillRunning := <-freqClock; !stillRunning {
			break
		}
	}
}

// Samples Interface for bytes RX or bytes TX depending on flag
func (iff MtInterface) SampleBytes(isRx bool) (sample InterfaceBandwidthSample, err error) {
	sample.Interface = iff
	sample.IsRx = isRx
	var bytes int64
	if isRx {
		bytes, err = iff.GetBytesRX()
	} else {
		bytes, err = iff.GetBytesTX()
	}
	// Capture the date as close as possible to return
	sample.Date = time.Now().Round(time.Second)

	// Now do err checking
	if err != nil {
		return InterfaceBandwidthSample{}, err
	}
	sample.ByteCount = bytes

	return sample, nil
}
