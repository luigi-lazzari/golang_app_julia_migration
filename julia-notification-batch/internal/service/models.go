package service

import "fmt"

type NewsExternalItem struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Channel     string `json:"channel"`
}

type NewsItem struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Channel     string `json:"channel"`
}

func (n NewsItem) String() string {
	return fmt.Sprintf("-----------------------------------\n"+
		"ID:          %s\n"+
		"Channel:     %s\n"+
		"Description: %s\n"+
		"-----------------------------------",
		n.ID, n.Channel, n.Description)
}
