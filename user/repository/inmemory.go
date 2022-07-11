package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/semirm-dev/faceit/user"
	"time"
)

type inmemory struct {
	accounts []*user.Account
}

func NewAccountInmemory() *inmemory {
	return &inmemory{}
}

func (repo *inmemory) AddAccount(ctx context.Context, account *user.Account) (*user.Account, error) {
	account.Id = uuid.New().String()
	account.CreatedAt = time.Now().UTC()
	account.UpdatedAt = time.Now().UTC()

	repo.accounts = append(repo.accounts, account)

	return account, nil
}

func (repo *inmemory) ModifyAccount(ctx context.Context, id string, account *user.Account) (*user.Account, error) {
	acc := repo.getById(id)
	if acc != nil {
		acc.Firstname = account.Firstname
		acc.Lastname = account.Lastname
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
	return repo.accounts, nil
}

func (repo *inmemory) GetById(ctx context.Context, id string) (*user.Account, error) {
	return repo.getById(id), nil
}

func (repo *inmemory) GetByEmail(ctx context.Context, email string) (*user.Account, error) {
	return repo.getByEmail(email), nil
}

func (repo *inmemory) getById(id string) *user.Account {
	for _, acc := range repo.accounts {
		if acc.Id == id {
			return acc
		}
	}

	return nil
}

func (repo *inmemory) getByEmail(email string) *user.Account {
	for _, acc := range repo.accounts {
		if acc.Email == email {
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
