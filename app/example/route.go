package example

import (
	"github.com/go-gorote/auth"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (r *appRouter) RegisterRouter(router fiber.Router) {
	router.Get("/health", r.controller.healthHandler)
	r.ws(router.Group("/ws"))
}

func (r *appRouter) ws(router fiber.Router) {
	router.Get(
		"/:id",
		gorote.IsWsMiddleware(),
		gorote.ValidationMiddleware(&WsConn{}),
		gorote.JWTProtectedRSA(&auth.JwtClaims{}, r.publicKey),
		websocket.New(r.controller.websocketHandler),
	)
}
