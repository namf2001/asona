package accounts

import (
	"asona/internal/model"
	"context"
)

func (i impl) Create(ctx context.Context, account model.Account) (model.Account, error) {
	return account, nil
}
