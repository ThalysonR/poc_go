package transport

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thalysonr/poc_go/common/errors"
	"github.com/thalysonr/poc_go/user/internal/app/model"
	"github.com/thalysonr/poc_go/user/internal/app/service"
	"github.com/thalysonr/poc_go/user/internal/config"
)

type HttpServer struct {
	app         *fiber.App
	cfg         config.Config
	userService service.UserService
}

func NewHttpServer(userService service.UserService) *HttpServer {
	return &HttpServer{
		userService: userService,
	}
}

func (h *HttpServer) ConfigChanged(cfg config.Config) bool {
	return !reflect.DeepEqual(h.cfg.Server.Http, cfg.Server.Http)
}

func (h *HttpServer) Start(cfg config.Config) error {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	})
	h.cfg = cfg
	h.loadRoutes(app)
	h.app = app
	return app.Listen(fmt.Sprintf(":%d", cfg.Server.Http.Port))
}

func (h *HttpServer) Stop() error {
	if h.app != nil {
		err := h.app.Shutdown()
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
