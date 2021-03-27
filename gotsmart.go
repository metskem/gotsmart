package main

import (
	"bufio"
	"fmt"
	"github.com/metskem/gotsmart/conf"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/metskem/gotsmart/crc16"
	"github.com/metskem/gotsmart/dsmr"
	dsmrprometheus "github.com/metskem/gotsmart/dsmr/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tarm/serial"
)

type frameupdate struct {
	Frame string
	Time  time.Time
}

func (f *frameupdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Last-Modified", f.Time.Format(http.TimeFormat))
	w.Write([]byte(f.Frame))
}

func (f *frameupdate) Update(frame string) {
	f.Frame = strings.Replace(frame, "\r", "", -1)
	f.Time = time.Now()
}

func (f *frameupdate) Process(br *bufio.Reader, collector *dsmrprometheus.DSMRCollector) {
	for {
		if b, err := br.Peek(1); err == nil {
			if string(b) != "/" {
				fmt.Printf("Ignoring garbage character: %c\n", b)
				br.ReadByte()
				continue
			}
		} else {
			continue
		}
		frame, err := br.ReadBytes('!')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		bcrc, err := br.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		// Check CRC
		mcrc := strings.ToUpper(strings.TrimSpace(string(bcrc)))
		crc := fmt.Sprintf("%04X", crc16.Checksum(frame))
		if mcrc != crc {
			fmt.Printf("CRC mismatch: %q != %q\n", mcrc, crc)
			continue
		}
		f.Update(string(frame))
		dsmrFrame, err := dsmr.ParseFrame(string(frame))
		if err != nil {
			log.Printf("could not parse frame: %v\n", err)
			continue
		}
		collector.Update(dsmrFrame)
	}
}

func main() {
	conf.Init()
	serialConfig := &serial.Config{Name: *conf.DeviceFlag, Baud: *conf.BaudFlag, Size: byte(*conf.BitsFlag), Parity: conf.Parity}
	p, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatal(err)
	}

	br := bufio.NewReader(p)
	collector := &dsmrprometheus.DSMRCollector{}
	prometheus.MustRegister(collector)
	f := &frameupdate{}
	go f.Process(br, collector)

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", f)
	err = http.ListenAndServe(*conf.AddrFlag, nil)
	if err != nil {
		log.Fatal(err)
	}
}
