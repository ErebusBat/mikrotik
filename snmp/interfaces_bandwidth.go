package snmp

import (
	"time"

	. "github.com/ErebusBat/mikrotik/core"
)

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
			delta := prevSample.Delta(sample)
			results <- delta
		}

		// Block/Check that the channel is still open, if not then exit
		if _, stillRunning := <-freqClock; !stillRunning {
			break
		}
	}
}

// Samples Interface for bytes RX or bytes TX depending on flag
func (iff MtInterface) SampleBytes(isRx bool) (sample InterfaceBandwidthSample, err error) {
	// log.Printf("SampleBytes %t %#v\n", isRx, iff.FullName())
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

	// log.Printf("  Got: %#v\n", sample)
	return sample, nil
}
