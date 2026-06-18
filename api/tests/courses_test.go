package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"ultrathreads/internal/domain"

	"github.com/gin-gonic/gin"
)

type courseResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Color       string `json:"color"`
	ImageURL    string `json:"imageUrl"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Published   bool   `json:"published"`
}

type offerResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	SchoolID    string   `json:"schoolId"`
	PackageIDs  []string `json:"packages"`
	Price       struct {
		Value    uint   `json:"value"`
		Currency string `json:"currency"`
	} `json:"price"`
}

func (s *APITestSuite) TestGetAllCourses() {
	router := gin.New()
	s.handler.InitRoutes(router.Group("/api"))
	r := s.Require()

	req, _ := http.NewRequest("GET", "/api/v1/courses", nil)
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var respCourses struct {
		Data []courseResponse `json:"data"`
	}

	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respCourses)
	s.NoError(err)

	r.Equal(1, len(respCourses.Data))
}

func (s *APITestSuite) TestGetCourseById() {
	router := gin.New()
	s.handler.InitRoutes(router.Group("/api"))
	r := s.Require()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/courses/%d", school.Courses[0].ID), nil)
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	// Get Unpublished Course
	router = gin.New()
	s.handler.InitRoutes(router.Group("/api"))
	r = s.Require()

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/courses/%d", school.Courses[1].ID), nil)
	req.Header.Set("Content-type", "application/json")

	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestGetCourseOffers() {
	router := gin.New()
	s.handler.InitRoutes(router.Group("/api"))
	r := s.Require()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/courses/%d/offers", school.Courses[0].ID), nil)
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var respOffers struct {
		Data []offerResponse `json:"data"`
	}

	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respOffers)
	s.NoError(err)

	r.Equal(1, len(respOffers.Data))
	r.Equal(offers[0].(domain.Offer).Name, respOffers.Data[0].Name)
	r.Equal(offers[0].(domain.Offer).Description, respOffers.Data[0].Description)
	r.Equal(offers[0].(domain.Offer).Price.Value, respOffers.Data[0].Price.Value)
	r.Equal(offers[0].(domain.Offer).Price.Currency, respOffers.Data[0].Price.Currency)
}
