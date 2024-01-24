package api

import (
	"github.com/jellydator/ttlcache/v3"
	hc "github.com/kazhuravlev/example-pow-guard/pkg/hashcash"
)

const msgWelcome = `================================================================================
Workflow:
- Send the solved and valid (not older that 3 seconds) solution of hashcash 
  challenge.
- Hit the enter as many times as you want to get a quote.
- Type 'quit' to close connection.
================================================================================
`

func (s *Service) GetWelcomeMessage() string {
	return msgWelcome
}

func (s *Service) GetQuote() string {
	quote := s.quotes.Get()

	return ">>> " + quote
}

func (s *Service) VerifyChallenge(result string) bool {
	if s.challenges.Has(result) {
		return false
	}

	isValid := hc.New(20, 8, "", s.challengeTTL).Check(result)
	s.challenges.Set(result, struct{}{}, ttlcache.DefaultTTL)

	return isValid
}
