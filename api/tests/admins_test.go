package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"ultrathreads/internal/domain"
)

func (s *APITestSuite) TestAdminCreateCourse() {
	router := gin.New()
	s.handler.InitRoutes(router.Group("/api"))
	r := s.Require()

	// populate DB data
	var id uint = 100
	var schoolID uint = 1
	adminEmail, password := "testAdmin@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	err = s.db.WithContext(context.Background()).Create(&domain.Admin{
		ID:       id,
		Email:    adminEmail,
		Password: passwordHash,
		SchoolID: schoolID,
	}).Error
	s.NoError(err)

	jwt, err := s.getJwt(id)
	s.NoError(err)

	adminCourseName := "admin course test name"

	name := fmt.Sprintf(`{"name":"%s"}`, adminCourseName)

	req, _ := http.NewRequest("POST", "/api/v1/admins/courses", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)
}

func (s *APITestSuite) TestAdminGetAllCourses() {
	router := gin.New()
	s.handler.InitRoutes(router.Group("/api"))
	r := s.Require()

	var id uint = 101
	var schoolID uint = 1
	adminEmail, password := "testAdmin@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	err = s.db.WithContext(context.Background()).Create(&domain.Admin{
		ID:       id,
		Email:    adminEmail,
		Password: passwordHash,
		SchoolID: schoolID,
	}).Error
	s.NoError(err)

	jwt, err := s.getJwt(id)
	s.NoError(err)

	req, _ := http.NewRequest("GET", "/api/v1/admins/courses", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)
}

func (s *APITestSuite) TestAdminGetCourseByID() {
	router := gin.New()
	s.handler.InitRoutes(router.Group("/api"))
	r := s.Require()

	var id uint = 102
	var schoolID uint = 1
	adminEmail, password := "testAdmin@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	err = s.db.WithContext(context.Background()).Create(&domain.Admin{
		ID:       id,
		Email:    adminEmail,
		Password: passwordHash,
		SchoolID: schoolID,
	}).Error
	s.NoError(err)

	jwt, err := s.getJwt(id)
	s.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/admins/courses/%d", school.Courses[0].ID), nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)
}
