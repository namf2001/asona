package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"asona/internal/model"
	"asona/internal/repository"
	"asona/internal/repository/users"

	pkgerrors "github.com/pkg/errors"
)

// RegisterInput represents the input data for new user registration.
type RegisterInput struct {
	Name     string
	Username string
	Email    string
	OTP      string
	Password string
}

func (i impl) RegisterStep1SendOTP(ctx context.Context, email string) error {
	// 1. Check if email already registered
	existingUser, _ := i.repo.User().GetByEmail(ctx, email)
	if existingUser.ID != 0 {
		return pkgerrors.WithStack(ErrUserAlreadyExists)
	}

	// 2. Clear old OTPs
	err := i.repo.VerificationToken().DeleteAllForIdentifier(ctx, email)
	if err != nil {
		return err
	}

	// 3. Generate new OTP
	otp, err := generateOTP(6)
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	// 4. Save to DB
	err = i.repo.VerificationToken().Create(ctx, model.VerificationToken{
		Identifier: email,
		Token:      otp,
		Expires:    time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return err
	}

	// 5. Send Email
	subject := "Verify your registration"
	body := fmt.Sprintf("Your verification code is: <b>%s</b><br><br>It will expire in 15 minutes.", otp)
	return i.mail.SendMail([]string{email}, subject, body)
}

func (i impl) RegisterStep2VerifyOTP(ctx context.Context, email, otp string) error {
	_, err := i.repo.VerificationToken().GetValidToken(ctx, email, otp)
	if err != nil {
		return pkgerrors.WithStack(ErrWrongOTP)
	}
	return nil
}

func (i impl) RegisterStep3Complete(ctx context.Context, input RegisterInput) (model.User, string, error) {
	// 1. Verify OTP again
	_, err := i.repo.VerificationToken().GetValidToken(ctx, input.Email, input.OTP)
	if err != nil {
		return model.User{}, "", pkgerrors.WithStack(ErrWrongOTP)
	}

	// 2. Hash password
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, "", pkgerrors.WithStack(fmt.Errorf("failed to hash password: %w", err))
	}

	// 3. Create user + account
	now := time.Now()
	user := model.User{
		Name:          input.Name,
		Username:      input.Username, // using the username sent from frontend
		Email:         input.Email,
		Password:      string(hashedPasswordBytes),
		EmailVerified: &now,
	}

	var createdUser model.User
	var sessionToken string

	err = i.repo.DoInTx(ctx, func(ctx context.Context, txRepo repository.Registry) error {
		var txErr error
		createdUser, txErr = txRepo.User().Create(ctx, user)
		if txErr != nil {
			if errors.Is(txErr, users.ErrEmailAlreadyExists) {
				return pkgerrors.WithStack(ErrUserAlreadyExists)
			}
			if errors.Is(txErr, users.ErrUsernameAlreadyExists) {
				return pkgerrors.WithStack(ErrUsernameAlreadyExists)
			}
			return pkgerrors.WithStack(fmt.Errorf("failed to register user: %w", txErr))
		}

		_, txErr = txRepo.Account().Create(ctx, model.Account{
			UserID:            createdUser.ID,
			Provider:          "credentials",
			ProviderAccountID: input.Email,
		})
		if txErr != nil {
			return pkgerrors.WithStack(fmt.Errorf("failed to insert credential account: %w", txErr))
		}

		return nil
	}, nil)
	if err != nil {
		return model.User{}, "", err
	}

	// Issue session outside the transaction using the standard helper
	sessionToken, err = i.issueSession(ctx, createdUser, "", "")
	if err != nil {
		return model.User{}, "", err
	}

	// 4. Delete OTP after successful registration
	err = i.repo.VerificationToken().DeleteAllForIdentifier(ctx, input.Email)
	if err != nil {
		// Just log this, no need to fail the whole process if delete fails
	}

	return createdUser, sessionToken, nil
}

func generateOTP(length int) (string, error) {
	const charset = "0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b), nil
}
