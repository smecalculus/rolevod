package pool

import (
	"context"
	"log/slog"

	"smecalculus/rolevod/lib/data"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/rev"

	"smecalculus/rolevod/internal/chnl"

	"smecalculus/rolevod/app/sig"
)

type ID = id.ADT
type Rev = rev.ADT
type Title = string

type Spec struct {
	Title  string
	SupID  id.ADT
	DepIDs []sig.ID
}

type Ref struct {
	ID    id.ADT
	Rev   rev.ADT
	Title string
}

type Snap struct {
	ID    id.ADT
	Title string
	Subs  []Ref
}

type Root struct {
	ID    id.ADT
	Rev   rev.ADT
	Title string
	SupID id.ADT
}

// Port
type API interface {
	Create(Spec) (Root, error)
	Retrieve(id.ADT) (Snap, error)
	RetreiveRefs() ([]Ref, error)
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

type service struct {
	pools    Repo
	operator data.Operator
	log      *slog.Logger
}

func newService(pools Repo, operator data.Operator, l *slog.Logger) *service {
	name := slog.String("name", "poolService")
	return &service{pools, operator, l.With(name)}
}

func (s *service) Create(spec Spec) (_ Root, err error) {
	ctx := context.Background()
	s.log.Debug("creation started", slog.Any("spec", spec))
	root := Root{
		ID:    id.New(),
		Rev:   rev.Initial(),
		Title: spec.Title,
		SupID: spec.SupID,
	}
	s.operator.Explicit(ctx, func(ds data.Source) error {
		err = s.pools.Insert(ds, root)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("creation failed")
		return Root{}, err
	}
	s.log.Debug("creation succeeded", slog.Any("id", root.ID))
	return root, nil
}

func (s *service) Retrieve(rid id.ADT) (snap Snap, err error) {
	ctx := context.Background()
	s.operator.Implicit(ctx, func(ds data.Source) {
		snap, err = s.pools.SelectByID(ds, rid)
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("id", rid))
		return Snap{}, err
	}
	return snap, nil
}

func (s *service) RetreiveRefs() (refs []Ref, err error) {
	ctx := context.Background()
	s.operator.Implicit(ctx, func(ds data.Source) {
		refs, err = s.pools.SelectAll(ds)
	})
	if err != nil {
		s.log.Error("retrieval failed")
		return nil, err
	}
	return refs, nil
}

// Port
type Repo interface {
	Insert(data.Source, Root) error
	SelectByID(data.Source, id.ADT) (Snap, error)
	SelectAll(data.Source) ([]Ref, error)
	Transfer(source data.Source, giver id.ADT, taker id.ADT, pids []chnl.ID) error
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	ConvertRootToRef func(Root) Ref
)
