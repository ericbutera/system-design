package models

import "time"

type URL struct {
	Long      string     `json:"long" binding:"required,url,min=1,max=2048"`
	Slug      string     `json:"slug" binding:"max=8"`
	ExpiresAt *time.Time `json:"expires_at"`
	UserID    uint
}

type URLStats struct {
	Slug      string `json:"slug"`
	ViewCount uint   `json:"view_count"`
}
