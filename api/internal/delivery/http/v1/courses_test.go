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

func TestHandler_getAllCourses(t *testing.T) {
	tests := []struct {
		name         string
		school       domain.School
		statusCode   int
		responseBody string
	}{
		{
			name: "ok",
			school: domain.School{
				ID: 1,
				Courses: []domain.Course{
					{ID: 1, Name: "Course 1", Published: true},
					{ID: 2, Name: "Course 2", Published: false},
					{ID: 3, Name: "Course 3", Published: true},
				},
			},
			statusCode:   200,
			responseBody: `{"data":[{"id":1,"name":"Course 1","code":"","description":"","color":"","imageUrl":"","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","published":true},{"id":3,"name":"Course 3","code":"","description":"","color":"","imageUrl":"","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","published":true}],"count":0}`,
		},
		{
			name: "no published courses",
			school: domain.School{
				ID: 1,
				Courses: []domain.Course{
					{ID: 1, Name: "Course 1", Published: false},
				},
			},
			statusCode:   200,
			responseBody: `{"data":[],"count":0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			services := &service.Services{}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/courses", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.getAllCourses)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/courses", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_getCourseById(t *testing.T) {
	type mockBehavior func(r *mock_service.MockModules, courseId uint)

	tests := []struct {
		name         string
		courseId     string
		school       domain.School
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			courseId: "1",
			school: domain.School{
				ID: 1,
				Courses: []domain.Course{
					{ID: 1, Name: "Course 1", Published: true},
				},
			},
			mockBehavior: func(r *mock_service.MockModules, courseId uint) {
				r.EXPECT().GetPublishedByCourseId(context.Background(), courseId).Return([]domain.Module{
					{ID: 1, Name: "Module 1", Position: 1, Lessons: []domain.Lesson{
						{ID: 1, Name: "Lesson 1", Position: 1, Published: true},
					}},
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"course":{"id":1,"name":"Course 1","code":"","description":"","color":"","imageUrl":"","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","published":true},"modules":[{"id":1,"name":"Module 1","position":1,"lessons":[{"id":1,"name":"Lesson 1","position":1}]}]}`,
		},
		{
			name:         "empty id",
			courseId:     "",
			school:       domain.School{ID: 1},
			mockBehavior: func(r *mock_service.MockModules, courseId uint) {},
			statusCode:   404,
			responseBody: `404 page not found`,
		},
		{
			name:         "invalid course id",
			courseId:     "abc",
			school:       domain.School{ID: 1},
			mockBehavior: func(r *mock_service.MockModules, courseId uint) {},
			statusCode:   400,
			responseBody: `{"message":"invalid course id"}`,
		},
		{
			name:         "course not found",
			courseId:     "999",
			school:       domain.School{ID: 1},
			mockBehavior: func(r *mock_service.MockModules, courseId uint) {},
			statusCode:   400,
			responseBody: `{"message":"not found"}`,
		},
		{
			name:     "service error",
			courseId: "1",
			school: domain.School{
				ID: 1,
				Courses: []domain.Course{
					{ID: 1, Name: "Course 1", Published: true},
				},
			},
			mockBehavior: func(r *mock_service.MockModules, courseId uint) {
				r.EXPECT().GetPublishedByCourseId(context.Background(), courseId).Return(nil, errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			modules := mock_service.NewMockModules(c)
			tt.mockBehavior(modules, 1)

			services := &service.Services{Modules: modules}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/courses/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.getCourseById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/courses/%s", tt.courseId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_getCourseOffers(t *testing.T) {
	type mockBehavior func(r *mock_service.MockOffers, courseId uint)

	tests := []struct {
		name         string
		courseId     string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			courseId: "1",
			mockBehavior: func(r *mock_service.MockOffers, courseId uint) {
				r.EXPECT().GetByCourse(context.Background(), courseId).Return([]domain.Offer{
					{ID: 1, Name: "Offer 1", Price: domain.Price{Value: 100, Currency: "USD"}},
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"data":[{"id":1,"name":"Offer 1","description":"","benefits":null,"schoolId":0,"packages":null,"price":{"value":100,"currency":"USD"},"paymentMethod":{"usesProvider":false,"provider":""}}],"count":0}`,
		},
		{
			name:         "empty id",
			courseId:     "",
			mockBehavior: func(r *mock_service.MockOffers, courseId uint) {},
			statusCode:   400,
			responseBody: `{"message":"empty id param"}`,
		},
		{
			name:         "invalid id",
			courseId:     "abc",
			mockBehavior: func(r *mock_service.MockOffers, courseId uint) {},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
		{
			name:     "service error",
			courseId: "1",
			mockBehavior: func(r *mock_service.MockOffers, courseId uint) {
				r.EXPECT().GetByCourse(context.Background(), courseId).Return(nil, errors.New("service error"))
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
			tt.mockBehavior(offers, 1)

			services := &service.Services{Offers: offers}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/courses/:id/offers", handler.getCourseOffers)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/courses/%s/offers", tt.courseId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}
