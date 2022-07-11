package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/semirm-dev/faceit/user"
	"time"
)

type inmemory struct {
	Accounts []*user.Account
}

func NewAccountInmemory() *inmemory {
	return &inmemory{}
}

func (repo *inmemory) AddAccount(ctx context.Context, account *user.Account) (*user.Account, error) {
	account.Id = uuid.New().String()
	account.CreatedAt = time.Now().UTC()
	account.UpdatedAt = time.Now().UTC()

	repo.Accounts = append(repo.Accounts, account)

	return account, nil
}

func (repo *inmemory) ModifyAccount(ctx context.Context, id string, account *user.Account) (*user.Account, error) {
	acc := repo.getById(id)
	if acc != nil {
		acc.FirstName = account.FirstName
		acc.LastName = account.LastName
		acc.Nickname = account.Nickname
		acc.Country = account.Country
		acc.UpdatedAt = time.Now().UTC()
	}

	return acc, nil
}

func (repo *inmemory) ChangePassword(ctx context.Context, id, password string) error {
	acc := repo.getById(id)
	if acc != nil {
		acc.Password = password
		acc.UpdatedAt = time.Now().UTC()
	}

	return nil
}

func (repo *inmemory) DeleteAccount(ctx context.Context, id string) error {
	for i, acc := range repo.Accounts {
		if acc.Id == id {
			copy(repo.Accounts[i:], repo.Accounts[i+1:])
			repo.Accounts[len(repo.Accounts)-1] = nil
			repo.Accounts = repo.Accounts[:len(repo.Accounts)-1]
			break
		}
	}

	return nil
}

func (repo *inmemory) GetAccountsByFilter(ctx context.Context, filter *user.Filter) ([]*user.Account, error) {
	accounts := repo.Accounts

	if filter.Country != "" {
		accounts = repo.getByCountry(filter.Country)
	}

	return accounts, nil
}

func (repo *inmemory) GetById(ctx context.Context, id string) (*user.Account, error) {
	return repo.getById(id), nil
}

func (repo *inmemory) GetByEmail(ctx context.Context, email string) (*user.Account, error) {
	return repo.getByEmail(email), nil
}

func (repo *inmemory) getById(id string) *user.Account {
	for _, acc := range repo.Accounts {
		if acc.Id == id {
			return acc
		}
	}

	return nil
}

func (repo *inmemory) getByEmail(email string) *user.Account {
	for _, acc := range repo.Accounts {
		if acc.Email == email {
			return acc
		}
	}

	return nil
}

func (repo *inmemory) getByCountry(country string) []*user.Account {
	var accounts []*user.Account

	for _, acc := range repo.Accounts {
		if acc.Country == country {
			accounts = append(accounts, acc)
		}
	}

	return accounts
}
