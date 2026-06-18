package public

import (
	"errors"
	"net/http"
	"strconv"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/handler"
	"ultrathreads/internal/service/dto"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initCoursesRoutes(api *gin.RouterGroup) {
	courses := api.Group("/courses", h.setSchoolFromRequest)
	{
		courses.GET("", h.getAllCourses)
		courses.GET("/:id", h.getCourseById)
		courses.GET("/:id/offers", h.getCourseOffers)
	}
}

// @Summary  Get All Courses
// @Tags courses
// @Description  get all courses
// @ModuleID getAllCourses
// @Accept  json
// @Produce  json
// @Success 200 {object} handler.DataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /courses [get]
func (h *Handler) getAllCourses(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		handler.NewResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	// Return only published courses
	courses := make([]domain.Course, 0)

	for _, course := range school.Courses {
		if course.Published {
			courses = append(courses, course)
		}
	}

	c.JSON(http.StatusOK, handler.DataResponse{Data: dto.CoursesToResponse(courses)})
}

type getCourseByIdResponse struct {
	Course  dto.CourseResponse   `json:"course"`
	Modules []dto.ModuleResponse `json:"modules"`
}

type module struct {
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
	Position uint     `json:"position"`
	Lessons  []lesson `json:"lessons"`
}

type lesson struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Position uint   `json:"position"`
}

func newGetCourseByIdResponse(course domain.Course, courseModules []domain.Module) getCourseByIdResponse {
	modules := make([]module, len(courseModules))

	for i := range courseModules {
		modules[i].ID = courseModules[i].ID
		modules[i].Name = courseModules[i].Name
		modules[i].Position = courseModules[i].Position
		modules[i].Lessons = toLessons(courseModules[i].Lessons)
	}

	return getCourseByIdResponse{
		Course:  course,
		Modules: modules,
	}
}

func toLessons(lessons []domain.Lesson) []lesson {
	out := make([]lesson, 0)

	for _, l := range lessons {
		if l.Published {
			out = append(out, lesson{
				ID:       l.ID,
				Name:     l.Name,
				Position: l.Position,
			})
		}
	}

	return out
}

// @Summary Get Course By ModuleID
// @Tags courses
// @Description  get course by id
// @ModuleID getCourseById
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {object} domain.Course
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /courses/{id} [get]
func (h *Handler) getCourseById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handler.NewResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		handler.NewResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	course, err := studentGetSchoolCourse(school, id)
	if err != nil {
		handler.NewResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	modules, err := h.services.Modules.GetPublishedByCourseId(c.Request.Context(), course.ID)
	if err != nil {
		handler.NewResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, newGetCourseByIdResponse(course, modules))
}

func studentGetSchoolCourse(school domain.School, courseId string) (domain.Course, error) {
	var searchedCourse domain.Course

	courseIDUint, err := strconv.ParseUint(courseId, 10, 64)
	if err != nil {
		return domain.Course{}, errors.New("invalid course id")
	}

	for _, course := range school.Courses {
		if course.Published && course.ID == uint(courseIDUint) {
			searchedCourse = course
		}
	}

	if searchedCourse.ID == 0 {
		return domain.Course{}, errors.New("not found")
	}

	return searchedCourse, nil
}

// @Summary Get Course Offers
// @Tags courses
// @Description  get course offers
// @ModuleID getCourseOffers
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {object} handler.DataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /courses/{id}/offers [get]
func (h *Handler) getCourseOffers(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handler.NewResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	courseIDUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		handler.NewResponse(c, http.StatusBadRequest, "invalid id param")

		return
	}

	offers, err := h.services.Offers.GetByCourse(c.Request.Context(), uint(courseIDUint))
	if err != nil {
		handler.NewResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, handler.DataResponse{Data: offers})
}
