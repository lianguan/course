package public

import (
	"errors"
	"net/http"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/service/dto"

	"github.com/gin-gonic/gin"
)

// @Summary Get PromoCode By Code
// @Tags promocodes
// @Description  get promocode by code
// @ModuleID getPromo
// @Accept  json
// @Produce  json
// @Param code path string true "code"
// @Success 200 {object} domain.PromoCode
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /promocodes/{code} [get]
func (h *Handler) getPromo(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		handler.NewResponse(c, http.StatusBadRequest, "empty code param")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		handler.NewResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	promocode, err := h.services.PromoCodes.GetByCode(c.Request.Context(), school.ID, code)
	if err != nil {
		if errors.Is(err, domain.ErrPromoNotFound) {
			handler.NewResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		handler.NewResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dto.PromoCodeToResponse(promocode))
}
