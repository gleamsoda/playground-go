package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"playground/config"
	"playground/domain"
	mock_domain "playground/domain/mock"
	"playground/internal"
	"playground/internal/token"
	mock_token "playground/internal/token/mock"
)

func TestUserUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_domain.NewMockUserRepository(ctrl)
	u := &UserUsecase{
		userRepo: mockUserRepo,
	}
	ctx := context.Background()
	args := domain.CreateUserInputParams{
		Username: "testuser",
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password",
	}
	usr := &domain.User{ID: 1}

	mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(usr, nil)

	got, err := u.Create(ctx, args)

	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, got, "Expected user")
}

func TestUserUsecase_GetByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_domain.NewMockUserRepository(ctrl)
	u := &UserUsecase{
		userRepo: mockUserRepo,
	}
	ctx := context.Background()
	usr := &domain.User{ID: 1}

	mockUserRepo.EXPECT().GetByUsername(ctx, "username").Return(usr, nil)

	got, err := u.GetByUsername(ctx, "username")

	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, got, "Expected user")
}

func TestUserUsecase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_domain.NewMockUserRepository(ctrl)
	mockSessionRepo := mock_domain.NewMockSessionRepository(ctrl)
	mockTokenMaker := mock_token.NewMockMaker(ctrl)
	cfg := config.Config{
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 1 * time.Hour,
	}
	u := &UserUsecase{
		userRepo:    mockUserRepo,
		sessionRepo: mockSessionRepo,
		tokenMaker:  mockTokenMaker,
		cfg:         cfg,
	}
	ctx := context.Background()
	args := domain.LoginUserInputParams{
		Username: "testuser",
		Password: "password",
	}
	pw, _ := internal.HashPassword("password")
	usr := &domain.User{
		ID:             1,
		HashedPassword: pw,
	}
	aPayload := &token.Payload{}
	rPayload := &token.Payload{}

	mockUserRepo.EXPECT().GetByUsername(ctx, args.Username).Return(usr, nil)
	mockTokenMaker.EXPECT().CreateToken(usr.ID, cfg.AccessTokenDuration).Return(gomock.Any().String(), aPayload, nil)
	mockTokenMaker.EXPECT().CreateToken(usr.ID, cfg.RefreshTokenDuration).Return(gomock.Any().String(), rPayload, nil)
	mockSessionRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	got, err := u.Login(ctx, args)

	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, got, "Expected params")
}

func TestUserUsecase_RenewAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mock_domain.NewMockSessionRepository(ctrl)
	mockTokenMaker := mock_token.NewMockMaker(ctrl)

	cfg := config.Config{
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 1 * time.Hour,
	}
	u := &UserUsecase{
		sessionRepo: mockSessionRepo,
		tokenMaker:  mockTokenMaker,
		cfg:         cfg,
	}
	ctx := context.Background()
	args := "refresh_token"
	rPayload := &token.Payload{
		ID:     uuid.New(),
		UserID: 123,
	}

	t.Run("Session is blocked", func(t *testing.T) {
		sess := &domain.Session{
			UserID:       123,
			IsBlocked:    true,
			RefreshToken: "refresh_token",
			ExpiresAt:    time.Now().Add(cfg.RefreshTokenDuration),
		}

		mockTokenMaker.EXPECT().VerifyToken(args).Return(rPayload, nil)
		mockSessionRepo.EXPECT().Get(ctx, rPayload.ID).Return(sess, nil)

		_, err := u.RenewAccessToken(ctx, args)

		assert.EqualError(t, err, "blocked session")
	})

	t.Run("Incorrect session", func(t *testing.T) {
		sess := &domain.Session{
			UserID:       456,
			IsBlocked:    false,
			RefreshToken: "refresh_token",
			ExpiresAt:    time.Now().Add(cfg.RefreshTokenDuration),
		}

		mockTokenMaker.EXPECT().VerifyToken(args).Return(rPayload, nil)
		mockSessionRepo.EXPECT().Get(ctx, rPayload.ID).Return(sess, nil)

		_, err := u.RenewAccessToken(ctx, args)

		assert.EqualError(t, err, "incorrect session user")
	})

	t.Run("Session token mismatched", func(t *testing.T) {
		sess := &domain.Session{
			UserID:       123,
			IsBlocked:    false,
			RefreshToken: "incorrect_token",
			ExpiresAt:    time.Now().Add(cfg.RefreshTokenDuration),
		}

		mockTokenMaker.EXPECT().VerifyToken(args).Return(rPayload, nil)
		mockSessionRepo.EXPECT().Get(ctx, rPayload.ID).Return(sess, nil)

		_, err := u.RenewAccessToken(ctx, args)

		assert.EqualError(t, err, "mismatched session token")
	})

	t.Run("Session expired", func(t *testing.T) {
		sess := &domain.Session{
			UserID:       123,
			IsBlocked:    false,
			RefreshToken: "refresh_token",
			ExpiresAt:    time.Now().Add(-cfg.RefreshTokenDuration),
		}

		mockTokenMaker.EXPECT().VerifyToken(args).Return(rPayload, nil)
		mockSessionRepo.EXPECT().Get(ctx, rPayload.ID).Return(sess, nil)

		_, err := u.RenewAccessToken(ctx, args)

		assert.EqualError(t, err, "expired session")
	})

	t.Run("Success", func(t *testing.T) {
		sess := &domain.Session{
			UserID:       123,
			IsBlocked:    false,
			RefreshToken: "refresh_token",
			ExpiresAt:    time.Now().Add(cfg.RefreshTokenDuration),
		}
		aPayload := &token.Payload{}

		mockTokenMaker.EXPECT().VerifyToken(args).Return(rPayload, nil)
		mockSessionRepo.EXPECT().Get(ctx, rPayload.ID).Return(sess, nil)
		mockTokenMaker.EXPECT().CreateToken(rPayload.UserID, cfg.AccessTokenDuration).Return(gomock.Any().String(), aPayload, nil)

		got, err := u.RenewAccessToken(ctx, args)

		assert.NoError(t, err, "Expected no error")
		assert.NotNil(t, got, "Expected params")
	})
}
