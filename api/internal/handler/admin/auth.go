package admin

import (
	"errors"
	"net/http"

	"ultrathreads/internal/domain"

	"github.com/gin-gonic/gin"
)

type signInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,max=64"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

// @Summary Admin SignIn
// @Tags admins-auth
// @Description admin sign in
// @ModuleID adminSignIn
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/sign-in [post]
func (h *Handler) adminSignIn(c *gin.Context) {
	var inp signInInput
	if err := c.BindJSON(&inp); err != nil {
		handler.NewResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		handler.NewResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	res, err := h.services.Admins.SignIn(c.Request.Context(), domain.SchoolSignInInput{
		Email:    inp.Email,
		Password: inp.Password,
		SchoolID: school.ID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			handler.NewResponse(c, http.StatusUnauthorized, err.Error())
		} else {
			handler.NewResponse(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// @Summary Admin Refresh Tokens
// @Tags admins-auth
// @Description admin refresh tokens
// @Accept  json
// @Produce  json
// @Param input body refreshInput true "refresh info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/auth/refresh [post]
func (h *Handler) adminRefresh(c *gin.Context) {
	var inp refreshInput
	if err := c.BindJSON(&inp); err != nil {
		handler.NewResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		handler.NewResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	res, err := h.services.Admins.RefreshTokens(c.Request.Context(), school.ID, inp.Token)
	if err != nil {
		handler.NewResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
