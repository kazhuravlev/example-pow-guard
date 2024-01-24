package facade

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"
)

func (s *Service) handle(ctx context.Context, conn net.Conn) error {
	s.log.Info("handle conn", slog.String("client", conn.RemoteAddr().String()))

	if _, err := conn.Write([]byte(s.api.GetWelcomeMessage())); err != nil {
		return fmt.Errorf("write welcome message: %w", err)
	}

	// Close connection at the end of fn.
	defer func() {
		if err := conn.Close(); err != nil {
			s.log.Error("close connection", slog.String("err", err.Error()))
		}
	}()

	var isVerified bool

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context was done: %w", ctx.Err())
		default:
			// NOTE(zhuravlev): here is slowread/slowwrite attack is possible.
			// NOTE(zhuravlev): copy deadlines from context to conn.
			netData, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				return fmt.Errorf("read request: %w", err)
			}

			var response string
			message := strings.TrimSpace(netData)
			switch message {
			case "quit", "exit", "stop":
				return nil
			case "":
				if isVerified {
					response = s.api.GetQuote() + "\n"
				} else {
					response = s.api.GetWelcomeMessage()
				}
			default:
				if isVerified {
					response = "You already have a bless from Chuck! Don't be impudent. Chuck doesn't like this.\n"
				} else {
					isVerified = s.api.VerifyChallenge(message)
					if isVerified {
						response = "Chuck bless you!\n"
					} else {
						response = "Chuck memorized your movements...\n"
					}
				}
			}

			if _, err := conn.Write([]byte(response)); err != nil {
				return fmt.Errorf("write response: %w", err)
			}
		}
	}
}
