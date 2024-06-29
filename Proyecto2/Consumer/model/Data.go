package model

import "time"

type Data struct {
    Texto      string    `json:"texto"`
    Pais       string    `json:"pais"`
    CreatedAt  time.Time `json:"created_at"`
}
