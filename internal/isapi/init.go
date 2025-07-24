package isapi

import (
	"github.com/hamza-farouk/go2rtc/internal/streams"
	"github.com/hamza-farouk/go2rtc/pkg/core"
	"github.com/hamza-farouk/go2rtc/pkg/isapi"
)

func Init() {
	streams.HandleFunc("isapi", func(source string) (core.Producer, error) {
		return isapi.Dial(source)
	})
}
