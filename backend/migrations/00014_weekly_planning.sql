-- +goose Up
-- Migrate planning system from monthly to weekly

-- Drop dependent tables that reference planificaciones or turnos
DROP TABLE IF EXISTS compensatory_days;
DROP TABLE IF EXISTS shift_swap_history;
DROP TABLE IF EXISTS shift_swap_requests;
DROP TABLE IF EXISTS turnos;
DROP TABLE IF EXISTS planificacion_dotacion;
DROP TABLE IF EXISTS planificacion_sectores;

-- Recreate planificaciones with semana instead of mes
DROP TABLE IF EXISTS planificaciones CASCADE;

CREATE TABLE planificaciones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    semana INTEGER NOT NULL CHECK (semana BETWEEN 1 AND 53),
    anio INTEGER NOT NULL,
    nombre VARCHAR(200) NOT NULL,
    estado VARCHAR(20) NOT NULL DEFAULT 'BORRADOR' CHECK (estado IN ('BORRADOR', 'PUBLICADO', 'CERRADO')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(semana, anio)
);

-- Recreate turnos with dia_semana instead of dia
CREATE TABLE turnos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    planificacion_id UUID NOT NULL REFERENCES planificaciones(id) ON DELETE CASCADE,
    empleado_id UUID NOT NULL REFERENCES employees(id),
    dia_semana INTEGER NOT NULL CHECK (dia_semana BETWEEN 1 AND 7),
    turno VARCHAR(20) NOT NULL DEFAULT 'MANANA' CHECK (turno IN ('MANANA', 'TARDE', 'VESPERTINO', 'NOCHE')),
    sector VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(planificacion_id, empleado_id, dia_semana, turno)
);

CREATE INDEX idx_turnos_planificacion ON turnos(planificacion_id);
CREATE INDEX idx_turnos_empleado ON turnos(empleado_id);

-- Recreate planificacion_sectores
CREATE TABLE planificacion_sectores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    planificacion_id UUID NOT NULL REFERENCES planificaciones(id) ON DELETE CASCADE,
    nombre VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(planificacion_id, nombre)
);

-- Recreate planificacion_dotacion
CREATE TABLE planificacion_dotacion (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    planificacion_id UUID NOT NULL REFERENCES planificaciones(id) ON DELETE CASCADE,
    sector VARCHAR(50) NOT NULL DEFAULT '',
    tipo_empleado VARCHAR(20) NOT NULL,
    turno VARCHAR(20) NOT NULL,
    cantidad_minima INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(planificacion_id, sector, tipo_empleado, turno)
);

-- Recreate compensatory_days (references turnos)
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

-- Recreate shift_swap tables (references planificaciones and turnos)
CREATE TABLE IF NOT EXISTS shift_swap_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    planificacion_id UUID NOT NULL REFERENCES planificaciones(id) ON DELETE CASCADE,
    turno_solicitante_id UUID NOT NULL REFERENCES turnos(id),
    turno_destino_id UUID NOT NULL REFERENCES turnos(id),
    solicitante_id UUID NOT NULL REFERENCES employees(id),
    destino_id UUID NOT NULL REFERENCES employees(id),
    estado VARCHAR(30) NOT NULL DEFAULT 'PENDIENTE_RESPUESTA'
        CHECK (estado IN ('PENDIENTE_RESPUESTA','PENDIENTE_APROBACION','APROBADO','RECHAZADO','CANCELADO')),
    aprobado_por UUID REFERENCES auth_users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS shift_swap_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    swap_request_id UUID NOT NULL REFERENCES shift_swap_requests(id) ON DELETE CASCADE,
    accion VARCHAR(30) NOT NULL,
    actor_id UUID NOT NULL REFERENCES auth_users(id),
    detalle TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_swap_requests_solicitante ON shift_swap_requests(solicitante_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_destino ON shift_swap_requests(destino_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_estado ON shift_swap_requests(estado);
CREATE INDEX IF NOT EXISTS idx_swap_history_request ON shift_swap_history(swap_request_id);

-- +goose Down
-- Restore monthly schema

-- Drop new weekly tables
DROP TABLE IF EXISTS shift_swap_history;
DROP TABLE IF EXISTS shift_swap_requests;
DROP TABLE IF EXISTS compensatory_days;
DROP TABLE IF EXISTS planificacion_dotacion;
DROP TABLE IF EXISTS planificacion_sectores;
DROP TABLE IF EXISTS turnos;
DROP TABLE IF EXISTS planificaciones CASCADE;

-- Recreate monthly planificaciones
CREATE TABLE planificaciones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mes INTEGER NOT NULL CHECK (mes BETWEEN 1 AND 12),
    anio INTEGER NOT NULL,
    nombre VARCHAR(200) NOT NULL,
    estado VARCHAR(20) NOT NULL DEFAULT 'BORRADOR' CHECK (estado IN ('BORRADOR', 'PUBLICADO', 'CERRADO')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(mes, anio)
);

-- Recreate monthly turnos
CREATE TABLE turnos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    planificacion_id UUID NOT NULL REFERENCES planificaciones(id) ON DELETE CASCADE,
    empleado_id UUID NOT NULL REFERENCES employees(id),
    dia INTEGER NOT NULL CHECK (dia BETWEEN 1 AND 31),
    turno VARCHAR(20) NOT NULL DEFAULT 'MANANA' CHECK (turno IN ('MANANA', 'TARDE', 'VESPERTINO', 'NOCHE')),
    sector VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(planificacion_id, empleado_id, dia, turno)
);

CREATE INDEX idx_turnos_planificacion ON turnos(planificacion_id);
CREATE INDEX idx_turnos_empleado ON turnos(empleado_id);

-- Recreate planificacion_sectores
CREATE TABLE planificacion_sectores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    planificacion_id UUID NOT NULL REFERENCES planificaciones(id) ON DELETE CASCADE,
    nombre VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(planificacion_id, nombre)
);

-- Recreate planificacion_dotacion
CREATE TABLE planificacion_dotacion (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    planificacion_id UUID NOT NULL REFERENCES planificaciones(id) ON DELETE CASCADE,
    sector VARCHAR(50) NOT NULL DEFAULT '',
    tipo_empleado VARCHAR(20) NOT NULL,
    turno VARCHAR(20) NOT NULL,
    cantidad_minima INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(planificacion_id, sector, tipo_empleado, turno)
);

-- Recreate compensatory_days
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

-- Recreate shift_swap tables
CREATE TABLE IF NOT EXISTS shift_swap_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    planificacion_id UUID NOT NULL REFERENCES planificaciones(id) ON DELETE CASCADE,
    turno_solicitante_id UUID NOT NULL REFERENCES turnos(id),
    turno_destino_id UUID NOT NULL REFERENCES turnos(id),
    solicitante_id UUID NOT NULL REFERENCES employees(id),
    destino_id UUID NOT NULL REFERENCES employees(id),
    estado VARCHAR(30) NOT NULL DEFAULT 'PENDIENTE_RESPUESTA'
        CHECK (estado IN ('PENDIENTE_RESPUESTA','PENDIENTE_APROBACION','APROBADO','RECHAZADO','CANCELADO')),
    aprobado_por UUID REFERENCES auth_users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS shift_swap_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    swap_request_id UUID NOT NULL REFERENCES shift_swap_requests(id) ON DELETE CASCADE,
    accion VARCHAR(30) NOT NULL,
    actor_id UUID NOT NULL REFERENCES auth_users(id),
    detalle TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_swap_requests_solicitante ON shift_swap_requests(solicitante_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_destino ON shift_swap_requests(destino_id);
CREATE INDEX IF NOT EXISTS idx_swap_requests_estado ON shift_swap_requests(estado);
CREATE INDEX IF NOT EXISTS idx_swap_history_request ON shift_swap_history(swap_request_id);
