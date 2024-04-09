package usecase

import (
	"banner_service/internal/entity"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestUseCase_GetUserBanner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockrepository(ctrl)
	useCase := New(mockRepo)

	tests := []struct {
		name          string
		tag           uint64
		feature       uint64
		limit         uint64
		offset        uint64
		useLastVer    bool
		timer         time.Time
		repoResult    entity.UserBannerResponse
		repoError     error
		expected      entity.UserBannerResponse
		expectedError error
	}{
		{
			name:          "6 mins ago result",
			tag:           1,
			feature:       1,
			limit:         10,
			offset:        0,
			useLastVer:    false,
			timer:         time.Now().Add(-6 * time.Minute),
			repoResult:    entity.UserBannerResponse{},
			repoError:     nil,
			expected:      entity.UserBannerResponse{},
			expectedError: nil,
		},
		{
			name:       "actual result since 5 min ago",
			tag:        1,
			feature:    1,
			limit:      10,
			offset:     0,
			useLastVer: false,
			timer:      time.Now(),
			repoResult: entity.UserBannerResponse{
				Title: "some_title",
				Text:  "some_text",
				URL:   "some_url",
			},
			repoError: nil,
			expected: entity.UserBannerResponse{
				Title: "some_title",
				Text:  "some_text",
				URL:   "some_url",
			},
			expectedError: nil,
		},
		{
			name:       "success",
			tag:        1,
			feature:    1,
			limit:      10,
			offset:     0,
			useLastVer: true,
			timer:      time.Now(),
			repoResult: entity.UserBannerResponse{
				Title: "some_title",
				Text:  "some_text",
				URL:   "some_url",
			},
			repoError: nil,
			expected: entity.UserBannerResponse{
				Title: "some_title",
				Text:  "some_text",
				URL:   "some_url",
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mockRepo.EXPECT().GetUserBanner(gomock.Any(), test.tag, test.feature, test.useLastVer).Return(test.repoResult, test.repoError)
			result, err := useCase.GetUserBanner(context.Background(), test.tag, test.feature, test.useLastVer)
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected error %v, got %v", test.expectedError, err)
			}
			if result.Title != test.expected.Title {
				t.Errorf("expected title %s, got %s", test.expected.Title, result.Title)
			}
			if result.Text != test.expected.Text {
				t.Errorf("expected text %s, got %s", test.expected.Text, result.Text)
			}
			if result.URL != test.expected.URL {
				t.Errorf("expected url %s, got %s", test.expected.URL, result.URL)
			}
		})
	}
}
