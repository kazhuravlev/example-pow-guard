package randstring

func (s *Store) Get() string {
	idx := s.rng.Intn(len(s.strings))
	return s.strings[idx]
}
