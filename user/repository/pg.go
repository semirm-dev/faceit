package repository

import (
	"context"
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
	Email     string
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

	acc.Firstname = account.FirstName
	acc.Lastname = account.LastName
	acc.Nickname = account.Nickname
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

	acc.Password = password

	return repo.db.Save(&acc).Error
}

func (repo *pgDb) DeleteAccount(ctx context.Context, id string) error {
	var acc *Account
	if err := repo.db.Where("id = ?", id).Find(&acc).Error; err != nil {
		return err
	}

	return repo.db.Delete(acc).Error
}

func (repo *pgDb) GetAccountsByFilter(ctx context.Context, filter *user.Filter) ([]*user.Account, error) {
	var accounts []*Account

	repo.db.Scopes(
		byCountry(repo.db, accounts, filter.Country),
		paginate(repo.db, accounts, &db.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		})).Find(&accounts)

	return entitiesToAccounts(accounts), nil
}

func (repo *pgDb) GetById(ctx context.Context, id string) (*user.Account, error) {
	var acc *Account
	if err := repo.db.Where("id = ?", id).Find(&acc).Error; err != nil {
		return nil, err
	}
	return entityToAccount(acc), nil
}

func (repo *pgDb) GetByEmail(ctx context.Context, email string) (*user.Account, error) {
	var acc *Account
	if err := repo.db.Where("email = ?", email).Find(&acc).Error; err != nil {
		return nil, err
	}
	return entityToAccount(acc), nil
}

func paginate(db *gorm.DB, model interface{}, pagination *db.Pagination) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(model).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func byCountry(db *gorm.DB, model interface{}, country string) func(db *gorm.DB) *gorm.DB {
	db.Model(model)

	return func(db *gorm.DB) *gorm.DB {
		if country != "" {
			db = db.Where("country = ?", country)
		}

		return db
	}
}

func accountToEntity(account *user.Account) *Account {
	return &Account{
		Firstname: account.FirstName,
		Lastname:  account.LastName,
		Nickname:  account.Nickname,
		Password:  account.Password,
		Email:     account.Email,
		Country:   account.Country,
	}
}

func entityToAccount(acc *Account) *user.Account {
	return &user.Account{
		Id:        acc.Id.String(),
		FirstName: acc.Firstname,
		LastName:  acc.Lastname,
		Nickname:  acc.Nickname,
		Password:  acc.Password,
		Email:     acc.Email,
		Country:   acc.Country,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
		DeletedAt: acc.DeletedAt.Time,
	}
}

func entitiesToAccounts(accs []*Account) []*user.Account {
	var accounts []*user.Account

	for _, acc := range accs {
		accounts = append(accounts, entityToAccount(acc))
	}

	return accounts
}

func emailExists(email string, db *gorm.DB) bool {
	var acc *Account
	db.Where("email = ?", email).Find(&acc)
	return acc.Email != ""
}
