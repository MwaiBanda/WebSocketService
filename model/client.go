package model

type Client struct {
	ID   string `json:"id"`
	BoardID   string `json:"boardId"`
	IP   string `json:"ip"`
	User User `json:"user"`
	Send func(int, []byte) `json:"-"`
}

func (c *Client) SetBoardId(boardId string) {
	c.BoardID = boardId
}