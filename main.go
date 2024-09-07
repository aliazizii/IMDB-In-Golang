package main

import (
	"github.com/aliazizii/IMDB-In-Golang/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//cfg, err := config.Read()
	//fmt.Println(err)
	//fmt.Println(cfg.DB)
	//
	//app := echo.New()
	//logger, err := zap.NewDevelopment()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	////cfg := config.Default
	//
	//db, err := database.New(cfg.DB)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = db.AutoMigrate(&model.Movie{}, &model.User{}, &model.Comment{}, &model.Vote{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// add admin
	//// admin should be in config
	//db.Save(&model.User{
	//	Username: "ali",
	//	Password: auth.Hash("1234"),
	//	Role:     auth.AdminRoleCode,
	//})
	//
	//u := user.NewSQL(db)
	//userHandler := handler.User{
	//	Store:     u,
	//	JwtSecret: cfg.Secret,
	//	Logger:    logger,
	//}
	//
	//m := movie.NewSQL(db)
	//movieHandler := handler.Movie{
	//	Store:  m,
	//	Logger: logger,
	//}
	//
	//c := comment.NewSQL(db)
	//commentHandler := handler.Comment{
	//	Store:  c,
	//	Logger: logger,
	//}
	//
	//v := vote.NewSQL(db)
	//voteHandler := handler.Vote{
	//	Store:  v,
	//	Logger: logger,
	//}
	//
	//// unrestricted endpoints
	//app.POST("/signup", userHandler.SignUp)
	//app.POST("/login", userHandler.Login)
	//app.GET("/movies", movieHandler.AllMovies)
	//app.GET("/comments", movieHandler.Comments)
	//app.GET("/movie/:id", movieHandler.Movie)
	//
	//// restricted endpoints. requires authorization
	//adminArea := app.Group("/admin")
	//userArea := app.Group("/user")
	//
	//jwtConfig := echojwt.Config{
	//	SigningKey: []byte(cfg.Secret),
	//	NewClaimsFunc: func(c echo.Context) jwt.Claims {
	//		return &auth.JwtClaim{}
	//	},
	//}
	//adminArea.Use(echojwt.WithConfig(jwtConfig))
	//userArea.Use(echojwt.WithConfig(jwtConfig))
	//
	//// admin area routing
	//adminArea.POST("/movie", movieHandler.AddMovie)
	//adminArea.PUT("/movie/:id", movieHandler.UpdateMovie)
	//adminArea.DELETE("/movie/:id", movieHandler.DeleteMovie)
	//adminArea.PUT("/comment/:id", commentHandler.UpdateComment)
	//adminArea.DELETE("/comment/:id", commentHandler.DeleteComment)
	//
	////user area routing
	//userArea.POST("/comment", commentHandler.Comment)
	//userArea.POST("/vote", voteHandler.Vote)
	//
	//if err := app.Start(":1234"); err != nil {
	//	log.Fatal("cannot start the http server")
	//}

	cmd.Execute()
}
