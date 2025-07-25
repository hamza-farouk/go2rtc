//go:build linux

package v4l2

import (
	"errors"
	"net/url"
	"strings"

	"github.com/hamza-farouk/go2rtc/pkg/core"
	"github.com/hamza-farouk/go2rtc/pkg/h264/annexb"
	"github.com/hamza-farouk/go2rtc/pkg/v4l2/device"
	"github.com/pion/rtp"
)

type Producer struct {
	core.Connection
	dev *device.Device
}

func Open(rawURL string) (*Producer, error) {
	// Example (ffmpeg source compatible):
	// v4l2:device?video=/dev/video0&input_format=mjpeg&video_size=1280x720
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := u.Query()

	dev, err := device.Open(query.Get("video"))
	if err != nil {
		return nil, err
	}

	codec := &core.Codec{
		ClockRate:   90000,
		PayloadType: core.PayloadTypeRAW,
	}

	var width, height, pixFmt uint32

	if wh := strings.Split(query.Get("video_size"), "x"); len(wh) == 2 {
		codec.FmtpLine = "width=" + wh[0] + ";height=" + wh[1]
		width = uint32(core.Atoi(wh[0]))
		height = uint32(core.Atoi(wh[1]))
	}

	switch query.Get("input_format") {
	case "yuyv422":
		if codec.FmtpLine == "" {
			return nil, errors.New("v4l2: invalid video_size")
		}
		codec.Name = core.CodecRAW
		codec.FmtpLine += ";colorspace=422"
		pixFmt = device.V4L2_PIX_FMT_YUYV
	case "nv12":
		if codec.FmtpLine == "" {
			return nil, errors.New("v4l2: invalid video_size")
		}
		codec.Name = core.CodecRAW
		codec.FmtpLine += ";colorspace=420mpeg2" // maybe 420jpeg
		pixFmt = device.V4L2_PIX_FMT_NV12
	case "mjpeg":
		codec.Name = core.CodecJPEG
		pixFmt = device.V4L2_PIX_FMT_MJPEG
	case "h264":
		codec.Name = core.CodecH264
		pixFmt = device.V4L2_PIX_FMT_H264
	case "hevc":
		codec.Name = core.CodecH265
		pixFmt = device.V4L2_PIX_FMT_HEVC
	default:
		return nil, errors.New("v4l2: invalid input_format")
	}

	if err = dev.SetFormat(width, height, pixFmt); err != nil {
		return nil, err
	}

	if fps := core.Atoi(query.Get("framerate")); fps > 0 {
		if err = dev.SetParam(uint32(fps)); err != nil {
			return nil, err
		}
	}

	medias := []*core.Media{
		{
			Kind:      core.KindVideo,
			Direction: core.DirectionRecvonly,
			Codecs:    []*core.Codec{codec},
		},
	}
	return &Producer{
		Connection: core.Connection{
			ID:         core.NewID(),
			FormatName: "v4l2",
			Medias:     medias,
		},
		dev: dev,
	}, nil
}

func (c *Producer) Start() error {
	if err := c.dev.StreamOn(); err != nil {
		return err
	}

	var bitstream bool
	switch c.Medias[0].Codecs[0].Name {
	case core.CodecH264, core.CodecH265:
		bitstream = true
	}

	for {
		buf, err := c.dev.Capture()
		if err != nil {
			return err
		}

		c.Recv += len(buf)

		if len(c.Receivers) == 0 {
			continue
		}

		if bitstream {
			buf = annexb.EncodeToAVCC(buf)
		}

		pkt := &rtp.Packet{
			Header:  rtp.Header{Timestamp: core.Now90000()},
			Payload: buf,
		}
		c.Receivers[0].WriteRTP(pkt)
	}
}

func (c *Producer) Stop() error {
	_ = c.Connection.Stop()
	return errors.Join(c.dev.StreamOff(), c.dev.Close())
}
