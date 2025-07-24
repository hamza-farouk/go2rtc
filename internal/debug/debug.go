package debug

import (
	"github.com/hamza-farouk/go2rtc/internal/api"
)

func Init() {
	api.HandleFunc("api/stack", stackHandler)
}
