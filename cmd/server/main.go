package main

import (
	"employee-management/api/delivery/httphandler"
	"employee-management/api/middleware"
	"employee-management/api/middleware/swagger"
	"employee-management/api/usecase"
	"employee-management/db"
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	// connect to db
	conn, err := db.Connect()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERROR] Failed to connect to db: %+v\n", err)
	}

	fmt.Println("connect", conn)

	// New gin server
	r := gin.New()

	// inject middlewares
	// newLogger
	// recover
	// swagger editor

	logger, _ := zap.NewProduction()
	r.Use(middleware.JSONMiddleware())

	/*  Add a ginzap middleware, which:
	    - Logs all requests, like a combined access and error log.
	    - Logs to stdout.
		- RFC3339 with UTC time format.
	*/

	// Host Swagger middleware
	r.Use(gin.WrapH(swagger.Middleware()))

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	/* Logs all panic to error log - stack means whether output the stack info. */
	r.Use(ginzap.RecoveryWithZap(logger, true))

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// employee endpoints
	employeeUsecase := usecase.NewEmployeeUsecase(conn)
	httphandler.NewEmployeeHandler(r, employeeUsecase)

	// Start the server
	_ = r.Run(":8080")
}
