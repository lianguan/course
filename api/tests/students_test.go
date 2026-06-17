package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"ultrathreads/internal/domain"
	"ultrathreads/pkg/email"
)

const (
	verificationCode = "CODE1234"
)

func (s *APITestSuite) TestStudentSignUp() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))

	r := s.Require()

	name, studentEmail, password := "Test Student", "test@test.com", "qwerty123"
	signUpData := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s"}`, name, studentEmail, password)

	s.mocks.otpGenerator.On("RandomSecret", 8).Return(verificationCode)
	s.mocks.emailSender.On("Send", email.SendEmailInput{
		To:      studentEmail,
		Subject: "Спасибо за регистрацию, Test Student!",
		Body: fmt.Sprintf(`<h1>Спасибо за регистрацию!</h1>
<br>
<p>Чтобы подтвердить свой аккаунт, <a href="https://workshop.ultrathreads.com/verification?code=%s">переходи по ссылке</a>.</p>`, verificationCode),
	}).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/students/sign-up", bytes.NewBuffer([]byte(signUpData)))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Referer", "https://workshop.ultrathreads.com/")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	var student domain.Student
	err := s.db.Where("email = ?", studentEmail).First(&student).Error
	s.NoError(err)

	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	r.Equal(name, student.Name)
	r.Equal(passwordHash, student.Password)
	r.Equal(false, student.Verification.Verified)
	r.Equal(verificationCode, student.Verification.Code)
}

func (s *APITestSuite) TestStudentSignInNotVerified() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentEmail, password := "test2@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	err = s.db.Create(&domain.Student{
		Email:    studentEmail,
		Password: passwordHash,
	}).Error
	s.NoError(err)

	signUpData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, studentEmail, password)
	req, _ := http.NewRequest("POST", "/api/v1/students/sign-in", bytes.NewBuffer([]byte(signUpData)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestStudentVerify() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentEmail, password := "test3@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	err = s.db.Create(&domain.Student{
		Email:        studentEmail,
		Password:     passwordHash,
		Verification: domain.Verification{Code: "CODE4321"},
	}).Error
	s.NoError(err)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/students/verify/%s", "CODE4321"), nil)
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var student domain.Student
	err = s.db.Where("email = ?", studentEmail).First(&student).Error
	s.NoError(err)

	r.Equal(true, student.Verification.Verified)
	r.Equal("", student.Verification.Code)
}

func (s *APITestSuite) TestStudentSignInVerified() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentEmail, password := "test4@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	err = s.db.Create(&domain.Student{
		Email:        studentEmail,
		Password:     passwordHash,
		SchoolID:     school.ID,
		Verification: domain.Verification{Verified: true},
	}).Error
	s.NoError(err)

	signUpData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, studentEmail, password)

	req, _ := http.NewRequest("POST", "/api/v1/students/sign-in", bytes.NewBuffer([]byte(signUpData)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)
}

func (s *APITestSuite) TestStudentGetPaidLessonsWithoutPurchase() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentEmail, password := "test4@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	student := domain.Student{
		Email:        studentEmail,
		Password:     passwordHash,
		SchoolID:     school.ID,
		Verification: domain.Verification{Verified: true},
	}
	err = s.db.Create(&student).Error
	s.NoError(err)

	jwt, err := s.getJwt(student.ID)
	s.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/students/modules/%d/content", modules[1].(domain.Module).ID), nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)
}

func (s *APITestSuite) TestStudentGetModuleOffers() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentEmail, password := "test4@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	student := domain.Student{
		Email:        studentEmail,
		Password:     passwordHash,
		SchoolID:     school.ID,
		Verification: domain.Verification{Verified: true},
	}
	err = s.db.Create(&student).Error
	s.NoError(err)

	jwt, err := s.getJwt(student.ID)
	s.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/students/modules/%d/offers", modules[1].(domain.Module).ID), nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

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

func (s *APITestSuite) TestStudentCreateOrderWithoutPromocode() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentEmail, password := "test4@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	student := domain.Student{
		Email:        studentEmail,
		Password:     passwordHash,
		SchoolID:     school.ID,
		Verification: domain.Verification{Verified: true},
	}
	err = s.db.Create(&student).Error
	s.NoError(err)

	jwt, err := s.getJwt(student.ID)
	s.NoError(err)

	orderData := fmt.Sprintf(`{"offerId":"%d"}`, offers[0].(domain.Offer).ID)

	req, _ := http.NewRequest("POST", "/api/v1/students/orders", bytes.NewBuffer([]byte(orderData)))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var order domain.Order
	err = s.db.Where("student_id = ?", student.ID).First(&order).Error
	s.NoError(err)

	r.Equal(offers[0].(domain.Offer).Price.Value, order.Amount)
	r.Equal(offers[0].(domain.Offer).Price.Currency, order.Currency)
}

func (s *APITestSuite) TestStudentCreateOrderWrongOffer() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentEmail, password := "test4@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	student := domain.Student{
		Email:        studentEmail,
		Password:     passwordHash,
		SchoolID:     school.ID,
		Verification: domain.Verification{Verified: true},
	}
	err = s.db.Create(&student).Error
	s.NoError(err)

	jwt, err := s.getJwt(student.ID)
	s.NoError(err)

	orderData := fmt.Sprintf(`{"offerId":"%d"}`, student.ID)

	req, _ := http.NewRequest("POST", "/api/v1/students/orders", bytes.NewBuffer([]byte(orderData)))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestStudentCreateOrderWithPromocode() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentEmail, password := "test4@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	student := domain.Student{
		Email:        studentEmail,
		Password:     passwordHash,
		SchoolID:     school.ID,
		Verification: domain.Verification{Verified: true},
	}
	err = s.db.Create(&student).Error
	s.NoError(err)

	jwt, err := s.getJwt(student.ID)
	s.NoError(err)

	orderData := fmt.Sprintf(`{"offerId":"%d", "promoId": "%d"}`,
		offers[0].(domain.Offer).ID, promocodes[0].(domain.PromoCode).ID)

	req, _ := http.NewRequest("POST", "/api/v1/students/orders", bytes.NewBuffer([]byte(orderData)))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var order domain.Order
	err = s.db.Where("student_id = ?", student.ID).First(&order).Error
	s.NoError(err)

	offerPrice := offers[0].(domain.Offer).Price.Value
	promocodeDiscount := promocodes[0].(domain.PromoCode).DiscountPercentage
	orderPrice := (offerPrice * uint(100-promocodeDiscount)) / 100

	r.Equal(orderPrice, order.Amount)
	r.Equal(offers[0].(domain.Offer).Price.Currency, order.Currency)
}

func (s *APITestSuite) TestStudentCreateOrderWrongPromo() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentEmail, password := "test4@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	student := domain.Student{
		Email:        studentEmail,
		Password:     passwordHash,
		SchoolID:     school.ID,
		Verification: domain.Verification{Verified: true},
	}
	err = s.db.Create(&student).Error
	s.NoError(err)

	jwt, err := s.getJwt(student.ID)
	s.NoError(err)

	orderData := fmt.Sprintf(`{"offerId":"%d", "promoId": "%d"}`,
		offers[0].(domain.Offer).ID, student.ID)

	req, _ := http.NewRequest("POST", "/api/v1/students/orders", bytes.NewBuffer([]byte(orderData)))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) getJwt(userId uint) (string, error) {
	return s.tokenManager.NewJWT(fmt.Sprintf("%d", userId), time.Hour)
}
