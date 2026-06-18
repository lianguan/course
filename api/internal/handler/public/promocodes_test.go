package public

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"ultrathreads/internal/domain"
	"ultrathreads/internal/service"
	mock_service "ultrathreads/internal/service/mocks"
)

func TestHandler_getPromocode(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPromoCodes, schoolId uint, code string, promocode domain.PromoCode)

	schoolId := uint(1)

	promocode := domain.PromoCode{
		Code:               "GOGOGO25",
		DiscountPercentage: 25,
	}

	setResponseBody := func(promocode domain.PromoCode) string {
		body, _ := json.Marshal(promocode)

		return string(body)
	}

	tests := []struct {
		name         string
		code         string
		schoolId     uint
		promocode    domain.PromoCode
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:      "ok",
			code:      "GOGOGO25",
			schoolId:  schoolId,
			promocode: promocode,
			mockBehavior: func(r *mock_service.MockPromoCodes, schoolId uint, code string, promocode domain.PromoCode) {
				r.EXPECT().GetByCode(context.Background(), schoolId, code).Return(promocode, nil)
			},
			statusCode:   200,
			responseBody: setResponseBody(promocode),
		},
		{
			name:      "empty code",
			code:      "",
			schoolId:  schoolId,
			promocode: promocode,
			mockBehavior: func(r *mock_service.MockPromoCodes, schoolId uint, code string, promocode domain.PromoCode) {
			},
			statusCode:   404,
			responseBody: `404 page not found`,
		},
		{
			name:      "service error",
			code:      "GOGOGO25",
			schoolId:  schoolId,
			promocode: promocode,
			mockBehavior: func(r *mock_service.MockPromoCodes, schoolId uint, code string, promocode domain.PromoCode) {
				r.EXPECT().GetByCode(context.Background(), schoolId, code).Return(promocode, errors.New("failed to get promocode"))
			},
			statusCode:   500,
			responseBody: `{"message":"failed to get promocode"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			s := mock_service.NewMockPromoCodes(c)
			tt.mockBehavior(s, tt.schoolId, tt.code, tt.promocode)

			services := &service.Services{PromoCodes: s}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/promocodes/:code", func(c *gin.Context) {
				c.Set(schoolCtx, domain.School{
					ID: schoolId,
				})
			}, handler.getPromo)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/promocodes/%s", tt.code), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}
