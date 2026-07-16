package domain

import "time"

type TandaTangan struct {
    ID             string    `json:"id"`
    KaryawanID     string    `json:"karyawan_id"`
    URLTandaTangan string    `json:"url_tanda_tangan"`
    HashTandaTangan *string  `json:"hash_tanda_tangan"`
    DiunggahPada   time.Time `json:"diunggah_pada"`
    DiperbaruiPada time.Time `json:"diperbarui_pada"`
}

type CreateTTDRequest struct {
    URLTandaTangan string `json:"url_tanda_tangan"`
}