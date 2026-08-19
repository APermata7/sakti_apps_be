package handler

import (
    "log"
    "runtime/debug"

    "sakti_apps_be/internal/repository"
    "sakti_apps_be/internal/utils"

    "github.com/gofiber/fiber/v2"
    "github.com/jackc/pgx/v5/pgxpool"
)

func UploadFile(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "File not found",
        })
    }

    src, err := file.Open()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to open file",
        })
    }
    defer src.Close()

    url, err := utils.UploadFile(src, file.Filename)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "url":     url,
    })
}

func UploadImage(c *fiber.Ctx) error {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("UploadImage panic: %v", r)
            log.Printf("Stack: %s", debug.Stack())
        }
    }()

    log.Println("UploadImage: start")

    role, ok := c.Locals("role").(string)
    if !ok || role == "" {
        log.Println("UploadImage: role not found")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "success": false,
            "message": "Unauthorized",
        })
    }

    userID, ok := c.Locals("user_id").(string)
    if !ok || userID == "" {
        log.Println("UploadImage: user_id not found")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "success": false,
            "message": "Unauthorized",
        })
    }

    karyawanID := c.FormValue("karyawan_id")
    log.Printf("UploadImage: role=%s, userID=%s, karyawanID=%s", role, userID, karyawanID)

    if !(role == "admin" && karyawanID != "") {
        karyawanID = userID
    }
    log.Printf("UploadImage: final karyawanID=%s", karyawanID)

    file, err := c.FormFile("image")
    if err != nil {
        log.Printf("UploadImage: FormFile error: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Image not found",
        })
    }
    log.Printf("UploadImage: file=%s, size=%d", file.Filename, file.Size)

    src, err := file.Open()
    if err != nil {
        log.Printf("UploadImage: Open error: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to open image",
        })
    }
    defer src.Close()

    log.Println("UploadImage: uploading to Cloudinary")
    url, err := utils.UploadImage(src, file.Filename)
    if err != nil {
        log.Printf("UploadImage: Cloudinary error: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Gagal upload: " + err.Error(),
        })
    }
    log.Printf("UploadImage: Cloudinary URL: %s", url)

    dbPool, ok := c.Locals("db").(*pgxpool.Pool)
    if !ok || dbPool == nil {
        log.Printf("UploadImage: Database connection not found")
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Database connection not found",
        })
    }
    log.Println("UploadImage: database connection OK")

    karyawanRepo := repository.NewKaryawanRepo(dbPool)
    err = karyawanRepo.UpdateFotoURL(c.Context(), karyawanID, url)
    if err != nil {
        log.Printf("UploadImage: UpdateFotoURL error: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Gagal simpan ke database: " + err.Error(),
        })
    }
    log.Println("UploadImage: success")

    return c.JSON(fiber.Map{
        "success":     true,
        "url":         url,
        "karyawan_id": karyawanID,
    })
}

func UploadLogoHandler(c *fiber.Ctx) error {
    log.Println("UploadLogoHandler: start")

    file, err := c.FormFile("image")
    if err != nil {
        log.Printf("UploadLogoHandler: FormFile error: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Image not found",
        })
    }
    log.Printf("UploadLogoHandler: file=%s, size=%d", file.Filename, file.Size)

    src, err := file.Open()
    if err != nil {
        log.Printf("UploadLogoHandler: Open error: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to open image",
        })
    }
    defer src.Close()

    log.Println("UploadLogoHandler: uploading to Cloudinary logos folder")
    url, err := utils.UploadLogo(src, file.Filename)
    if err != nil {
        log.Printf("UploadLogoHandler: Cloudinary error: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Gagal upload logo: " + err.Error(),
        })
    }
    log.Printf("UploadLogoHandler: Cloudinary URL: %s", url)

    return c.JSON(fiber.Map{
        "success": true,
        "url":     url,
    })
}