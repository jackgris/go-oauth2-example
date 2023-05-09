package main

import (
	"bytes"
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

func main() {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// client memory store
	clientStore := store.NewClientStore()

	manager.MapClientStorage(clientStore)

	// create the default authorization server
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	app := fiber.New()
	app.Get("/token", Token(srv))
	app.Get("/credentials", Credentials(clientStore))
	app.Get("/protected", ValidateToken(Protected, srv))
	app.Get("/", Home)
	log.Fatal(app.Listen(":3000"))
}

func Token(srv *server.Server) fiber.Handler {
	// first, create a Handlerfunc because go-oauth2/oauth2 library was
	// created to be used with the standard net/http library
	var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		_ = srv.HandleTokenRequest(w, r)
	}

	// wraps net/http HandlerFunc to fiber handler
	return func(c *fiber.Ctx) error {
		c.Request()
		handler := fasthttpadaptor.NewFastHTTPHandler(h)
		handler(c.Context())
		return nil
	}
}

func Credentials(clientStore *store.ClientStore) fiber.Handler {
	handler := func(c *fiber.Ctx) error {
		clientId := uuid.New().String()[:8]
		clientSecret := uuid.New().String()[:8]
		err := clientStore.Set(clientId, &models.Client{
			ID:     clientId,
			Secret: clientSecret,
			Domain: "http://localhost:9094",
		})
		if err != nil {
			return err
		}
		r := map[string]string{"CLIENT_ID": clientId, "CLIENT_SECRET": clientSecret}
		c.Status(http.StatusOK)
		return c.JSON(r)
	}

	return handler
}

func Home(c *fiber.Ctx) error {
	return c.SendString("Hello, I'm not protected 👋!")
}

func Protected(c *fiber.Ctx) error {
	return c.SendString("I'm protected 👋!")
}

func NotAllowed(c *fiber.Ctx) error {
	return c.SendString("You Shall Not Pass!")
}

func ValidateToken(hand fiber.Handler, srv *server.Server) fiber.Handler {

	handler := func(c *fiber.Ctx) error {

		ctx := context.TODO()
		method := c.Method()
		url := c.OriginalURL()
		body := c.Body()
		r, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
		if err != nil {
			return NotAllowed(c)
		}
		_, err = srv.ValidationBearerToken(r)
		if err != nil {
			return NotAllowed(c)
		}
		return hand(c)
	}

	return handler
}
