package accounts

import (
	"asona/internal/model"
	"context"
)

func (i impl) GetByProvider(ctx context.Context, provider, providerAccountID string) (model.Account, error) {
	return model.Account{}, nil
}
