package repository

import "ARIGATOBANK/domain"

type Database interface {
	InsertBalanceOperation(balanceOperationTableName, userBalanceTableName string, bo *domain.BalanceOperation) error
	GetBalanceOperation(balanceOperationTableName string, bo *domain.BalanceOperation) (*domain.BalanceOperation, error)
}
