package transport

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thalysonr/poc_go/common/errors"
	"github.com/thalysonr/poc_go/common/log"
	"github.com/thalysonr/poc_go/user/internal/app/model"
	"github.com/thalysonr/poc_go/user/internal/app/service"
)

type HttpServer struct {
	app         *fiber.App
	userService service.UserService
}

func NewHttpServer(userService service.UserService) *HttpServer {
	return &HttpServer{
		userService: userService,
	}
}

func (h *HttpServer) Start() error {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	})
	h.loadRoutes(app)
	h.app = app
	return app.Listen(":3000")
}

func (h *HttpServer) Stop() error {
	if h.app != nil {
		log.GetLogger().Debug("Shutting down http...")
		err := h.app.Shutdown()
		log.GetLogger().Debug("Http shut down successfully")
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
///////                       AUXILIARY FUNCTIONS                        ///////
////////////////////////////////////////////////////////////////////////////////

func (h *HttpServer) loadRoutes(app *fiber.App) {
	app.Get("/user", func(c *fiber.Ctx) error {
		users, err := h.userService.FindAll(c.Context())
		if err != nil {
			return err
		}
		return c.JSON(users)
	})
	app.Post("/user", func(c *fiber.Ctx) error {
		var user model.User
		err := json.Unmarshal(c.Body(), &user)
		if err != nil {
			return err
		}
		id, err := h.userService.Create(c.Context(), user)
		if err != nil {
			return err
		}
		return c.JSON(id)
	})
}

func errorHandler(c *fiber.Ctx, err error) error {
	// Default 500 statuscode
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	if _, ok := err.(*errors.ValidationError); ok {
		// Override status code if fiber.Error type
		code = 400
		msg = "Bad Request"
	}
	// Set Content-Type: text/plain; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	// Return statuscode with error message
	return c.Status(code).SendString(msg)
}
