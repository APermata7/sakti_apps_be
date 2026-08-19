package handler

import (
    "log"
    "strconv"

    "github.com/gofiber/fiber/v2"

    "sakti_apps_be/internal/domain"
    "sakti_apps_be/internal/usecase"
)

type PresensiHandler struct {
    PresensiUsecase *usecase.PresensiUsecase
}

func NewPresensiHandler(presensiUsecase *usecase.PresensiUsecase) *PresensiHandler {
    return &PresensiHandler{PresensiUsecase: presensiUsecase}
}

func (h *PresensiHandler) CheckIn(c *fiber.Ctx) error {
    log.Println("CheckIn dimulai")

    karyawanID := c.Locals("user_id").(string)
    log.Printf("KaryawanID: %s", karyawanID)

    var req domain.CheckInRequest
    if err := c.BodyParser(&req); err != nil {
        log.Printf("Error parsing request: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Request tidak valid",
        })
    }
    log.Printf("Request: %+v", req)

    if req.SelfieURL == "" || req.Latitude == 0 || req.Longitude == 0 {
        log.Printf("Validasi gagal: SelfieURL=%s, Latitude=%f, Longitude=%f", req.SelfieURL, req.Latitude, req.Longitude)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Selfie URL, latitude, dan longitude wajib diisi",
        })
    }

    log.Println("Memanggil usecase CheckIn")
    resp, err := h.PresensiUsecase.CheckIn(c.Context(), karyawanID, req)
    log.Printf("Hasil usecase - resp: %+v, err: %v", resp, err)

    if err != nil {
        log.Printf("Error usecase: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": err.Error(),
        })
    }

    log.Println("CheckIn berhasil")
    return c.JSON(fiber.Map{
        "success": true,
        "data":    resp,
        "message": "Presensi berhasil",
    })
}

func (h *PresensiHandler) CheckOut(c *fiber.Ctx) error {
    log.Println("CheckOut dimulai")

    karyawanID := c.Locals("user_id").(string)
    log.Printf("KaryawanID: %s", karyawanID)

    var req domain.CheckOutRequest
    if err := c.BodyParser(&req); err != nil {
        log.Printf("Error parsing request: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Request tidak valid",
        })
    }
    log.Printf("Request: %+v", req)

    if req.SelfieURL == "" {
        log.Printf("SelfieURL kosong")
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Selfie URL wajib diisi untuk check-out",
        })
    }

    if req.Latitude == 0 || req.Longitude == 0 {
        log.Printf("Validasi koordinat gagal: Latitude=%f, Longitude=%f", req.Latitude, req.Longitude)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Latitude dan longitude wajib diisi",
        })
    }

    log.Println("Memanggil usecase CheckOut")
    resp, err := h.PresensiUsecase.CheckOut(c.Context(), karyawanID, req)
    log.Printf("Hasil usecase - resp: %+v, err: %v", resp, err)

    if err != nil {
        log.Printf("Error usecase: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": err.Error(),
        })
    }

    message := ""
    if resp.Lembur {
        message = "Lembur " + strconv.FormatFloat(resp.JamLembur, 'f', 1, 64) + " jam tercatat"
    }

    log.Println("CheckOut berhasil")
    return c.JSON(fiber.Map{
        "success": true,
        "data":    resp,
        "message": message,
    })
}

func (h *PresensiHandler) GetToday(c *fiber.Ctx) error {
    log.Println("GetToday dimulai")

    karyawanID := c.Locals("user_id").(string)
    log.Printf("KaryawanID: %s", karyawanID)

    resp, err := h.PresensiUsecase.GetToday(c.Context(), karyawanID)
    log.Printf("Hasil GetToday - resp: %+v, err: %v", resp, err)

    if err != nil {
        log.Printf("Error GetToday: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": err.Error(),
        })
    }

    log.Println("GetToday berhasil")
    return c.JSON(fiber.Map{
        "success": true,
        "data":    resp,
    })
}

func (h *PresensiHandler) GetHistory(c *fiber.Ctx) error {
    log.Println("GetHistory dimulai")

    karyawanID := c.Locals("user_id").(string)
    log.Printf("KaryawanID: %s", karyawanID)

    startDate := c.Query("start_date")
    endDate := c.Query("end_date")
    status := c.Query("status")
    limit, _ := strconv.Atoi(c.Query("limit", "30"))
    page, _ := strconv.Atoi(c.Query("page", "1"))

    if limit <= 0 {
        limit = 30
    }
    if page <= 0 {
        page = 1
    }

    log.Printf("Query params - startDate: %s, endDate: %s, status: %s, limit: %d, page: %d", startDate, endDate, status, limit, page)

    items, total, err := h.PresensiUsecase.GetHistory(c.Context(), karyawanID, startDate, endDate, status, limit, page)
    log.Printf("Hasil GetHistory - items: %d, total: %d, err: %v", len(items), total, err)

    if err != nil {
        log.Printf("Error GetHistory: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": err.Error(),
        })
    }

    totalPages := (total + limit - 1) / limit

    log.Println("GetHistory berhasil")
    return c.JSON(fiber.Map{
        "success": true,
        "data": fiber.Map{
            "items": items,
            "meta": fiber.Map{
                "total":       total,
                "page":        page,
                "limit":       limit,
                "total_pages": totalPages,
            },
        },
    })
}

func (h *PresensiHandler) UpdateAlasanTerlambat(c *fiber.Ctx) error {
    log.Println("UpdateAlasanTerlambat dimulai")

    karyawanID := c.Locals("user_id").(string)
    log.Printf("KaryawanID: %s", karyawanID)

    var req struct {
        AlasanTerlambat string `json:"alasan_terlambat"`
    }
    if err := c.BodyParser(&req); err != nil {
        log.Printf("Error parsing request: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Request tidak valid",
        })
    }
    log.Printf("Request: %+v", req)

    if req.AlasanTerlambat == "" {
        log.Printf("AlasanTerlambat kosong")
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Alasan terlambat wajib diisi",
        })
    }

    err := h.PresensiUsecase.UpdateAlasanTerlambat(c.Context(), karyawanID, req.AlasanTerlambat)
    log.Printf("Hasil UpdateAlasanTerlambat - err: %v", err)

    if err != nil {
        log.Printf("Error UpdateAlasanTerlambat: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": err.Error(),
        })
    }

    log.Println("UpdateAlasanTerlambat berhasil")
    return c.JSON(fiber.Map{
        "success": true,
        "message": "Alasan terlambat berhasil disimpan",
    })
}