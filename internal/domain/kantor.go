package domain

import "time"

type Kantor struct {
	ID             string    `json:"id"`
	Nama           string    `json:"nama"`
	Lintang        float64   `json:"lintang"`
	Bujur          float64   `json:"bujur"`
	Radius         int       `json:"radius"`
	Alamat         string    `json:"alamat"`
	DibuatPada     time.Time `json:"dibuat_pada"`
	DiperbaruiPada time.Time `json:"diperbarui_pada"`
}