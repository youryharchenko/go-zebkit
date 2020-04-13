package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

	//"log"

	//"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

var secretKey string

func main() {
	dir := flag.String("dir", ".", "static dir")
	port := flag.String("port", "1323", "port")
	debug := flag.Bool("debug", false, "debug")
	flag.Parse()

	log.Info("Port:", *port)
	log.Info("Root:", *dir)

	var ok bool
	secretKey, ok = os.LookupEnv("JWTSECRET")
	if !ok {
		secretKey = "secret key"
	}

	e := echo.New()
	e.Debug = *debug

	e.Logger.SetLevel(log.DEBUG)
	// Log all requests
	e.Use(middleware.Logger())
	// Recover from panic
	e.Use(middleware.Recover())
	//
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper:      JWTSkipper,
		ErrorHandler: JWTErrorHandler,
		SigningKey:   []byte(secretKey),
	}))
	//
	e.Use(middleware.Static(*dir))

	e.POST("/login", Login)
	e.GET("/auth-test", AuthTest)

	//e.POST("/api", API)

	e.Logger.Fatal(e.Start(":" + *port))
}

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Login -
func Login(ctx echo.Context) (err error) {

	login := new(LoginRequest)
	if err = ctx.Bind(login); err != nil {
		return
	}

	ctx.Logger().Info("Login -- User:", login.Username)

	// Throws unauthorized error
	if login.Username != "admin" || login.Password != "admin" {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["name"] = login.Username
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"token": t,
	})

}

// AuthTest -
func AuthTest(ctx echo.Context) (err error) {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)

	return ctx.JSON(http.StatusOK, map[string]string{
		"authed": name,
	})
}

// API -
func API(ctx echo.Context) (err error) {
	return
}

// JWTSkipper -
func JWTSkipper(ctx echo.Context) bool {
	p := ctx.Path()
	m := ctx.Request().Method

	switch m {
	case "GET":
		switch {
		case p == "":
			return true
		case p == "/":
			return true
		case strings.HasPrefix(p, "/zebkit"):
			return true
		case strings.HasPrefix(p, "/main"):
			return true
		case strings.HasPrefix(p, "/rs/"):
			return true
		}
	case "POST":
		switch {
		case p == "/login":
			return true
		}
	}
	//log.Info("Unknown - method:", m, ", path:", p)
	return false
}

// JWTErrorHandler -
func JWTErrorHandler(err error) error {
	return &echo.HTTPError{
		Code:    http.StatusUnauthorized,
		Message: "invalid or expired jwt: " + err.Error(),
	}
}
