package funcionesArreglos

import (
	"bytes"
	"fmt"
	"modulo/config"
	"modulo/structs"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"math/rand"

	"github.com/go-mail/mail"
)

func NormalicarCedulas(cedula string) (int64, error) {
	cedula = strings.TrimSpace(cedula)
	cedulaSin__Espacios := strings.ReplaceAll(cedula, ".", "")
	cedula = strings.ReplaceAll(cedulaSin__Espacios, ",", "")
	cedula__formateada, err := strconv.ParseInt(cedula, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error al normaliza cedula %v:%v", cedula, err)
	}
	return cedula__formateada, nil
}

func ConvertirFechas(fecha string) (time.Time, error) {
	formato__fecha := "01-02-06"
	return time.Parse(formato__fecha, fecha)
}

func CumpleañosActuales(
	Nombre__funcionarios map[int64]string,
	Informacion__funcionarios map[int64]structs.DatosFuncionarios) []int64 {
	fecha__actual := time.Now()
	var cumpleaños []int64
	for doc, info := range Informacion__funcionarios {
		if info.FechaNacimiento.Month() == fecha__actual.Month() &&
			info.FechaNacimiento.Day() == fecha__actual.Day() {
			cumpleaños = append(cumpleaños, doc)
		}
	}
	return cumpleaños
}

func CorreosSinRegistrar(
	Informacion__funcionarios map[int64]structs.DatosFuncionarios) error {
	for _, inf := range Informacion__funcionarios {
		if strings.TrimSpace(inf.CorreoFuncionario) == "" {
			fmt.Printf("funcionari@ %v %v sin correo electronico\n", inf.NombreFuncionario, inf.ApellidoFuncionario)
			continue
		}
	}
	return nil
}

func EnviarCalendarioCumpleañosMes(
	Informacion__funcionarios map[int64]structs.DatosFuncionarios) {
	var Meses__año = map[time.Month]string{
		time.January: "enero", time.February: "febrero",
		time.March: "marzo", time.April: "abril",
		time.May: "mayo", time.June: "junio",
		time.July: "julio", time.August: "agosto",
		time.September: "septiembre", time.October: "octubre",
		time.November: "noviembre", time.December: "diciembre",
	}
	var mes__titulo string
	// guarda los dias ordenados junto con nombre + apellido
	var diasOrdenados []structs.PlantillaCumpleanos
	// guarda toda la informacion para poder usarla en la plantilla
	var lista__funcionarios []string
	for _, inf := range Informacion__funcionarios {
		mes__titulo = Meses__año[inf.FechaNacimiento.Month()]
		dia__cumpleaños := inf.FechaNacimiento.Day()
		descripcion__plantilla := strconv.Itoa(inf.FechaNacimiento.Day()) + " " +
			inf.NombreFuncionario + " " + inf.ApellidoFuncionario
		// guarda en diasOrdenados los dias y descipcion__plantilla
		diasOrdenados = append(diasOrdenados, structs.PlantillaCumpleanos{
			Dias:        dia__cumpleaños,
			Descripcion: descripcion__plantilla,
		})
	}
	// funcion para organizar los dias
	sort.Slice(diasOrdenados, func(i, j int) bool {
		return diasOrdenados[i].Dias < diasOrdenados[j].Dias
	})
	// guarda en lista__funcionarios dias ordenados + nombre + apellido
	for _, inf := range diasOrdenados {
		lista__funcionarios = append(lista__funcionarios, inf.Descripcion)
	}
	funcionarios := structs.PlantillaCumpleanos{
		Titulo:       "Calendario cumpleaños" + " " + mes__titulo,
		Funcionarios: lista__funcionarios,
	}
	tmlp, err := template.ParseGlob("plantillasMes/*.html")
	if err != nil {
		fmt.Printf("error al analizar archivos: %v", err)
	}
	var informacion__plantilla bytes.Buffer
	if err := tmlp.Execute(&informacion__plantilla, funcionarios); err != nil {
		fmt.Printf("error al ejecutar la plantilla:%v", err)
	}
	datos__correo := mail.NewMessage()
	datos__correo.SetHeader("From", "gygculpleanos@gmail.com")
	datos__correo.SetHeader("To", "cm1094871@gmail.com")
	datos__correo.SetHeader("Subject", "Calendario cumpleaños"+" "+mes__titulo)
	datos__correo.SetBody("text/html", informacion__plantilla.String())
	datos := mail.NewDialer("smtp.gmail.com", 587, "gygculpleanos@gmail.com", "hzit rqsa dpwd vebc")
	if err := datos.DialAndSend(datos__correo); err != nil {
		fmt.Printf("error enviando correo: %v", err)
	} else {
		fmt.Println("correo enviado satisfactoriamente")
	}
}

func NotificarCumpleañosFuncionario(
	Nombre__funcionarios, Apellido__funcionario,
	Descripcion__trabajador map[int64]string,
	Genero__funcionario map[int64]string,
	Edad__trabajador map[int64]int,
	Correo__trabajador map[int64]string,
	Informacion__funcionarios map[int64]structs.DatosFuncionarios,
) {
	dataConfig := config.CargaConfig()
	cumpleaños := CumpleañosActuales(Nombre__funcionarios, Informacion__funcionarios)
	for _, doc := range cumpleaños {
		// Variables globales para usar posteriormente
		var tmlp *template.Template
		var err error
		var ListaPlantillas []*template.Template
		var plantillaElegida *template.Template

		// segun genero, envia plantilla correspondiente
		switch Genero__funcionario[doc] {
		case "MASCULINO":
			if tmlp, err = template.ParseGlob("plantillasHombre/*.html"); err != nil {
				fmt.Printf("error al analizar archivos")
			}
			// Se obtienen todas las plantillas dentro del archivo
			for _, savePlantillas := range tmlp.Templates() {
				savePlantillas.Name()
				ListaPlantillas = append(ListaPlantillas, savePlantillas)
			}
			// plantilla aleatorias
			indice := rand.Intn(len(ListaPlantillas))
			plantillaElegida = ListaPlantillas[indice]
		case "FEMENINO":
			if tmlp, err = template.ParseGlob("plantillasMujer/*.html"); err != nil {
				fmt.Printf("error al analizar archivos")
			}
			// Se obtienen todas las plantillas dentro del archivo
			for _, savePlantillas := range tmlp.Templates() {
				savePlantillas.Name()
				ListaPlantillas = append(ListaPlantillas, savePlantillas)
			}
			// plantilla aleatorias
			indice := rand.Intn(len(ListaPlantillas))
			plantillaElegida = ListaPlantillas[indice]
		}
		// datos a enviar a plantilla
		info__funcionario := structs.PlantillaCumpleanos{
			Titulo:      "Feliz cumpleaños",
			Nombre:      Nombre__funcionarios[doc],
			Descripcion: Descripcion__trabajador[doc],
			Edad:        Edad__trabajador[doc],
		}
		var informacion__plantilla bytes.Buffer
		// ejecucion de plantilla
		if err := tmlp.ExecuteTemplate(&informacion__plantilla, plantillaElegida.Name(), info__funcionario); err != nil {
			fmt.Printf("error al ejecutar la plantilla:%v", err)
		}
		// cargue de informacion a correos
		datos__correo := mail.NewMessage()
		datos__correo.SetHeader("From", dataConfig.BaseEmail)
		datos__correo.SetHeader("To", Correo__trabajador[doc]) // Añadir a Cristian y Deysy como comprobante
		// Llamado funcion correos sin registrar
		CorreosSinRegistrar(Informacion__funcionarios)
		datos__correo.SetHeader("Subject", "RE: Feliz cumpleaños"+" "+Nombre__funcionarios[doc])
		// Contenido principal del mensaje
		datos__correo.SetBody("text/html", informacion__plantilla.String())
		// dialer SMTP configurado para poder enviar correos electronicos
		datos := mail.NewDialer("smtp.gmail.com", 587, dataConfig.BaseEmail, dataConfig.PasswordBaseEmail)
		// Envio correo con la informacion correspondiente
		if err := datos.DialAndSend(datos__correo); err != nil {
			fmt.Printf("error enviando correo:%v", err)
		} else {
			fmt.Println("correo enviado satisfactoriamente")
		}
	}
}

func AñosEnEmpresa(fecha__ingreso time.Time) int {
	fecha__actual := time.Now()
	años := fecha__actual.Year() - fecha__ingreso.Year()

	if fecha__actual.Month() < fecha__ingreso.Month() ||
		fecha__actual.Month() == fecha__ingreso.Month() && fecha__actual.Day() < fecha__ingreso.Day() {
		años--
	}
	return años
}

func NotificarAniversarioFuncionario(
	Informacion__aniversarios map[int64]structs.DatosAniversarios,
) {
	dataConfig := config.CargaConfig()
	for _, inf := range Informacion__aniversarios {
		fecha__actual := time.Now()

		if fecha__actual.Month() == inf.FechaIngreso.Month() &&
			fecha__actual.Day() == inf.FechaIngreso.Day() {

			// Variables globales para usar posteriormente
			var tmlpAniversarios *template.Template
			var err error
			var ListaPlantillasAniverarios []*template.Template
			var plantillaElegidaAniversarios *template.Template

			// segun genero, envia plantilla correspondiente
			switch inf.GeneroAniversarios {
			case "MASCULINO":
				if tmlpAniversarios, err = template.ParseGlob("plantillasHombreAniversarios/*.html"); err != nil {
					fmt.Printf("error al analizar archivos")
				}
				// Se obtienen todas las plantillas dentro del archivo
				for _, savePlantillas := range tmlpAniversarios.Templates() {
					savePlantillas.Name()
					ListaPlantillasAniverarios = append(ListaPlantillasAniverarios, savePlantillas)
				}
				// plantilla aleatorias
				indice := rand.Intn(len(ListaPlantillasAniverarios))
				plantillaElegidaAniversarios = ListaPlantillasAniverarios[indice]
			case "FEMENINO":
				if tmlpAniversarios, err = template.ParseGlob("plantillasMujerAniversarios/*.html"); err != nil {
					fmt.Printf("error al analizar archivos")
				}
				// Se obtienen todas las plantillas dentro del archivo
				for _, savePlantillas := range tmlpAniversarios.Templates() {
					savePlantillas.Name()
					ListaPlantillasAniverarios = append(ListaPlantillasAniverarios, savePlantillas)
				}
				// plantilla aleatorias
				indice := rand.Intn(len(ListaPlantillasAniverarios))
				plantillaElegidaAniversarios = ListaPlantillasAniverarios[indice]
			}
			// datos a enviar a plantilla
			info__aniversario := structs.PlantillaAniversarios{
				Titulo:      "Feliz aniversario",
				Nombre:      inf.NombreAniversario,
				Descripcion: inf.Descripcion,
				Tiempo:      AñosEnEmpresa(inf.FechaIngreso),
			}

			var informacion__plantilla__aniversarios bytes.Buffer
			// ejecucion de plantilla
			if err := tmlpAniversarios.ExecuteTemplate(&informacion__plantilla__aniversarios, plantillaElegidaAniversarios.Name(), info__aniversario); err != nil {
				fmt.Printf("error al ejecutar la plantilla:%v", err)
			}

			// cargue de informacion a correos
			datos__correo := mail.NewMessage()
			datos__correo.SetHeader("From", dataConfig.BaseEmail)
			datos__correo.SetHeader("To", inf.CorreoAniversario)
			datos__correo.SetHeader("Subject", "RE: Feliz aniversario"+" "+inf.NombreAniversario+""+
				inf.ApellidoAniversario)
			datos__correo.SetBody("text/html", informacion__plantilla__aniversarios.String())
			// dialer SMTP configurado para poder enviar correos electronicos
			datos := mail.NewDialer("smtp.gmail.com", 587, dataConfig.BaseEmail, dataConfig.PasswordBaseEmail)
			// Envio correo con la informacion correspondiente
			if err := datos.DialAndSend(datos__correo); err != nil {
				fmt.Printf("error enviando correo:%v", err)
			} else {
				fmt.Println("correo enviado satisfactoriamente")
			}
		}
	}
}
