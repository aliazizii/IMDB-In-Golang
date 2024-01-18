package main

import (
	"github.com/aliazizii/IMDB-In-Golang/internal/auth"
	"github.com/aliazizii/IMDB-In-Golang/internal/handler"
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/comment"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/movie"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/user"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/vote"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func main() {
	app := echo.New()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open(sqlite.Open("test10.db"), new(gorm.Config))
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&model.Movie{}, &model.User{}, &model.Comment{}, &model.Vote{})

	// Add Admin
	// Admin should be in config
	db.Save(&model.User{
		Username: "ali",
		Password: auth.Hash("1234"),
		Role:     auth.AdminRoleCode,
	})

	u := user.NewSQL(db)
	userHandler := handler.User{
		Store: u,
		// It should be in config
		JwtSecret: "secret",
		Logger:    logger,
	}

	m := movie.NewSQL(db)
	movieHandler := handler.Movie{
		Store:  m,
		Logger: logger,
	}

	c := comment.NewSQL(db)
	commentHandler := handler.Comment{
		Store:  c,
		Logger: logger,
	}

	v := vote.NewSQL(db)
	voteHandler := handler.Vote{
		Store:  v,
		Logger: logger,
	}

	// unrestricted endpoints
	app.POST("/signup", userHandler.SignUp)
	app.POST("/login", userHandler.Login)
	app.GET("/movies", movieHandler.AllMovies)
	app.GET("/comments", movieHandler.Comments)
	app.GET("/movie/:id", movieHandler.Movie)

	// restricted endpoints. requires authorization
	adminArea := app.Group("/admin")
	userArea := app.Group("/user")

	jwtConfig := echojwt.Config{
		SigningKey: []byte("secret"),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &auth.JwtClaim{}
		},
	}
	adminArea.Use(echojwt.WithConfig(jwtConfig))
	userArea.Use(echojwt.WithConfig(jwtConfig))

	// admin area routing
	adminArea.POST("/movie", movieHandler.AddMovie)
	adminArea.PUT("/movie/:id", movieHandler.UpdateMovie)
	adminArea.DELETE("/movie/:id", movieHandler.DeleteMovie)
	adminArea.PUT("/comment/:id", commentHandler.UpdateComment)
	adminArea.DELETE("/comment/:id", commentHandler.DeleteComment)

	//user area routing
	userArea.POST("/comment", commentHandler.Comment)
	userArea.POST("vote", voteHandler.Vote)

	if err := app.Start(":1234"); err != nil {
		log.Fatal("cannot start the http server")
	}
}
