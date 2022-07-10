package repository

import (
	"context"
	"github.com/semirm-dev/faceit/user"
)

type inmemory struct {
	accounts []*user.Account
}

func NewAccountInmemory() *inmemory {
	return &inmemory{}
}

func (repo *inmemory) AddAccount(ctx context.Context, account *user.Account) (*user.Account, error) {
	account.Id = len(repo.accounts) + 1

	repo.accounts = append(repo.accounts, account)

	return account, nil
}

func (repo *inmemory) GetAccountsByFilter(ctx context.Context, filter *user.Filter) ([]*user.Account, error) {
	var accounts []*user.Account

	if filter.Id > 0 {
		for _, acc := range repo.accounts {
			if acc.Id == filter.Id {
				accounts = append(accounts, acc)
			}
		}
	}

	if filter.Nickname != "" {
		for _, acc := range repo.accounts {
			if acc.Nickname == filter.Nickname {
				accounts = append(accounts, acc)
			}
		}
	}

	if filter.Id == 0 && filter.Nickname == "" {
		accounts = repo.accounts
	}

	return accounts, nil
}
