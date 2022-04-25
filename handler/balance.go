package handler

import (
	"ARIGATOBANK/domain"
	"ARIGATOBANK/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type BalanceRequest struct {
	UserId        string `json:"user_id"`
	Timestamp     string `json:"timestamp"`
	Difference    int    `json:"difference"`
	IsTransaction bool   `json:"is_transaction"`
	CurrencyCode  string `json:"currency_code"`
}

func (br *BalanceRequest) ToBalanceOperation() (*domain.BalanceOperation, error) {
	bo, err := domain.NewBalanceOperation(br.UserId, br.Timestamp, br.Difference, br.IsTransaction)
	if err != nil {
		return nil, err
	}

	id := uuid.NewString()
	bo.Id = id
	return bo, nil
}

type BalanceResponse struct {
	OperationId string `json:"operation_id"`
	UserId      string `json:"user_id"`
	Timestamp   string `json:"timestamp"`
	Result      bool   `json:"result"`
	Message     string `json:"message"`
}

func NewBalanceResponse(operationId string, userId string, timestamp string, result bool, message string) *BalanceResponse {
	return &BalanceResponse{OperationId: operationId, UserId: userId, Timestamp: timestamp, Result: result, Message: message}
}

type Balance struct {
	DB repository.Database
}

const (
	balanceOperationTableName = "balance_operation"
	userBalanceTableName      = "user_balance"
)

func (b *Balance) Update(c *gin.Context) {
	balanceRequestList := make([]BalanceRequest, 0)
	err := c.Bind(balanceRequestList)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	balanceResponseList := make([]*BalanceResponse, 0)
	for _, balanceRequest := range balanceRequestList {
		bo, err := balanceRequest.ToBalanceOperation()
		if err != nil {
			balanceResponse := NewBalanceResponse(bo.Id, bo.UserId, fmt.Sprintf("%s", bo.Timestamp), false, err.Error())
			balanceResponseList = append(balanceResponseList, balanceResponse)
		}

		// check if balance operation is existing
		checkBO, err := b.DB.GetBalanceOperation(balanceOperationTableName, bo)
		if err != nil {
			balanceResponse := NewBalanceResponse(bo.Id, bo.UserId, fmt.Sprintf("%s", bo.Timestamp), false, err.Error())
			balanceResponseList = append(balanceResponseList, balanceResponse)
		}
		if checkBO != nil {
			balanceResponse := NewBalanceResponse(bo.Id, bo.UserId, fmt.Sprintf("%s", bo.Timestamp), true, "duplicated operation")
			balanceResponseList = append(balanceResponseList, balanceResponse)
		}

		// insert balance operation
		err = b.DB.InsertBalanceOperation(balanceOperationTableName, userBalanceTableName, bo)
		if err != nil {
			balanceResponse := NewBalanceResponse(bo.Id, bo.UserId, fmt.Sprintf("%s", bo.Timestamp), true, err.Error())
			balanceResponseList = append(balanceResponseList, balanceResponse)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"operation_result": balanceResponseList,
	})
}
