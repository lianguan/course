package service

import (
	"context"

	"ultrathreads/internal/domain"
)

type ModulesService struct {
	repo        ModulesRepository
	contentRepo LessonContentRepository
}

func NewModulesService(repo ModulesRepository, contentRepo LessonContentRepository) *ModulesService {
	return &ModulesService{repo: repo, contentRepo: contentRepo}
}

func (s *ModulesService) GetPublishedByCourseId(ctx context.Context, courseID uint) ([]domain.Module, error) {
	modules, err := s.repo.GetPublishedByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	for i := range modules {
		modules[i].SortLessons()
	}

	return modules, nil
}

func (s *ModulesService) GetByCourseId(ctx context.Context, courseID uint) ([]domain.Module, error) {
	modules, err := s.repo.GetByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	for i := range modules {
		modules[i].SortLessons()
	}

	return modules, nil
}

func (s *ModulesService) GetById(ctx context.Context, moduleID uint) (domain.Module, error) {
	module, err := s.repo.GetPublishedByID(ctx, moduleID)
	if err != nil {
		return module, err
	}

	module.SortLessons()

	return module, nil
}

func (s *ModulesService) GetWithContent(ctx context.Context, moduleID uint) (domain.Module, error) {
	module, err := s.repo.GetByID(ctx, moduleID)
	if err != nil {
		return module, err
	}

	// 使用 domain 方法获取已发布的课时ID
	lessonIDs := module.GetPublishedLessonIDs()
	publishedLessons := module.GetPublishedLessons()

	module.Lessons = publishedLessons

	content, err := s.contentRepo.GetByLessons(ctx, lessonIDs)
	if err != nil {
		return module, err
	}

	for i := range module.Lessons {
		for _, lessonContent := range content {
			if module.Lessons[i].ID == lessonContent.LessonID {
				module.Lessons[i].Content = lessonContent.Content
			}
		}
	}

	module.SortLessons()

	return module, nil
}

func (s *ModulesService) GetByPackages(ctx context.Context, packageIDs []uint) ([]domain.Module, error) {
	modules, err := s.repo.GetByPackages(ctx, packageIDs)
	if err != nil {
		return nil, err
	}

	for i := range modules {
		modules[i].SortLessons()
	}

	return modules, nil
}

func (s *ModulesService) GetByLesson(ctx context.Context, lessonID uint) (domain.Module, error) {
	return s.repo.GetByLesson(ctx, lessonID)
}

func (s *ModulesService) Create(ctx context.Context, inp domain.CreateModuleInput) (uint, error) {
	module := domain.NewModule(inp.Name, inp.Position, inp.CourseID, inp.SchoolID)

	return s.repo.Create(ctx, *module)
}

func (s *ModulesService) Update(ctx context.Context, inp domain.UpdateModuleInput) error {
	return s.repo.Update(ctx, inp)
}

func (s *ModulesService) Delete(ctx context.Context, schoolID, moduleID uint) error {
	module, err := s.repo.GetByID(ctx, moduleID)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, schoolID, moduleID); err != nil {
		return err
	}

	// 使用 domain 方法获取课时ID
	lessonIDs := module.GetLessonIDs()

	return s.contentRepo.DeleteContent(ctx, schoolID, lessonIDs)
}

func (s *ModulesService) DeleteByCourse(ctx context.Context, schoolID, courseID uint) error {
	modules, err := s.repo.GetPublishedByCourseID(ctx, courseID)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteByCourse(ctx, schoolID, courseID); err != nil {
		return err
	}

	// 使用 domain 方法获取所有课时ID
	lessonIDs := make([]uint, 0)
	for i := range modules {
		lessonIDs = append(lessonIDs, modules[i].GetLessonIDs()...)
	}

	return s.contentRepo.DeleteContent(ctx, schoolID, lessonIDs)
}
