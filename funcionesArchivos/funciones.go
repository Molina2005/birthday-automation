package funcionesArchivos

import (
	"errors"
	"fmt"
	"log"
	"modulo/funcionesArreglos"
	"modulo/structs"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kardianos/service"
	"github.com/xuri/excelize/v2"
)

func LeerArchivoFuncionariosgyg(Nombre__archivo string) ([][]string, [][]string, error) {
	archivo__funcionrios, err := excelize.OpenFile(Nombre__archivo)
	if err != nil {
		return nil, nil, fmt.Errorf("error al abrir el archivo de funcionarios: %v", err)
	}
	defer archivo__funcionrios.Close()

	nombreHoja__Funcionarios := archivo__funcionrios.GetSheetName(0)
	nombreHoja__aniversarios := archivo__funcionrios.GetSheetName(1)
	rows__funcionarios, err := archivo__funcionrios.GetRows(nombreHoja__Funcionarios)
	if err != nil {
		return nil, nil, fmt.Errorf("error al leer las filas de funcionariosgyg: %v", err)
	}
	// lectura hoja aniversarios
	rows__aniversarios, err := archivo__funcionrios.GetRows(nombreHoja__aniversarios)
	if err != nil {
		return nil, nil, fmt.Errorf("error al leer las filas de aniversarios: %v", err)
	}

	return rows__funcionarios, rows__aniversarios, nil
}

func ProcesarDatosFuncionarios(rows__funcionarios [][]string) (
	map[int64]struct{},
	map[int64]time.Time,
	map[int64]string,
	map[int64]string,
	map[int64]string,
	map[int64]string,
	map[int64]int,
	map[int64]string,
	map[int64]string,
	map[int64]structs.DatosFuncionarios) {
	Cedula__funcionarios := make(map[int64]struct{})
	FechaDe__nacimiento := make(map[int64]time.Time)
	Nombre__funcionarios := make(map[int64]string)
	Apellido__funcionario := make(map[int64]string)
	Genero__funcionario := make(map[int64]string)
	Lugar__trabajo := make(map[int64]string)
	Edad__trabajador := make(map[int64]int)
	Correo__trabajador := make(map[int64]string)
	Descripcion__trabajador := make(map[int64]string)
	Informacion__funcionarios := make(map[int64]structs.DatosFuncionarios)

	for i, rows := range rows__funcionarios {
		if i > 1 {
			doc, err := funcionesArreglos.NormalicarCedulas(rows[0])
			if err != nil {
				log.Fatal(errors.New("error al normalizar cedula %v: %v"), doc, err)
			}
			Cedula__funcionarios[doc] = struct{}{}

			fecha, err := funcionesArreglos.ConvertirFechas(rows[1])
			if err != nil {
				log.Fatal(errors.New("error la convertir fecha %v: %v"), fecha, err)
			}
			FechaDe__nacimiento[doc] = fecha

			Nombre__funcionarios[doc] = rows[2]
			Apellido__funcionario[doc] = rows[3]
			Genero__funcionario[doc] = rows[4]
			Lugar__trabajo[doc] = rows[5]
			if Edad__convertida, err := strconv.Atoi(rows[6]); err == nil {
				Edad__trabajador[doc] = Edad__convertida
			}
			Correo__trabajador[doc] = rows[7]
			Descripcion__trabajador[doc] = rows[8]

			Informacion__funcionarios[doc] = structs.DatosFuncionarios{
				FechaNacimiento:     FechaDe__nacimiento[doc],
				NombreFuncionario:   Nombre__funcionarios[doc],
				ApellidoFuncionario: Apellido__funcionario[doc],
				CorreoFuncionario:   Correo__trabajador[doc],
			}
		}
	}
	return Cedula__funcionarios, FechaDe__nacimiento,
		Nombre__funcionarios, Apellido__funcionario,
		Genero__funcionario, Lugar__trabajo,
		Edad__trabajador, Correo__trabajador,
		Descripcion__trabajador, Informacion__funcionarios
}

func ProcesarDatosAniversarios(rows__aniversarios [][]string) map[int64]structs.DatosAniversarios {
	Cedula__aniversarios := make(map[int64]struct{})
	Fecha__ingreso := make(map[int64]time.Time)
	Nombre__anivesario := make(map[int64]string)
	Apellido__aniversarios := make(map[int64]string)
	Genero__aniversarios := make(map[int64]string)
	Correo__aniversarios := make(map[int64]string)
	Descripcion__aniversario := make(map[int64]string)
	Informacion__aniversarios := make(map[int64]structs.DatosAniversarios)

	for i, rows := range rows__aniversarios {
		if i > 1 {
			doc, err := funcionesArreglos.NormalicarCedulas(rows[0])
			if err != nil {
				log.Fatalf("error al normalizar cedula %v: %v", doc, err)
			}
			Cedula__aniversarios[doc] = struct{}{}

			fecha, err := funcionesArreglos.ConvertirFechas(rows[1])
			if err != nil {
				log.Fatalf("error la convertir fecha %v: %v", fecha, err)
			}
			Fecha__ingreso[doc] = fecha
			Nombre__anivesario[doc] = rows[2]
			Apellido__aniversarios[doc] = rows[3]
			Genero__aniversarios[doc] = rows[4]
			Correo__aniversarios[doc] = rows[5]
			Descripcion__aniversario[doc] = rows[6]

			Informacion__aniversarios[doc] = structs.DatosAniversarios{
				FechaIngreso:       Fecha__ingreso[doc],
				NombreAniversario:  Nombre__anivesario[doc],
				ApellidoAniversario: Apellido__aniversarios[doc],
				GeneroAniversarios: Genero__aniversarios[doc],
				CorreoAniversario:  Correo__aniversarios[doc],
				Descripcion:        Descripcion__aniversario[doc],
			}
		}
	}
	return Informacion__aniversarios
}

func NotificaCambiosArchivoPrincipal() {
	// Creacion de nuevo observador
	observador, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("error al crear nuevo observador: %v", err)
	}
	defer observador.Close()
	// Se agrega ruta de los excel
	err = observador.Add("C:\\Users\\cristian.gonzalez\\go\\src\\cumpleañosgyg\\archivo\\")
	if err != nil {
		fmt.Printf("error de ruta %v", err)
	}
	// pase observador por parametro a funcion EscuchaEventos
	EscuchaEventos(observador)
}

// Escucha de eventos
func EscuchaEventos(observador *fsnotify.Watcher) {
	for {
		select {
		// recibe eventos de cambio en los archivos excel
		case evento, ok := <-observador.Events:
			if !ok {
				return
			}
			log.Println("evento:", evento)
			// operacion que se realizo en el evento:write,remove,rename,chmod
			if evento.Has(fsnotify.Write) {
				log.Println("archivo modificado:", evento.Name) // nombre del evento
			}
			// canal de errores del observador
		case err, ok := <-observador.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

// Instalacion programa como servicio
// maneja todo lo relacionados con mi servicio: Start, Stop, Run
type Programa struct {
	InfoFuncionarios      map[int64]structs.DatosFuncionarios
	NombreFuncionarios    map[int64]string
	ApellidoFuncionario   map[int64]string
	DescripcionTrabajador map[int64]string
	GeneroFuncionario     map[int64]string
	EdadTrabajador        map[int64]int
	CorreoTrabajador      map[int64]string
}

// se arranca el servicio
func (p *Programa) Start(service.Service) error {
	go p.run()
	return nil
}

// ejecuta continuamente mientras el servidor este activo
func (p *Programa) run() {
	for {
		// funciones con logica necesaria, las cuales se van a ejecutar
		funcionesArreglos.EnviarCalendarioCumpleañosMes(p.InfoFuncionarios)
		funcionesArreglos.NotificarCumpleañosFuncionario(p.NombreFuncionarios,
			p.ApellidoFuncionario, p.DescripcionTrabajador,
			p.GeneroFuncionario, p.EdadTrabajador, p.CorreoTrabajador,
			p.InfoFuncionarios)
		// tiempo para que vuelva a repetir la ejecicion de nuevo
		time.Sleep(20 * time.Second)
	}
}

// sin funcionamiento, se deja solo por que la libreria kardianos obliga a utilizarlo
func (p *Programa) Stop(service.Service) error {
	return nil
}
