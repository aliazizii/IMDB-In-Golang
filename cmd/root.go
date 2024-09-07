package cmd

import (
	"github.com/aliazizii/IMDB-In-Golang/cmd/migrate"
	"github.com/aliazizii/IMDB-In-Golang/cmd/serve"
	"github.com/aliazizii/IMDB-In-Golang/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func Execute() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	root := &cobra.Command{
		Use:   "imdb",
		Short: "short",
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(cfg.Logger.Level)

	serve.Register(root, cfg)
	migrate.Register(root, cfg)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
