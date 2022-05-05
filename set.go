package main

type VertexSet struct {
	list map[*Vertex]struct{}
}

func (s *VertexSet) Has(v *Vertex) bool {
	for item := range s.list {
		if v.stationID == item.stationID {
			return true
		}
	}
	return false
}

func (s *VertexSet) Add(v *Vertex) {
	if !s.Has(v) {
		s.list[v] = struct{}{}
	}
}

func (s *VertexSet) Remove(v *Vertex) {
	for item := range s.list {
		if v.stationID == item.stationID {
			delete(s.list, item)
		}
	}
}

func (s *VertexSet) Clear() {
	s.list = make(map[*Vertex]struct{})
}

func (s *VertexSet) Size() int {
	return len(s.list)
}

func NewVertexSet() *VertexSet {
	s := &VertexSet{}
	s.list = make(map[*Vertex]struct{})
	return s
}

func (s *VertexSet) AddMulti(list ...*Vertex) {
	for _, v := range list {
		s.Add(v)
	}
}

func (s *VertexSet) GetList() []*Vertex {
	keys := make([]*Vertex, 0, len(s.list))
	for k := range s.list {
		keys = append(keys, k)
	}

	return keys
}

type FilterFunc func(v *Vertex) bool

func (s *VertexSet) Filter(P FilterFunc) *VertexSet {
	res := NewVertexSet()
	for v := range s.list {
		if P(v) == false {
			continue
		}
		res.Add(v)
	}
	return res
}

func (s *VertexSet) Union(s2 *VertexSet) *VertexSet {
	res := NewVertexSet()
	for v := range s.list {
		res.Add(v)
	}

	for v := range s2.list {
		res.Add(v)
	}
	return res
}

func (s *VertexSet) Intersect(s2 *VertexSet) *VertexSet {
	res := NewVertexSet()
	for v := range s.list {
		if s2.Has(v) == false {
			continue
		}
		res.Add(v)
	}
	return res
}

func (s *VertexSet) Difference(s2 *VertexSet) *VertexSet {
	res := NewVertexSet()
	for v := range s.list {
		if s2.Has(v) {
			continue
		}
		res.Add(v)
	}
	return res
}
