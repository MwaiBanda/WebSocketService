package controller

type Client struct {
	id string
	deviceId string
	// Buffered func of outbound messages.
	send func(int, []byte)
}
