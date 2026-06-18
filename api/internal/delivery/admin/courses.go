package admin

import (
	"net/http"
	"strconv"

	"ultrathreads/internal/domain"

	"github.com/gin-gonic/gin"
)

type createCourseInput struct {
	Name string `json:"name" binding:"required"`
}

// @Summary Admin Create New Courses
// @Security AdminAuth
// @Tags admins-courses
// @Description admin create new course
// @ModuleID adminCreateCourse
// @Accept  json
// @Produce  json
// @Param input body createCourseInput true "course info"
// @Success 200 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses [post]
func (h *Handler) adminCreateCourse(c *gin.Context) {
	var inp createCourseInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	id, err := h.services.Courses.Create(c.Request.Context(), school.ID, inp.Name)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusCreated, idResponse{id})
}

// @Summary Admin Get All Courses
// @Security AdminAuth
// @Tags admins-courses
// @Description admin get all courses
// @ModuleID adminGetAllCourses
// @Accept  json
// @Produce  json
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses [get]
func (h *Handler) adminGetAllCourses(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	courses, err := h.services.Admins.GetCourses(c.Request.Context(), school.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	response := make([]domain.Course, len(courses))
	if courses != nil {
		response = courses
	}

	c.JSON(http.StatusOK, dataResponse{Data: response})
}

type adminGetCourseByIdResponse struct {
	Course  domain.Course   `json:"course"`
	Modules []domain.Module `json:"modules"`
}

// @Summary Admin Get Course By ID
// @Security AdminAuth
// @Tags admins-courses
// @Description admin get course by id
// @ModuleID adminGetCourseById
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {object} domain.Course
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses/{id} [get]
func (h *Handler) adminGetCourseById(c *gin.Context) {
	id, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	course, err := h.services.Admins.GetCourseById(c.Request.Context(), school.ID, id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	modules, err := h.services.Modules.GetByCourseId(c.Request.Context(), course.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, adminGetCourseByIdResponse{
		Course:  course,
		Modules: modules,
	})
}

type updateCourseInput struct {
	Name        *string `json:"name"`
	ImageURL    *string `json:"imageUrl"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	Published   *bool   `json:"published"`
}

// @Summary Admin Update Course
// @Security AdminAuth
// @Tags admins-courses
// @Description admin update course
// @ModuleID adminUpdateCourse
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Param input body updateCourseInput true "course update info"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses/{id} [put]
func (h *Handler) adminUpdateCourse(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	courseID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id param")

		return
	}

	var inp updateCourseInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Courses.Update(c.Request.Context(), domain.UpdateCourseInput{
		ID:          uint(courseID),
		SchoolID:    school.ID,
		Name:        inp.Name,
		Description: inp.Description,
		ImageURL:    inp.ImageURL,
		Color:       inp.Color,
		Published:   inp.Published,
	}); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

// @Summary Admin Delete Course
// @Security AdminAuth
// @Tags admins-courses
// @Description admin delete course
// @ModuleID adminDeleteCourse
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses/{id} [delete]
func (h *Handler) adminDeleteCourse(c *gin.Context) {
	id, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Courses.Delete(c.Request.Context(), school.ID, id); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}
