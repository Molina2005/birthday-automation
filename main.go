package main

import (
	"fmt"
	"log"
	"modulo/funcionesArchivos"
	"modulo/funcionesArreglos"
	"os"
	"time"

	"github.com/kardianos/service"
	"github.com/robfig/cron/v3"
)

func main() {
	loc, _ := time.LoadLocation("America/Bogota")            // Zona horaria para fechas y horas
	horario__establecido := cron.New(cron.WithLocation(loc)) // Programacion tareas automaticas en un horario establecido

	// Ejecucion funcion lectura archivo funcionariosGyG
	rows__funcionarios, rows__aniversarios, err := funcionesArchivos.LeerArchivoFuncionariosgyg("../cumpleañosgyg/archivo/funcionariosgyg.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	// Ejecucion funcion procesar Datos funcionarios
	Cedula__funcionarios, FechaDe__nacimiento,
		Nombre__Funcionarios, Apellido__funcionario,
		Genero__funcionario, Lugar__trabajo,
		Edad__trabajador, Correo__trabajador,
		Descripcion__trabajador,
		Informacion__funcionarios := funcionesArchivos.ProcesarDatosFuncionarios(rows__funcionarios)
	fmt.Println(Cedula__funcionarios, Edad__trabajador, Correo__trabajador, Descripcion__trabajador, Apellido__funcionario, FechaDe__nacimiento, Nombre__Funcionarios, Genero__funcionario, Lugar__trabajo)
	// Ejecucion funcion procesar datos funcionarios
	Informacion__aniversarios := funcionesArchivos.ProcesarDatosAniversarios(rows__aniversarios)
	// Ejecucion funcion cumpleaños actuales
	funcionesArreglos.CumpleañosActuales(
		Nombre__Funcionarios,
		Informacion__funcionarios)
	// Ejecucion funciones Envio de plantillas
	horario__establecido.AddFunc("0 8 1 * *", func() { // envio de plantilla 01 de cada mes, 8 am
		informacion := Informacion__funcionarios
		funcionesArreglos.EnviarCalendarioCumpleañosMes(informacion)
	})
	// Ejecucion funcion cumpleaños individuales de funcionarios
	horario__establecido.AddFunc("0 8 * * *", func() {
		funcionesArreglos.NotificarCumpleañosFuncionario(
			Nombre__Funcionarios, Apellido__funcionario,
			Descripcion__trabajador, Genero__funcionario,
			Edad__trabajador, Correo__trabajador, Informacion__funcionarios)
	})
	// Ejecucion funcion aniversarios funcionarios
	horario__establecido.AddFunc("0 8 * * *", func() {
		funcionesArreglos.NotificarAniversarioFuncionario(Informacion__aniversarios)
	})
	horario__establecido.Start()

	// Ejecucion funcion para poder vigilar cambios archivo excel
	funcionesArchivos.NotificaCambiosArchivoPrincipal()

	// parte de configuracion servicio!estudiar linea por linea
	config__Servicio := &service.Config{
		Name:        "GoServiceTest",
		DisplayName: "Prueba de servicio Go",
		Description: "Este es un servicio Go de automatizacion",
	}
	prg := &funcionesArchivos.Programa{}
	s, err := service.New(prg, config__Servicio)
	if err != nil {
		log.Fatal(err)
	}
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	// soporte a comandos: install/start/stop/uninstall/statun
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "install":
			err = s.Install()
		case "uninstall":
			err = s.Uninstall()
		case "start":
			err = service.Control(s, "start")
		case "stop":
			err = service.Control(s, "stop")
		case "restart":
			err = service.Control(s, "restart")
		case "status":
			status, err := s.Status()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("estatus:", status)
			return
		default:
			fmt.Println("Uso: [install | uninstall | start | stop | restart | status]")
			return
		}
		if err != nil {
			logger.Error(err)
		} else {
			logger.Infof("Comando %s ejecutado correctamente", cmd)
			return
		}
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
