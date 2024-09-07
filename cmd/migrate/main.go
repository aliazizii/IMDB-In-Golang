package migrate

import (
	"github.com/aliazizii/IMDB-In-Golang/config"
	"github.com/aliazizii/IMDB-In-Golang/internal/auth"
	"github.com/aliazizii/IMDB-In-Golang/internal/database"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main(cfg *config.Config) {
	db, err := database.New(cfg.DB)
	if err != nil {
		logrus.Fatal(err)
	}

	err = db.AutoMigrate(&model.Movie{}, &model.User{}, &model.Comment{}, &model.Vote{})
	if err != nil {
		logrus.Fatal(err)
	}

	db.Save(&model.User{
		Username: cfg.Admin.Username,
		Password: auth.Hash(cfg.Admin.Password),
		Role:     auth.AdminRoleCode,
	})
}

func Register(root *cobra.Command, cfg *config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "migrate",
			Short: "database migration",
			Run: func(_ *cobra.Command, _ []string) {
				main(cfg)
			},
		},
	)
}
