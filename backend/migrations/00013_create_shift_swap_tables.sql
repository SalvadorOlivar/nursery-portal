-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shift_swap_history;
DROP TABLE IF EXISTS shift_swap_requests;
-- +goose StatementEnd
