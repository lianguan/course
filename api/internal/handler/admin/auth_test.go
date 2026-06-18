package admin

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

func TestHandler_adminSignIn(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAdmins, input domain.SchoolSignInInput)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		body         string
		school       domain.School
		input        domain.SchoolSignInInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:   "ok",
			body:   `{"email":"admin@example.com","password":"password123"}`,
			school: school,
			input: domain.SchoolSignInInput{
				Email:    "admin@example.com",
				Password: "password123",
				SchoolID: school.ID,
			},
			mockBehavior: func(r *mock_service.MockAdmins, input domain.SchoolSignInInput) {
				r.EXPECT().SignIn(context.Background(), input).Return(domain.Tokens{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"accessToken":"access-token","refreshToken":"refresh-token"}`,
		},
		{
			name:         "invalid input",
			body:         `{wrong}`,
			school:       school,
			mockBehavior: func(r *mock_service.MockAdmins, input domain.SchoolSignInInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:   "user not found",
			body:   `{"email":"admin@example.com","password":"password123"}`,
			school: school,
			input: domain.SchoolSignInInput{
				Email:    "admin@example.com",
				Password: "password123",
				SchoolID: school.ID,
			},
			mockBehavior: func(r *mock_service.MockAdmins, input domain.SchoolSignInInput) {
				r.EXPECT().SignIn(context.Background(), input).Return(domain.Tokens{}, domain.ErrUserNotFound)
			},
			statusCode:   401,
			responseBody: `{"message":"user doesn't exists"}`,
		},
		{
			name:   "service error",
			body:   `{"email":"admin@example.com","password":"password123"}`,
			school: school,
			input: domain.SchoolSignInInput{
				Email:    "admin@example.com",
				Password: "password123",
				SchoolID: school.ID,
			},
			mockBehavior: func(r *mock_service.MockAdmins, input domain.SchoolSignInInput) {
				r.EXPECT().SignIn(context.Background(), input).Return(domain.Tokens{}, errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			admins := mock_service.NewMockAdmins(c)
			tt.mockBehavior(admins, tt.input)

			services := &service.Services{Admins: admins}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/admins/sign-in", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminSignIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/admins/sign-in", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminRefresh(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAdmins, schoolID uint, token string)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		body         string
		school       domain.School
		token        string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:   "ok",
			body:   `{"token":"refresh-token"}`,
			school: school,
			token:  "refresh-token",
			mockBehavior: func(r *mock_service.MockAdmins, schoolID uint, token string) {
				r.EXPECT().RefreshTokens(context.Background(), schoolID, token).Return(domain.Tokens{
					AccessToken:  "new-access-token",
					RefreshToken: "new-refresh-token",
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"accessToken":"new-access-token","refreshToken":"new-refresh-token"}`,
		},
		{
			name:         "invalid input",
			body:         `{wrong}`,
			school:       school,
			mockBehavior: func(r *mock_service.MockAdmins, schoolID uint, token string) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:   "service error",
			body:   `{"token":"refresh-token"}`,
			school: school,
			token:  "refresh-token",
			mockBehavior: func(r *mock_service.MockAdmins, schoolID uint, token string) {
				r.EXPECT().RefreshTokens(context.Background(), schoolID, token).Return(domain.Tokens{}, errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			admins := mock_service.NewMockAdmins(c)
			tt.mockBehavior(admins, tt.school.ID, tt.token)

			services := &service.Services{Admins: admins}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/admins/auth/refresh", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminRefresh)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/admins/auth/refresh", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminCreateCourse(t *testing.T) {
	type mockBehavior func(r *mock_service.MockCourses, schoolID uint, name string)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		body         string
		school       domain.School
		name_        string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:   "ok",
			body:   `{"name":"New Course"}`,
			school: school,
			name_:  "New Course",
			mockBehavior: func(r *mock_service.MockCourses, schoolID uint, name string) {
				r.EXPECT().Create(context.Background(), schoolID, name).Return(uint(1), nil)
			},
			statusCode:   201,
			responseBody: `{"id":1}`,
		},
		{
			name:         "invalid input",
			body:         `{wrong}`,
			school:       school,
			mockBehavior: func(r *mock_service.MockCourses, schoolID uint, name string) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:   "service error",
			body:   `{"name":"New Course"}`,
			school: school,
			name_:  "New Course",
			mockBehavior: func(r *mock_service.MockCourses, schoolID uint, name string) {
				r.EXPECT().Create(context.Background(), schoolID, name).Return(uint(0), errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			courses := mock_service.NewMockCourses(c)
			tt.mockBehavior(courses, tt.school.ID, tt.name_)

			services := &service.Services{Courses: courses}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/admins/courses", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminCreateCourse)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/admins/courses", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminGetAllCourses(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAdmins, schoolID uint)

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
			mockBehavior: func(r *mock_service.MockAdmins, schoolID uint) {
				r.EXPECT().GetCourses(context.Background(), schoolID).Return([]domain.Course{
					{ID: 1, Name: "Course 1"},
					{ID: 2, Name: "Course 2"},
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"data":[{"id":1,"name":"Course 1","code":"","description":"","color":"","imageUrl":"","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","published":false},{"id":2,"name":"Course 2","code":"","description":"","color":"","imageUrl":"","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","published":false}],"count":0}`,
		},
		{
			name:   "service error",
			school: school,
			mockBehavior: func(r *mock_service.MockAdmins, schoolID uint) {
				r.EXPECT().GetCourses(context.Background(), schoolID).Return(nil, errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			admins := mock_service.NewMockAdmins(c)
			tt.mockBehavior(admins, tt.school.ID)

			services := &service.Services{Admins: admins}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/admins/courses", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminGetAllCourses)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/admins/courses", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminGetCourseById(t *testing.T) {
	type mockBehavior func(a *mock_service.MockAdmins, m *mock_service.MockModules, schoolID, courseID uint)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		courseId     uint
		school       domain.School
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			courseId: 1,
			school:   school,
			mockBehavior: func(a *mock_service.MockAdmins, m *mock_service.MockModules, schoolID, courseID uint) {
				a.EXPECT().GetCourseById(context.Background(), schoolID, courseID).Return(domain.Course{
					ID:   courseID,
					Name: "Course 1",
				}, nil)
				m.EXPECT().GetByCourseId(context.Background(), courseID).Return([]domain.Module{
					{ID: 1, Name: "Module 1"},
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"course":{"id":1,"name":"Course 1","code":"","description":"","color":"","imageUrl":"","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","published":false},"modules":[{"id":1,"name":"Module 1","position":0,"published":false,"courseId":0,"schoolId":0,"survey":{"title":"","questions":null,"required":false}}]}`,
		},
		{
			name:         "invalid id",
			courseId:     0,
			school:       school,
			mockBehavior: func(a *mock_service.MockAdmins, m *mock_service.MockModules, schoolID, courseID uint) {},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
		{
			name:     "service error",
			courseId: 1,
			school:   school,
			mockBehavior: func(a *mock_service.MockAdmins, m *mock_service.MockModules, schoolID, courseID uint) {
				a.EXPECT().GetCourseById(context.Background(), schoolID, courseID).Return(domain.Course{}, errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			admins := mock_service.NewMockAdmins(c)
			modules := mock_service.NewMockModules(c)
			tt.mockBehavior(admins, modules, tt.school.ID, tt.courseId)

			services := &service.Services{Admins: admins, Modules: modules}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/admins/courses/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminGetCourseById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/admins/courses/%d", tt.courseId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}
