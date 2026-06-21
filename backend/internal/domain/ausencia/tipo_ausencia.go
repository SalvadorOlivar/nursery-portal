package ausencia

import "fmt"

type TipoAusencia string

const (
	Vacaciones      TipoAusencia = "VACACIONES"
	Enfermedad      TipoAusencia = "ENFERMEDAD"
	Personal        TipoAusencia = "PERSONAL"
	DiaFavor        TipoAusencia = "DIA_FAVOR"
)

func (t TipoAusencia) IsValid() bool {
	switch t {
	case Vacaciones, Enfermedad, Personal, DiaFavor:
		return true
	}
	return false
}

func ParseTipoAusencia(s string) (TipoAusencia, error) {
	t := TipoAusencia(s)
	if !t.IsValid() {
		return "", fmt.Errorf("invalid leave type: %s", s)
	}
	return t, nil
}

type EstadoAusencia string

const (
	Pendiente EstadoAusencia = "PENDIENTE"
	Aprobado  EstadoAusencia = "APROBADO"
	Rechazado EstadoAusencia = "RECHAZADO"
)

func (e EstadoAusencia) IsValid() bool {
	switch e {
	case Pendiente, Aprobado, Rechazado:
		return true
	}
	return false
}

type MotivoCompensatorio string

const (
	DobleTurno       MotivoCompensatorio = "DOBLE_TURNO"
	DescansoLaborado MotivoCompensatorio = "DESCANSO_LABORADO"
)

func (m MotivoCompensatorio) IsValid() bool {
	switch m {
	case DobleTurno, DescansoLaborado:
		return true
	}
	return false
}
