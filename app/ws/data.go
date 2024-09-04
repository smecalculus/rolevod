package ws

type envRootData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type tpIntroData struct {
	EnvID string `db:"env_id"`
	TpID  string `db:"tp_id"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/app/dcl:To.*
var (
	// goverter:ignore Tps Exps
	dataToEnvRoot   func(envRootData) (EnvRoot, error)
	dataFromEnvRoot func(EnvRoot) envRootData
	dataToEnvRoots  func([]envRootData) ([]EnvRoot, error)
	dataFromTpIntro func(TpIntro) tpIntroData
)
