package router

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/paimon_bank/internal/customErr"
	"github.com/paimon_bank/internal/handler"
	"github.com/paimon_bank/internal/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

type RouteConfig struct {
	App                *fiber.App
	ImageHandler       *handler.ImageHandler
	BalanceHandler     *handler.BalanceHandler
	UserHandler        *handler.UserHandler
	TransactionHandler *handler.TranssactionHandler
}

var (
	httpRequestProm = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_histogram",
		Help:    "Histogram of the http request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10),
	}, []string{"path", "method", "status"})
)

func (c *RouteConfig) Setup() {

	//prome := fiberprometheus.New("paimon-bank-app")
	//prome.RegisterAt(c.App, "/metrics")
	//c.App.Use(prome.Middleware)

	//c.App.Use(FiberPrometheusMiddleware)
	metrics := c.App.Group("/metrics", FiberPrometheusMiddleware)
	metrics.Get("", adaptor.HTTPHandler(promhttp.Handler()))
	//c.App.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// Health check
	c.App.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "ok",
		})
	}, FiberPrometheusMiddleware)

	authMiddleware := middleware.NewAuthMiddleware()
	// promotheus middleware
	c.App.Post("/v1/user/register", FiberPrometheusMiddleware, c.UserHandler.Register)
	c.App.Post("/v1/user/login", FiberPrometheusMiddleware, c.UserHandler.Login)

	//c.App.Post("/v1/user/register", c.UserHandler.Register)
	//c.App.Post("/v1/user/login", c.UserHandler.Login, adaptor.HTTPHandler(promhttp.Handler()))
	//prometheus := c.App.Group("", FiberPrometheusMiddleware)
	//prometheus.Use(FiberPrometheusMiddleware)

	// Image
	image := c.App.Group("/v1/image", FiberPrometheusMiddleware, authMiddleware)
	image.Post("", c.ImageHandler.Upload)

	// Balance
	balance := c.App.Group("/v1/balance", FiberPrometheusMiddleware, authMiddleware)
	balance.Post("", c.BalanceHandler.Create)
	balance.Get("", c.BalanceHandler.ListBalance)
	balance.Get("/history", FiberPrometheusMiddleware, c.TransactionHandler.TransactionHistory)

	// transaction
	transaction := c.App.Group("/v1/transaction", FiberPrometheusMiddleware, authMiddleware)
	transaction.Post("", c.TransactionHandler.Create)

}

func FiberPrometheusMiddleware(ctx *fiber.Ctx) error {
	start := time.Now()
	method := ctx.Route().Method
	path := ctx.Route().Path
	err := ctx.Next()

	status := fiber.StatusInternalServerError
	if err != nil {
		if e, ok := err.(*fiber.Error); ok {
			// Get correct error code from fiber.Error type
			status = e.Code
		}

		if e, ok := err.(customErr.CustomError); ok {
			// This is a custom error, handle it accordingly
			status = e.Status()
		}
	} else {
		status = ctx.Response().StatusCode()
	}
	statusCode := strconv.Itoa(status)

	httpRequestProm.WithLabelValues(path, method, statusCode).Observe(float64(time.Since(start).Milliseconds()))

	return err
}
