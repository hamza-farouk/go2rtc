package expr

import (
	"errors"

	"github.com/hamza-farouk/go2rtc/internal/app"
	"github.com/hamza-farouk/go2rtc/internal/streams"
	"github.com/hamza-farouk/go2rtc/pkg/expr"
)

func Init() {
	log := app.GetLogger("expr")

	streams.RedirectFunc("expr", func(url string) (string, error) {
		v, err := expr.Eval(url[5:], nil)
		if err != nil {
			return "", err
		}

		log.Debug().Msgf("[expr] url=%s", url)

		if url = v.(string); url == "" {
			return "", errors.New("expr: result is empty")
		}

		return url, nil
	})
}
