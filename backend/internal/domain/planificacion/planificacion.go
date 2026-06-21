package planificacion

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Estado string

const (
	EstadoBorrador  Estado = "BORRADOR"
	EstadoPublicado Estado = "PUBLICADO"
	EstadoCerrado   Estado = "CERRADO"
)

func (e Estado) IsValid() bool {
	switch e {
	case EstadoBorrador, EstadoPublicado, EstadoCerrado:
		return true
	}
	return false
}

const DiasSemana = 7

type Planificacion struct {
	ID        string
	Semana    int
	Anio      int
	Nombre    string
	Estado    Estado
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewPlanificacionParams struct {
	Semana int
	Anio   int
	Nombre string
}

func NewPlanificacion(params NewPlanificacionParams) (*Planificacion, error) {
	if params.Semana < 1 || params.Semana > 53 {
		return nil, fmt.Errorf("semana must be between 1 and 53")
	}
	if params.Anio < 2020 || params.Anio > 2100 {
		return nil, fmt.Errorf("anio must be between 2020 and 2100")
	}
	if params.Nombre == "" {
		return nil, fmt.Errorf("nombre is required")
	}

	now := time.Now().UTC()
	return &Planificacion{
		ID:        uuid.New().String(),
		Semana:    params.Semana,
		Anio:      params.Anio,
		Nombre:    params.Nombre,
		Estado:    EstadoBorrador,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (p *Planificacion) Publicar() error {
	if p.Estado != EstadoBorrador {
		return fmt.Errorf("only plans in BORRADOR can be published")
	}
	p.Estado = EstadoPublicado
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Planificacion) Cerrar() error {
	if p.Estado != EstadoPublicado {
		return fmt.Errorf("only plans in PUBLICADO can be closed")
	}
	p.Estado = EstadoCerrado
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Planificacion) Dias() int {
	return DiasSemana
}
