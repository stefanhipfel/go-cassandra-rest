package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"sap/errorlog-rest-dataingestion/cassandra"
	"sap/errorlog-rest-dataingestion/errorLogMessages"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	defaultPort    = "3000"
	algorithmRS256 = "RS256"
)

func main() {

	CassandraSession := cassandra.Session
	defer CassandraSession.Close()

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		log.Printf("Warning, PORT not set. Defaulting to %+v", defaultPort)
		port = defaultPort
	}

	publicKey, err := GetVerificationKey()
	if err != nil {
		log.Panic("can't get verifyKey")
	}

	var DefaultJWTConfig = middleware.JWTConfig{
		Skipper:       middleware.DefaultSkipper,
		SigningMethod: algorithmRS256,
		ContextKey:    "user",
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		AuthScheme:    "Bearer",
		SigningKey:    publicKey,
		Claims:        jwt.MapClaims{},
	}
	e := echo.New()

	//Custom Error Handler
	//e.HTTPErrorHandler = customHTTPErrorHandler
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	//The server runs each handler in a separate goroutine so that it can serve multiple requests simultaneously
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	r := e.Group("/v1")
	e.Use(middleware.JWTWithConfig(DefaultJWTConfig))
	r.Use(checkClaimsHandler)
	r.GET("/ErrorMessages/", errorLogMessages.Get)
	e.Logger.Fatal(e.Start(":" + port))

}

// checkClaimsHandler middleware checks jwt claims.
func checkClaimsHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		zid := claims["zid"]
		if zid != "" {
			c.Set("tenant", zid)
			return next(c)
		}
		return echo.NewHTTPError(http.StatusBadRequest, "Tenant-Id is missing")
	}
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}
