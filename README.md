# mikrotik
Golang Bandwidth monitoring utility for Mikrotik RouterOS Devices

This library reads [RouterOS](http://routerboard.com/) devices (primarily Mikrotik routerboads).
Currently only SNMP and the functions to support bandwidth monitoring are supported.

This library is built using [gosnmp](https://github.com/alouca/gosnmp).  However it should be abstracted enough that if 
one were so inclined you could implement a [raw API](http://wiki.mikrotik.com/wiki/Manual:API) client.

Currently there only exists the `bandwidth` tool (See [bandwidth.go](https://github.com/ErebusBat/mikrotik/blob/master/cmd/bandwidth.go)) which consumes the library and reports interface statistics at a given interval.

## Installation ##

```sh
go get -u -v github.com/ErebusBat/mikrotik/

# Optional... compile the include tool / sample
go build cmd/*.go
```

## Usage ##

The tool tries to have intelligent defaults, so if your SNMP community is `public` and the interface name you want to monitor is the standard `ether1` then you just need to specify a host:

```
$ ./bandwidth -h 192.168.0.1 

# Help
$ ./bandwidth --help
Usage of ./bandwidth:
  -c="public": SNMP Community Name
  -h="127.0.0.7": Mikrotik IP
  -i="ether1": Mikrotik Interface Name
  -list=false: Lists all known interfaces and exits
  -s=1s: Sample Interval
```

## Example ##

If you want to use the library in your own tools then you can do something like

```go
package main

import (
	"log"
	"time"

	"github.com/ErebusBat/mikrotik/snmp"
)

func main() {
	rb := snmp.Connect("10.0.1.250", "public")
	iface, err := rb.FindInterfaceByName("UPSTREAM")
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	sampleChannel := iface.MonitorBandwidth(time.Second * 5)
	for sample := range sampleChannel {
		log.Printf("%s tx/rx %s/%s",
			iface.FullName(),
			sample.TX().BitsString(),
			sample.RX().BitsString(),
		)
	}
}

```

## Contributing ##

Contributions welcome! Please fork the repository and open a pull request
with your changes.

## License ##

This is free software, licensed under the Apache License, Version 2.0.
