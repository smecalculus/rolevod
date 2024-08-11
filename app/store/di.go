package store

import (
	"go.uber.org/fx"

	ce "smecalculus/rolevod/app/core/env"
	se "smecalculus/rolevod/app/store/env"
)

var Module = fx.Module("app/store",
	fx.Provide(
		// fx.Private,
		fx.Annotate(se.NewRepoPgx, fx.As(new(ce.Repo))),
	),
)
