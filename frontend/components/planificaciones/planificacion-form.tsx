'use client'

import { useMemo, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { toast } from 'sonner'
import { useCreatePlanificacion } from '@/features/planificaciones/hooks/use-planificaciones'

const monthNames = [
  'Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio',
  'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre',
]

const shortMonths = ['ene', 'feb', 'mar', 'abr', 'may', 'jun', 'jul', 'ago', 'sep', 'oct', 'nov', 'dic']

interface WeekOption {
  semana: number
  anio: number
  monday: Date
  label: string
}

function getISOWeek(date: Date): { week: number; year: number } {
  const d = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()))
  const dayNum = d.getUTCDay() || 7
  d.setUTCDate(d.getUTCDate() + 4 - dayNum)
  const yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1))
  const weekNo = Math.ceil(((d.getTime() - yearStart.getTime()) / 86400000 + 1) / 7)
  return { week: weekNo, year: d.getUTCFullYear() }
}

function getWeeksInMonth(mes: number, anio: number): WeekOption[] {
  const monthStart = new Date(anio, mes - 1, 1)
  const monthEnd = new Date(anio, mes, 0)

  const firstMonday = new Date(monthStart)
  const dow = firstMonday.getDay()
  firstMonday.setDate(firstMonday.getDate() - (dow === 0 ? 6 : dow - 1))

  const lastMonday = new Date(monthEnd)
  const dow2 = lastMonday.getDay()
  lastMonday.setDate(lastMonday.getDate() - (dow2 === 0 ? 6 : dow2 - 1))

  const weeks: WeekOption[] = []
  const current = new Date(firstMonday)

  while (current <= lastMonday) {
    const iso = getISOWeek(current)
    const weekEnd = new Date(current)
    weekEnd.setDate(weekEnd.getDate() + 6)

    let extra = ''
    if (current.getMonth() !== mes - 1 || weekEnd.getMonth() !== mes - 1) {
      const other = weekEnd.getMonth() !== mes - 1 ? weekEnd.getMonth() : current.getMonth()
      extra = ` (incluye ${shortMonths[other]})`
    }

    weeks.push({
      semana: iso.week,
      anio: iso.year,
      monday: new Date(current),
      label: `Semana del lun ${current.getDate()}/${current.getMonth() + 1}${extra}`,
    })

    current.setDate(current.getDate() + 7)
  }

  return weeks
}

export function PlanificacionForm() {
  const router = useRouter()
  const createMutation = useCreatePlanificacion()

  const today = new Date()
  const [mes, setMes] = useState(String(today.getMonth() + 1))
  const [anio, setAnio] = useState(String(today.getFullYear()))
  const [semana, setSemana] = useState<number | null>(null)
  const [nombre, setNombre] = useState('')
  const [error, setError] = useState('')

  const weeks = useMemo(() => {
    return getWeeksInMonth(Number(mes), Number(anio))
  }, [mes, anio])

  const selectedWeek = weeks.find((w) => w.semana === semana)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')

    if (semana === null) {
      setError('Seleccioná una semana')
      return
    }

    const selected = weeks.find((w) => w.semana === semana)
    const generatedName = nombre || `Planificación ${selected?.label ?? `Semana ${semana} - ${anio}`}`

    try {
      await createMutation.mutateAsync({
        semana,
        anio: Number(anio),
        nombre: generatedName,
      })
      toast.success('Planificación creada correctamente')
      router.push('/planificaciones')
      router.refresh()
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Error al crear planificación'
      setError(message)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6 max-w-lg">
      {error && (
        <div className="bg-destructive/10 text-destructive text-sm p-3 rounded-md">
          {error}
        </div>
      )}

      <div className="space-y-2">
        <Label htmlFor="nombre">Nombre</Label>
        <Input
          id="nombre"
          placeholder={selectedWeek?.label ?? `Planificación Semana ${semana ?? ''} - ${anio}`}
          value={nombre}
          onChange={(e) => setNombre(e.target.value)}
        />
        <p className="text-xs text-muted-foreground">Si se deja vacío, se generará automáticamente.</p>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="mes">Mes</Label>
          <Select value={mes} onValueChange={(v) => { if (v) { setMes(v); setSemana(null) } }}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {monthNames.map((name, i) => (
                <SelectItem key={i + 1} value={String(i + 1)}>{name}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-2">
          <Label htmlFor="anio">Año</Label>
          <Input
            id="anio"
            type="number"
            value={anio}
            onChange={(e) => { setAnio(e.target.value); setSemana(null) }}
            min={2020}
            max={2100}
            required
          />
        </div>
      </div>

      {weeks.length > 0 && (
        <div className="space-y-2">
          <Label>Semana</Label>
          <div className="flex flex-col gap-1.5">
            {weeks.map((w) => (
              <Button
                key={`${w.semana}-${w.anio}`}
                type="button"
                variant={semana === w.semana ? 'default' : 'outline'}
                className="justify-start"
                onClick={() => setSemana(w.semana)}
              >
                {w.label}
              </Button>
            ))}
          </div>
        </div>
      )}

      <div className="flex gap-4">
        <Button type="submit" disabled={createMutation.isPending || semana === null}>
          Crear planificación
        </Button>
        <Button type="button" variant="outline" onClick={() => router.back()}>
          Cancelar
        </Button>
      </div>
    </form>
  )
}