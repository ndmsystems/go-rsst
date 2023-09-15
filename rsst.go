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

func (s *svc) Process(authLevel rsstApi.AuthorizeLevel, in []rsstApi.Info) {
	for i := range in {
		info := in[i]

		if fn, ok := s.handlers[info.ID]; ok && authorizeRequest(authLevel, info.ID) {
			fn(&info)
			in[i] = info
		}
	}
}

func authorizeRequest(authLevel rsstApi.AuthorizeLevel, id uint16) bool {
	value := (id / 0x1000)
	if value == 0x1 || value%0x2 == 0 {
		return true
	}
	if authLevel >= rsstApi.TokenTrust {
		return true
	}
	return false
}
