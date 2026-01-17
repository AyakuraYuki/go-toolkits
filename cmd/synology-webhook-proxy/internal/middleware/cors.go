package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc:  func(_ string) bool { return true },
		MaxAge:           2 * time.Hour,
	})
}
