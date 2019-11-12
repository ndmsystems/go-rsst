package rsst

import (
	rsstApi "github.com/tdx/go/api/rsst"
)

type svc struct {
	handlers map[uint16]rsstApi.Handler
}

// New ...
func New() rsstApi.Rsst {
	return &svc{
		handlers: make(map[uint16]rsstApi.Handler),
	}
}

func (s *svc) AddHandler(id uint16, f rsstApi.Handler) {
	s.handlers[id] = f
}

func (s *svc) Process(in []rsstApi.Info) {
	for _, info := range in {
		if fn, ok := s.handlers[info.ID]; ok {
			fn(&info)
		}
	}
}
