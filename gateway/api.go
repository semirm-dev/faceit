package gateway

import (
	"github.com/semirm-dev/faceit/internal/grpc"
	pbUser "github.com/semirm-dev/faceit/user/proto"
)

type api struct {
	rpcClient pbUser.AccountManagementClient
}

func NewApi(accAddr string) *api {
	conn := grpc.CreateClientConnection(accAddr)
	client := pbUser.NewAccountManagementClient(conn)

	return &api{
		rpcClient: client,
	}
}
