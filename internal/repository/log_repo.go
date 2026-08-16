package repository

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type LogRepo struct {
    DB *pgxpool.Pool
}

func NewLogRepo(db *pgxpool.Pool) *LogRepo {
    return &LogRepo{DB: db}
}

func (r *LogRepo) CreateLog(ctx context.Context, karyawanID, action, detail string) error {
    query := `
        INSERT INTO log_aktivitas (karyawan_id, action, detail, dibuat_pada)
        VALUES ($1, $2, $3, NOW())
    `
    _, err := r.DB.Exec(ctx, query, karyawanID, action, detail)
    return err
}

func (r *LogRepo) GetLogs(ctx context.Context, limit, offset int) ([]map[string]interface{}, int, error) {
    var total int
    countQuery := `SELECT COUNT(*) FROM log_aktivitas`
    err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    query := `
        SELECT id, karyawan_id, action, detail, dibuat_pada
        FROM log_aktivitas
        ORDER BY dibuat_pada DESC
        LIMIT $1 OFFSET $2
    `
    rows, err := r.DB.Query(ctx, query, limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var logs []map[string]interface{}
    for rows.Next() {
        var id, karyawanID, action, detail string
        var dibuatPada time.Time
        err := rows.Scan(&id, &karyawanID, &action, &detail, &dibuatPada)
        if err != nil {
            return nil, 0, err
        }
        logs = append(logs, map[string]interface{}{
            "id":          id,
            "karyawan_id": karyawanID,
            "action":      action,
            "detail":      detail,
            "dibuat_pada": dibuatPada,
        })
    }
    return logs, total, nil
}