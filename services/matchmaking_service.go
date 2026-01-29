package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/ruangsawala/backend/repositories"
)

type MatchmakingService struct {
	MatchmakingRepo *repositories.MatchmakingRepository
	UserRepo        *repositories.UserInterestRepository
}

type MatchResult struct {
	UserAID int
	UserBID int
	Score   int
	Matched bool
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

func (s *MatchmakingService) FindMatch(ctx context.Context) (*MatchResult, error) {
	candidates, err := s.MatchmakingRepo.PopRandomCandidates(ctx, 2)
	if err != nil {
		return nil, err
	}

	if len(candidates) < 2 {
		if len(candidates) == 1 {
			userID, _ := strconv.Atoi(candidates[0])
			s.MatchmakingRepo.AddToPool(ctx, userID)
		}
		return nil, errors.New("not enough users in pool")
	}

	userAID, _ := strconv.Atoi(candidates[0])
	userBID, _ := strconv.Atoi(candidates[1])

	interestsA, err := s.UserRepo.GetUserInterests(userAID)
	if err != nil {
		s.MatchmakingRepo.AddToPoolMulti(ctx, []int{userAID, userBID})
		return nil, err
	}

	interestsB, err := s.UserRepo.GetUserInterests(userBID)
	if err != nil {
		s.MatchmakingRepo.AddToPoolMulti(ctx, []int{userAID, userBID})
		return nil, err
	}

	interestNamesA := make([]string, len(interestsA))
	for i, interest := range interestsA {
		interestNamesA[i] = interest.Name
	}

	interestNamesB := make([]string, len(interestsB))
	for i, interest := range interestsB {
		interestNamesB[i] = interest.Name
	}

	score := s.CalculateMatchScore(interestNamesA, interestNamesB)

	result := &MatchResult{
		UserAID: userAID,
		UserBID: userBID,
		Score:   score,
		Matched: score >= 2,
	}

	if !result.Matched {
		s.MatchmakingRepo.AddToPoolMulti(ctx, []int{userAID, userBID})
	}

	return result, nil
}
