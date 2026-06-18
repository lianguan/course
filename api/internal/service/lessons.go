package service

import (
	"context"
	"errors"

	"ultrathreads/internal/domain"
	"gorm.io/gorm"
)

type LessonsService struct {
	repo        ModulesRepository
	contentRepo LessonContentRepository
}

func NewLessonsService(repo ModulesRepository, contentRepo LessonContentRepository) *LessonsService {
	return &LessonsService{repo: repo, contentRepo: contentRepo}
}

func (s *LessonsService) Create(ctx context.Context, inp domain.AddLessonInput) (uint, error) {
	lesson := domain.NewLesson(inp.Name, inp.Position, inp.SchoolID)

	if err := s.repo.AddLesson(ctx, inp.SchoolID, inp.ModuleID, *lesson); err != nil {
		return 0, err
	}

	return lesson.ID, nil
}

func (s *LessonsService) GetById(ctx context.Context, lessonID uint) (domain.Lesson, error) {
	module, err := s.repo.GetByLesson(ctx, lessonID)
	if err != nil {
		return domain.Lesson{}, err
	}

	var lesson domain.Lesson

	for _, l := range module.Lessons {
		if l.ID == lessonID {
			lesson = l
		}
	}

	content, err := s.contentRepo.GetByLesson(ctx, lessonID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lesson, nil
		}

		return lesson, err
	}

	lesson.Content = content.Content

	return lesson, nil
}

func (s *LessonsService) Update(ctx context.Context, inp domain.UpdateLessonInput) error {
	if inp.Name != "" || inp.Position != nil || inp.Published != nil {
		if err := s.repo.UpdateLesson(ctx, domain.UpdateLessonInput{
			ID:        inp.ID,
			Name:      inp.Name,
			Position:  inp.Position,
			Published: inp.Published,
			SchoolID:  inp.SchoolID,
		}); err != nil {
			return err
		}
	}

	if inp.Content != "" {
		if err := s.contentRepo.Update(ctx, inp.SchoolID, inp.ID, inp.Content); err != nil {
			return err
		}
	}

	return nil
}

func (s *LessonsService) Delete(ctx context.Context, schoolID, id uint) error {
	return s.repo.DeleteLesson(ctx, schoolID, id)
}

func (s *LessonsService) DeleteContent(ctx context.Context, schoolID uint, lessonIDs []uint) error {
	return s.contentRepo.DeleteContent(ctx, schoolID, lessonIDs)
}
