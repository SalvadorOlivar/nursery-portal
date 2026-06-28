import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function isoWeekToDate(anio: number, semana: number, diaSemana: number): Date {
  const jan4 = new Date(anio, 0, 4)
  const daysSinceMonday = jan4.getDay() === 0 ? 6 : jan4.getDay() - 1
  const mondayWeek1 = new Date(anio, 0, 4 - daysSinceMonday)
  const result = new Date(mondayWeek1)
  result.setDate(mondayWeek1.getDate() + (semana - 1) * 7 + (diaSemana - 1))
  return result
}

export function getWeekRange(anio: number, semana: number): { start: Date; end: Date } {
  const start = isoWeekToDate(anio, semana, 1)
  const end = isoWeekToDate(anio, semana, 7)
  return { start, end }
}

export function getMonthFromWeek(anio: number, semana: number): number {
  return isoWeekToDate(anio, semana, 4).getMonth()
}

export function getYearFromWeek(anio: number, semana: number): number {
  return isoWeekToDate(anio, semana, 4).getFullYear()
}

export function isDateInRange(date: Date, start: Date, end: Date): boolean {
  return date >= start && date <= end
}

export function addMonths(date: Date, months: number): Date {
  const d = new Date(date)
  d.setMonth(d.getMonth() + months)
  return d
}

export function subMonths(date: Date, months: number): Date {
  return addMonths(date, -months)
}
