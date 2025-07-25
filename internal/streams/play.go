package streams

import (
	"errors"
	"time"

	"github.com/hamza-farouk/go2rtc/pkg/core"
)

func (s *Stream) Play(urlOrProd any) error {
	s.mu.Lock()
	for _, producer := range s.producers {
		if producer.state == stateInternal && producer.conn != nil {
			_ = producer.conn.Stop()
		}
	}
	s.mu.Unlock()

	var source string
	var src core.Producer

	switch urlOrProd.(type) {
	case string:
		if source = urlOrProd.(string); source == "" {
			return nil
		}
	case core.Producer:
		src = urlOrProd.(core.Producer)
	}

	for _, producer := range s.producers {
		if producer.conn == nil {
			continue
		}

		cons, ok := producer.conn.(core.Consumer)
		if !ok {
			continue
		}

		if src == nil {
			var err error
			if src, err = GetProducer(source); err != nil {
				return err
			}
		}

		if !matchMedia(src, cons) {
			continue
		}

		s.AddInternalProducer(src)

		go func() {
			_ = src.Start()
			s.RemoveProducer(src)
		}()

		return nil
	}

	for _, producer := range s.producers {
		// start new client
		dst, err := GetProducer(producer.url)
		if err != nil {
			continue
		}

		// check if client support consumer interface
		cons, ok := dst.(core.Consumer)
		if !ok {
			_ = dst.Stop()
			continue
		}

		// start new producer
		if src == nil {
			if src, err = GetProducer(source); err != nil {
				return err
			}
		}

		if !matchMedia(src, cons) {
			_ = dst.Stop()
			continue
		}

		s.AddInternalProducer(src)
		s.AddInternalConsumer(cons)

		go func() {
			_ = dst.Start()
			_ = src.Stop()
			s.RemoveInternalConsumer(cons)
		}()

		go func() {
			_ = src.Start()
			// little timeout before stop dst, so the buffer can be transferred
			time.Sleep(time.Second)
			_ = dst.Stop()
			s.RemoveProducer(src)
		}()

		return nil
	}

	return errors.New("can't find consumer")
}

func (s *Stream) AddInternalProducer(conn core.Producer) {
	producer := &Producer{conn: conn, state: stateInternal, url: "internal"}
	s.mu.Lock()
	s.producers = append(s.producers, producer)
	s.mu.Unlock()
}

func (s *Stream) AddInternalConsumer(conn core.Consumer) {
	s.mu.Lock()
	s.consumers = append(s.consumers, conn)
	s.mu.Unlock()
}

func (s *Stream) RemoveInternalConsumer(conn core.Consumer) {
	s.mu.Lock()
	for i, consumer := range s.consumers {
		if consumer == conn {
			s.consumers = append(s.consumers[:i], s.consumers[i+1:]...)
			break
		}
	}
	s.mu.Unlock()
}

func matchMedia(prod core.Producer, cons core.Consumer) bool {
	for _, consMedia := range cons.GetMedias() {
		for _, prodMedia := range prod.GetMedias() {
			if prodMedia.Direction != core.DirectionRecvonly {
				continue
			}

			prodCodec, consCodec := prodMedia.MatchMedia(consMedia)
			if prodCodec == nil {
				continue
			}

			track, err := prod.GetTrack(prodMedia, prodCodec)
			if err != nil {
				log.Warn().Err(err).Msg("[streams] can't get track")
				continue
			}

			if err = cons.AddTrack(consMedia, consCodec, track); err != nil {
				log.Warn().Err(err).Msg("[streams] can't add track")
				continue
			}

			return true
		}
	}

	return false
}
