package api

import (
	"github.com/jellydator/ttlcache/v3"
	"github.com/kazhuravlev/example-pow-guard/pkg/randstring"
	"log/slog"
	"time"
)

type Service struct {
	log    *slog.Logger
	quotes *randstring.Store
	// TTL for challenge. Old challenges is not allowed.
	challengeTTL time.Duration
	challenges   *ttlcache.Cache[string, struct{}]
}

func New(logger *slog.Logger, quotes *randstring.Store) (*Service, error) {
	const challengeTTL = 30 * time.Second

	// NOTE: this lib is looks a bit awkward. So it used only for mvp.
	cache := ttlcache.New[string, struct{}](
		ttlcache.WithTTL[string, struct{}](challengeTTL),
	)

	go cache.Start()

	return &Service{
		log:          logger,
		quotes:       quotes,
		challengeTTL: challengeTTL,
		challenges:   cache,
	}, nil
}
