'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog'
import {
  useSwapRequests,
  useAcceptSwapRequest,
  useRejectSwapRequest,
  useApproveSwapRequest,
  useCancelSwapRequest,
  useSwapHistory,
} from '@/features/intercambio/hooks/use-intercambio'
import { useMe } from '@/features/auth/hooks/use-auth'
import { useEmployees } from '@/features/employees/hooks/use-employees'
import { usePlanificaciones } from '@/features/planificaciones/hooks/use-planificaciones'
import { toast } from 'sonner'
import type { ShiftSwapRequest } from '@/types/intercambio'

const estadoLabels: Record<string, string> = {
  PENDIENTE_RESPUESTA: 'Esperando respuesta',
  PENDIENTE_APROBACION: 'Esperando aprobación',
  APROBADO: 'Aprobado',
  RECHAZADO: 'Rechazado',
  CANCELADO: 'Cancelado',
}

const estadoColors: Record<string, 'default' | 'secondary' | 'outline' | 'destructive'> = {
  PENDIENTE_RESPUESTA: 'secondary',
  PENDIENTE_APROBACION: 'secondary',
  APROBADO: 'default',
  RECHAZADO: 'destructive',
  CANCELADO: 'outline',
}

function HistoryDialog({ swapId, open, onOpenChange }: { swapId: string; open: boolean; onOpenChange: (v: boolean) => void }) {
  const { data, isLoading } = useSwapHistory(swapId)

  const accionLabels: Record<string, string> = {
    SOLICITADO: 'Solicitado',
    ACEPTADO: 'Aceptado por destino',
    RECHAZADO: 'Rechazado',
    APROBADO: 'Aprobado por supervisor',
    CANCELADO: 'Cancelado',
    EJECUTADO: 'Swap ejecutado',
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Historial del intercambio</DialogTitle>
          <DialogDescription>Secuencia de eventos del intercambio de turnos</DialogDescription>
        </DialogHeader>
        {isLoading ? (
          <div className="text-sm text-muted-foreground">Cargando historial...</div>
        ) : (
          <div className="space-y-2">
            {data?.data?.map((entry) => (
              <div key={entry.id} className="flex items-center gap-2 text-sm">
                <Badge variant="outline">{accionLabels[entry.accion] ?? entry.accion}</Badge>
                <span className="text-muted-foreground">
                  {new Date(entry.created_at).toLocaleString('es-AR')}
                </span>
              </div>
            ))}
            {!data?.data?.length && (
              <div className="text-sm text-muted-foreground">Sin historial disponible.</div>
            )}
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}

function SwapRequestRow({ req, canApprove }: { req: ShiftSwapRequest; canApprove: boolean }) {
  const [historyOpen, setHistoryOpen] = useState(false)
  const { data: meData } = useMe()
  const { data: employeesData } = useEmployees()
  const { data: planifData } = usePlanificaciones()
  const acceptMutation = useAcceptSwapRequest()
  const rejectMutation = useRejectSwapRequest()
  const approveMutation = useApproveSwapRequest()
  const cancelMutation = useCancelSwapRequest()

  const employees = employeesData?.data ?? []
  const solicitante = employees.find((e) => e.id === req.solicitante_id)
  const destino = employees.find((e) => e.id === req.destino_id)
  const planif = planifData?.data?.find((p) => p.id === req.planificacion_id)

  const userEmployeeID = meData?.user?.employee_id
  const isDestino = userEmployeeID === req.destino_id
  const isSolicitante = userEmployeeID === req.solicitante_id

  async function handleAccept() {
    try {
      await acceptMutation.mutateAsync(req.id)
      toast.success('Intercambio aceptado')
    } catch { toast.error('Error al aceptar') }
  }

  async function handleReject() {
    try {
      await rejectMutation.mutateAsync(req.id)
      toast.success('Intercambio rechazado')
    } catch { toast.error('Error al rechazar') }
  }

  async function handleApprove() {
    try {
      await approveMutation.mutateAsync(req.id)
      toast.success('Intercambio aprobado')
    } catch { toast.error('Error al aprobar') }
  }

  async function handleCancel() {
    try {
      await cancelMutation.mutateAsync(req.id)
      toast.success('Intercambio cancelado')
    } catch { toast.error('Error al cancelar') }
  }

  return (
    <>
      <TableRow>
        <TableCell className="font-medium">
          {solicitante ? `${solicitante.apellido}, ${solicitante.nombre}` : req.solicitante_id}
        </TableCell>
        <TableCell>
          {destino ? `${destino.apellido}, ${destino.nombre}` : req.destino_id}
        </TableCell>
        <TableCell className="text-xs text-muted-foreground">
          {planif?.nombre ?? req.planificacion_id}
        </TableCell>
        <TableCell>
          <Badge variant={estadoColors[req.estado] ?? 'secondary'}>
            {estadoLabels[req.estado] ?? req.estado}
          </Badge>
        </TableCell>
        <TableCell className="text-right">
          <div className="flex justify-end gap-1 flex-wrap">
            {isDestino && req.estado === 'PENDIENTE_RESPUESTA' && (
              <>
                <Button variant="outline" size="sm" className="text-green-600" onClick={handleAccept} disabled={acceptMutation.isPending}>
                  Aceptar
                </Button>
                <Button variant="outline" size="sm" className="text-destructive" onClick={handleReject} disabled={rejectMutation.isPending}>
                  Rechazar
                </Button>
              </>
            )}
            {canApprove && req.estado === 'PENDIENTE_APROBACION' && (
              <>
                <Button variant="outline" size="sm" className="text-green-600" onClick={handleApprove} disabled={approveMutation.isPending}>
                  Aprobar
                </Button>
                <Button variant="outline" size="sm" className="text-destructive" onClick={handleReject} disabled={rejectMutation.isPending}>
                  Rechazar
                </Button>
              </>
            )}
            {isSolicitante && (req.estado === 'PENDIENTE_RESPUESTA' || req.estado === 'PENDIENTE_APROBACION') && (
              <Button variant="outline" size="sm" onClick={handleCancel} disabled={cancelMutation.isPending}>
                Cancelar
              </Button>
            )}
            <Button variant="ghost" size="sm" onClick={() => setHistoryOpen(true)}>
              Historial
            </Button>
          </div>
        </TableCell>
      </TableRow>
      <HistoryDialog swapId={req.id} open={historyOpen} onOpenChange={setHistoryOpen} />
    </>
  )
}

export function SwapRequestList() {
  const { data: meData } = useMe()
  const { data, isLoading, isError } = useSwapRequests()
  const canApprove = meData?.user.role === 'ADMIN' || meData?.user.role === 'SUPERVISOR'

  if (isLoading) {
    return <div className="text-center py-8 text-muted-foreground">Cargando solicitudes de intercambio...</div>
  }

  if (isError) {
    return <div className="text-center py-8 text-destructive">Error al cargar solicitudes de intercambio.</div>
  }

  const requests = data?.data ?? []

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Intercambio de Turnos</CardTitle>
        <Link href="/intercambio/new">
          <Button>Nuevo intercambio</Button>
        </Link>
      </CardHeader>
      <CardContent>
        {requests.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            No hay solicitudes de intercambio.
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Solicitante</TableHead>
                <TableHead>Destino</TableHead>
                <TableHead>Planificación</TableHead>
                <TableHead>Estado</TableHead>
                <TableHead className="text-right">Acciones</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {requests.map((req) => (
                <SwapRequestRow key={req.id} req={req} canApprove={canApprove} />
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>
    </Card>
  )
}
