package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/awhdesmond/revolut-user-service/pkg/common"
	"go.uber.org/zap"
)

var (
	ErrUsernameContainsNonLetters = errors.New("username contains non letters")
	ErrInvalidDOB                 = errors.New("invalid date of birth")
	ErrFutureDOBUsed              = errors.New("a date of birth in the future is used")
)

type Service interface {
	Upsert(ctx context.Context, username, dob string) error
	Read(ctx context.Context, username string) (string, error)
}

type service struct {
	store  Store
	logger *zap.Logger
}

func NewService(store Store) Service {
	return &service{store: store}
}

func (svc *service) Upsert(ctx context.Context, username, dob string) error {
	if !common.StringContainsOnlyLetters(username) {
		return ErrUsernameContainsNonLetters
	}

	dobDt, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return ErrInvalidDOB
	}

	today := time.Now()
	if dobDt.After(today) {
		return ErrFutureDOBUsed
	}

	return svc.store.Upsert(ctx, username, dobDt)
}

// Read retrieves a user generates a Hello Birthday message based on
// the user's birthday
func (svc *service) Read(ctx context.Context, username string) (string, error) {
	user, err := svc.store.Read(ctx, username)
	if err != nil {
		return "", err
	}

	numDaysToBirthday := user.CalcDaysToBirthday()
	if numDaysToBirthday == 0 {
		return fmt.Sprintf("Hello, %s! Happy birthday!", username), nil
	}

	return fmt.Sprintf("Hello, %s! Your birthday is in %d days(s)", username, numDaysToBirthday), nil
}
