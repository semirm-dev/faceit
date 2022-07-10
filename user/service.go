package user

import (
	"context"
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
	Id       int
	Nickname string
}

// AccountRepository communicates to data store with user accounts
type AccountRepository interface {
	AddAccount(ctx context.Context, account *Account) (*Account, error)
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
	account, err := svc.repo.AddAccount(ctx, protoReqToUserAccount(req))
	if err != nil {
		return nil, err
	}

	go func(id int) {
		if pubErr := svc.pub.Publish(event.AccountCreated, id); pubErr != nil {
			logrus.Error(pubErr)
		}
	}(account.Id)

	return userAccountToProto(account), nil
}

// GetAccountsByFilter will get user accounts based on given filters
func (svc *accountService) GetAccountsByFilter(ctx context.Context, req *pbUser.GetAccountsByFilterRequest) (*pbUser.AccountsResponse, error) {
	accounts, err := svc.repo.GetAccountsByFilter(ctx, &Filter{
		Id:       int(req.Id),
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
		Id:       int64(account.Id),
		Nickname: account.Nickname,
	}
}

func protoReqToUserAccount(pbAccount *pbUser.AccountRequest) *Account {
	return &Account{
		Nickname: pbAccount.Nickname,
	}
}
