package msg

type props struct {
	Protocol protocol
	Mapping  mapping
}

type protocol struct {
	Modes []string
	Http  http
}

type mapping struct {
	Modes []string
}

type http struct {
	Port int
}
