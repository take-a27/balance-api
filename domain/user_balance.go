package domain

type UserBalance struct {
	Id     string `json:"id"`
	Amount int    `json:"amount"`
}

func NewUserBalance(userId string) *UserBalance {
	return &UserBalance{Id: userId, Amount: 0}
}
