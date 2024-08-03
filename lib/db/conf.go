package db

type props struct {
	Protocol protocol
	Mapping  mapping
}

type protocol struct {
	Mode     string
	Postgres postgres
}

type mapping struct {
	Modes []string
}

type postgres struct {
	Url string
}
