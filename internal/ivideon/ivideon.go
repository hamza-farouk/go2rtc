package ivideon

import (
	"github.com/hamza-farouk/go2rtc/internal/streams"
	"github.com/hamza-farouk/go2rtc/pkg/ivideon"
)

func Init() {
	streams.HandleFunc("ivideon", ivideon.Dial)
}
