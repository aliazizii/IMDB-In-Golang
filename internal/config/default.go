package config

var Default = Config{
	DB: Database{
		Host:     "localhost",
		Port:     "5432",
		Name:     "imdb",
		User:     "postgres",
		Password: "your_password",
	},
	Secret: "secret",
}
