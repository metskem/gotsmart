GotSmart
========

GotSmart collects information from the Dutch Slimme Meter (translated as Smart Meter) and exports them as Prometheus metrics.
For extensive details on the metrics see the [docs from Netbeheer Nederland](https://www.netbeheernederland.nl/_upload/Files/Slimme_meter_15_a727fce1f1.pdf)

Setup
-----

### Build

```sh
cd gotsmart
go get ./...
make 
```

### Run

Specify the serial device that is connected with the Smart Meter, and a non-default listen address:port.

```sh
gotsmart -device /dev/ttyS0 -listen-address :8082
```

If you want to run it as a daemon, use the file ``/etc/systemd/system/gotsmart.service`` and fill it with this (example, fixed port and device):
```sh
[Unit]
Description=GotSmart, reading your smart meter and providing a prometheus exporter endpoint
Documentation=https://github.com/metskem/gotsmart
After=network-online.target

[Service]
User=pi
Restart=on-failure

ExecStart=/usr/local/bin/gotsmart -listen-address :8082 -device /dev/ttyUSB0

[Install]
WantedBy=multi-user.target
```

Usage
-----

By default gotsmart reads device /dev/ttyAMA0, listens on address:port 0.0.0.0:8080 and exposes the metrics under `/metrics`.
