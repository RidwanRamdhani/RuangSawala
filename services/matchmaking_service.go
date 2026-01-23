package services

import (
	"context"
	"errors"

	"github.com/ruangsawala/backend/repositories"
)

type MatchmakingService struct {
	MatchmakingRepo *repositories.MatchmakingRepository
	UserRepo        *repositories.UserInterestRepository
}

func (s *MatchmakingService) StartSearching(ctx context.Context, userID int) error {
	interests, err := s.UserRepo.GetUserInterests(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if len(interests) == 0 {
		return errors.New("user has no interests")
	}

	return s.MatchmakingRepo.AddToPool(ctx, userID)
}

func (s *MatchmakingService) CancelSearching(ctx context.Context, userID int) error {
	_, err := s.UserRepo.GetUserInterests(userID)
	if err != nil {
		return errors.New("user not found")
	}
	return s.MatchmakingRepo.RemoveFromPool(ctx, userID)
}

func (s *MatchmakingService) CalculateMatchScore(interestsA, interestsB []string) int {
	if len(interestsA) == 0 || len(interestsB) == 0 {
		return 0
	}

	score := 0
	m := make(map[string]bool)
	for _, item := range interestsA {
		m[item] = true
	}
	for _, item := range interestsB {
		if m[item] {
			score++
		}
	}
	return score
}
