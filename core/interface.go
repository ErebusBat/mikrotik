package core

import (
	"fmt"
	"time"
)

type RbInterface interface {
	fmt.Stringer
	Bandwidther
	Routerboarder

	Index() int
	Name() string
	FullName() string
	Equals(rhs RbInterface) bool
}

type Bandwidther interface {
	GetBytesRX() (ifBytes int64, err error)
	GetBytesTX() (ifBytes int64, err error)
	SampleBytes(isRx bool) (sample InterfaceBandwidthSample, err error)
	MonitorBandwidth(sampleSize time.Duration) (resultsChan chan InterfaceRxTxSample)
}
