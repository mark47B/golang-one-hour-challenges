package main

type Stats struct {
    Requests uint64
    Errors   uint64
}

type StatsStore struct {
    stats atomic.Pointer[Stats]
}

func NewStatsStore() *StatsStore {
	stats := &Stats{
		Requests: 0,
		Errors: 0,
	}
	statsStore := &StatsStore{}
	statsStore.stats.Store(stats)
	return statsStore
}



func (s *StatsStore) IncRequests() {
	for {
		stats := s.stats.Load()
		if stats == nil {
			return 
		}

		statsCpy := *stats
		statsCpy.Requests++

		if s.stats.CompareAndSwap(stats, &statsCpy) {
			return
		}
	}
}

func (s *StatsStore) IncErrors() {
	for {
		stats := s.stats.Load()
		if stats == nil {
			return 
		}

		statsCpy := *stats
		statsCpy.Errors++

		if s.stats.CompareAndSwap(stats, &statsCpy) {
			return
		}
	}
}

func (s *StatsStore) Snapshot() Stats {
	stats := s.stats.Load()
	return *stats
}

func main(){
	cfgStore := NewStatsStore()
	_ = cfgStore
}
