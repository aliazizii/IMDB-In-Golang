package config

type Config struct {
	DB     Database
	Secret string
}

type Database struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}
