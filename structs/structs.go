package structs

import (
	"time"
)

type PlantillaCumpleanos struct {
	Titulo           string
	Fecha            int
	Nombre           string
	Apellido         string
	Descripcion      string
	Edad             int
	Funcionarios     []string
	Dias             int
}

type DatosFuncionarios struct {
	FechaNacimiento     time.Time
	NombreFuncionario   string
	ApellidoFuncionario string
	CorreoFuncionario   string
}
