-- +goose Up
-- Fix planificaciones incorrectly set to CERRADO by the old auto-close logic
-- that triggered on Sunday morning instead of after Monday starts.
-- Restore to PUBLICADO any plan whose week hasn't fully ended yet
-- (i.e., Monday after the plan's week > TODAY).
-- date_trunc('week', make_date(anio, 1, 4)) = Monday of ISO week 1
-- + semana * 7 days = Monday after the target plan's week
UPDATE planificaciones
SET estado = 'PUBLICADO', updated_at = NOW()
WHERE estado = 'CERRADO'
  AND (date_trunc('week', make_date(anio, 1, 4))::date + semana * 7) > CURRENT_DATE;

-- +goose Down
-- No rollback: the fix is data-corrective and non-reversible.
SELECT 1;
