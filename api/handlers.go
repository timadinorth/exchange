package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
	"github.com/timadinorth/bet-exchange/model"
	"github.com/timadinorth/bet-exchange/util"
	"golang.org/x/crypto/bcrypt"
)

type SignUpReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignUpResp struct {
	Username string `json:"username"`
}

// SignUp godoc
//
// @Summary 	Register user
// @Description Creates new user
// @Tags 		auth
// @Accept 		json
// @Produce 	json
// @Param username body string true "Username"
// @Param password body string true "Password plan text"
// @Success 	201 		{object} 	SignUpResp
// @Failure		400			{object}	util.HTTPError
// @Failure		401			{object}	util.HTTPError
// @Failure		500			{object}	util.HTTPError
// @Router      /signup [post]
func (s *Server) SignUp(c *fiber.Ctx) error {
	var req SignUpReq

	if err := c.BodyParser(&req); err != nil {
		return util.NewError(c, fiber.StatusBadRequest, err)
	}

	if err := s.validator.Struct(&req); err != nil {
		return util.NewError(c, fiber.StatusBadRequest, err)
	}

	user := model.User{
		Username: req.Username,
		Password: req.Password,
	}

	if savedUser, err := user.Save(s.DB); err != nil {
		return util.NewErrorStr(c, fiber.StatusBadRequest, "invalid username or password")
	} else {
		resp := SignUpResp{
			Username: savedUser.Username,
		}
		return c.Status(fiber.StatusCreated).JSON(&fiber.Map{"data": resp})
	}
}

type SigninReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SignIn godoc
//
// @Summary 	User signin
// @Description Create new session for user
// @Tags 		auth
// @Accept 		json
// @Produce 	json
// @Param username body string true "Username"
// @Param password body string true "Password plan text"
// @Success 	200
// @Failure		400			{object}	util.HTTPError
// @Failure		401			{object}	util.HTTPError
// @Failure		500			{object}	util.HTTPError
// @Router      /signin [post]
func (s *Server) SignIn(c *fiber.Ctx) error {
	var req SigninReq

	if err := c.BodyParser(&req); err != nil {
		return util.NewError(c, fiber.StatusBadRequest, err)
	}

	dbUser := model.User{}
	if err := dbUser.FindByUsername(s.DB, req.Username); err != nil {
		return util.NewErrorStr(c, fiber.StatusForbidden, "Invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(req.Password)); err != nil {
		return util.NewErrorStr(c, fiber.StatusForbidden, "Invalid username or password")
	}

	session, err := s.Session.Get(c)
	if err != nil {
		return util.NewError(c, fiber.StatusInternalServerError, err)
	}

	sessionToken := xid.New().String()
	session.Set("username", dbUser.Username)
	session.Set("token", sessionToken)
	if err := session.Save(); err != nil {
		return util.NewError(c, http.StatusInternalServerError, err)
	}
	return c.Status(fiber.StatusOK).SendString("")
}

// Signout godoc
//
// @Summary 	User signout
// @Description Delete user session
// @Tags 		auth
// @Produce 	json
// @Success 	200
// @Failure		401			{object}	util.HTTPError
// @Failure		500			{object}	util.HTTPError
// @Router      /signout [post]
func (s *Server) SignOut(c *fiber.Ctx) error {
	session, err := s.Session.Get(c)
	if err != nil {
		return util.NewError(c, fiber.StatusInternalServerError, err)
	}

	session.Delete("username") // TODO: check
	session.Delete("token")

	if err := session.Destroy(); err != nil {
		return util.NewError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).SendString("")
}

// CreateCategory godoc
//
// @Summary 	Add a category
// @Description Creates new category and assigns unique Id
// @Tags 		categories
// @Accept 		json
// @Produce 	json
// @Param name body string true "Category name"
// @Param icon body string true "Icon path"
// @Param type body string true "Category type"
// @Success 	201 		{object} 	model.Category
// @Failure		400			{object}	util.HTTPError
// @Failure		401			{object}	util.HTTPError
// @Failure		500			{object}	util.HTTPError
// @Router      /categories [post]
func (s *Server) CreateCategory(c *fiber.Ctx) error {
	var category model.Category

	if err := c.BodyParser(&category); err != nil {
		return util.NewError(c, fiber.StatusBadRequest, err)
	}

	s.DB.Create(&category)
	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{"data": category})
}

// ListCategories godoc
//
// @Summary 	Get categories
// @Description Returns list of all categories
// @Tags 		categories
// @Produce 	json
// @Success 	200 		{array} 	model.Category
// @Failure		401			{object}	util.HTTPError
// @Failure		500			{object}	util.HTTPError
// @Router      /categories [get]
func (s *Server) ListCategories(c *fiber.Ctx) error {
	var categories []model.Category
	s.DB.Find(&categories)
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": categories})
}

// func (s *Server) ListCompetitions(c *gin.Context) {
// 	categoryId := c.Param("category_id")
// 	c.JSON(http.StatusOK, gin.H{"data": categoryId}) // FIXME
// }
