package gateway

import (
	"github.com/gin-gonic/gin"
	pbUser "github.com/semirm-dev/faceit/user/proto"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type CreateAccount struct {
	Nickname string
}

func (api *api) CreateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *CreateAccount
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		account, err := api.rpcClient.AddAccount(c.Request.Context(), &pbUser.AccountRequest{
			Nickname: req.Nickname,
		})
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, account)
	}
}

func (api *api) GetAccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accId int
		accIdParam, ok := c.GetQuery("id")
		if ok {
			id, _ := strconv.Atoi(accIdParam)
			accId = id
		}

		accounts, err := api.rpcClient.GetAccountsByFilter(c.Request.Context(), &pbUser.GetAccountsByFilterRequest{
			Id: int64(accId),
		})
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, accounts)
	}
}
