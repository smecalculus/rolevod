package prog

import (
	"log/slog"

	pooltypedef "orglang/go-engine/pool/typedef"
	proctypedef "orglang/go-engine/proc/typedef"
)

type API interface {
	Create(Spec) error
}

type Spec struct {
	Pools []pooltypedef.DefSpec
	Procs []proctypedef.DefSpec
}

type service struct {
	log *slog.Logger
}

func newAPI() API {
	return new(service)
}

func newService(log *slog.Logger) *service {
	return &service{log: log}
}

func (s *service) Create(spec Spec) error {
	return nil
}
