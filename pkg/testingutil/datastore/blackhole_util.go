package datastore

func (s *Blackhole) Count() int {
	s.l.Lock()
	defer s.l.Unlock()

	return s.count
}

func (s *Blackhole) SetWriteError(err error) {
	s.l.Lock()
	defer s.l.Unlock()

	s.writeError = err
}
