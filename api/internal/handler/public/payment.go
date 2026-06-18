package public

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"ultrathreads/internal/domain"
	"ultrathreads/pkg/payment/fondy"
)

func (h *Handler) initCallbackRoutes(api *gin.RouterGroup) {
	callback := api.Group("/callback")
	{
		callback.POST("/fondy", h.handleFondyCallback)
	}
}

func (h *Handler) handleFondyCallback(c *gin.Context) {
	if c.Request.UserAgent() != fondy.UserAgent {
		handler.NewResponse(c, http.StatusForbidden, "forbidden")

		return
	}

	var inp fondy.Callback
	if err := c.BindJSON(&inp); err != nil {
		handler.NewResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Payments.ProcessTransaction(c.Request.Context(), inp); err != nil {
		if errors.Is(err, domain.ErrTransactionInvalid) {
			handler.NewResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		handler.NewResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}
