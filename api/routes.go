package api

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/timadinorth/bet-exchange/docs"
	"github.com/timadinorth/bet-exchange/util"
)

func (s *Server) Auth(c *fiber.Ctx) error {
	session, err := s.Session.Get(c)
	if err != nil {
		return util.NewError(c, fiber.StatusInternalServerError, err)
	}
	token := session.Get("token")
	if token != nil {
		return c.Next()
	}
	return c.SendStatus(fiber.StatusUnauthorized)
}

func (s *Server) RegisterRoutes() {
	app := s.Web
	app.Use(etag.New())
	app.Use(logger.New())
	app.Get("/docs/*", swagger.HandlerDefault)
	v1 := app.Group("/api/v1")
	auth := v1.Group("auth")
	auth.Post("/signup", s.SignUp)
	auth.Post("/signin", s.SignIn)
	v1.Use(s.Auth)
	auth.Post("/signout", s.SignOut)
	v1.Get("/categories", s.ListCategories)
	v1.Post("/categories", s.CreateCategory)
	// {
	// 	//
	// 	// 	v1.GET("/categories/:category_id/competitions", s.ListCompetitions)
	// }
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
}
