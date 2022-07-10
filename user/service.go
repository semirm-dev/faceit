package user

import (
	"context"
	"errors"
	"github.com/gobackpack/crypto"
	"github.com/semirm-dev/faceit/event"
	"github.com/semirm-dev/faceit/internal/grpc"
	pbUser "github.com/semirm-dev/faceit/user/proto"
	"github.com/sirupsen/logrus"
	grpcLib "google.golang.org/grpc"
)

const serviceName = "account management service"

// accountService will expose account management service via grpc
type accountService struct {
	pbUser.UnimplementedAccountManagementServer
	addr string
	repo AccountRepository
	pub  AccountPublisher
}

// Filter when querying data store for user accounts
type Filter struct {
	Id       string
	Nickname string
}

// AccountRepository communicates to data store with user accounts
type AccountRepository interface {
	AddAccount(ctx context.Context, account *Account) (*Account, error)
	ModifyAccount(ctx context.Context, id string, account *Account) (*Account, error)
	ChangePassword(ctx context.Context, id, newPassword string) error
	DeleteAccount(ctx context.Context, id string) error
	GetById(ctx context.Context, id string) (*Account, error)
	GetAccountsByFilter(ctx context.Context, filter *Filter) ([]*Account, error)
}

// AccountPublisher will publish event that corresponds to an account action
type AccountPublisher interface {
	Publish(event string, msg interface{}) error
}

func NewAccountService(addr string, repo AccountRepository, pub AccountPublisher) *accountService {
	return &accountService{
		addr: addr,
		repo: repo,
		pub:  pub,
	}
}

func (svc *accountService) ListenForConnections(ctx context.Context) {
	grpc.ListenForConnections(ctx, svc, svc.addr, serviceName)
}

func (svc *accountService) RegisterGrpcServer(server *grpcLib.Server) {
	pbUser.RegisterAccountManagementServer(server, svc)
}

// AddAccount will add new user account
func (svc *accountService) AddAccount(ctx context.Context, req *pbUser.AccountRequest) (*pbUser.AccountMessage, error) {
	argon2 := crypto.NewArgon2()
	hashed, err := argon2.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashed

	account, err := svc.repo.AddAccount(ctx, protoReqToUserAccount(req))
	if err != nil {
		return nil, err
	}

	go func(id string) {
		if pubErr := svc.pub.Publish(event.AccountCreated, id); pubErr != nil {
			logrus.Error(pubErr)
		}
	}(account.Id)

	return userAccountToProto(account), nil
}

func (svc *accountService) ModifyAccount(ctx context.Context, req *pbUser.AccountMessage) (*pbUser.AccountMessage, error) {
	account, err := svc.repo.ModifyAccount(ctx, req.Id, protoToUserAccount(req))
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}

	go func(id string) {
		if pubErr := svc.pub.Publish(event.AccountModified, id); pubErr != nil {
			logrus.Error(pubErr)
		}
	}(account.Id)

	return userAccountToProto(account), nil
}

func (svc *accountService) ChangePassword(ctx context.Context, req *pbUser.ChangePasswordRequest) (*pbUser.ChangePasswordResponse, error) {
	account, err := svc.repo.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}

	argon2 := crypto.NewArgon2()
	if !argon2.Validate(account.Password, req.OldPassword) {
		return nil, errors.New("invalid credentials")
	}

	hashed, err := argon2.Hash(req.NewPassword)
	if err != nil {
		return nil, err
	}

	if err := svc.repo.ChangePassword(ctx, req.Id, hashed); err != nil {
		return nil, err
	}

	go func(id string) {
		if pubErr := svc.pub.Publish(event.AccountModified, id); pubErr != nil {
			logrus.Error(pubErr)
		}
	}(account.Id)

	return &pbUser.ChangePasswordResponse{
		Success: true,
	}, nil
}

func (svc *accountService) DeleteAccount(ctx context.Context, req *pbUser.DeleteAccountRequest) (*pbUser.DeleteAccountResponse, error) {
	err := svc.repo.DeleteAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	go func(id string) {
		if pubErr := svc.pub.Publish(event.AccountDeleted, id); pubErr != nil {
			logrus.Error(pubErr)
		}
	}(req.Id)

	return &pbUser.DeleteAccountResponse{
		Success: true,
	}, nil
}

// GetAccountsByFilter will get user accounts based on given filters
func (svc *accountService) GetAccountsByFilter(ctx context.Context, req *pbUser.GetAccountsByFilterRequest) (*pbUser.AccountsResponse, error) {
	accounts, err := svc.repo.GetAccountsByFilter(ctx, &Filter{
		Id:       req.Id,
		Nickname: req.Nickname,
	})
	if err != nil {
		return nil, err
	}

	return &pbUser.AccountsResponse{
		Accounts: userAccountsToProto(accounts),
	}, nil
}

func userAccountsToProto(accounts []*Account) []*pbUser.AccountMessage {
	var userAccounts []*pbUser.AccountMessage

	for _, acc := range accounts {
		userAccounts = append(userAccounts, userAccountToProto(acc))
	}

	return userAccounts
}

func userAccountToProto(account *Account) *pbUser.AccountMessage {
	return &pbUser.AccountMessage{
		Id:        account.Id,
		FirstName: account.Firstname,
		LastName:  account.Lastname,
		Nickname:  account.Nickname,
		Password:  account.Password,
		Email:     account.Email,
		Country:   account.Country,
		CreatedAt: account.CreatedAt.String(),
		UpdatedAt: account.UpdatedAt.String(),
	}
}

func protoReqToUserAccount(pbAccount *pbUser.AccountRequest) *Account {
	return &Account{
		Firstname: pbAccount.FirstName,
		Lastname:  pbAccount.LastName,
		Nickname:  pbAccount.Nickname,
		Password:  pbAccount.Password,
		Email:     pbAccount.Email,
		Country:   pbAccount.Country,
	}
}

func protoToUserAccount(pbAccount *pbUser.AccountMessage) *Account {
	return &Account{
		Id:        pbAccount.Id,
		Firstname: pbAccount.FirstName,
		Lastname:  pbAccount.LastName,
		Nickname:  pbAccount.Nickname,
		Password:  pbAccount.Password,
		Email:     pbAccount.Email,
		Country:   pbAccount.Country,
	}
}
