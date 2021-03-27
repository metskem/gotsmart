package conf

import (
	"flag"
	"fmt"
	"github.com/tarm/serial"
	"log"
)

var (
	VersionTag string
	BuildTime  string

	Parity serial.Parity

	AddrFlag   = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	DeviceFlag = flag.String("device", "/dev/ttyAMA0", "Serial device to read P1 data from.")
	BaudFlag   = flag.Int("baud", 115200, "Baud rate (speed) to use.")
	BitsFlag   = flag.Int("bits", 8, "Number of databits.")
	ParityFlag = flag.String("parity", "none", "Parity the use (none/odd/even/mark/space).")
)

func Init() {
	flag.Parse()

	fmt.Printf("GotSmart (version %s, build time %s)\n", VersionTag, BuildTime)

	switch *ParityFlag {
	case "none":
		Parity = serial.ParityNone
	case "odd":
		Parity = serial.ParityOdd
	case "even":
		Parity = serial.ParityEven
	case "mark":
		Parity = serial.ParityMark
	case "space":
		Parity = serial.ParitySpace
	default:
		log.Fatal("Invalid Parity setting")
	}
}
