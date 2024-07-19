package credential

// JsTicketHandle js ticket 获取
type JsTicketHandle interface {
	// GetTicket 获取 ticket
	GetTicket(accessToken string) (ticket string, err error)
}
