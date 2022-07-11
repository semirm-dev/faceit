package gateway

import (
	"github.com/gin-gonic/gin"
	pbUser "github.com/semirm-dev/faceit/user/proto"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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

type ChangePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
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
			c.JSON(http.StatusBadRequest, err.Error())
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
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, account)
	}
}

func (api *api) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		var req *ChangePassword
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		resp, err := api.rpcClient.ChangePassword(c.Request.Context(), &pbUser.ChangePasswordRequest{
			Id:          idParam,
			OldPassword: req.OldPassword,
			NewPassword: req.NewPassword,
		})
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, resp)
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
		var page, limit int

		pageQuery, ok := c.GetQuery("page")
		if ok {
			p, err := strconv.Atoi(pageQuery)
			if err == nil {
				page = p
			}
		}

		limitQuery, ok := c.GetQuery("limit")
		if ok {
			l, err := strconv.Atoi(limitQuery)
			if err == nil {
				limit = l
			}
		}

		country, _ := c.GetQuery("country")

		resp, err := api.rpcClient.GetAccountsByFilter(c.Request.Context(), &pbUser.GetAccountsByFilterRequest{
			Page:    int64(page),
			Limit:   int64(limit),
			Country: country,
		})
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, resp.Accounts)
	}
}
