package planificacion

import "time"

func TurnoDate(semana, anio, diaSemana int) time.Time {
	jan4 := time.Date(anio, time.January, 4, 0, 0, 0, 0, time.UTC)

	jan4Weekday := jan4.Weekday()
	if jan4Weekday == 0 {
		jan4Weekday = 7
	}

	daysBack := int(jan4Weekday) - 1
	mondayOfWeek1 := jan4.AddDate(0, 0, -daysBack)

	monday := mondayOfWeek1.AddDate(0, 0, (semana-1)*7)
	return monday.AddDate(0, 0, diaSemana-1)
}
