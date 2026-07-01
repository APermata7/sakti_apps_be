CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.diperbarui_pada = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_karyawan_updated ON karyawan;
CREATE TRIGGER trg_karyawan_updated BEFORE UPDATE ON karyawan FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_presensi_updated ON presensi;
CREATE TRIGGER trg_presensi_updated BEFORE UPDATE ON presensi FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_pengajuan_cuti_updated ON pengajuan_cuti;
CREATE TRIGGER trg_pengajuan_cuti_updated BEFORE UPDATE ON pengajuan_cuti FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_sisa_cuti_updated ON sisa_cuti;
CREATE TRIGGER trg_sisa_cuti_updated BEFORE UPDATE ON sisa_cuti FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_tanda_tangan_updated ON tanda_tangan;
CREATE TRIGGER trg_tanda_tangan_updated BEFORE UPDATE ON tanda_tangan FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_token_fcm_updated ON token_fcm;
CREATE TRIGGER trg_token_fcm_updated BEFORE UPDATE ON token_fcm FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_kantor_updated ON kantor;
CREATE TRIGGER trg_kantor_updated BEFORE UPDATE ON kantor FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_libur_updated ON libur;
CREATE TRIGGER trg_libur_updated BEFORE UPDATE ON libur FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS trg_konfigurasi_kerja_updated ON konfigurasi_kerja;
CREATE TRIGGER trg_konfigurasi_kerja_updated BEFORE UPDATE ON konfigurasi_kerja FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();