package service

import (
	"context"
	"time"

	"ultrathreads/internal/domain"
)

type CoursesService struct {
	repo           CoursesRepository
	modulesService Modules
}

func NewCoursesService(repo CoursesRepository, modulesService Modules) *CoursesService {
	return &CoursesService{repo: repo, modulesService: modulesService}
}

func (s *CoursesService) Create(ctx context.Context, schoolID uint, name string) (uint, error) {
	now := time.Now().Unix()
	return s.repo.Create(ctx, schoolID, domain.Course{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	})
}

func (s *CoursesService) Update(ctx context.Context, inp domain.UpdateCourseInput) error {
	return s.repo.Update(ctx, inp)
}

func (s *CoursesService) Delete(ctx context.Context, schoolID, courseID uint) error {
	if err := s.repo.Delete(ctx, schoolID, courseID); err != nil {
		return err
	}

	return s.modulesService.DeleteByCourse(ctx, schoolID, courseID)
}
