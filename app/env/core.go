package env

// message
type EnvSpec struct {
	Id string `json:"id"`
}

// domain
type Env struct {
	Id string
}

// state
type envRoot struct {
	Id string
}

// port
type Api interface {
	Create(es EnvSpec) (Env, error)
}

// core
type service struct {
	repo repo
}

func (s *service) Create(es EnvSpec) (Env, error) {
	er := envRoot{}
	err := s.repo.Insert(er)
	if err != nil {
		return Env{}, err
	}
	return Env{}, nil
}

// port
type repo interface {
	Insert(er envRoot) error
}
