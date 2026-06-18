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

func TestHandler_adminCreateModule(t *testing.T) {
	type mockBehavior func(r *mock_service.MockModules, input service.CreateModuleInput)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		courseId     uint
		body         string
		school       domain.School
		input        service.CreateModuleInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			courseId: 1,
			body:     `{"name":"New Module","position":1}`,
			school:   school,
			input: service.CreateModuleInput{
				SchoolID: school.ID,
				CourseID: 1,
				Name:     "New Module",
				Position: 1,
			},
			mockBehavior: func(r *mock_service.MockModules, input service.CreateModuleInput) {
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
			mockBehavior: func(r *mock_service.MockModules, input service.CreateModuleInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:     "service error",
			courseId: 1,
			body:     `{"name":"New Module","position":1}`,
			school:   school,
			input: service.CreateModuleInput{
				SchoolID: school.ID,
				CourseID: 1,
				Name:     "New Module",
				Position: 1,
			},
			mockBehavior: func(r *mock_service.MockModules, input service.CreateModuleInput) {
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

			modules := mock_service.NewMockModules(c)
			tt.mockBehavior(modules, tt.input)

			services := &service.Services{Modules: modules}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/admins/courses/:id/modules", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminCreateModule)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/admins/courses/%d/modules", tt.courseId), strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminUpdateModule(t *testing.T) {
	type mockBehavior func(r *mock_service.MockModules, input service.UpdateModuleInput)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		moduleId     uint
		body         string
		school       domain.School
		input        service.UpdateModuleInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			moduleId: 1,
			body:     `{"name":"Updated Module","position":2}`,
			school:   school,
			input: service.UpdateModuleInput{
				ID:       1,
				SchoolID: school.ID,
				Name:     "Updated Module",
				Position: uintPtr(2),
			},
			mockBehavior: func(r *mock_service.MockModules, input service.UpdateModuleInput) {
				r.EXPECT().Update(context.Background(), input).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:         "invalid input",
			moduleId:     1,
			body:         `{wrong}`,
			school:       school,
			mockBehavior: func(r *mock_service.MockModules, input service.UpdateModuleInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:     "service error",
			moduleId: 1,
			body:     `{"name":"Updated Module","position":2}`,
			school:   school,
			input: service.UpdateModuleInput{
				ID:       1,
				SchoolID: school.ID,
				Name:     "Updated Module",
				Position: uintPtr(2),
			},
			mockBehavior: func(r *mock_service.MockModules, input service.UpdateModuleInput) {
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

			modules := mock_service.NewMockModules(c)
			tt.mockBehavior(modules, tt.input)

			services := &service.Services{Modules: modules}
			handler := Handler{services: services}

			r := gin.New()
			r.PUT("/admins/modules/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminUpdateModule)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/admins/modules/%d", tt.moduleId), strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminDeleteModule(t *testing.T) {
	type mockBehavior func(r *mock_service.MockModules, schoolID, moduleID uint)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		moduleId     uint
		school       domain.School
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			moduleId: 1,
			school:   school,
			mockBehavior: func(r *mock_service.MockModules, schoolID, moduleID uint) {
				r.EXPECT().Delete(context.Background(), schoolID, moduleID).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:     "service error",
			moduleId: 1,
			school:   school,
			mockBehavior: func(r *mock_service.MockModules, schoolID, moduleID uint) {
				r.EXPECT().Delete(context.Background(), schoolID, moduleID).Return(errors.New("service error"))
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
			tt.mockBehavior(modules, tt.school.ID, tt.moduleId)

			services := &service.Services{Modules: modules}
			handler := Handler{services: services}

			r := gin.New()
			r.DELETE("/admins/modules/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminDeleteModule)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/admins/modules/%d", tt.moduleId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminCreateLesson(t *testing.T) {
	type mockBehavior func(r *mock_service.MockLessons, input service.AddLessonInput)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		moduleId     uint
		body         string
		school       domain.School
		input        service.AddLessonInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			moduleId: 1,
			body:     `{"name":"New Lesson","position":1}`,
			school:   school,
			input: service.AddLessonInput{
				ModuleID: 1,
				SchoolID: school.ID,
				Name:     "New Lesson",
				Position: 1,
			},
			mockBehavior: func(r *mock_service.MockLessons, input service.AddLessonInput) {
				r.EXPECT().Create(context.Background(), input).Return(uint(1), nil)
			},
			statusCode:   201,
			responseBody: `{"id":1}`,
		},
		{
			name:         "invalid input",
			moduleId:     1,
			body:         `{wrong}`,
			school:       school,
			mockBehavior: func(r *mock_service.MockLessons, input service.AddLessonInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:     "service error",
			moduleId: 1,
			body:     `{"name":"New Lesson","position":1}`,
			school:   school,
			input: service.AddLessonInput{
				ModuleID: 1,
				SchoolID: school.ID,
				Name:     "New Lesson",
				Position: 1,
			},
			mockBehavior: func(r *mock_service.MockLessons, input service.AddLessonInput) {
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

			lessons := mock_service.NewMockLessons(c)
			tt.mockBehavior(lessons, tt.input)

			services := &service.Services{Lessons: lessons}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/admins/modules/:id/lessons", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminCreateLesson)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/admins/modules/%d/lessons", tt.moduleId), strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminGetLessonById(t *testing.T) {
	type mockBehavior func(r *mock_service.MockLessons, lessonID uint)

	tests := []struct {
		name         string
		lessonId     uint
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			lessonId: 1,
			mockBehavior: func(r *mock_service.MockLessons, lessonID uint) {
				r.EXPECT().GetById(context.Background(), lessonID).Return(domain.Lesson{
					ID:       lessonID,
					Name:     "Test Lesson",
					Content:  "Test Content",
					Position: 1,
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"id":1,"name":"Test Lesson","position":1,"published":false,"content":"Test Content","schoolId":0}`,
		},
		{
			name:         "invalid id",
			lessonId:     0,
			mockBehavior: func(r *mock_service.MockLessons, lessonID uint) {},
			statusCode:   400,
			responseBody: `{"message":"invalid id param"}`,
		},
		{
			name:     "service error",
			lessonId: 1,
			mockBehavior: func(r *mock_service.MockLessons, lessonID uint) {
				r.EXPECT().GetById(context.Background(), lessonID).Return(domain.Lesson{}, errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			lessons := mock_service.NewMockLessons(c)
			tt.mockBehavior(lessons, tt.lessonId)

			services := &service.Services{Lessons: lessons}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/admins/lessons/:id", handler.adminGetLessonById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/admins/lessons/%d", tt.lessonId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminUpdateLesson(t *testing.T) {
	type mockBehavior func(r *mock_service.MockLessons, input service.UpdateLessonInput)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		lessonId     uint
		body         string
		school       domain.School
		input        service.UpdateLessonInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			lessonId: 1,
			body:     `{"name":"Updated Lesson","content":"Updated Content"}`,
			school:   school,
			input: service.UpdateLessonInput{
				LessonID: 1,
				SchoolID: school.ID,
				Name:     "Updated Lesson",
				Content:  "Updated Content",
			},
			mockBehavior: func(r *mock_service.MockLessons, input service.UpdateLessonInput) {
				r.EXPECT().Update(context.Background(), input).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:         "invalid input",
			lessonId:     1,
			body:         `{wrong}`,
			school:       school,
			mockBehavior: func(r *mock_service.MockLessons, input service.UpdateLessonInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:     "service error",
			lessonId: 1,
			body:     `{"name":"Updated Lesson","content":"Updated Content"}`,
			school:   school,
			input: service.UpdateLessonInput{
				LessonID: 1,
				SchoolID: school.ID,
				Name:     "Updated Lesson",
				Content:  "Updated Content",
			},
			mockBehavior: func(r *mock_service.MockLessons, input service.UpdateLessonInput) {
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

			lessons := mock_service.NewMockLessons(c)
			tt.mockBehavior(lessons, tt.input)

			services := &service.Services{Lessons: lessons}
			handler := Handler{services: services}

			r := gin.New()
			r.PUT("/admins/lessons/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminUpdateLesson)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/admins/lessons/%d", tt.lessonId), strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_adminDeleteLesson(t *testing.T) {
	type mockBehavior func(r *mock_service.MockLessons, schoolID, lessonID uint)

	school := domain.School{ID: 1}

	tests := []struct {
		name         string
		lessonId     uint
		school       domain.School
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			lessonId: 1,
			school:   school,
			mockBehavior: func(r *mock_service.MockLessons, schoolID, lessonID uint) {
				r.EXPECT().Delete(context.Background(), schoolID, lessonID).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:     "service error",
			lessonId: 1,
			school:   school,
			mockBehavior: func(r *mock_service.MockLessons, schoolID, lessonID uint) {
				r.EXPECT().Delete(context.Background(), schoolID, lessonID).Return(errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			lessons := mock_service.NewMockLessons(c)
			tt.mockBehavior(lessons, tt.school.ID, tt.lessonId)

			services := &service.Services{Lessons: lessons}
			handler := Handler{services: services}

			r := gin.New()
			r.DELETE("/admins/lessons/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminDeleteLesson)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/admins/lessons/%d", tt.lessonId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func uintPtr(i uint) *uint {
	return &i
}
