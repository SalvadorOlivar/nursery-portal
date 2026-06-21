package planificacion

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPlanificacion(t *testing.T) {
	p, err := NewPlanificacion(NewPlanificacionParams{
		Semana: 24,
		Anio:   2026,
		Nombre: "Planificación Semana 24 2026",
	})
	require.NoError(t, err)
	assert.Equal(t, 24, p.Semana)
	assert.Equal(t, 2026, p.Anio)
	assert.Equal(t, "Planificación Semana 24 2026", p.Nombre)
	assert.Equal(t, EstadoBorrador, p.Estado)
	assert.NotEmpty(t, p.ID)
	assert.Equal(t, 7, p.Dias())
}

func TestNewPlanificacion_InvalidSemana(t *testing.T) {
	_, err := NewPlanificacion(NewPlanificacionParams{Semana: 54, Anio: 2026, Nombre: "test"})
	assert.Error(t, err)
}

func TestPlanificacion_Publicar(t *testing.T) {
	p, _ := NewPlanificacion(NewPlanificacionParams{Semana: 24, Anio: 2026, Nombre: "test"})
	err := p.Publicar()
	assert.NoError(t, err)
	assert.Equal(t, EstadoPublicado, p.Estado)
}

func TestPlanificacion_Publicar_Twice(t *testing.T) {
	p, _ := NewPlanificacion(NewPlanificacionParams{Semana: 24, Anio: 2026, Nombre: "test"})
	_ = p.Publicar()
	err := p.Publicar()
	assert.Error(t, err)
}

func TestPlanificacion_Cerrar(t *testing.T) {
	p, _ := NewPlanificacion(NewPlanificacionParams{Semana: 24, Anio: 2026, Nombre: "test"})
	_ = p.Publicar()
	err := p.Cerrar()
	assert.NoError(t, err)
	assert.Equal(t, EstadoCerrado, p.Estado)
}

func TestDias(t *testing.T) {
	assert.Equal(t, 7, DiasSemana)
}