package main

type Model struct {
	Id          int64  `column:"id"`
	Name        string `column:"name"`
	WalletMoney int64  `column:"wallet_money"`
}

func (m *Model) Table() string {
	return "suser_copy"
}

func main() {
}
