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
)

func TestHandler_studentCreateOrder(t *testing.T) {
	type mockBehavior func(r *mock_service.MockOrders, studentId, offerId, promoId uint)

	studentId := uint(1)
	offerId := uint(1)
	promoId := uint(1)
	orderId := uint(1)

	tests := []struct {
		name         string
		body         string
		studentId    uint
		offerId      uint
		promoId      uint
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:      "ok",
			body:      fmt.Sprintf(`{"offerId": "%d"}`, offerId),
			studentId: studentId,
			offerId:   offerId,
			mockBehavior: func(r *mock_service.MockOrders, studentId, offerId, promoId uint) {
				r.EXPECT().Create(context.Background(), studentId, offerId, promoId).Return(orderId, nil)
			},
			statusCode:   200,
			responseBody: fmt.Sprintf(`{"orderId":"%d"}`, orderId),
		},
		{
			name:      "ok w/ promocode",
			body:      fmt.Sprintf(`{"offerId": "%d", "promoId": "%d"}`, offerId, promoId),
			studentId: studentId,
			offerId:   offerId,
			promoId:   promoId,
			mockBehavior: func(r *mock_service.MockOrders, studentId, offerId, promoId uint) {
				r.EXPECT().Create(context.Background(), studentId, offerId, promoId).Return(orderId, nil)
			},
			statusCode:   200,
			responseBody: fmt.Sprintf(`{"orderId":"%d"}`, orderId),
		},
		{
			name:         "offerId missing",
			body:         fmt.Sprintf(`{"offerId": "", "promoId": "%d"}`, promoId),
			mockBehavior: func(r *mock_service.MockOrders, studentId, offerId, promoId uint) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:         "invalid offerId",
			body:         fmt.Sprintf(`{"offerId": "abc", "promoId": "%d"}`, promoId),
			mockBehavior: func(r *mock_service.MockOrders, studentId, offerId, promoId uint) {},
			statusCode:   400,
			responseBody: `{"message":"invalid offer id"}`,
		},
		{
			name:         "invalid promoId",
			body:         fmt.Sprintf(`{"offerId": "%d", "promoId": "abc"}`, offerId),
			mockBehavior: func(r *mock_service.MockOrders, studentId, offerId, promoId uint) {},
			statusCode:   400,
			responseBody: `{"message":"invalid promo id"}`,
		},
		{
			name:      "service error",
			body:      fmt.Sprintf(`{"offerId": "%d", "promoId": "%d"}`, offerId, promoId),
			studentId: studentId,
			offerId:   offerId,
			promoId:   promoId,
			mockBehavior: func(r *mock_service.MockOrders, studentId, offerId, promoId uint) {
				r.EXPECT().Create(context.Background(), studentId, offerId, promoId).Return(orderId, errors.New("failed to create order"))
			},
			statusCode:   500,
			responseBody: `{"message":"failed to create order"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			s := mock_service.NewMockOrders(c)
			tt.mockBehavior(s, tt.studentId, tt.offerId, tt.promoId)

			services := &service.Services{Orders: s}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.POST("/order", func(c *gin.Context) {
				c.Set(studentCtx, fmt.Sprintf("%d", tt.studentId))
			}, handler.studentCreateOrder)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/order",
				strings.NewReader(tt.body))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}

func TestHandler_studentGetModuleOffers(t *testing.T) {
	type mockBehavior func(r *mock_service.MockOffers, schoolId, moduleId uint, offers []domain.Offer)

	schoolId := uint(1)
	moduleId := uint(1)

	packageIds := []uint{1, 2}

	tests := []struct {
		name         string
		moduleId     uint
		schoolId     uint
		offers       []domain.Offer
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			moduleId: moduleId,
			schoolId: schoolId,
			offers: []domain.Offer{
				{
					Name:        "test offer",
					Description: "description",
					SchoolID:    schoolId,
					PackageIDs:  packageIds,
					Benefits: []string{
						"benefit 1",
						"benefit 2",
					},
					Price: domain.Price{
						Value:    6900,
						Currency: "USD",
					},
				},
			},
			mockBehavior: func(r *mock_service.MockOffers, schoolId, moduleId uint, offers []domain.Offer) {
				r.EXPECT().GetByModule(context.Background(), schoolId, moduleId).Return(offers, nil)
			},
			statusCode:   200,
			responseBody: `{"data":[{"id":0,"name":"test offer","description":"description","price":{"value":6900,"currency":"USD"},"benefits":["benefit 1","benefit 2"],"paymentMethod":{"usesProvider":false}}],"count":0}`,
		},
		{
			name:         "invalid module id",
			moduleId:     0,
			schoolId:     schoolId,
			mockBehavior: func(r *mock_service.MockOffers, schoolId, moduleId uint, offers []domain.Offer) {},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
		{
			name:     "service error",
			moduleId: moduleId,
			schoolId: schoolId,
			mockBehavior: func(r *mock_service.MockOffers, schoolId, moduleId uint, offers []domain.Offer) {
				r.EXPECT().GetByModule(context.Background(), schoolId, moduleId).Return(nil, errors.New("failed to get offers"))
			},
			statusCode:   500,
			responseBody: `{"message":"failed to get offers"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			s := mock_service.NewMockOffers(c)

			tt.mockBehavior(s, tt.schoolId, tt.moduleId, tt.offers)

			services := &service.Services{Offers: s}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/modules/:id/offers", func(c *gin.Context) {
				c.Set(schoolCtx, domain.School{
					ID: schoolId,
				})
			}, handler.studentGetModuleOffers)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/modules/%d/offers", tt.moduleId), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}

func TestHandler_studentGetModuleContent(t *testing.T) {
	type mockBehavior func(r *mock_service.MockStudents, schoolId, studentId, moduleId uint, content domain.ModuleContent)

	schoolId := uint(1)
	moduleId := uint(1)
	studentId := uint(1)

	tests := []struct {
		name         string
		moduleId     uint
		schoolId     uint
		studentId    uint
		content      domain.ModuleContent
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:      "ok",
			moduleId:  moduleId,
			schoolId:  schoolId,
			studentId: studentId,
			content: domain.ModuleContent{
				Lessons: []domain.Lesson{
					{
						Name:      "test lesson",
						Position:  0,
						Published: true,
						Content:   "content",
						SchoolID:  schoolId,
					},
				},
			},
			mockBehavior: func(r *mock_service.MockStudents, schoolId, studentId, moduleId uint, content domain.ModuleContent) {
				r.EXPECT().GetModuleContent(context.Background(), schoolId, studentId, moduleId).Return(content, nil)
			},
			statusCode:   200,
			responseBody: fmt.Sprintf(`{"lessons":[{"id":0,"name":"test lesson","position":0,"published":true,"content":"content","schoolId":%d}],"survey":{"title":"","questions":null,"required":false}}`, schoolId),
		},
		{
			name:      "invalid module id",
			moduleId:  0,
			schoolId:  schoolId,
			studentId: studentId,
			mockBehavior: func(r *mock_service.MockStudents, schoolId, studentId, moduleId uint, content domain.ModuleContent) {
			},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
		{
			name:      "module is not available",
			moduleId:  moduleId,
			schoolId:  schoolId,
			studentId: studentId,
			mockBehavior: func(r *mock_service.MockStudents, schoolId, studentId, moduleId uint, content domain.ModuleContent) {
				r.EXPECT().GetModuleContent(context.Background(), schoolId, studentId, moduleId).Return(content, domain.ErrModuleIsNotAvailable)
			},
			statusCode:   403,
			responseBody: fmt.Sprintf(`{"message":"%s"}`, domain.ErrModuleIsNotAvailable.Error()),
		},
		{
			name:      "service error",
			moduleId:  moduleId,
			schoolId:  schoolId,
			studentId: studentId,
			mockBehavior: func(r *mock_service.MockStudents, schoolId, studentId, moduleId uint, content domain.ModuleContent) {
				r.EXPECT().GetModuleContent(context.Background(), schoolId, studentId, moduleId).Return(content, errors.New("failed to get module"))
			},
			statusCode:   500,
			responseBody: `{"message":"failed to get module"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			s := mock_service.NewMockStudents(c)

			tt.mockBehavior(s, tt.schoolId, tt.studentId, tt.moduleId, tt.content)

			services := &service.Services{Students: s}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/modules/:id/content", func(c *gin.Context) {
				c.Set(schoolCtx, domain.School{
					ID: schoolId,
				})
				c.Set(studentCtx, fmt.Sprintf("%d", tt.studentId))
			}, handler.studentGetModuleContent)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/modules/%d/content", tt.moduleId), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}

func TestHandler_studentSetLessonFinished(t *testing.T) {
	type mockBehavior func(r *mock_service.MockStudents, studentId, lessonId uint)

	lessonId := uint(1)
	schoolId := uint(1)
	studentId := uint(1)

	tests := []struct {
		name         string
		lessonId     uint
		studentId    uint
		schoolId     uint
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:      "ok",
			lessonId:  lessonId,
			studentId: studentId,
			schoolId:  schoolId,
			mockBehavior: func(r *mock_service.MockStudents, studentId, lessonId uint) {
				r.EXPECT().SetLessonFinished(context.Background(), studentId, lessonId).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:      "invalid lesson id",
			lessonId:  0,
			schoolId:  schoolId,
			studentId: studentId,
			mockBehavior: func(r *mock_service.MockStudents, studentId, moduleId uint) {
			},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
		{
			name:      "module is not available",
			lessonId:  lessonId,
			schoolId:  schoolId,
			studentId: studentId,
			mockBehavior: func(r *mock_service.MockStudents, studentId, moduleId uint) {
				r.EXPECT().SetLessonFinished(context.Background(), studentId, moduleId).Return(domain.ErrModuleIsNotAvailable)
			},
			statusCode:   403,
			responseBody: fmt.Sprintf(`{"message":"%s"}`, domain.ErrModuleIsNotAvailable.Error()),
		},
		{
			name:      "service error",
			lessonId:  lessonId,
			schoolId:  schoolId,
			studentId: studentId,
			mockBehavior: func(r *mock_service.MockStudents, studentId, moduleId uint) {
				r.EXPECT().SetLessonFinished(context.Background(), studentId, moduleId).Return(errors.New("failed to update student lessons"))
			},
			statusCode:   500,
			responseBody: `{"message":"failed to update student lessons"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			s := mock_service.NewMockStudents(c)

			tt.mockBehavior(s, tt.studentId, tt.lessonId)

			services := &service.Services{Students: s}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.POST("/lessons/:id/finished", func(c *gin.Context) {
				c.Set(schoolCtx, domain.School{
					ID: schoolId,
				})
				c.Set(studentCtx, fmt.Sprintf("%d", tt.studentId))
			}, handler.studentSetLessonFinished)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/lessons/%d/finished", tt.lessonId), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}

func TestHandler_studentSignUp(t *testing.T) {
	type mockBehavior func(r *mock_service.MockStudents, input service.StudentSignUpInput)

	schoolId := uint(1)

	tests := []struct {
		name         string
		requestBody  string
		schoolId     uint
		serviceInput service.StudentSignUpInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:        "ok",
			requestBody: `{"name":"Vasya","email":"test@test.com","password":"qwerty123","registerSource":"test-course"}`,
			schoolId:    schoolId,
			serviceInput: service.StudentSignUpInput{
				Name:         "Vasya",
				Email:        "test@test.com",
				Password:     "qwerty123",
				SchoolID:     schoolId,
				SchoolDomain: "localhost",
			},
			mockBehavior: func(r *mock_service.MockStudents, input service.StudentSignUpInput) {
				r.EXPECT().SignUp(context.Background(), input).Return(nil)
			},
			statusCode: 201,
		},
		{
			name:         "missing name",
			requestBody:  `{"name":"","email":"test@test.com","password":"qwerty123","registerSource":"test-course"}`,
			schoolId:     schoolId,
			mockBehavior: func(r *mock_service.MockStudents, input service.StudentSignUpInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:         "invalid name",
			requestBody:  `{"name":"q","email":"test@test.com","password":"qwerty123","registerSource":"test-course"}`,
			schoolId:     schoolId,
			mockBehavior: func(r *mock_service.MockStudents, input service.StudentSignUpInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:         "missing email",
			requestBody:  `{"name":"Vasya","email":"","password":"qwerty123","registerSource":"test-course"}`,
			schoolId:     schoolId,
			mockBehavior: func(r *mock_service.MockStudents, input service.StudentSignUpInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:         "missing password",
			requestBody:  `{"name":"Vasya","email":"test@test.com","password":"","registerSource":"test-course"}`,
			schoolId:     schoolId,
			mockBehavior: func(r *mock_service.MockStudents, input service.StudentSignUpInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:         "password too short",
			requestBody:  `{"name":"Vasya","email":"test@test.com","password":"qwerty","registerSource":"test-course"}`,
			schoolId:     schoolId,
			mockBehavior: func(r *mock_service.MockStudents, input service.StudentSignUpInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			s := mock_service.NewMockStudents(c)

			tt.mockBehavior(s, tt.serviceInput)

			services := &service.Services{Students: s}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.POST("/sign-up", func(c *gin.Context) {
				c.Set(schoolCtx, domain.School{
					ID: tt.schoolId,
					Settings: domain.Settings{
						Domains: []string{"localhost"},
					},
				})
				c.Set(domainCtx, "localhost")
			}, handler.studentSignUp)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}
