package repository

import (
    "context"
    "errors"
    "strconv"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "sakti_apps_be/internal/domain"
)

type LiburRepo struct {
    DB *pgxpool.Pool
}

func NewLiburRepo(db *pgxpool.Pool) *LiburRepo {
    return &LiburRepo{DB: db}
}

func (r *LiburRepo) Create(ctx context.Context, libur *domain.Libur) error {
    query := `
        INSERT INTO libur (tanggal, nama, jenis, sumber, aktif, dibuat_pada, diperbarui_pada)
        VALUES ($1, $2, $3, $4, true, NOW(), NOW())
        RETURNING id
    `
    err := r.DB.QueryRow(ctx, query, libur.Tanggal, libur.Nama, libur.Jenis, libur.Sumber).Scan(&libur.ID)
    return err
}

func (r *LiburRepo) GetByID(ctx context.Context, id string) (*domain.Libur, error) {
    query := `
        SELECT id, tanggal, nama, jenis, aktif, sumber, dibuat_pada, diperbarui_pada
        FROM libur
        WHERE id = $1
    `
    var l domain.Libur
    err := r.DB.QueryRow(ctx, query, id).Scan(
        &l.ID, &l.Tanggal, &l.Nama, &l.Jenis, &l.Aktif,
        &l.Sumber, &l.DibuatPada, &l.DiperbaruiPada,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &l, nil
}

func (r *LiburRepo) GetByTanggal(ctx context.Context, tanggal string) (*domain.Libur, error) {
    query := `
        SELECT id, tanggal, nama, jenis, aktif, sumber, dibuat_pada, diperbarui_pada
        FROM libur
        WHERE tanggal = $1
    `
    t, err := time.Parse("2006-01-02", tanggal)
    if err != nil {
        return nil, err
    }

    var l domain.Libur
    err = r.DB.QueryRow(ctx, query, t).Scan(
        &l.ID, &l.Tanggal, &l.Nama, &l.Jenis, &l.Aktif,
        &l.Sumber, &l.DibuatPada, &l.DiperbaruiPada,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &l, nil
}

func (r *LiburRepo) GetAll(ctx context.Context, tahun int, jenis, sumber string, aktif *bool, limit, offset int) ([]domain.Libur, int, error) {
    var items []domain.Libur
    var total int

    baseQuery := `FROM libur WHERE 1=1`
    args := []interface{}{}
    argIdx := 1

    if tahun > 0 {
        baseQuery += ` AND EXTRACT(YEAR FROM tanggal) = $` + strconv.Itoa(argIdx)
        args = append(args, tahun)
        argIdx++
    }
    if jenis != "" {
        baseQuery += ` AND jenis = $` + strconv.Itoa(argIdx)
        args = append(args, jenis)
        argIdx++
    }
    if sumber != "" {
        baseQuery += ` AND sumber = $` + strconv.Itoa(argIdx)
        args = append(args, sumber)
        argIdx++
    }
    if aktif != nil {
        baseQuery += ` AND aktif = $` + strconv.Itoa(argIdx)
        args = append(args, *aktif)
        argIdx++
    }

    countQuery := `SELECT COUNT(*) ` + baseQuery
    err := r.DB.QueryRow(ctx, countQuery, args...).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    dataQuery := `
        SELECT id, tanggal, nama, jenis, aktif, sumber, dibuat_pada, diperbarui_pada
    ` + baseQuery + ` ORDER BY tanggal ASC LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

    finalArgs := append(args, limit, offset)
    rows, err := r.DB.Query(ctx, dataQuery, finalArgs...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    for rows.Next() {
        var l domain.Libur
        err := rows.Scan(
            &l.ID, &l.Tanggal, &l.Nama, &l.Jenis, &l.Aktif,
            &l.Sumber, &l.DibuatPada, &l.DiperbaruiPada,
        )
        if err != nil {
            return nil, 0, err
        }
        items = append(items, l)
    }

    return items, total, nil
}

func (r *LiburRepo) GetByTahun(ctx context.Context, tahun int) ([]domain.Libur, error) {
    query := `
        SELECT id, tanggal, nama, jenis, aktif, sumber, dibuat_pada, diperbarui_pada
        FROM libur
        WHERE EXTRACT(YEAR FROM tanggal) = $1
        ORDER BY tanggal ASC
    `
    rows, err := r.DB.Query(ctx, query, tahun)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []domain.Libur
    for rows.Next() {
        var l domain.Libur
        err := rows.Scan(
            &l.ID, &l.Tanggal, &l.Nama, &l.Jenis, &l.Aktif,
            &l.Sumber, &l.DibuatPada, &l.DiperbaruiPada,
        )
        if err != nil {
            return nil, err
        }
        items = append(items, l)
    }
    return items, nil
}

func (r *LiburRepo) GetByBulan(ctx context.Context, bulan string) ([]domain.Libur, error) {
    query := `
        SELECT id, tanggal, nama, jenis, aktif, sumber, dibuat_pada, diperbarui_pada
        FROM libur
        WHERE DATE_TRUNC('month', tanggal) = DATE_TRUNC('month', $1::date)
        ORDER BY tanggal ASC
    `
    rows, err := r.DB.Query(ctx, query, bulan+"-01")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []domain.Libur
    for rows.Next() {
        var l domain.Libur
        err := rows.Scan(
            &l.ID, &l.Tanggal, &l.Nama, &l.Jenis, &l.Aktif,
            &l.Sumber, &l.DibuatPada, &l.DiperbaruiPada,
        )
        if err != nil {
            return nil, err
        }
        items = append(items, l)
    }
    return items, nil
}

func (r *LiburRepo) Update(ctx context.Context, libur *domain.Libur) error {
    query := `
        UPDATE libur 
        SET nama = $2, jenis = $3, aktif = $4, diperbarui_pada = NOW()
        WHERE id = $1
    `
    _, err := r.DB.Exec(ctx, query, libur.ID, libur.Nama, libur.Jenis, libur.Aktif)
    return err
}

func (r *LiburRepo) Delete(ctx context.Context, id string) error {
    query := `DELETE FROM libur WHERE id = $1`
    _, err := r.DB.Exec(ctx, query, id)
    return err
}

func (r *LiburRepo) Toggle(ctx context.Context, id string) (bool, error) {
    var aktif bool
    query := `
        UPDATE libur
        SET aktif = NOT aktif, diperbarui_pada = NOW()
        WHERE id = $1
        RETURNING aktif
    `
    err := r.DB.QueryRow(ctx, query, id).Scan(&aktif)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return false, errors.New("hari libur tidak ditemukan")
        }
        return false, err
    }
    return aktif, nil
}