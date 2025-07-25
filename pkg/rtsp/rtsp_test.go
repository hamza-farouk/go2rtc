package rtsp

import (
	"testing"

	"github.com/hamza-farouk/go2rtc/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestURLParse(t *testing.T) {
	// https://github.com/hamza-farouk/WebRTC/issues/395
	base := "rtsp://::ffff:192.168.1.123/onvif/profile.1/"
	u, err := urlParse(base)
	assert.Empty(t, err)
	assert.Equal(t, "::ffff:192.168.1.123:", u.Host)

	// https://github.com/hamza-farouk/go2rtc/issues/208
	base = "rtsp://rtsp://turret2-cam.lan:554/stream1/"
	u, err = urlParse(base)
	assert.Empty(t, err)
	assert.Equal(t, "turret2-cam.lan:554", u.Host)
}

func TestBugSDP1(t *testing.T) {
	// https://github.com/hamza-farouk/WebRTC/issues/417
	s := `v=0
o=- 91674849066 1 IN IP4 192.168.1.123
s=RtspServer
i=live
t=0 0
a=control:*
a=range:npt=0-
m=video 0 RTP/AVP 96
c=IN IP4 0.0.0.0
s=RtspServer
i=live
a=control:track0
a=range:npt=0-
a=rtpmap:96 H264/90000
a=fmtp:96 packetization-mode=1;profile-level-id=42001E;sprop-parameter-sets=Z0IAHvQCgC3I,aM48gA==
a=control:track0
m=audio 0 RTP/AVP 97
c=IN IP4 0.0.0.0
s=RtspServer
i=live
a=control:track1
a=range:npt=0-
a=rtpmap:97 MPEG4-GENERIC/8000/1
a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3;config=1588
a=control:track1
`
	medias, err := UnmarshalSDP([]byte(s))
	assert.Nil(t, err)
	assert.NotNil(t, medias)
}

func TestBugSDP2(t *testing.T) {
	// https://github.com/hamza-farouk/WebRTC/issues/419
	s := `v=0
o=- 1675628282 1675628283 IN IP4 192.168.1.123
s=streamed by the RTSP server
t=0 0
m=video 0 RTP/AVP 96
a=rtpmap:96 H264/90000
a=control:track0
m=audio 0 RTP/AVP 8
a=rtpmap:0 pcma/8000/1
a=control:track1
a=framerate:25
a=range:npt=now-
a=fmtp:96 packetization-mode=1;profile-level-id=64001F;sprop-parameter-sets=Z0IAH5WoFAFuQA==,aM48gA==
`
	medias, err := UnmarshalSDP([]byte(s))
	assert.Nil(t, err)
	assert.NotNil(t, medias)
	assert.NotEqual(t, "", medias[0].Codecs[0].FmtpLine)
}

func TestBugSDP3(t *testing.T) {
	s := `v=0
o=- 1680614126554766 1 IN IP4 192.168.0.3
s=Session streamed by "preview"
t=0 0
a=tool:BC Streaming Media v202210012022.10.01
a=type:broadcast
a=control:*
a=range:npt=now-
a=x-qt-text-nam:Session streamed by "preview"
m=video 0 RTP/AVP 96
c=IN IP4 0.0.0.0
b=AS:8192
a=rtpmap:96 H264/90000
a=range:npt=now-
a=fmtp:96 packetization-mode=1;profile-level-id=640033;sprop-parameter-sets=Z2QAM6wVFKAoAPGQ,aO48sA==
a=recvonly
a=control:track1
m=audio 0 RTP/AVP 97
c=IN IP4 0.0.0.0
b=AS:8192
a=rtpmap:97 MPEG4-GENERIC/16000
a=fmtp:97 streamtype=5;profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3;config=1408;
a=recvonly
a=control:track2
m=audio 0 RTP/AVP 8
a=control:track3
a=rtpmap:8 PCMA/8000
a=sendonly`
	medias, err := UnmarshalSDP([]byte(s))
	assert.Nil(t, err)
	assert.Len(t, medias, 3)
}

func TestBugSDP4(t *testing.T) {
	s := `v=0
o=- 14665860 31787219 1 IN IP4 10.0.0.94
s=Session streamed by "MERCURY RTSP Server"
t=0 0
m=video 0 RTP/AVP 96
c=IN IP4 0.0.0.0
b=AS:4096
a=range:npt=0-
a=control:track1
a=rtpmap:96 H264/90000
a=fmtp:96 packetization-mode=1; profile-level-id=640016; sprop-parameter-sets=Z2QAFqzGoCgPaEAAAAMAQAAAB6E=,aOqPLA==
m=audio 0 RTP/AVP 8
a=rtpmap:8 PCMA/8000
a=control:track2
m=application/MERCURY 0 RTP/AVP smart/1/90000
a=rtpmap:95 MERCURY/90000
a=control:track3
`
	medias, err := UnmarshalSDP([]byte(s))
	assert.Nil(t, err)
	assert.Len(t, medias, 3)
}

func TestBugSDP5(t *testing.T) {
	s := `v=0
o=CV-RTSPHandler 1123412 0 IN IP4 192.168.1.22
s=Camera
c=IN IP4 192.168.1.22
t=0 0
a=charset:Shift_JIS
a=range:npt=0-
a=control:*
a=etag:1234567890
m=video 0 RTP/AVP 99
a=rtpmap:99 H264/90000
a=fmtp:99 profile-level-id=42A01E;packetization-mode=1;sprop-parameter-sets=Z0KgKedAPAET8uAIEAABd2AAK/IGAAADAC+vCAAAHc1lP//jAAADABfXhAAADuayn//wIA==,aN48gA==
a=control:trackID=1
a=sendonly
m=audio 0 RTP/AVP 127
a=rtpmap:127 mpeg4-generic/8000/1
a=fmtp:127 streamtype=5; profile-level-id=15; mode=AAC-hbr; sizeLength=13; indexLength=3; indexDeltalength=3; config=1588; CTSDeltaLength=0; DTSDeltaLength=0;
a=control:trackID=2
`
	medias, err := UnmarshalSDP([]byte(s))
	assert.Nil(t, err)
	assert.Len(t, medias, 2)
	assert.Equal(t, "recvonly", medias[0].Direction)
	assert.Equal(t, "recvonly", medias[1].Direction)
}

func TestBugSDP6(t *testing.T) {
	// https://github.com/hamza-farouk/go2rtc/issues/1278
	s := `v=0
o=- 3730506281693 1 IN IP4 172.20.0.215
s=IP camera Live streaming
i=stream1
t=0 0
a=tool:LIVE555 Streaming Media v2014.02.04
a=type:broadcast
a=control:*
a=range:npt=0-
a=x-qt-text-nam:IP camera Live streaming
a=x-qt-text-inf:stream1
m=video 0 RTP/AVP 26
c=IN IP4 172.20.0.215
b=AS:1500
a=x-bufferdelay:0.55000
a=x-dimensions:1280,960
a=control:track1
m=audio 0 RTP/AVP 0
c=IN IP4 172.20.0.215
b=AS:64
a=x-bufferdelay:0.55000
a=control:track2
m=application 0 RTP/AVP 107
c=IN IP4 172.20.0.215
b=AS:1
a=x-bufferdelay:0.55000
a=rtpmap:107 vnd.onvif.metadata/90000/500
a=control:track4
m=vana 0 RTP/AVP 108
c=IN IP4 172.20.0.215
b=AS:1
a=x-bufferdelay:0.55000
a=rtpmap:108 video.analysis/90000/500
a=control:track5
`
	medias, err := UnmarshalSDP([]byte(s))
	assert.Nil(t, err)
	assert.Len(t, medias, 4)
}

func TestBugSDP7(t *testing.T) {
	// https://github.com/hamza-farouk/go2rtc/issues/1426
	s := `v=0
o=- 1001 1 IN
s=VCP IPC Realtime stream
m=video 0 RTP/AVP 105
c=IN
a=control:rtsp://1.0.1.113/media/video2/video
a=rtpmap:105 H264/90000
a=fmtp:105 profile-level-id=640016; packetization-mode=1; sprop-parameter-sets=Z2QAFqw7UFAX/LCAAAH0AABOIEI=,aOqPLA==
a=recvonly
m=audio 0 RTP/AVP 0
c=IN
a=fmtp:0 RTCP=0
a=control:rtsp://1.0.1.113/media/video2/audio1
a=recvonly
m=audio 0 RTP/AVP 0
c=IN
a=control:rtsp://1.0.1.113/media/video2/backchannel
a=rtpmap:0 PCMA/8000
a=rtpmap:0 PCMU/8000
a=sendonly
m=application 0 RTP/AVP 107
c=IN
a=control:rtsp://1.0.1.113/media/video2/metadata
a=rtpmap:107 vnd.onvif.metadata/90000
a=fmtp:107 DecoderTag=h3c-v3 RTCP=0
a=recvonly
`
	medias, err := UnmarshalSDP([]byte(s))
	assert.Nil(t, err)
	assert.Len(t, medias, 4)
}

func TestHikvisionPCM(t *testing.T) {
	s := `v=0
o=- 1721969533379665 1721969533379665 IN IP4 192.168.1.12
s=Media Presentation
e=NONE
b=AS:5100
t=0 0
a=control:rtsp://192.168.1.12:554/Streaming/channels/101/
m=video 0 RTP/AVP 96
c=IN IP4 0.0.0.0
b=AS:5000
a=recvonly
a=x-dimensions:3200,1800
a=control:rtsp://192.168.1.12:554/Streaming/channels/101/trackID=1
a=rtpmap:96 H264/90000
a=fmtp:96 profile-level-id=420029; packetization-mode=1; sprop-parameter-sets=Z2QAM6wVFKAyAOP5f/AAEAAWyAAAH0AAB1MAIA==,aO48sA==
m=audio 0 RTP/AVP 11
c=IN IP4 0.0.0.0
b=AS:50
a=recvonly
a=control:rtsp://192.168.1.12:554/Streaming/channels/101/trackID=2
a=rtpmap:11 PCM/48000
a=Media_header:MEDIAINFO=494D4B4801030000040000010170011080BB0000007D000000000000000000000000000000000000;
a=appversion:1.0
`
	medias, err := UnmarshalSDP([]byte(s))
	assert.Nil(t, err)
	assert.Len(t, medias, 2)
	assert.Equal(t, core.CodecPCML, medias[1].Codecs[0].Name)
}
