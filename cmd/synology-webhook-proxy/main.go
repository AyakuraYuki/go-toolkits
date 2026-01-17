package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/AyakuraYuki/go-toolkits/cmd/synology-webhook-proxy/internal/dingtalk"
	"github.com/AyakuraYuki/go-toolkits/cmd/synology-webhook-proxy/internal/middleware"
	"github.com/AyakuraYuki/go-toolkits/pkg/cjson"
	"github.com/AyakuraYuki/go-toolkits/pkg/signals"
)

const (
	port = ":6789"
)

func init() {
	cjson.RegisterFuzzyDecoders()
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	engine.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.Cors())

	engine.NoRoute(func(c *gin.Context) { c.AbortWithStatus(http.StatusNotFound) })

	engine.POST("/dingtalk", SendDingTalkMessage)

	hs := &http.Server{
		Addr:         port,
		Handler:      engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		err := hs.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	signals.GraceStop(func() {
		_ = hs.Shutdown(context.Background())
	})
}

func SendDingTalkMessage(c *gin.Context) {
	token := c.PostForm("token")
	secret := c.PostForm("secret")
	title := c.PostForm("title")
	text := c.PostForm("text")

	dingtalk.Notify(&dingtalk.Message{
		Title:  title,
		Text:   text,
		Token:  token,
		Secret: secret,
	})

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "ok"})
}
