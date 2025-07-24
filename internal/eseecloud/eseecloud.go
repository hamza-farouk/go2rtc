package eseecloud

import (
	"github.com/hamza-farouk/go2rtc/internal/streams"
	"github.com/hamza-farouk/go2rtc/pkg/eseecloud"
)

func Init() {
	streams.HandleFunc("eseecloud", eseecloud.Dial)
}
