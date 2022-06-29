package provider

import (
	"sync"

	"golang.org/x/exp/maps"
)

type fetchStats struct {
	l       sync.Mutex
	fetches map[string]int
}

func (s *fetchStats) track(name string) {
	s.l.Lock()
	defer s.l.Unlock()

	if s.fetches == nil {
		s.fetches = make(map[string]int)
	}

	count := s.fetches[name]
	count++
	s.fetches[name] = count
}

func (s *fetchStats) get() map[string]int {
	s.l.Lock()
	defer s.l.Unlock()

	return maps.Clone(s.fetches)
}

var stats fetchStats

func FetchStats() map[string]int {
	return stats.get()
}
