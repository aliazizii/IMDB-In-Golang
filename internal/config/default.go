package config

var Default = Config{
	DB: Database{
		Host:     "localhost",
		Port:     "5432",
		Name:     "postgres",
		User:     "postgres",
		Password: "demo",
	},
	Secret: "secret",
}
