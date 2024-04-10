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

	// Создание фейкового провайдера времени для теста
	fakeTime := time.Now().Add(-6 * time.Minute)

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := NewMockrepository(ctrl)
	useCase := New(mockRepo)

	type banner struct {
		Banner entity.Content
		Timer  time.Time
	}

	tests := []struct {
		name                string
		tag                 uint64
		feature             uint64
		useLastVer          bool
		isUserRequest       bool
		RepoBannerResult    banner
		isActiveBanner      bool
		repoError           error
		CacheBannerExpected banner
		expectedError       error
	}{
		{
			name:          "admin request with no active banner",
			tag:           1,
			feature:       1,
			useLastVer:    true,
			isUserRequest: false,
			RepoBannerResult: banner{
				Timer: time.Now(),
				Banner: entity.Content{
					Title: "some_title",
					Text:  "some_text",
					URL:   "some_url",
				},
			},
			isActiveBanner: false,
			repoError:      nil,
			CacheBannerExpected: banner{
				Timer: time.Now(),
				Banner: entity.Content{
					Title: "some_title",
					Text:  "some_text",
					URL:   "some_url",
				},
			},
			expectedError: nil,
		},
		{
			name:          "6 mins ago result but use last version",
			tag:           1,
			feature:       1,
			useLastVer:    true,
			isUserRequest: false,
			RepoBannerResult: banner{
				Timer: fakeTime,
				Banner: entity.Content{
					Title: "some_title",
					Text:  "some_text",
					URL:   "some_url",
				},
			},
			isActiveBanner: true,
			repoError:      nil,
			CacheBannerExpected: banner{
				Timer: fakeTime,
				Banner: entity.Content{
					Title: "some_title",
					Text:  "some_text",
					URL:   "some_url",
				},
			},
			expectedError: nil,
		},
		{
			name:                "user request with no active banner",
			tag:                 1,
			feature:             1,
			useLastVer:          true,
			isUserRequest:       true,
			RepoBannerResult:    banner{},
			isActiveBanner:      false,
			repoError:           nil,
			CacheBannerExpected: banner{},
			expectedError:       nil,
		},
		{
			name:          "admin request with active banner",
			tag:           1,
			feature:       1,
			useLastVer:    true,
			isUserRequest: false,
			RepoBannerResult: banner{
				Timer: time.Now(),
				Banner: entity.Content{
					Title: "some_title",
					Text:  "some_text",
					URL:   "some_url",
				},
			},
			isActiveBanner: true,
			repoError:      nil,
			CacheBannerExpected: banner{
				Timer: time.Now(),
				Banner: entity.Content{
					Title: "some_title",
					Text:  "some_text",
					URL:   "some_url",
				},
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mockRepo.EXPECT().GetUserBanner(gomock.Any(), test.tag, test.feature).Return(test.RepoBannerResult.Banner, test.isActiveBanner, test.repoError).AnyTimes()
			result, err := useCase.GetUserBanner(context.Background(), test.tag, test.feature, test.useLastVer, test.isUserRequest)
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected error %v, got %v", test.expectedError, err)
			}
			if result.Title != test.CacheBannerExpected.Banner.Title {
				t.Errorf("expected title %s, got %s", test.CacheBannerExpected.Banner.Title, result.Title)
			}
			if result.Text != test.CacheBannerExpected.Banner.Text {
				t.Errorf("expected text %s, got %s", test.CacheBannerExpected.Banner.Text, result.Text)
			}
			if result.URL != test.CacheBannerExpected.Banner.URL {
				t.Errorf("expected url %s, got %s", test.CacheBannerExpected.Banner.URL, result.URL)
			}
		})
	}
}
