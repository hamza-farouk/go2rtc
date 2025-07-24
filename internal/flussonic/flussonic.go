package flussonic

import (
	"github.com/hamza-farouk/go2rtc/internal/streams"
	"github.com/hamza-farouk/go2rtc/pkg/flussonic"
)

func Init() {
	streams.HandleFunc("flussonic", flussonic.Dial)
}
