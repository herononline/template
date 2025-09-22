package example

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type controller interface {
	healthHandler(ctx *fiber.Ctx) error
	websocketHandler(ctx *websocket.Conn)
}

// healthHandler godoc
// @Summary      Check service health
// @Description  Returns information about the service health
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{} "Service health information"
// @Failure      500 {object} map[string]string "Service health check failed"
// @Router       /health [get]
func (c *appController) healthHandler(ctx *fiber.Ctx) error {
	res, err := c.service.health()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (con *appController) websocketHandler(ctx *websocket.Conn) {
	req := ctx.Locals("validatedData").(*WsConn)
	clientID := req.ID

	wsMux.Lock()
	Cws[clientID] = ctx
	wsMux.Unlock()
	log.Printf("Novo cliente conectado: %v", clientID)

	defer func() {
		wsMux.Lock()
		delete(Cws, clientID)
		wsMux.Unlock()
		log.Printf("Cliente desconectado: %v", clientID)
		ctx.Close()
	}()

	welcomeMsg := map[string]any{
		"type":    "welcome",
		"message": "Bem-vindo ao WebSocket seguro!",
		"sub":     clientID,
	}
	if err := ctx.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Erro ao enviar mensagem de boas-vindas: %v", err)
		return
	}

	for {
		var msg map[string]any
		err := ctx.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Erro na leitura: %v", err)
			}
			break
		}

		log.Printf("Mensagem recebida de %v: %v", clientID, msg)
	}
}
