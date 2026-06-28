'use client'

import Link from 'next/link'
import { ArrowRight, Check, Clock3, Plus, X } from 'lucide-react'
import { toast } from 'sonner'
import { Card, EmptyState, PageHeader, StatsGrid, StatusBadge, fullName } from '@/components/hospital/hospital-ui'
import { useMe } from '@/features/auth/hooks/use-auth'
import { useEmployees } from '@/features/employees/hooks/use-employees'
import {
  useAcceptSwapRequest,
  useApproveSwapRequest,
  useCancelSwapRequest,
  useRejectSwapRequest,
  useSwapRequests,
} from '@/features/intercambio/hooks/use-intercambio'
import type { ShiftSwapRequest } from '@/types/intercambio'
import { ApiError } from '@/lib/api/client'

const statusLabels: Record<string, string> = {
  PENDIENTE_RESPUESTA: 'Pendiente',
  PENDIENTE_APROBACION: 'Pendiente',
  APROBADO: 'Aprobado',
  RECHAZADO: 'Rechazado',
  CANCELADO: 'Cancelado',
}

function statusTone(status: string) {
  if (status === 'APROBADO') return 'success' as const
  if (status === 'RECHAZADO' || status === 'CANCELADO') return 'danger' as const
  return 'warn' as const
}

export function SwapRequestList() {
  const { data: meData } = useMe()
  const { data, isLoading, isError } = useSwapRequests()
  const { data: employeesData } = useEmployees()
  const currentEmployeeId = meData?.user.employee_id
  const canApprove = meData?.user.role === 'ADMIN' || meData?.user.role === 'SUPERVISOR'
  const requests = data?.data ?? []
  const employees = employeesData?.data ?? []
  const today = new Date().toISOString().slice(0, 10)
  const currentMonth = today.slice(0, 7)
  const pending = requests.filter((req) => req.estado === 'PENDIENTE_RESPUESTA' || req.estado === 'PENDIENTE_APROBACION')
  const completed = requests.filter((req) => req.estado === 'APROBADO')

  if (isLoading) {
    return <div className="np-empty">Cargando solicitudes de intercambio...</div>
  }

  if (isError) {
    return <div className="np-empty text-[var(--danger)]">Error al cargar solicitudes de intercambio.</div>
  }

  return (
    <div className="np-page">
      <PageHeader
        title="Intercambios de turnos"
        subtitle="Solicitudes entre funcionarios con revision de supervision."
        actions={
          <Link href="/intercambio/new" className="np-btn np-btn-primary">
            <Plus className="size-4" />
            <span className="np-action-text">Solicitar intercambio</span>
          </Link>
        }
      />

      <StatsGrid
        items={[
          { label: 'Solicitudes activas', value: pending.length, highlight: true },
          { label: 'Pendientes de aprobacion', value: requests.filter((req) => req.estado === 'PENDIENTE_APROBACION').length },
          { label: 'Completados hoy', value: completed.filter((req) => req.updated_at?.slice(0, 10) === today).length },
          { label: 'Historial mes', value: requests.filter((req) => req.created_at?.slice(0, 7) === currentMonth).length },
        ]}
      />

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1fr)_360px]">
        <section className="space-y-3">
          {requests.length === 0 ? (
            <Card>
              <EmptyState>No hay solicitudes de intercambio.</EmptyState>
            </Card>
          ) : (
            requests.map((request) => (
              <SwapCard key={request.id} request={request} employees={employees} canApprove={canApprove} currentEmployeeId={currentEmployeeId} />
            ))
          )}
        </section>

        <Card title="Historial completado">
          {completed.length === 0 ? (
            <EmptyState>Sin intercambios completados.</EmptyState>
          ) : (
            <ol className="space-y-4">
              {completed.slice(0, 6).map((request) => {
                const employee = employees.find((item) => item.id === request.solicitante_id)
                return (
                  <li key={request.id} className="relative border-l border-[var(--border)] pl-4">
                    <span className="absolute -left-[5px] top-1 size-2.5 rounded-full bg-[var(--success)]" />
                    <div className="text-sm font-[510]">{fullName(employee)}</div>
                    <div className="text-[0.78rem] text-[var(--muted-foreground)]">{new Date(request.updated_at).toLocaleString('es-UY')}</div>
                  </li>
                )
              })}
            </ol>
          )}
        </Card>
      </div>
    </div>
  )
}

function SwapCard({
  request,
  employees,
  canApprove,
  currentEmployeeId,
}: {
  request: ShiftSwapRequest
  employees: { id: string; nombre: string; apellido: string }[]
  canApprove: boolean
  currentEmployeeId?: string
}) {
  const acceptMutation = useAcceptSwapRequest()
  const approveMutation = useApproveSwapRequest()
  const rejectMutation = useRejectSwapRequest()
  const cancelMutation = useCancelSwapRequest()
  const isSolicitante = currentEmployeeId === request.solicitante_id
  const isDestino = currentEmployeeId === request.destino_id
  const solicitante = employees.find((employee) => employee.id === request.solicitante_id)
  const destino = employees.find((employee) => employee.id === request.destino_id)

  async function handleAccept() {
    try {
      await acceptMutation.mutateAsync(request.id)
      toast.success('Intercambio aceptado')
    } catch (error) {
      toast.error(error instanceof ApiError ? error.message : 'Error al aceptar intercambio')
    }
  }

  async function handleApprove() {
    try {
      await approveMutation.mutateAsync(request.id)
      toast.success('Intercambio aprobado')
    } catch (error) {
      toast.error(error instanceof ApiError ? error.message : 'Error al aprobar intercambio')
    }
  }

  async function handleReject() {
    try {
      await rejectMutation.mutateAsync(request.id)
      toast.success('Intercambio rechazado')
    } catch (error) {
      toast.error(error instanceof ApiError ? error.message : 'Error al rechazar intercambio')
    }
  }

  async function handleCancel() {
    try {
      await cancelMutation.mutateAsync(request.id)
      toast.success('Intercambio cancelado')
    } catch (error) {
      toast.error(error instanceof ApiError ? error.message : 'Error al cancelar intercambio')
    }
  }

  return (
    <article className="np-card">
      <div className="np-card-body">
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div className="space-y-4">
            <div className="flex flex-wrap items-center gap-3">
              <div>
                <div className="text-sm font-[510]">{fullName(solicitante)}</div>
                <div className="text-[0.78rem] text-[var(--muted-foreground)]">Solicitante</div>
              </div>
              <ArrowRight className="size-4 text-[var(--muted-foreground)]" />
              <div>
                <div className="text-sm font-[510]">{fullName(destino)}</div>
                <div className="text-[0.78rem] text-[var(--muted-foreground)]">Receptor</div>
              </div>
            </div>

            <div className="flex flex-wrap items-center gap-2 text-sm">
              <span className="np-badge">Turno origen</span>
              <ArrowRight className="size-4 text-[var(--muted-foreground)]" />
              <span className="np-badge">Turno destino</span>
            </div>

            <div className="flex items-center gap-2 text-[0.8rem] text-[var(--muted-foreground)]">
              <Clock3 className="size-4" />
              {new Date(request.created_at).toLocaleDateString('es-UY')}
            </div>
          </div>

          <div className="flex flex-col items-end gap-3">
            <StatusBadge tone={statusTone(request.estado)}>{statusLabels[request.estado] ?? request.estado}</StatusBadge>
            {request.estado === 'PENDIENTE_RESPUESTA' && isDestino && (
              <div className="flex gap-2">
                <button type="button" className="np-btn np-btn-sm" onClick={handleAccept} disabled={acceptMutation.isPending}>
                  <Check className="size-4 text-[var(--success)]" />
                  <span className="np-action-text">Aceptar</span>
                </button>
                <button type="button" className="np-btn np-btn-sm" onClick={handleReject} disabled={rejectMutation.isPending}>
                  <X className="size-4 text-[var(--danger)]" />
                  <span className="np-action-text">Rechazar</span>
                </button>
              </div>
            )}
            {request.estado === 'PENDIENTE_RESPUESTA' && isSolicitante && (
              <button type="button" className="np-btn np-btn-sm" onClick={handleCancel} disabled={cancelMutation.isPending}>
                <X className="size-4 text-[var(--muted-foreground)]" />
                <span className="np-action-text">Cancelar</span>
              </button>
            )}
            {request.estado === 'PENDIENTE_APROBACION' && canApprove && (
              <div className="flex gap-2">
                <button type="button" className="np-btn np-btn-sm" onClick={handleApprove} disabled={approveMutation.isPending}>
                  <Check className="size-4 text-[var(--success)]" />
                  <span className="np-action-text">Aprobar</span>
                </button>
              </div>
            )}
            {request.estado === 'PENDIENTE_APROBACION' && isSolicitante && (
              <button type="button" className="np-btn np-btn-sm" onClick={handleCancel} disabled={cancelMutation.isPending}>
                <X className="size-4 text-[var(--muted-foreground)]" />
                <span className="np-action-text">Cancelar</span>
              </button>
            )}
          </div>
        </div>
      </div>
    </article>
  )
}
