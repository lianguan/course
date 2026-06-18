package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
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

func TestHandler_adminCreateOffer(t *testing.T) {
	type mockBehavior func(r *mock_service.MockOffers, input service.CreateOfferInput)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		body         string
		school       domain.School
		input        service.CreateOfferInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:   "ok",
			body:   `{"name":"Test Offer","description":"Test","benefits":["b1"],"packages":[1],"price":{"value":100,"currency":"USD"},"paymentMethod":{"usesProvider":false}}`,
			school: school,
			input: service.CreateOfferInput{
				SchoolID:    school.ID,
				Name:        "Test Offer",
				Description: "Test",
				Benefits:    []string{"b1"},
				Price:       domain.Price{Value: 100, Currency: "USD"},
				PaymentMethod: domain.PaymentMethod{
					UsesProvider: false,
				},
				Packages: []uint{1},
			},
			mockBehavior: func(r *mock_service.MockOffers, input service.CreateOfferInput) {
				r.EXPECT().Create(context.Background(), input).Return(uint(1), nil)
			},
			statusCode:   201,
			responseBody: `{"id":1}`,
		},
		{
			name:         "invalid input",
			body:         `{wrong}`,
			school:       school,
			mockBehavior: func(r *mock_service.MockOffers, input service.CreateOfferInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:   "service error",
			body:   `{"name":"Test Offer","description":"Test","benefits":["b1"],"packages":[1],"price":{"value":100,"currency":"USD"},"paymentMethod":{"usesProvider":false}}`,
			school: school,
			input: service.CreateOfferInput{
				SchoolID:    school.ID,
				Name:        "Test Offer",
				Description: "Test",
				Benefits:    []string{"b1"},
				Price:       domain.Price{Value: 100, Currency: "USD"},
				PaymentMethod: domain.PaymentMethod{
					UsesProvider: false,
				},
				Packages: []uint{1},
			},
			mockBehavior: func(r *mock_service.MockOffers, input service.CreateOfferInput) {
				r.EXPECT().Create(context.Background(), input).Return(uint(0), errors.New("service error"))
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
			tt.mockBehavior(offers, tt.input)

			services := &service.Services{Offers: offers}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/admins/offers", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminCreateOffer)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/admins/offers", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminGetAllOffers(t *testing.T) {
	type mockBehavior func(o *mock_service.MockOffers, p *mock_service.MockPackages, schoolID uint)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		school       domain.School
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:   "ok",
			school: school,
			mockBehavior: func(o *mock_service.MockOffers, p *mock_service.MockPackages, schoolID uint) {
				o.EXPECT().GetAll(context.Background(), schoolID).Return([]domain.Offer{
					{ID: 1, Name: "Offer 1", Price: domain.Price{Value: 100, Currency: "USD"}},
				}, nil)
				p.EXPECT().GetByIds(context.Background(), gomock.Any()).Return([]domain.Package{}, nil)
			},
			statusCode:   200,
			responseBody: `{"data":[{"id":1,"name":"Offer 1","description":"","benefits":null,"packages":[],"price":{"value":100,"currency":"USD"},"paymentMethod":{"usesProvider":false,"provider":""}}],"count":0}`,
		},
		{
			name:   "service error",
			school: school,
			mockBehavior: func(o *mock_service.MockOffers, p *mock_service.MockPackages, schoolID uint) {
				o.EXPECT().GetAll(context.Background(), schoolID).Return(nil, errors.New("service error"))
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
			packages := mock_service.NewMockPackages(c)
			tt.mockBehavior(offers, packages, tt.school.ID)

			services := &service.Services{Offers: offers, Packages: packages}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/admins/offers", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminGetAllOffers)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/admins/offers", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminGetOfferById(t *testing.T) {
	type mockBehavior func(o *mock_service.MockOffers, p *mock_service.MockPackages, offerID uint)

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
			mockBehavior: func(o *mock_service.MockOffers, p *mock_service.MockPackages, offerID uint) {
				o.EXPECT().GetById(context.Background(), offerID).Return(domain.Offer{
					ID:          offerID,
					Name:        "Test Offer",
					Description: "Test",
					Price:       domain.Price{Value: 100, Currency: "USD"},
				}, nil)
				p.EXPECT().GetByIds(context.Background(), gomock.Any()).Return([]domain.Package{}, nil)
			},
			statusCode:   200,
			responseBody: `{"id":1,"name":"Test Offer","description":"Test","benefits":null,"packages":[],"price":{"value":100,"currency":"USD"},"paymentMethod":{"usesProvider":false,"provider":""}}`,
		},
		{
			name:         "invalid id",
			offerId:      0,
			mockBehavior: func(o *mock_service.MockOffers, p *mock_service.MockPackages, offerID uint) {},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
		{
			name:    "service error",
			offerId: 1,
			mockBehavior: func(o *mock_service.MockOffers, p *mock_service.MockPackages, offerID uint) {
				o.EXPECT().GetById(context.Background(), offerID).Return(domain.Offer{}, errors.New("service error"))
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
			packages := mock_service.NewMockPackages(c)
			tt.mockBehavior(offers, packages, tt.offerId)

			services := &service.Services{Offers: offers, Packages: packages}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/admins/offers/:id", handler.adminGetOfferById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/admins/offers/%d", tt.offerId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminUpdateOffer(t *testing.T) {
	type mockBehavior func(r *mock_service.MockOffers, input service.UpdateOfferInput)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		offerId      uint
		body         string
		school       domain.School
		input        service.UpdateOfferInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:    "ok",
			offerId: 1,
			body:    `{"name":"Updated Offer","description":"Updated"}`,
			school:  school,
			input: service.UpdateOfferInput{
				ID:          1,
				SchoolID:    school.ID,
				Name:        "Updated Offer",
				Description: "Updated",
			},
			mockBehavior: func(r *mock_service.MockOffers, input service.UpdateOfferInput) {
				r.EXPECT().Update(context.Background(), input).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:         "invalid input",
			offerId:      1,
			body:         `{wrong}`,
			school:       school,
			mockBehavior: func(r *mock_service.MockOffers, input service.UpdateOfferInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:    "service error",
			offerId: 1,
			body:    `{"name":"Updated Offer","description":"Updated"}`,
			school:  school,
			input: service.UpdateOfferInput{
				ID:          1,
				SchoolID:    school.ID,
				Name:        "Updated Offer",
				Description: "Updated",
			},
			mockBehavior: func(r *mock_service.MockOffers, input service.UpdateOfferInput) {
				r.EXPECT().Update(context.Background(), input).Return(errors.New("service error"))
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
			tt.mockBehavior(offers, tt.input)

			services := &service.Services{Offers: offers}
			handler := Handler{services: services}

			r := gin.New()
			r.PUT("/admins/offers/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminUpdateOffer)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/admins/offers/%d", tt.offerId), strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminDeleteOffer(t *testing.T) {
	type mockBehavior func(r *mock_service.MockOffers, schoolID, offerID uint)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		offerId      uint
		school       domain.School
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:    "ok",
			offerId: 1,
			school:  school,
			mockBehavior: func(r *mock_service.MockOffers, schoolID, offerID uint) {
				r.EXPECT().Delete(context.Background(), schoolID, offerID).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:    "service error",
			offerId: 1,
			school:  school,
			mockBehavior: func(r *mock_service.MockOffers, schoolID, offerID uint) {
				r.EXPECT().Delete(context.Background(), schoolID, offerID).Return(errors.New("service error"))
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
			tt.mockBehavior(offers, tt.school.ID, tt.offerId)

			services := &service.Services{Offers: offers}
			handler := Handler{services: services}

			r := gin.New()
			r.DELETE("/admins/offers/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminDeleteOffer)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/admins/offers/%d", tt.offerId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}
