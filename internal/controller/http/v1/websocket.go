package v1

import (
	"github.com/gin-gonic/gin"
	"lemon_be/internal/usecase"
	"lemon_be/pkg/logger"
	"net/http"
)

type WebbsocketRoutes struct {
	w usecase.Websocket
	l logger.Interface
}

func NewWebsocketRoutes(h *gin.RouterGroup, w usecase.Websocket, l logger.Interface) {

	r := &WebbsocketRoutes{w, l}

	h.GET("/ws", r.websocketHandlerRoute)

}

// websocketHandler handler saat buka koneksi websocketc
func (r *WebbsocketRoutes) websocketHandlerRoute(c *gin.Context) {
	err := r.w.WebsocketHandler(c.Writer, c.Request, c.Request.Context())
	if err != nil {
		if err == usecase.WebsocketUnauthorizedError {
			//w.WriteHeader(http.StatusUnauthorized)
			ErrorResponse(c, http.StatusUnauthorized, "Websocket connection unauthorized")
			return
		}

		ErrorResponse(c, http.StatusInternalServerError, "Websocket service error")
		return
	}

}
