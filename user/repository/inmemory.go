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

func (repo *inmemory) ModifyAccount(ctx context.Context, id int, account *user.Account) (*user.Account, error) {
	acc := repo.getById(id)
	if acc != nil {
		acc.Nickname = account.Nickname
	}

	return acc, nil
}

func (repo *inmemory) DeleteAccount(ctx context.Context, id int) error {
	for i, acc := range repo.accounts {
		if acc.Id == id {
			copy(repo.accounts[i:], repo.accounts[i+1:])
			repo.accounts[len(repo.accounts)-1] = nil
			repo.accounts = repo.accounts[:len(repo.accounts)-1]
			break
		}
	}

	return nil
}

func (repo *inmemory) GetAccountsByFilter(ctx context.Context, filter *user.Filter) ([]*user.Account, error) {
	var accounts []*user.Account

	if filter.Id > 0 {
		acc := repo.getById(filter.Id)
		if acc != nil {
			accounts = append(accounts, acc)
		}
	}

	if filter.Nickname != "" {
		acc := repo.getByNickname(filter.Nickname)
		if acc != nil {
			accounts = append(accounts, acc)
		}
	}

	if filter.Id == 0 && filter.Nickname == "" {
		accounts = repo.accounts
	}

	return accounts, nil
}

func (repo *inmemory) getById(id int) *user.Account {
	for _, acc := range repo.accounts {
		if acc.Id == id {
			return acc
		}
	}

	return nil
}

func (repo *inmemory) getByNickname(nickname string) *user.Account {
	for _, acc := range repo.accounts {
		if acc.Nickname == nickname {
			return acc
		}
	}

	return nil
}
