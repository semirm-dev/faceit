package gateway

import (
	"github.com/gin-gonic/gin"
	pbUser "github.com/semirm-dev/faceit/user/proto"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CreateAccount struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Country   string `json:"country"`
}

type ModifyAccount struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Country   string `json:"country"`
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
			FirstName: req.Firstname,
			LastName:  req.Lastname,
			Nickname:  req.Nickname,
			Password:  req.Password,
			Email:     req.Email,
			Country:   req.Country,
		})
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, account)
	}
}

func (api *api) ModifyAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		var req *ModifyAccount
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		account, err := api.rpcClient.ModifyAccount(c.Request.Context(), &pbUser.AccountMessage{
			Id:        idParam,
			FirstName: req.Firstname,
			LastName:  req.Lastname,
			Nickname:  req.Nickname,
			Email:     req.Email,
			Country:   req.Country,
		})
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, account)
	}
}

func (api *api) DeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		resp, err := api.rpcClient.DeleteAccount(c.Request.Context(), &pbUser.DeleteAccountRequest{
			Id: idParam,
		})
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func (api *api) GetAccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam, ok := c.GetQuery("id")
		if !ok {
			idParam = ""
		}

		resp, err := api.rpcClient.GetAccountsByFilter(c.Request.Context(), &pbUser.GetAccountsByFilterRequest{
			Id: idParam,
		})
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, resp.Accounts)
	}
}
