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

	tests := []struct {
		name       string
		tag        uint64
		feature    uint64
		useLastVer bool
		//timer            time.Time
		RepoBannerResult    CacheBanner
		repoError           error
		CacheBannerExpected CacheBanner
		expectedError       error
	}{
		{
			name:       "6 mins ago result",
			tag:        1,
			feature:    1,
			useLastVer: false,
			RepoBannerResult: CacheBanner{
				Timer:  fakeTime,
				Banner: entity.Content{},
			},
			repoError: nil,
			CacheBannerExpected: CacheBanner{
				Timer:  fakeTime,
				Banner: entity.Content{},
			},
			expectedError: nil,
		},
		{
			name:       "6 mins ago result but use last version",
			tag:        5,
			feature:    66,
			useLastVer: true,
			RepoBannerResult: CacheBanner{
				Timer: fakeTime,
				Banner: entity.Content{
					Title: "some_title",
					Text:  "some_text",
					URL:   "some_url",
				},
			},
			repoError: nil,
			CacheBannerExpected: CacheBanner{
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
			name:       "success",
			tag:        5,
			feature:    66,
			useLastVer: true,
			RepoBannerResult: CacheBanner{
				Timer: time.Now(),
				Banner: entity.Content{
					Title: "some_title",
					Text:  "some_text",
					URL:   "some_url",
				},
			},
			repoError: nil,
			CacheBannerExpected: CacheBanner{
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
			mockRepo.EXPECT().GetUserBanner(gomock.Any(), test.tag, test.feature).Return(test.RepoBannerResult.Banner, test.repoError).AnyTimes()
			result, err := useCase.GetUserBanner(context.Background(), test.tag, test.feature, test.useLastVer)
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
