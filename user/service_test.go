package user_test

import (
	"context"
	"github.com/semirm-dev/faceit/user"
	pbUser "github.com/semirm-dev/faceit/user/proto"
	"github.com/semirm-dev/faceit/user/repository"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"strings"
	"testing"
)

const (
	bufSize = 1024 * 1024
	addr    = "8001"
)

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	srv := grpc.NewServer()

	pbUser.RegisterAccountManagementServer(srv, user.NewAccountService(
		addr,
		repository.NewAccountInmemory(),
		&mockPublisher{},
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

func TestAccountService_AddAccount(t *testing.T) {
	
}

func TestAccountService_ModifyAccount(t *testing.T) {

}

func TestAccountService_ChangePassword(t *testing.T) {

}

func TestAccountService_DeleteAccount(t *testing.T) {

}

func TestAccountService_GetAccountsByFilter(t *testing.T) {

}
