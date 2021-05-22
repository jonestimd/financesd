package domain

type idSet struct {
	ids map[int64]struct{}
}

func newIDSet() *idSet {
	return &idSet{ids: make(map[int64]struct{})}
}

func (s *idSet) Add(id int64) {
	if _, ok := s.ids[id]; !ok {
		s.ids[id] = struct{}{}
	}
}

func (s *idSet) Values() []int64 {
	ids := make([]int64, len(s.ids))
	i := 0
	for id := range s.ids {
		ids[i] = id
		i++
	}
	return ids
}
