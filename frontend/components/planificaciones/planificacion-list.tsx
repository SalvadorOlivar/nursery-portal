'use client'

import { useMemo } from 'react'
import Link from 'next/link'
import { CalendarPlus, CheckCircle2, Star } from 'lucide-react'
import { EmptyState, PageHeader, StatusBadge } from '@/components/hospital/hospital-ui'
import { useMe } from '@/features/auth/hooks/use-auth'
import { usePlanificaciones } from '@/features/planificaciones/hooks/use-planificaciones'
import { CollapsibleSection } from '@/components/ui/collapsible-section'
import { getWeekRange, getMonthFromWeek, getYearFromWeek, subMonths } from '@/lib/utils'
import type { Planificacion } from '@/types/planificacion'

const monthsEs = ['enero', 'febrero', 'marzo', 'abril', 'mayo', 'junio', 'julio', 'agosto', 'septiembre', 'octubre', 'noviembre', 'diciembre']

const estadoLabel: Record<string, string> = {
  PUBLICADO: 'Publicado',
  CERRADO: 'Cerrado',
  BORRADOR: 'Borrador',
}

const estadoTone: Record<string, 'success' | 'accent' | 'warn'> = {
  PUBLICADO: 'success',
  CERRADO: 'accent',
  BORRADOR: 'warn',
}

interface MonthGroup {
  label: string
  sortKey: number
  plans: Planificacion[]
}

function buildMonthGroups(planList: Planificacion[], sortAsc: boolean): MonthGroup[] {
  const map = new Map<string, MonthGroup>()
  for (const plan of planList) {
    const month = getMonthFromWeek(plan.anio, plan.semana)
    const year = getYearFromWeek(plan.anio, plan.semana)
    const label = `${monthsEs[month]} ${year}`
    const sortKey = year * 12 + month
    if (!map.has(label)) {
      map.set(label, { label, sortKey, plans: [] })
    }
    map.get(label)!.plans.push(plan)
  }
  for (const group of map.values()) {
    group.plans.sort((a, b) =>
      sortAsc
        ? a.anio * 100 + a.semana - (b.anio * 100 + b.semana)
        : b.anio * 100 + b.semana - (a.anio * 100 + a.semana),
    )
  }
  return Array.from(map.values()).sort((a, b) =>
    sortAsc ? a.sortKey - b.sortKey : b.sortKey - a.sortKey,
  )
}

function PlanTable({ groups }: { groups: MonthGroup[] }) {
  return (
    <div className="space-y-3">
      {groups.map((group) => (
        <div key={group.label}>
          <h4 className="text-xs font-semibold uppercase tracking-wider text-[var(--muted-foreground)] mb-2">{group.label}</h4>
          <div className="np-table-wrap">
            <table className="np-table">
              <thead>
                <tr>
                  <th>Nombre</th>
                  <th>Periodo</th>
                  <th>Dias</th>
                  <th>Estado</th>
                </tr>
              </thead>
              <tbody>
                {group.plans.map((plan) => (
                  <tr key={plan.id}>
                    <td className="font-[510]">
                      <Link href={`/planificaciones/${plan.id}`} className="hover:underline">
                        {plan.nombre}
                      </Link>
                    </td>
                    <td>Semana {plan.semana} de {plan.anio}</td>
                    <td>{plan.dias}</td>
                    <td>
                      <StatusBadge tone={estadoTone[plan.estado]}>
                        {estadoLabel[plan.estado]}
                      </StatusBadge>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      ))}
    </div>
  )
}

export function PlanificacionList() {
  const { data: plansData, isLoading, isError } = usePlanificaciones()
  const { data: meData } = useMe()

  const plans = plansData?.data ?? []
  const canEdit = meData?.user.role === 'ADMIN' || meData?.user.role === 'SUPERVISOR'

  const grouped = useMemo((): {
    current: Planificacion | null
    proximasGroups: MonthGroup[]
    recientesGroups: MonthGroup[]
    anterioresGroups: MonthGroup[]
  } => {
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const monthAgo = subMonths(today, 1)

    let currentPlan: Planificacion | null = null
    const proximas: Planificacion[] = []
    const recientes: Planificacion[] = []
    const anteriores: Planificacion[] = []

    for (const plan of plans) {
      const { start, end } = getWeekRange(plan.anio, plan.semana)
      if (today >= start && today <= end) {
        if (!currentPlan || plan.estado === 'PUBLICADO') {
          currentPlan = plan
        }
      } else if (end < today) {
        if (end >= monthAgo) {
          recientes.push(plan)
        } else {
          anteriores.push(plan)
        }
      } else {
        proximas.push(plan)
      }
    }

    const currentId = currentPlan?.id

    return {
      current: currentPlan,
      proximasGroups: buildMonthGroups(proximas.filter((p) => p.id !== currentId), true),
      recientesGroups: buildMonthGroups(recientes.filter((p) => p.id !== currentId), false),
      anterioresGroups: buildMonthGroups(anteriores.filter((p) => p.id !== currentId), false),
    }
  }, [plans])

  if (isLoading) {
    return <div className="np-empty">Cargando planificaciones...</div>
  }

  if (isError) {
    return <div className="np-empty text-[var(--danger)]">Error al cargar planificaciones. Verifica que el servidor este corriendo.</div>
  }

  return (
    <div className="np-page">
      <PageHeader
        title="Planificaciones"
        subtitle="Gestiona las planificaciones semanales del personal."
        actions={
          <>
            <StatusBadge tone={plans.some((plan) => plan.estado === 'PUBLICADO') ? 'success' : 'warn'}>
              <CheckCircle2 className="size-3.5" />
              {plans.some((plan) => plan.estado === 'PUBLICADO') ? 'Semana publicada' : 'Sin publicacion'}
            </StatusBadge>
            {canEdit && (
              <Link href="/planificaciones/new" className="np-btn np-btn-primary">
                <CalendarPlus className="size-4" />
                <span className="np-action-text">Nueva planificacion</span>
              </Link>
            )}
          </>
        }
      />

      {plans.length === 0 ? (
        <EmptyState>No hay planificaciones configuradas.</EmptyState>
      ) : (
        <div className="space-y-4">
          {grouped.current && (
            <div className="np-card border-l-4" style={{ borderLeftColor: 'var(--success)' }}>
              <div className="np-card-body">
                <div className="flex items-start justify-between gap-4">
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <Star className="size-3.5" style={{ color: 'var(--success)' }} />
                      <span className="text-xs font-semibold uppercase tracking-wider" style={{ color: 'var(--success)' }}>
                        Planificacion actual
                      </span>
                      <StatusBadge tone={estadoTone[grouped.current.estado]}>{estadoLabel[grouped.current.estado]}</StatusBadge>
                    </div>
                    <Link href={`/planificaciones/${grouped.current.id}`} className="font-[510] text-base hover:underline">
                      {grouped.current.nombre}
                    </Link>
                    <p className="text-sm" style={{ color: 'var(--muted-foreground)' }}>
                      Semana {grouped.current.semana} de {grouped.current.anio}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          )}

          {grouped.proximasGroups.length > 0 && (
            <CollapsibleSection
              title="Proximas"
              count={grouped.proximasGroups.reduce((s, g) => s + g.plans.length, 0)}
              defaultOpen
            >
              <PlanTable groups={grouped.proximasGroups} />
            </CollapsibleSection>
          )}

          {grouped.recientesGroups.length > 0 && (
            <CollapsibleSection
              title="Recientes"
              count={grouped.recientesGroups.reduce((s, g) => s + g.plans.length, 0)}
              defaultOpen
            >
              <PlanTable groups={grouped.recientesGroups} />
            </CollapsibleSection>
          )}

          {grouped.anterioresGroups.length > 0 && (
            <CollapsibleSection
              title="Anteriores"
              count={grouped.anterioresGroups.reduce((s, g) => s + g.plans.length, 0)}
              defaultOpen={false}
            >
              <PlanTable groups={grouped.anterioresGroups} />
            </CollapsibleSection>
          )}
        </div>
      )}
    </div>
  )
}
