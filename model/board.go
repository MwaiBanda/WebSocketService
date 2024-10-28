package model

import (
	"strconv"
)


type Board struct {
	ID          int     `json:"id"`
	Title       string    `json:"title"`
	Prayers    []Prayer  `json:"prayers"`
	Boards	 []BoardMetadata   `json:"boards"`
	Clients    []Client `json:"clients"`
}

type BoardMetadata struct {
	ID          int     `json:"id"`
	Title       string    `json:"title"`
}

type BoardEvent struct {
	ID		  int     `json:"id"`
}
type BoardData struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Prayers    []Prayer  `json:"prayers"`
	Boards	 []BoardMetadata   `json:"boards"`
}

func (b *Board) GetBoardData() BoardData {
	return BoardData{
		ID: b.ID,
		Title: b.Title,
		Prayers: b.Prayers,
		Boards: b.Boards,
	}
}
func (b *Board) SetCanReceiveMessages(client *Client, canReceiveMessages bool) {
	for _, c := range b.Clients {
		if c.ID == client.ID {
		   c.SetCanReceiveMessages(canReceiveMessages)
		}
	}
}

func (b *Board) AddClient(client Client) {
	if !b.HasClient(client) {
	client.SetBoardId(strconv.Itoa(b.ID))
	b.Clients = append(b.Clients, client)
	}
}

func (b *Board) RemoveClient(client Client) {
	filtered := []Client{}
	for _, c := range b.Clients {
		if c.ID != client.ID {
			filtered = append(filtered, c)
		}
	}
	b.Clients = filtered
}

func (b *Board) HasClient(client Client) bool {
	for _, c := range b.Clients {
		if c.DeviceId == client.DeviceId {
			return true
		}
	}
	return false
}

func (b *Board) GetClientIndex(client Client) int {
	for i, c := range b.Clients {
		if c.DeviceId == client.DeviceId {
			return i
		}
	}
	return -1
}