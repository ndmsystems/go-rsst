package rsst

import (
	rsstApi "github.com/ndmsystems/go-rsst/api"
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
	for i := range in {
		info := in[i]
		if fn, ok := s.handlers[info.ID]; ok {
			fn(&info)
			in[i] = info
		}
	}
}
