package main

import (
	"github.com/hamza-farouk/go2rtc/internal/alsa"
	"github.com/hamza-farouk/go2rtc/internal/api"
	"github.com/hamza-farouk/go2rtc/internal/api/ws"
	"github.com/hamza-farouk/go2rtc/internal/app"
	"github.com/hamza-farouk/go2rtc/internal/bubble"
	"github.com/hamza-farouk/go2rtc/internal/debug"
	"github.com/hamza-farouk/go2rtc/internal/doorbird"
	"github.com/hamza-farouk/go2rtc/internal/dvrip"
	"github.com/hamza-farouk/go2rtc/internal/echo"
	"github.com/hamza-farouk/go2rtc/internal/eseecloud"
	"github.com/hamza-farouk/go2rtc/internal/exec"
	"github.com/hamza-farouk/go2rtc/internal/expr"
	"github.com/hamza-farouk/go2rtc/internal/ffmpeg"
	"github.com/hamza-farouk/go2rtc/internal/flussonic"
	"github.com/hamza-farouk/go2rtc/internal/gopro"
	"github.com/hamza-farouk/go2rtc/internal/hass"
	"github.com/hamza-farouk/go2rtc/internal/hls"
	"github.com/hamza-farouk/go2rtc/internal/homekit"
	"github.com/hamza-farouk/go2rtc/internal/http"
	"github.com/hamza-farouk/go2rtc/internal/isapi"
	"github.com/hamza-farouk/go2rtc/internal/ivideon"
	"github.com/hamza-farouk/go2rtc/internal/mjpeg"
	"github.com/hamza-farouk/go2rtc/internal/mp4"
	"github.com/hamza-farouk/go2rtc/internal/mpegts"
	"github.com/hamza-farouk/go2rtc/internal/nest"
	"github.com/hamza-farouk/go2rtc/internal/ngrok"
	"github.com/hamza-farouk/go2rtc/internal/onvif"
	"github.com/hamza-farouk/go2rtc/internal/ring"
	"github.com/hamza-farouk/go2rtc/internal/roborock"
	"github.com/hamza-farouk/go2rtc/internal/rtmp"
	"github.com/hamza-farouk/go2rtc/internal/rtsp"
	"github.com/hamza-farouk/go2rtc/internal/srtp"
	"github.com/hamza-farouk/go2rtc/internal/streams"
	"github.com/hamza-farouk/go2rtc/internal/tapo"
	"github.com/hamza-farouk/go2rtc/internal/v4l2"
	"github.com/hamza-farouk/go2rtc/internal/webrtc"
	"github.com/hamza-farouk/go2rtc/internal/webtorrent"
	"github.com/hamza-farouk/go2rtc/internal/wyoming"
	"github.com/hamza-farouk/go2rtc/internal/yandex"
	"github.com/hamza-farouk/go2rtc/pkg/shell"
)

func main() {
	app.Version = "1.9.9"

	// 1. Core modules: app, api/ws, streams

	app.Init() // init config and logs

	api.Init() // init API before all others
	ws.Init()  // init WS API endpoint

	streams.Init() // streams module

	// 2. Main sources and servers

	rtsp.Init()   // rtsp source, RTSP server
	webrtc.Init() // webrtc source, WebRTC server

	// 3. Main API

	mp4.Init()   // MP4 API
	hls.Init()   // HLS API
	mjpeg.Init() // MJPEG API

	// 4. Other sources and servers

	hass.Init()       // hass source, Hass API server
	onvif.Init()      // onvif source, ONVIF API server
	webtorrent.Init() // webtorrent source, WebTorrent module
	wyoming.Init()

	// 5. Other sources

	rtmp.Init()     // rtmp source
	exec.Init()     // exec source
	ffmpeg.Init()   // ffmpeg source
	echo.Init()     // echo source
	ivideon.Init()  // ivideon source
	http.Init()     // http/tcp source
	dvrip.Init()    // dvrip source
	tapo.Init()     // tapo source
	isapi.Init()    // isapi source
	mpegts.Init()   // mpegts passive source
	roborock.Init() // roborock source
	homekit.Init()  // homekit source
	ring.Init()     // ring source
	nest.Init()     // nest source
	bubble.Init()   // bubble source
	expr.Init()     // expr source
	gopro.Init()    // gopro source
	doorbird.Init() // doorbird source
	v4l2.Init()     // v4l2 source
	alsa.Init()     // alsa source
	flussonic.Init()
	eseecloud.Init()
	yandex.Init()

	// 6. Helper modules

	ngrok.Init() // ngrok module
	srtp.Init()  // SRTP server
	debug.Init() // debug API

	// 7. Go

	shell.RunUntilSignal()
}
