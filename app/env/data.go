package env

type envRootData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	// goverter:ignore Tps Exps
	dataToEnvRoot   func(envRootData) (EnvRoot, error)
	dataFromEnvRoot func(EnvRoot) envRootData
)
