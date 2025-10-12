package structs

import (
	"time"
)

type PlantillaCumpleanos struct {
	Titulo       string
	Nombre       string
	Descripcion  string
	Edad         int
	Funcionarios []string
	Dias         int
}

type PlantillaAniversarios struct {
	Titulo      string
	Nombre      string
	Tiempo      int
	Descripcion string
}

type DatosFuncionarios struct {
	FechaNacimiento     time.Time
	NombreFuncionario   string
	ApellidoFuncionario string
	CorreoFuncionario   string
}

type DatosAniversarios struct {
	FechaIngreso        time.Time
	NombreAniversario   string
	ApellidoAniversario string
	GeneroAniversarios  string
	CorreoAniversario   string
	Descripcion         string
}

type Config struct {
	BaseEmail         string
	PasswordBaseEmail string
}

