package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/service"
	mock_service "ultrathreads/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func TestHandler_getOffer(t *testing.T) {
	type mockBehavior func(r *mock_service.MockOffers, id uint)

	tests := []struct {
		name         string
		offerId      uint
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:    "ok",
			offerId: 1,
			mockBehavior: func(r *mock_service.MockOffers, id uint) {
				r.EXPECT().GetById(context.Background(), id).Return(domain.Offer{
					ID:          id,
					Name:        "Test Offer",
					Description: "Test Description",
					Price: domain.Price{
						Value:    100,
						Currency: "USD",
					},
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"id":1,"name":"Test Offer","description":"Test Description","benefits":null,"schoolId":0,"packages":null,"price":{"value":100,"currency":"USD"},"paymentMethod":{"usesProvider":false,"provider":""}}`,
		},
		{
			name:         "invalid id",
			offerId:      0,
			mockBehavior: func(r *mock_service.MockOffers, id uint) {},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
		{
			name:    "offer not found",
			offerId: 1,
			mockBehavior: func(r *mock_service.MockOffers, id uint) {
				r.EXPECT().GetById(context.Background(), id).Return(domain.Offer{}, domain.ErrPromoNotFound)
			},
			statusCode:   400,
			responseBody: `{"message":"promocode doesn't exists"}`,
		},
		{
			name:    "service error",
			offerId: 1,
			mockBehavior: func(r *mock_service.MockOffers, id uint) {
				r.EXPECT().GetById(context.Background(), id).Return(domain.Offer{}, errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			offers := mock_service.NewMockOffers(c)
			tt.mockBehavior(offers, tt.offerId)

			services := &service.Services{Offers: offers}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/offers/:id", handler.getOffer)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/offers/%d", tt.offerId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_getSchoolSettings(t *testing.T) {
	tests := []struct {
		name         string
		school       domain.School
		statusCode   int
		responseBody string
	}{
		{
			name: "ok",
			school: domain.School{
				ID:          1,
				Name:        "Test School",
				Subtitle:    "Test Subtitle",
				Description: "Test Description",
				Settings: domain.Settings{
					Domains: []string{"example.com"},
					Color:   "#ffffff",
				},
			},
			statusCode:   200,
			responseBody: `{"name":"Test School","subtitle":"Test Subtitle","description":"Test Description","settings":{"color":"#ffffff","domains":["example.com"],"contactInfo":{"businessName":"","registrationNumber":"","address":"","email":"","phone":""},"pages":{"confidential":"","serviceAgreement":"","newsletterConsent":""},"showPaymentImages":false,"logo":"","googleAnalyticsCode":"","fondy":{"merchantId":"","merchantPassword":"","connected":false},"sendpulse":{"id":"","secret":"","listId":"","connected":false},"disableRegistration":false}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			services := &service.Services{}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/settings", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.getSchoolSettings)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/settings", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}
