package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

	//"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	acc "github.com/youryharchenko/accountancy"
)

var secretKey string
var driver = "sqlite3"
var connection = ""

func init() {

}

func main() {
	dir := flag.String("dir", ".", "static dir")
	db := flag.String("db", ".", "db connection")
	query := flag.String("query", ".", "query folder")
	port := flag.String("port", "1323", "port")
	debug := flag.Bool("debug", false, "debug")
	flag.Parse()

	log.Info("Port:", *port)
	log.Info("Root:", *dir)
	log.Info("DB:", *db)

	resp, err := initAll(*db, *query)
	if err != nil {
		panic(err)
	}
	log.Info("initAll:" + resp)

	var ok bool
	secretKey, ok = os.LookupEnv("JWTSECRET")
	if !ok {
		secretKey = "secret key"
	}

	connection = *db

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
	e.POST("/api/run", RunAPI)

	e.Logger.Fatal(e.Start(":" + *port))
}

// Request -
type Request struct {
	Command string `json:"command"`
	Service string `json:"service"`
}

// DB -
type DB struct {
	Driver     string `json:"driver"`
	Connection string `json:"connection"`
	Show       bool   `json:"show"`
}

// RunQuery -
type RunQuery struct {
	Request Request                `json:"request"`
	DB      DB                     `json:"db"`
	Body    map[string]interface{} `json:"body"`
}

// RunAPI -
func RunAPI(ctx echo.Context) (err error) {
	runQuery := new(RunQuery)
	if err = ctx.Bind(runQuery); err != nil {
		return
	}

	if len(runQuery.DB.Driver) == 0 {
		runQuery.DB.Driver = driver
	}
	if len(runQuery.DB.Connection) == 0 {
		runQuery.DB.Connection = connection
	}

	buf, err := json.MarshalIndent(runQuery, "", "  ")
	if err != nil {
		return
	}

	ctx.Logger().Info("RunAPI:", string(buf))

	response, err := acc.Run(string(buf), nil, nil)
	if err != nil {
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal([]byte(response), &m)
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, m)
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
		case strings.HasPrefix(p, "/favicon"):
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
