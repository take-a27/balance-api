package domain

import (
	"fmt"
	"time"
)

type BalanceOperation struct {
	Id        string
	UserId    string
	Timestamp time.Time

	// a difference between current amount and previous amount
	Difference int

	// if this operation is transaction or not
	IsTransaction bool
}

func NewBalanceOperation(userId string, timestamp string, difference int, isTransaction bool) (*BalanceOperation, error) {
	if userId == "" || timestamp == "" || difference == 0 {
		return nil, fmt.Errorf("required parameter missing. either of user id, timestamp, difference")
	}

	parsedTimestamp, err := time.Parse("2006/01/02 15:04:05", timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %s", err)
	}

	return &BalanceOperation{UserId: userId, Timestamp: parsedTimestamp,
		Difference: difference, IsTransaction: isTransaction}, nil
}
