package randstring

import (
	"errors"
	"fmt"
	"math/rand"
)

var ErrBadConfiguration = errors.New("bad configuration")

type Store struct {
	strings []string
	rng     *rand.Rand
}

func New(rng *rand.Rand, in []string) (*Store, error) {
	if len(in) == 0 {
		return nil, fmt.Errorf("at least one string is required: %w", ErrBadConfiguration)
	}

	return &Store{
		strings: in,
		rng:     rng,
	}, nil
}
