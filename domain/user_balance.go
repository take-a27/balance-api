package domain

type UserBalance struct {
	Id     string `json:"id"`
	Amount int    `json:"amount"`
}

func NewUserBalance(id string) *UserBalance {
	return &UserBalance{Id: id}
}
