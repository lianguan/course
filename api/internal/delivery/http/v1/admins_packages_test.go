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

func TestHandler_adminCreatePackage(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPackages, input service.CreatePackageInput)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		courseId     uint
		body         string
		school       domain.School
		input        service.CreatePackageInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			courseId: 1,
			body:     `{"name":"New Package","modules":[1,2,3]}`,
			school:   school,
			input: service.CreatePackageInput{
				SchoolID: school.ID,
				CourseID: 1,
				Name:     "New Package",
				Modules:  []uint{1, 2, 3},
			},
			mockBehavior: func(r *mock_service.MockPackages, input service.CreatePackageInput) {
				r.EXPECT().Create(context.Background(), input).Return(uint(1), nil)
			},
			statusCode:   201,
			responseBody: `{"id":1}`,
		},
		{
			name:         "invalid input",
			courseId:     1,
			body:         `{wrong}`,
			school:       school,
			mockBehavior: func(r *mock_service.MockPackages, input service.CreatePackageInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:     "service error",
			courseId: 1,
			body:     `{"name":"New Package","modules":[1,2,3]}`,
			school:   school,
			input: service.CreatePackageInput{
				SchoolID: school.ID,
				CourseID: 1,
				Name:     "New Package",
				Modules:  []uint{1, 2, 3},
			},
			mockBehavior: func(r *mock_service.MockPackages, input service.CreatePackageInput) {
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

			packages := mock_service.NewMockPackages(c)
			tt.mockBehavior(packages, tt.input)

			services := &service.Services{Packages: packages}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/admins/courses/:id/packages", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminCreatePackage)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/admins/courses/%d/packages", tt.courseId), strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminGetAllPackages(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPackages, courseID uint)

	tests := []struct {
		name         string
		courseId     uint
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			courseId: 1,
			mockBehavior: func(r *mock_service.MockPackages, courseID uint) {
				r.EXPECT().GetByCourse(context.Background(), courseID).Return([]domain.Package{
					{ID: 1, Name: "Package 1"},
					{ID: 2, Name: "Package 2"},
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"data":[{"id":1,"name":"Package 1"},{"id":2,"name":"Package 2"}],"count":0}`,
		},
		{
			name:         "invalid id",
			courseId:     0,
			mockBehavior: func(r *mock_service.MockPackages, courseID uint) {},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
		{
			name:     "service error",
			courseId: 1,
			mockBehavior: func(r *mock_service.MockPackages, courseID uint) {
				r.EXPECT().GetByCourse(context.Background(), courseID).Return(nil, errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			packages := mock_service.NewMockPackages(c)
			tt.mockBehavior(packages, tt.courseId)

			services := &service.Services{Packages: packages}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/admins/courses/:id/packages", handler.adminGetAllPackages)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/admins/courses/%d/packages", tt.courseId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminGetPackageById(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPackages, packageID uint)

	tests := []struct {
		name         string
		packageId    uint
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:      "ok",
			packageId: 1,
			mockBehavior: func(r *mock_service.MockPackages, packageID uint) {
				r.EXPECT().GetById(context.Background(), packageID).Return(domain.Package{
					ID:   packageID,
					Name: "Test Package",
					Modules: []domain.Module{
						{ID: 1, Name: "Module 1"},
					},
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"id":1,"name":"Test Package","modules":[{"id":1,"name":"Module 1"}]}`,
		},
		{
			name:         "invalid id",
			packageId:    0,
			mockBehavior: func(r *mock_service.MockPackages, packageID uint) {},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
		{
			name:      "service error",
			packageId: 1,
			mockBehavior: func(r *mock_service.MockPackages, packageID uint) {
				r.EXPECT().GetById(context.Background(), packageID).Return(domain.Package{}, errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			packages := mock_service.NewMockPackages(c)
			tt.mockBehavior(packages, tt.packageId)

			services := &service.Services{Packages: packages}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/admins/packages/:id", handler.adminGetPackageById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/admins/packages/%d", tt.packageId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminUpdatePackage(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPackages, input service.UpdatePackageInput)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		packageId    uint
		body         string
		school       domain.School
		input        service.UpdatePackageInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:      "ok",
			packageId: 1,
			body:      `{"name":"Updated Package","modules":[1,2]}`,
			school:    school,
			input: service.UpdatePackageInput{
				ID:       1,
				SchoolID: school.ID,
				Name:     "Updated Package",
				Modules:  []uint{1, 2},
			},
			mockBehavior: func(r *mock_service.MockPackages, input service.UpdatePackageInput) {
				r.EXPECT().Update(context.Background(), input).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:         "invalid input",
			packageId:    1,
			body:         `{wrong}`,
			school:       school,
			mockBehavior: func(r *mock_service.MockPackages, input service.UpdatePackageInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "service error",
			packageId: 1,
			body:      `{"name":"Updated Package","modules":[1,2]}`,
			school:    school,
			input: service.UpdatePackageInput{
				ID:       1,
				SchoolID: school.ID,
				Name:     "Updated Package",
				Modules:  []uint{1, 2},
			},
			mockBehavior: func(r *mock_service.MockPackages, input service.UpdatePackageInput) {
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

			packages := mock_service.NewMockPackages(c)
			tt.mockBehavior(packages, tt.input)

			services := &service.Services{Packages: packages}
			handler := Handler{services: services}

			r := gin.New()
			r.PUT("/admins/packages/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminUpdatePackage)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/admins/packages/%d", tt.packageId), strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminDeletePackage(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPackages, schoolID, packageID uint)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		packageId    uint
		school       domain.School
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:      "ok",
			packageId: 1,
			school:    school,
			mockBehavior: func(r *mock_service.MockPackages, schoolID, packageID uint) {
				r.EXPECT().Delete(context.Background(), schoolID, packageID).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:      "service error",
			packageId: 1,
			school:    school,
			mockBehavior: func(r *mock_service.MockPackages, schoolID, packageID uint) {
				r.EXPECT().Delete(context.Background(), schoolID, packageID).Return(errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			packages := mock_service.NewMockPackages(c)
			tt.mockBehavior(packages, tt.school.ID, tt.packageId)

			services := &service.Services{Packages: packages}
			handler := Handler{services: services}

			r := gin.New()
			r.DELETE("/admins/packages/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminDeletePackage)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/admins/packages/%d", tt.packageId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}
