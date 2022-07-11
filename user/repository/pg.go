package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/semirm-dev/faceit/internal/db"
	"github.com/semirm-dev/faceit/user"
	"gorm.io/gorm"
	"math"
	"time"
)

type Account struct {
	Id        uuid.UUID `gorm:"primarykey"`
	Firstname string
	Lastname  string
	Nickname  string
	Password  string
	Email     string `gorm:"uniqueIndex"`
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type pgDb struct {
	db *gorm.DB
}

func NewPgDb(db *gorm.DB) *pgDb {
	db.AutoMigrate(&Account{})

	return &pgDb{
		db: db,
	}
}

func (repo *pgDb) AddAccount(ctx context.Context, account *user.Account) (*user.Account, error) {
	acc := accountToEntity(account)
	acc.Id = uuid.New()

	if err := repo.db.Create(&acc).Error; err != nil {
		return nil, err
	}

	account.Id = acc.Id.String()
	account.CreatedAt = acc.CreatedAt
	account.UpdatedAt = acc.UpdatedAt

	return account, nil
}

func (repo *pgDb) ModifyAccount(ctx context.Context, id string, account *user.Account) (*user.Account, error) {
	var acc *Account
	if err := repo.db.Where("id = ?", id).Find(&acc).Error; err != nil {
		return nil, err
	}
	if acc.Id.String() == "" {
		return nil, errors.New("account not found")
	}

	acc.Firstname = account.Firstname
	acc.Lastname = account.Lastname
	acc.Nickname = account.Nickname
	acc.Email = account.Email
	acc.Country = account.Country

	if err := repo.db.Save(&acc).Error; err != nil {
		return nil, err
	}

	return entityToAccount(acc), nil
}

func (repo *pgDb) ChangePassword(ctx context.Context, id, password string) error {
	var acc *Account
	if err := repo.db.Where("id = ?", id).Find(&acc).Error; err != nil {
		return err
	}
	if acc.Id.String() == "" {
		return errors.New("account not found")
	}

	acc.Password = password
	if err := repo.db.Save(&acc).Error; err != nil {
		return err
	}

	return nil
}

func (repo *pgDb) DeleteAccount(ctx context.Context, id string) error {
	var acc *Account
	if err := repo.db.Where("id = ?", id).Find(&acc).Error; err != nil {
		return err
	}
	if acc.Id.String() == "" {
		return errors.New("account not found")
	}

	if err := repo.db.Delete(acc).Error; err != nil {
		return err
	}

	return nil
}

func (repo *pgDb) GetAccountsByFilter(ctx context.Context, filter *user.Filter) ([]*user.Account, error) {
	var accounts []*user.Account

	if filter.Id != "" {

	}

	if filter.Nickname != "" {

	}

	if filter.Id == "" && filter.Nickname == "" {

	}

	return accounts, nil
}

func (repo *pgDb) GetById(ctx context.Context, id string) (*user.Account, error) {
	return nil, nil
}

func paginate(value interface{}, pagination *db.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func accountToEntity(account *user.Account) *Account {
	return &Account{
		Firstname: account.Firstname,
		Lastname:  account.Lastname,
		Nickname:  account.Nickname,
		Password:  account.Password,
		Email:     account.Email,
		Country:   account.Country,
	}
}

func entityToAccount(acc *Account) *user.Account {
	return &user.Account{
		Id:        acc.Id.String(),
		Firstname: acc.Firstname,
		Lastname:  acc.Lastname,
		Nickname:  acc.Nickname,
		Password:  acc.Password,
		Email:     acc.Email,
		Country:   acc.Country,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
}
