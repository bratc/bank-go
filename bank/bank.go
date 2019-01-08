package bank

type BankHandler struct {
	service Service
}

func NewHandler(service Service) *BankHandler {
	return &BankHandler{service}
}

func (h *BankHandler) NewClient(balance int) (Client, error) {
	client := Client{0, "", "", ""}
	client, err := h.service.CreateClient(client)
	if err != nil {
		return client, err
	}
	transaction := Transaction{0, 0, client.Id, balance}
	_, err = h.service.CreateTransaction(transaction)
	return client, err
}

func (h *BankHandler) NewTransaction(from_client_id int, to_client_id int, amount int) (Transaction, error) {
	transaction := Transaction{0, from_client_id, to_client_id, amount}
	transaction, err := h.service.CreateTransaction(transaction)
	return transaction, err
}

func (h *BankHandler) CheckBalance(client_id int) (Balance, error) {
	return h.service.GetBalance(client_id)
}
