package user_test

import (
	"context"
	"github.com/semirm-dev/faceit/user"
	pbUser "github.com/semirm-dev/faceit/user/proto"
	"github.com/semirm-dev/faceit/user/repository"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"strings"
	"testing"
	"time"
)

const (
	bufSize = 1024 * 1024
	addr    = "8001"
)

var (
	lis       *bufconn.Listener
	repo      = repository.NewAccountInmemory()
	publisher = &mockPublisher{events: make(map[string]interface{})}
)

func init() {
	lis = bufconn.Listen(bufSize)
	srv := grpc.NewServer()

	pbUser.RegisterAccountManagementServer(srv, user.NewAccountService(
		addr,
		repo,
		publisher,
		&mockPwdHash{}))

	go func() {
		if err := srv.Serve(lis); err != nil {
			logrus.Fatalf("grpc server failed: %v", err)
		}
	}()
}

type mockPublisher struct {
	events map[string]interface{}
}

func (pub *mockPublisher) Publish(event string, msg interface{}) error {
	pub.events[event] = msg
	return nil
}

type mockPwdHash struct{}

func (pwdHash *mockPwdHash) Hash(plain string) (string, error) {
	return plain + "-hashed", nil
}
func (pwdHash *mockPwdHash) Validate(hashed, plain string) bool {
	h := strings.TrimRight(hashed, "-hashed")

	return h == plain
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func grpcConn(addr string) *grpc.ClientConn {
	ctx := context.Background()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(bufDialer),
	}

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		logrus.Fatal(err)
	}

	return conn
}

func grpcClient() pbUser.AccountManagementClient {
	conn := grpcConn(addr)
	return pbUser.NewAccountManagementClient(conn)
}

func TestAccountService_AddAccount_Valid_Returns_Success(t *testing.T) {
	// given
	repo.Accounts = nil
	publisher.events = make(map[string]interface{})

	rpcClient := grpcClient()
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	accountReq := &pbUser.AccountRequest{
		FirstName: "user 1",
		LastName:  "user 1",
		Nickname:  "user_1",
		Password:  "pwd123",
		Email:     "user1@mail.com",
		Country:   "country1",
	}

	// when
	resp, err := rpcClient.AddAccount(rootCtx, accountReq)

	// then
	assert.Nil(t, err)
	assert.NotEmpty(t, resp.Id)
	assert.Equal(t, "user1@mail.com", resp.Email)
	assert.Equal(t, "pwd123-hashed", resp.Password)
	assert.NotNil(t, resp.CreatedAt)
	assert.Equal(t, 1, len(repo.Accounts))

	publishedMsg := publisher.events["account_created"]
	assert.NotNil(t, publishedMsg)
}

func TestAccountService_AddAccount_ExistingEmail_Returns_Fail(t *testing.T) {
	// given
	repo.Accounts = []*user.Account{
		{
			Id:        "123",
			FirstName: "user 1",
			LastName:  "user 1",
			Nickname:  "user_1",
			Password:  "pwd123",
			Email:     "user1@mail.com",
			Country:   "country1",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
	}
	publisher.events = make(map[string]interface{})

	rpcClient := grpcClient()
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	accountReq := &pbUser.AccountRequest{
		FirstName: "user 1",
		LastName:  "user 1",
		Nickname:  "user_1",
		Password:  "pwd123",
		Email:     "user1@mail.com",
		Country:   "country1",
	}

	// when
	resp, err := rpcClient.AddAccount(rootCtx, accountReq)

	// then
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "email already exists")

	publishedMsg := publisher.events["account_created"]
	assert.Nil(t, publishedMsg)
}

func TestAccountService_ModifyAccount_Valid_Returns_Success(t *testing.T) {
	repo.Accounts = []*user.Account{
		{
			Id:        "123",
			FirstName: "user 1",
			LastName:  "user 1",
			Nickname:  "user_1",
			Password:  "pwd123",
			Email:     "user1@mail.com",
			Country:   "country1",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
	}
	publisher.events = make(map[string]interface{})

	rpcClient := grpcClient()
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	accountReq := &pbUser.AccountMessage{
		Id:        "123",
		FirstName: "user 1 changed",
		LastName:  "user 1",
		Nickname:  "user_1",
		Password:  "pwd123 changed",
		Email:     "user1@mail.com changed",
		Country:   "country1",
	}

	resp, err := rpcClient.ModifyAccount(rootCtx, accountReq)

	assert.Nil(t, err)
	assert.Equal(t, "user 1 changed", resp.FirstName)
	assert.Equal(t, "pwd123", resp.Password)      // shouldnt be changed
	assert.Equal(t, "user1@mail.com", resp.Email) // shouldnt be changed

	publishedMsg := publisher.events["account_modified"]
	assert.NotNil(t, publishedMsg)
}

func TestAccountService_ModifyAccount_NoAccount_Returns_Fail(t *testing.T) {
	repo.Accounts = []*user.Account{
		{
			Id:        "123",
			FirstName: "user 1",
			LastName:  "user 1",
			Nickname:  "user_1",
			Password:  "pwd123",
			Email:     "user1@mail.com",
			Country:   "country1",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
	}
	publisher.events = make(map[string]interface{})

	rpcClient := grpcClient()
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	accountReq := &pbUser.AccountMessage{
		Id:        "12345",
		FirstName: "user 1 changed",
		LastName:  "user 1",
		Nickname:  "user_1",
		Password:  "pwd123 changed",
		Email:     "user1@mail.com changed",
		Country:   "country1",
	}

	resp, err := rpcClient.ModifyAccount(rootCtx, accountReq)

	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "account not found")

	publishedMsg := publisher.events["account_modified"]
	assert.Nil(t, publishedMsg)
}

func TestAccountService_ChangePassword_ValidPassword_Returns_Success(t *testing.T) {
	repo.Accounts = []*user.Account{
		{
			Id:        "123",
			FirstName: "user 1",
			LastName:  "user 1",
			Nickname:  "user_1",
			Password:  "pwd123",
			Email:     "user1@mail.com",
			Country:   "country1",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
	}
	publisher.events = make(map[string]interface{})

	rpcClient := grpcClient()
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	accountReq := &pbUser.ChangePasswordRequest{
		Id:          "123",
		OldPassword: "pwd123",
		NewPassword: "pwd12345",
	}

	resp, err := rpcClient.ChangePassword(rootCtx, accountReq)

	assert.Nil(t, err)
	assert.True(t, resp.Success)

	publishedMsg := publisher.events["account_modified"]
	assert.NotNil(t, publishedMsg)
}

func TestAccountService_ChangePassword_InvalidPassword_Returns_Fail(t *testing.T) {
	repo.Accounts = []*user.Account{
		{
			Id:        "123",
			FirstName: "user 1",
			LastName:  "user 1",
			Nickname:  "user_1",
			Password:  "pwd123",
			Email:     "user1@mail.com",
			Country:   "country1",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
	}
	publisher.events = make(map[string]interface{})

	rpcClient := grpcClient()
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	accountReq := &pbUser.ChangePasswordRequest{
		Id:          "123",
		OldPassword: "invalid",
		NewPassword: "pwd12345",
	}

	resp, err := rpcClient.ChangePassword(rootCtx, accountReq)

	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAccountService_ChangePassword_NoAccount_Returns_Fail(t *testing.T) {
	repo.Accounts = nil
	publisher.events = make(map[string]interface{})

	rpcClient := grpcClient()
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	accountReq := &pbUser.ChangePasswordRequest{
		Id:          "123",
		OldPassword: "pwd123",
		NewPassword: "pwd12345",
	}

	resp, err := rpcClient.ChangePassword(rootCtx, accountReq)

	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "account not found")
}

func TestAccountService_DeleteAccount(t *testing.T) {
	repo.Accounts = []*user.Account{
		{
			Id:        "123",
			FirstName: "user 1",
			LastName:  "user 1",
			Nickname:  "user_1",
			Password:  "pwd123",
			Email:     "user1@mail.com",
			Country:   "country1",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
	}
	publisher.events = make(map[string]interface{})

	rpcClient := grpcClient()
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	accountReq := &pbUser.DeleteAccountRequest{
		Id: "123",
	}

	resp, err := rpcClient.DeleteAccount(rootCtx, accountReq)

	assert.Nil(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 0, len(repo.Accounts))

	publishedMsg := publisher.events["account_deleted"]
	assert.NotNil(t, publishedMsg)
}

func TestAccountService_DeleteAccount_NoAccount_Returns_Fail(t *testing.T) {
	repo.Accounts = nil
	publisher.events = make(map[string]interface{})

	rpcClient := grpcClient()
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	accountReq := &pbUser.DeleteAccountRequest{
		Id: "123",
	}

	resp, err := rpcClient.DeleteAccount(rootCtx, accountReq)

	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "account not found")
}

func TestAccountService_GetAccountsByFilter(t *testing.T) {
	repo.Accounts = []*user.Account{
		{
			Id:        "1",
			FirstName: "user 1",
			LastName:  "user 1",
			Nickname:  "user_1",
			Password:  "pwd123",
			Email:     "user1@mail.com",
			Country:   "country1",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
		{
			Id:        "2",
			FirstName: "user 2",
			LastName:  "user 2",
			Nickname:  "user_2",
			Password:  "pwd123",
			Email:     "user2@mail.com",
			Country:   "country2",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: time.Time{},
		},
	}
	publisher.events = make(map[string]interface{})

	rpcClient := grpcClient()
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	accountReq := &pbUser.GetAccountsByFilterRequest{
		Page:    0,
		Limit:   0,
		Country: "country1",
	}

	resp, err := rpcClient.GetAccountsByFilter(rootCtx, accountReq)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp.Accounts))

	acc := resp.Accounts[0]
	assert.NotNil(t, acc)
	assert.Equal(t, "1", acc.Id)
}
