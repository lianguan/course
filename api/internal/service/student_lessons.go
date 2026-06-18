package service

import (
	"context"
)

type StudentLessonsService struct {
	repo StudentLessonsRepository
}

func NewStudentLessonsService(repo StudentLessonsRepository) *StudentLessonsService {
	return &StudentLessonsService{
		repo: repo,
	}
}

func (s *StudentLessonsService) AddFinished(ctx context.Context, studentID, lessonID uint) error {
	return s.repo.AddFinished(ctx, studentID, lessonID)
}

func (s *StudentLessonsService) SetLastOpened(ctx context.Context, studentID, lessonID uint) error {
	return s.repo.SetLastOpened(ctx, studentID, lessonID)
}
