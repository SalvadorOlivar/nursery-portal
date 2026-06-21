-- +goose Up
CREATE TABLE leave_requests (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    fecha_inicio DATE NOT NULL,
    fecha_fin   DATE NOT NULL,
    tipo        VARCHAR(20) NOT NULL CHECK (tipo IN ('VACACIONES', 'ENFERMEDAD', 'PERSONAL', 'DIA_FAVOR')),
    estado      VARCHAR(20) NOT NULL DEFAULT 'PENDIENTE' CHECK (estado IN ('PENDIENTE', 'APROBADO', 'RECHAZADO')),
    motivo      TEXT NOT NULL DEFAULT '',
    aprobado_por UUID REFERENCES auth_users(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_leave_requests_employee_id ON leave_requests(employee_id);
CREATE INDEX idx_leave_requests_estado ON leave_requests(estado);
CREATE INDEX idx_leave_requests_fechas ON leave_requests(fecha_inicio, fecha_fin);

CREATE TABLE compensatory_days (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    fecha_origen DATE NOT NULL,
    motivo      VARCHAR(30) NOT NULL CHECK (motivo IN ('DOBLE_TURNO', 'DESCANSO_LABORADO')),
    turno_id    UUID REFERENCES turnos(id) ON DELETE SET NULL,
    descripcion TEXT NOT NULL DEFAULT '',
    utilizado   BOOLEAN NOT NULL DEFAULT false,
    fecha_uso   DATE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_compensatory_days_employee_id ON compensatory_days(employee_id);

-- +goose Down
DROP TABLE IF EXISTS compensatory_days;
DROP TABLE IF EXISTS leave_requests;
