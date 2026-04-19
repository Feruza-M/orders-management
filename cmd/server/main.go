package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"order-management-app/internal/db"
	"order-management-app/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := db.ConfigFromEnv()
	pool, err := db.NewPool(cfg)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer pool.Close()

	h := handlers.New(pool)
	r := gin.Default()
	r.Use(corsMiddleware())

	r.Static("/app", "./web")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/app/")
	})

	api := r.Group("/api")
	{
		api.GET("/health", h.Health)
		api.GET("/orders", h.ListOrders)
		api.GET("/orders/:id", h.GetOrder)
		api.POST("/orders", h.CreateOrder)
		api.PUT("/orders/:id", h.UpdateOrder)
		api.DELETE("/orders/:id", h.DeleteOrder)
	}

	srv := &http.Server{
		Addr:              envOrDefault("APP_ADDR", ":8080"),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
