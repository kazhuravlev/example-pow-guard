package facade

import (
	"context"
	"fmt"
	"github.com/kazhuravlev/example-pow-guard/internal/api"
	"log/slog"
	"net"
	"strconv"
	"sync"
)

type Service struct {
	log  *slog.Logger
	api  *api.Service
	port int
	wg   *sync.WaitGroup
}

func New(log *slog.Logger, apiInst *api.Service, port int) (*Service, error) {
	return &Service{
		log:  log,
		api:  apiInst,
		port: port,
		wg:   new(sync.WaitGroup),
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	address := ":" + strconv.Itoa(s.port)
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		return fmt.Errorf("listen tcp addr: %w", err)
	}

	// Stop listener on ctx.Done
	go func() {
		<-ctx.Done()

		if err := listener.Close(); err != nil {
			s.log.Error("close tcp listener", slog.String("err", err.Error()))
		}
	}()

	// Handle new connections.
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				s.log.Error("accept tcp conn", slog.String("err", err.Error()))
				return
			}

			s.wg.Add(1)
			// NOTE(zhuravlev): add a pool that limits the number of connections.
			go func() {
				defer s.wg.Done()

				if err := s.handle(ctx, conn); err != nil {
					s.log.Error("handle connection", slog.String("err", err.Error()))
				}
			}()
		}
	}()

	return nil
}

// Wait will wait that all connections was closed.
func (s *Service) Wait() {
	s.wg.Wait()
}
