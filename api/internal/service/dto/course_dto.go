package dto

import "ultrathreads/internal/domain"

// ===== Course DTOs =====

type CreateCourseRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCourseRequest struct {
	Name        *string `json:"name"`
	ImageURL    *string `json:"imageUrl"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	Published   *bool   `json:"published"`
}

type CourseResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Color       string `json:"color"`
	ImageURL    string `json:"imageUrl"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
	Published   bool   `json:"published"`
}

func CourseToResponse(c domain.Course) CourseResponse {
	return CourseResponse{
		ID:          c.ID,
		Name:        c.Name,
		Code:        c.Code,
		Description: c.Description,
		Color:       c.Color,
		ImageURL:    c.ImageURL,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Published:   c.Published,
	}
}

func CoursesToResponse(courses []domain.Course) []CourseResponse {
	res := make([]CourseResponse, len(courses))
	for i, c := range courses {
		res[i] = CourseToResponse(c)
	}
	return res
}

// ===== Module DTOs =====

type CreateModuleRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateModuleRequest struct {
	Name      *string `json:"name"`
	Position  *uint   `json:"position"`
	Published *bool   `json:"published"`
}

type ModuleResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Position  uint   `json:"position"`
	Published bool   `json:"published"`
	CourseID  uint   `json:"courseId"`
}

func ModuleToResponse(m domain.Module) ModuleResponse {
	return ModuleResponse{
		ID:        m.ID,
		Name:      m.Name,
		Position:  m.Position,
		Published: m.Published,
		CourseID:  m.CourseID,
	}
}

func ModulesToResponse(modules []domain.Module) []ModuleResponse {
	res := make([]ModuleResponse, len(modules))
	for i, m := range modules {
		res[i] = ModuleToResponse(m)
	}
	return res
}

// ===== Lesson DTOs =====

type CreateLessonRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateLessonRequest struct {
	Name      *string `json:"name"`
	Content   *string `json:"content"`
	Position  *uint   `json:"position"`
	Published *bool   `json:"published"`
}

type LessonResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Position  uint   `json:"position"`
	Published bool   `json:"published"`
}

func LessonToResponse(l domain.Lesson) LessonResponse {
	return LessonResponse{
		ID:        l.ID,
		Name:      l.Name,
		Position:  l.Position,
		Published: l.Published,
	}
}

func LessonsToResponse(lessons []domain.Lesson) []LessonResponse {
	res := make([]LessonResponse, len(lessons))
	for i, l := range lessons {
		res[i] = LessonToResponse(l)
	}
	return res
}

// ===== Package DTOs =====

type CreatePackageRequest struct {
	Name    string `json:"name" binding:"required"`
	Modules []uint `json:"modules"`
}

type UpdatePackageRequest struct {
	Name    *string `json:"name"`
	Modules []uint  `json:"modules"`
}

type PackageResponse struct {
	ID       uint             `json:"id"`
	Name     string           `json:"name"`
	CourseID uint             `json:"courseId"`
	Modules  []ModuleResponse `json:"modules"`
}

func PackageToResponse(p domain.Package) PackageResponse {
	return PackageResponse{
		ID:       p.ID,
		Name:     p.Name,
		CourseID: p.CourseID,
		Modules:  ModulesToResponse(p.Modules),
	}
}

func PackagesToResponse(packages []domain.Package) []PackageResponse {
	res := make([]PackageResponse, len(packages))
	for i, p := range packages {
		res[i] = PackageToResponse(p)
	}
	return res
}

// ===== ModuleContent DTOs =====

type ModuleContentResponse struct {
	Lessons []LessonResponse `json:"lessons"`
}

func ModuleContentToResponse(content domain.ModuleContent) ModuleContentResponse {
	return ModuleContentResponse{
		Lessons: LessonsToResponse(content.Lessons),
	}
}
