# Arquitectura del proyecto

**Nurse Portal** — Gestión de turnos para enfermería. Monorepo con frontend Next.js + backend Go hexagonal + PostgreSQL.

## Frontend (Next.js 16 / React / TypeScript)

- `app/*/page.tsx` — App Router pages (13 rutas: planificaciones, employees, intercambio, leave-requests, sectores, dotacion, etc.)
- `components/` — React components: `auth/AppShell`, `planificaciones/`, `employees/`, `ausencia/`, `intercambio/`, `ui/` (shadcn + `collapsible-section.tsx`)
- `features/*/hooks/` — React Query hooks (use-auth, use-planificaciones, use-employees, use-ausencia, use-intercambio)
- `lib/api/` — Cliente HTTP + 5 módulos API (auth, employees, planificaciones, ausencia, intercambio)
- `lib/utils.ts` — Funciones utilitarias: `cn()`, `isoWeekToDate()`, `getWeekRange()`, helpers de fecha
- `types/` — Interfaces TypeScript

## Backend (Go, hexagonal + CQRS)

- `internal/adapters/http/` — Handlers REST + middleware auth (chi router)
- `internal/application/services/` — 5 service facades (Auth, Employee, Planificacion, Ausencia, Intercambio)
- `internal/application/commands/` — CQRS command handlers (employee, planificacion, turno, leave, intercambio, compensatory)
- `internal/application/queries/` — CQRS query handlers (planificacion, leave, intercambio, compensatory)
- `internal/domain/` — Entidades de negocio (employee, turno, planificacion, ausencia, intercambio, auth)
- `internal/ports/repositories.go` — Interfaces de repositorio
- `internal/adapters/repository/postgres/` — Implementaciones PostgreSQL (8 repos)

## DB

PostgreSQL 16 con migrations goose.

---

## Regla: Actualizar de arquitectura

Como paso final de cada tarea:

1. Si hiciste cambios que afecten la estructura del proyecto (nuevos componentes, handlers, servicios, repos, rutas, flujos), actualiza `docs/arquitectura.md` y tambien este archivo `AGENTS.md` y su descripcion de la arquitectura.
2. Mantén sincronizadas ambas versiones (detallada y compacta).
3. El diagrama usa sintaxis `architecture` de Mermaid. Refleja las capas: Frontend (pages → components → hooks → API) → Backend (handlers → services → CQRS → ports → repos) → DB.
