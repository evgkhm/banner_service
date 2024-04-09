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
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := NewMockrepository(ctrl)
	useCase := New(mockRepo)

	tests := []struct {
		name          string
		tag           uint64
		feature       uint64
		useLastVer    bool
		timer         time.Time
		repoResult    entity.Content
		repoError     error
		expected      entity.Content
		expectedError error
	}{
		{
			name:          "6 mins ago result",
			tag:           1,
			feature:       1,
			useLastVer:    false,
			timer:         time.Now().Add(-6 * time.Minute),
			repoResult:    entity.Content{},
			repoError:     nil,
			expected:      entity.Content{},
			expectedError: nil,
		},
		{
			name:       "success",
			tag:        1,
			feature:    1,
			useLastVer: true,
			timer:      time.Now(),
			repoResult: entity.Content{
				Title: "some_title",
				Text:  "some_text",
				URL:   "some_url",
			},
			repoError: nil,
			expected: entity.Content{
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

			mockRepo.EXPECT().GetUserBanner(gomock.Any(), test.tag, test.feature).Return(test.repoResult, test.repoError)
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
