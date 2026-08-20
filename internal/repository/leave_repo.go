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

type LeaveRepo struct {
    DB *pgxpool.Pool
}

func NewLeaveRepo(db *pgxpool.Pool) *LeaveRepo {
    return &LeaveRepo{DB: db}
}

func (r *LeaveRepo) Create(ctx context.Context, leave *domain.PengajuanCuti) error {
    query := `
        INSERT INTO pengajuan_cuti (
            karyawan_id, sub_tipe, tanggal_mulai, tanggal_selesai,
            total_hari, alasan, status, back_date, mengurangi_cuti, 
            langsung_approve, langsung_final,
            judul_dokumen, difinalisasi_oleh, tanggal_difinalisasi,
            dibuat_pada, diperbarui_pada
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW())
        RETURNING id
    `

    err := r.DB.QueryRow(ctx, query,
        leave.KaryawanID,
        leave.SubTipe,
        leave.TanggalMulai,
        leave.TanggalSelesai,
        leave.TotalHari,
        leave.Alasan,
        leave.Status,
        leave.BackDate,
        leave.MengurangiCuti,
        leave.LangsungApprove,
        leave.LangsungFinal,
        leave.JudulDokumen,
        leave.DifinalisasiOleh,
        leave.TanggalDifinalisasi,
    ).Scan(&leave.ID)

    return err
}

func (r *LeaveRepo) CreateBalance(ctx context.Context, balance *domain.SisaCuti) error {
    query := `
        INSERT INTO sisa_cuti (karyawan_id, tahun, jumlah_cuti, telah_dilaksanakan, akan_dilaksanakan, sisa_cuti)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
    err := r.DB.QueryRow(ctx, query,
        balance.KaryawanID,
        balance.Tahun,
        balance.JumlahCuti,
        balance.TelahDilaksanakan,
        balance.AkanDilaksanakan,
        balance.SisaCuti,
    ).Scan(&balance.ID)
    return err
}

func (r *LeaveRepo) GetByID(ctx context.Context, id string) (*domain.PengajuanCuti, error) {
    query := `
        SELECT id, karyawan_id, sub_tipe, tanggal_mulai, tanggal_selesai,
               total_hari, alasan, status, back_date, mengurangi_cuti, 
               langsung_approve, langsung_final,
               judul_dokumen, disetujui_oleh, tanggal_disetujui, difinalisasi_oleh,
               tanggal_difinalisasi, url_pdf, alasan_batal, tanggal_dibatalkan,
               alasan_ditolak, tanggal_ditolak,
               dibuat_pada, diperbarui_pada
        FROM pengajuan_cuti
        WHERE id = $1
    `

    var l domain.PengajuanCuti
    err := r.DB.QueryRow(ctx, query, id).Scan(
        &l.ID,
        &l.KaryawanID,
        &l.SubTipe,
        &l.TanggalMulai,
        &l.TanggalSelesai,
        &l.TotalHari,
        &l.Alasan,
        &l.Status,
        &l.BackDate,
        &l.MengurangiCuti,
        &l.LangsungApprove,
        &l.LangsungFinal,
        &l.JudulDokumen,
        &l.DisetujuiOleh,
        &l.TanggalDisetujui,
        &l.DifinalisasiOleh,
        &l.TanggalDifinalisasi,
        &l.URLPDF,
        &l.AlasanBatal,
        &l.TanggalDibatalkan,
        &l.AlasanDitolak,
        &l.TanggalDitolak,
        &l.DibuatPada,
        &l.DiperbaruiPada,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &l, nil
}

func (r *LeaveRepo) GetByKaryawanID(ctx context.Context, karyawanID string, status string, limit, offset int) ([]domain.PengajuanCuti, int, error) {
    var items []domain.PengajuanCuti
    var total int

    baseQuery := `FROM pengajuan_cuti WHERE karyawan_id = $1`
    args := []interface{}{karyawanID}
    argIdx := 2

    if status != "" {
        baseQuery += ` AND status = $` + strconv.Itoa(argIdx)
        args = append(args, status)
        argIdx++
    }

    countQuery := `SELECT COUNT(*) ` + baseQuery
    err := r.DB.QueryRow(ctx, countQuery, args...).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    dataQuery := `
        SELECT id, karyawan_id, sub_tipe, tanggal_mulai, tanggal_selesai,
               total_hari, alasan, status, back_date, mengurangi_cuti, 
               langsung_approve, langsung_final,
               judul_dokumen, disetujui_oleh, tanggal_disetujui, difinalisasi_oleh,
               tanggal_difinalisasi, url_pdf, alasan_batal, tanggal_dibatalkan,
               alasan_ditolak, tanggal_ditolak,
               dibuat_pada, diperbarui_pada
    ` + baseQuery + ` ORDER BY dibuat_pada DESC LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

    finalArgs := append(args, limit, offset)
    rows, err := r.DB.Query(ctx, dataQuery, finalArgs...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    for rows.Next() {
        var l domain.PengajuanCuti
        err := rows.Scan(
            &l.ID,
            &l.KaryawanID,
            &l.SubTipe,
            &l.TanggalMulai,
            &l.TanggalSelesai,
            &l.TotalHari,
            &l.Alasan,
            &l.Status,
            &l.BackDate,
            &l.MengurangiCuti,
            &l.LangsungApprove,
            &l.LangsungFinal,
            &l.JudulDokumen,
            &l.DisetujuiOleh,
            &l.TanggalDisetujui,
            &l.DifinalisasiOleh,
            &l.TanggalDifinalisasi,
            &l.URLPDF,
            &l.AlasanBatal,
            &l.TanggalDibatalkan,
            &l.AlasanDitolak,
            &l.TanggalDitolak,
            &l.DibuatPada,
            &l.DiperbaruiPada,
        )
        if err != nil {
            return nil, 0, err
        }
        items = append(items, l)
    }

    return items, total, nil
}

func (r *LeaveRepo) UpdateStatus(ctx context.Context, id, status string) error {
    query := `UPDATE pengajuan_cuti SET status = $1, diperbarui_pada = NOW() WHERE id = $2`
    _, err := r.DB.Exec(ctx, query, status, id)
    return err
}

func (r *LeaveRepo) Approve(ctx context.Context, id, managerID string) error {
    query := `
        UPDATE pengajuan_cuti 
        SET status = 'disetujui', disetujui_oleh = $1, tanggal_disetujui = NOW(),
            diperbarui_pada = NOW()
        WHERE id = $2
    `
    _, err := r.DB.Exec(ctx, query, managerID, id)
    return err
}

func (r *LeaveRepo) Reject(ctx context.Context, id, managerID, alasan string) error {
    query := `
        UPDATE pengajuan_cuti 
        SET status = 'ditolak', 
            disetujui_oleh = $1, 
            tanggal_disetujui = NOW(),
            alasan_ditolak = $2,
            tanggal_ditolak = NOW(),
            diperbarui_pada = NOW()
        WHERE id = $3
    `
    _, err := r.DB.Exec(ctx, query, managerID, alasan, id)
    return err
}

func (r *LeaveRepo) UpdateAlasanBatal(ctx context.Context, id, alasan string) error {
    query := `
        UPDATE pengajuan_cuti 
        SET alasan_batal = $1,
            tanggal_dibatalkan = NOW(),
            diperbarui_pada = NOW()
        WHERE id = $2
    `
    _, err := r.DB.Exec(ctx, query, alasan, id)
    return err
}

func (r *LeaveRepo) Finalize(ctx context.Context, id, hrdID string) error {
    query := `
        UPDATE pengajuan_cuti 
        SET status = 'disetujui',
            difinalisasi_oleh = $1, 
            tanggal_difinalisasi = NOW(),
            diperbarui_pada = NOW()
        WHERE id = $2 AND difinalisasi_oleh IS NULL
    `
    result, err := r.DB.Exec(ctx, query, hrdID, id)
    if err != nil {
        return err
    }

    rowsAffected := result.RowsAffected()
    if rowsAffected == 0 {
        return errors.New("tidak ada data yang diupdate")
    }

    return nil
}

func (r *LeaveRepo) UpdatePDFURL(ctx context.Context, id, url string) error {
    query := `UPDATE pengajuan_cuti SET url_pdf = $1, diperbarui_pada = NOW() WHERE id = $2`
    _, err := r.DB.Exec(ctx, query, url, id)
    return err
}

func (r *LeaveRepo) GetBalance(ctx context.Context, karyawanID string, year int) (*domain.SisaCuti, error) {
    query := `
        SELECT id, karyawan_id, tahun, jumlah_cuti, telah_dilaksanakan, 
               akan_dilaksanakan, sisa_cuti, dibuat_pada, diperbarui_pada
        FROM sisa_cuti
        WHERE karyawan_id = $1 AND tahun = $2
    `

    var b domain.SisaCuti
    err := r.DB.QueryRow(ctx, query, karyawanID, year).Scan(
        &b.ID,
        &b.KaryawanID,
        &b.Tahun,
        &b.JumlahCuti,
        &b.TelahDilaksanakan,
        &b.AkanDilaksanakan,
        &b.SisaCuti,
        &b.DibuatPada,
        &b.DiperbaruiPada,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &b, nil
}

func (r *LeaveRepo) GetBalanceWithCarryOver(ctx context.Context, karyawanID string, year int) (*domain.SisaCuti, error) {
    now := time.Now()
    bulan := int(now.Month())

    balance, err := r.GetBalance(ctx, karyawanID, year)
    if err != nil {
        return nil, err
    }

    if balance == nil {
        balance = &domain.SisaCuti{
            KaryawanID:        karyawanID,
            Tahun:             year,
            JumlahCuti:        12,
            TelahDilaksanakan: 0,
            AkanDilaksanakan:  0,
            SisaCuti:          12,
        }
    }

    if bulan >= 1 && bulan <= 3 {
        balanceTahunLalu, err := r.GetBalance(ctx, karyawanID, year-1)
        if err == nil && balanceTahunLalu != nil && balanceTahunLalu.SisaCuti > 0 {
            balance.JumlahCuti = 12 + balanceTahunLalu.SisaCuti
        }
    }

    return balance, nil
}

func (r *LeaveRepo) GetLiburByDate(ctx context.Context, tanggal time.Time) (*domain.Libur, error) {
    query := `
        SELECT id, nama, tanggal, aktif, dibuat_pada, diperbarui_pada
        FROM libur
        WHERE DATE(tanggal) = DATE($1)
        LIMIT 1
    `

    var l domain.Libur
    err := r.DB.QueryRow(ctx, query, tanggal).Scan(
        &l.ID,
        &l.Nama,
        &l.Tanggal,
        &l.Aktif,
        &l.DibuatPada,
        &l.DiperbaruiPada,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &l, nil
}

func (r *LeaveRepo) UpdateAkanDilaksanakan(ctx context.Context, karyawanID string, tahun int) error {
    query := `
        UPDATE sisa_cuti 
        SET akan_dilaksanakan = (
            SELECT COALESCE(SUM(total_hari), 0)
            FROM pengajuan_cuti
            WHERE karyawan_id = $1
              AND mengurangi_cuti = true
              AND EXTRACT(YEAR FROM tanggal_mulai) = $2
              AND (
                  status = 'menunggu'
                  OR (
                      status = 'disetujui'
                      AND difinalisasi_oleh IS NULL
                  )
              )
        ),
        diperbarui_pada = NOW()
        WHERE karyawan_id = $1 AND tahun = $2
    `
    _, err := r.DB.Exec(ctx, query, karyawanID, tahun)
    return err
}

func (r *LeaveRepo) UpdateBalance(ctx context.Context, karyawanID string, tahun int) error {
    query := `
        UPDATE sisa_cuti 
        SET telah_dilaksanakan = (
            SELECT COALESCE(SUM(total_hari), 0)
            FROM pengajuan_cuti
            WHERE karyawan_id = $1
              AND status = 'disetujui'
              AND difinalisasi_oleh IS NOT NULL
              AND mengurangi_cuti = true
              AND EXTRACT(YEAR FROM tanggal_mulai) = $2
        ),
        akan_dilaksanakan = (
            SELECT COALESCE(SUM(total_hari), 0)
            FROM pengajuan_cuti
            WHERE karyawan_id = $1
              AND mengurangi_cuti = true
              AND EXTRACT(YEAR FROM tanggal_mulai) = $2
              AND (
                  status = 'menunggu'
                  OR (
                      status = 'disetujui'
                      AND difinalisasi_oleh IS NULL
                  )
              )
        ),
        sisa_cuti = jumlah_cuti - (
            SELECT COALESCE(SUM(total_hari), 0)
            FROM pengajuan_cuti
            WHERE karyawan_id = $1
              AND status = 'disetujui'
              AND difinalisasi_oleh IS NOT NULL
              AND mengurangi_cuti = true
              AND EXTRACT(YEAR FROM tanggal_mulai) = $2
        ),
        diperbarui_pada = NOW()
        WHERE karyawan_id = $1 AND tahun = $2
    `
    _, err := r.DB.Exec(ctx, query, karyawanID, tahun)
    return err
}

func (r *LeaveRepo) GetActiveLeaves(ctx context.Context, karyawanID string) ([]domain.PengajuanCuti, error) {
    query := `
        SELECT id, karyawan_id, sub_tipe, tanggal_mulai, tanggal_selesai,
               total_hari, alasan, status, back_date, mengurangi_cuti, 
               langsung_approve, langsung_final,
               judul_dokumen, disetujui_oleh, tanggal_disetujui, difinalisasi_oleh,
               tanggal_difinalisasi, url_pdf, alasan_batal, tanggal_dibatalkan,
               alasan_ditolak, tanggal_ditolak,
               dibuat_pada, diperbarui_pada
        FROM pengajuan_cuti
        WHERE karyawan_id = $1 AND status NOT IN ('ditolak', 'dibatalkan')
    `

    rows, err := r.DB.Query(ctx, query, karyawanID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []domain.PengajuanCuti
    for rows.Next() {
        var l domain.PengajuanCuti
        err := rows.Scan(
            &l.ID,
            &l.KaryawanID,
            &l.SubTipe,
            &l.TanggalMulai,
            &l.TanggalSelesai,
            &l.TotalHari,
            &l.Alasan,
            &l.Status,
            &l.BackDate,
            &l.MengurangiCuti,
            &l.LangsungApprove,
            &l.LangsungFinal,
            &l.JudulDokumen,
            &l.DisetujuiOleh,
            &l.TanggalDisetujui,
            &l.DifinalisasiOleh,
            &l.TanggalDifinalisasi,
            &l.URLPDF,
            &l.AlasanBatal,
            &l.TanggalDibatalkan,
            &l.AlasanDitolak,
            &l.TanggalDitolak,
            &l.DibuatPada,
            &l.DiperbaruiPada,
        )
        if err != nil {
            return nil, err
        }
        items = append(items, l)
    }
    return items, nil
}

func (r *LeaveRepo) GetActiveLeavesByDate(ctx context.Context, tanggal time.Time) ([]string, error) {
    query := `
        SELECT DISTINCT karyawan_id
        FROM pengajuan_cuti
        WHERE status = 'disetujui'
        AND tanggal_mulai <= $1 AND tanggal_selesai >= $1
        AND mengurangi_cuti = true
    `
    rows, err := r.DB.Query(ctx, query, tanggal)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ids []string
    for rows.Next() {
        var id string
        if err := rows.Scan(&id); err != nil {
            continue
        }
        ids = append(ids, id)
    }
    return ids, nil
}

func (r *LeaveRepo) GetAllLeaves(ctx context.Context, atasanID string, role string, status string, subTipe string, startDate string, endDate string, limit int, offset int) ([]domain.LeaveWithKaryawanResponse, int, error) {
    var items []domain.LeaveWithKaryawanResponse
    var total int

    query := `
        FROM pengajuan_cuti pc
        JOIN karyawan k ON pc.karyawan_id = k.id
        LEFT JOIN sisa_cuti sc ON pc.karyawan_id = sc.karyawan_id AND sc.tahun = EXTRACT(YEAR FROM NOW())
        WHERE 1=1
        AND k.role != 'admin'
    `
    args := []interface{}{}
    argIdx := 1

    if role == "atasan" && atasanID != "" {
        query += ` AND k.atasan_langsung_id = $` + strconv.Itoa(argIdx)
        args = append(args, atasanID)
        argIdx++
    }

    if status != "" {
        query += ` AND pc.status = $` + strconv.Itoa(argIdx)
        args = append(args, status)
        argIdx++
    }

    if subTipe != "" {
        query += ` AND pc.sub_tipe = $` + strconv.Itoa(argIdx)
        args = append(args, subTipe)
        argIdx++
    }

    if startDate != "" {
        query += ` AND pc.tanggal_mulai >= $` + strconv.Itoa(argIdx)
        args = append(args, startDate)
        argIdx++
    }

    if endDate != "" {
        query += ` AND pc.tanggal_selesai <= $` + strconv.Itoa(argIdx)
        args = append(args, endDate)
        argIdx++
    }

    if role == "hrd" {
        query += ` AND pc.difinalisasi_oleh IS NULL AND pc.status = 'disetujui'`
    }

    countQuery := `SELECT COUNT(*) ` + query
    err := r.DB.QueryRow(ctx, countQuery, args...).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    dataQuery := `
        SELECT pc.id, pc.karyawan_id, pc.sub_tipe, pc.tanggal_mulai, pc.tanggal_selesai,
               pc.total_hari, pc.alasan, pc.status, pc.back_date, pc.mengurangi_cuti,
               pc.langsung_approve, pc.langsung_final, pc.judul_dokumen,
               pc.disetujui_oleh, pc.tanggal_disetujui, pc.difinalisasi_oleh,
               pc.tanggal_difinalisasi, pc.url_pdf, pc.alasan_batal, pc.tanggal_dibatalkan,
               pc.alasan_ditolak, pc.tanggal_ditolak,
               pc.dibuat_pada, pc.diperbarui_pada,
               k.nama_lengkap, k.divisi, k.unit, k.role,
               COALESCE(sc.sisa_cuti, 12) as sisa_cuti
    ` + query + ` ORDER BY pc.dibuat_pada DESC LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

    finalArgs := append(args, limit, offset)
    rows, err := r.DB.Query(ctx, dataQuery, finalArgs...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    for rows.Next() {
        var l domain.LeaveWithKaryawanResponse
        err := rows.Scan(
            &l.ID, &l.KaryawanID, &l.SubTipe, &l.TanggalMulai, &l.TanggalSelesai,
            &l.TotalHari, &l.Alasan, &l.Status, &l.BackDate, &l.MengurangiCuti,
            &l.LangsungApprove, &l.LangsungFinal, &l.JudulDokumen,
            &l.DisetujuiOleh, &l.TanggalDisetujui, &l.DifinalisasiOleh,
            &l.TanggalDifinalisasi, &l.URLPDF, &l.AlasanBatal, &l.TanggalDibatalkan,
            &l.AlasanDitolak, &l.TanggalDitolak,
            &l.DibuatPada, &l.DiperbaruiPada,
            &l.KaryawanNama, &l.KaryawanDivisi, &l.KaryawanUnit, &l.KaryawanRole,
            &l.SisaCuti,
        )
        if err != nil {
            return nil, 0, err
        }
        items = append(items, l)
    }

    return items, total, nil
}

func (r *LeaveRepo) GetPendingLeavesByAtasan(ctx context.Context, atasanID string, limit int, offset int) ([]domain.LeaveWithKaryawanResponse, int, error) {
    return r.GetAllLeaves(ctx, atasanID, "atasan", "menunggu", "", "", "", limit, offset)
}