package main

type VertexSet struct {
	list map[*Vertex]struct{}
}

func NewVertexSet() *VertexSet {
	s := &VertexSet{}
	s.list = make(map[*Vertex]struct{})
	return s
}

func (s *VertexSet) has(v *Vertex) bool {
	for item := range s.list {
		if v.stationID == item.stationID {
			return true
		}
	}
	return false
}

func (s *VertexSet) add(v *Vertex) {
	if !s.has(v) {
		s.list[v] = struct{}{}
	}
}

func (s *VertexSet) remove(v *Vertex) {
	for item := range s.list {
		if v.stationID == item.stationID {
			delete(s.list, item)
		}
	}
}

func (s *VertexSet) clear() {
	s.list = make(map[*Vertex]struct{})
}

func (s *VertexSet) size() int {
	return len(s.list)
}

func (s *VertexSet) addMany(list ...*Vertex) {
	for _, v := range list {
		s.add(v)
	}
}

func (s *VertexSet) getList() []*Vertex {
	keys := make([]*Vertex, 0, len(s.list))
	for k := range s.list {
		keys = append(keys, k)
	}

	return keys
}

func (s *VertexSet) union(s2 *VertexSet) *VertexSet {
	res := NewVertexSet()
	for v := range s.list {
		res.add(v)
	}

	for v := range s2.list {
		res.add(v)
	}
	return res
}
