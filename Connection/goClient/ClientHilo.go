package main

import (
	"bufio"
	"fmt"
	"net"
	"paquete/resources"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var comando string   //Variable global para almacenar los comandos ingresados por el usuario.
var socketC net.Conn //Variable global para mantener la conexión con el servidor.
func main() {

	fmt.Println("***********************************")
	fmt.Println("|| CLIENTE J&J* OPERTATIVOS 2024.1 ||")
	fmt.Println("***********************************")

	var tcpAddress *net.TCPAddr
	var ip string
	var puerto string

	for {
		var err error
		fmt.Println("INGRESE LA DIRECCION IP A LA QUE SE VA A CONECTAR")
		fmt.Scanln(&ip)
		fmt.Println("INGRESE EL PUERTO AL QUE SE VA A CONECTAR")
		fmt.Scanln(&puerto)

		if ip == "" || puerto == "" { //Verificacion necesaria para que no reviente por espacios en blanco
			fmt.Println("~~~~~~~~DIRECCION IP O PUERTO NO VALIDO~~~~~~~~")
			continue
		}

		tcpAddress, err = net.ResolveTCPAddr("tcp4", ip+":"+puerto)
		if err == nil {
			break
		} //Verificacion para que no reviente por ip invalida
		fmt.Println("~~~~~~~~DIRECCION IP O PUERTO NO VALIDO~~~~~~~~")
	}

	socketC, _ = net.DialTCP("tcp", nil, tcpAddress) //Establecer conexion TCP
	fmt.Println("Client# Se ha establecido la conexion [", socketC.RemoteAddr(), "]")

	time.Sleep(1 * time.Second) //Para evitar por si el usuario por error se detiene.
	fmt.Println("INGRESE USUARIO")
	resources.EnviaRespuesta(&socketC)
	fmt.Println("INGRESE CONTRASENA")
	resources.EnviaRespuesta(&socketC)

	mensaje := resources.LeeDatos(socketC)    //Lee datos del servidor a traves del socket
	mensajeVec := strings.Split(mensaje, ":") //Divide el mensaje en partes mas pequeñas utilizando el separador ":"
	num := strings.Fields(mensajeVec[1])      //Accede al segundo elemento del Slice
	i, _ := strconv.Atoi(num[0])              //Lo convierte en un entero

	for mensajeVec[0] == "F" && i > 0 {
		fmt.Println("usuario o contrasena incorrecta, intentos: " + mensajeVec[1])
		fmt.Println("INGRESE USUARIO")
		resources.EnviaRespuesta(&socketC)
		fmt.Println("INGRESE CONTRASENA")
		resources.EnviaRespuesta(&socketC)

		mensaje = resources.LeeDatos(socketC)
		mensajeVec = strings.Split(mensaje, ":")
		num = strings.Fields(mensajeVec[1])
		i, _ = strconv.Atoi(num[0])
	}
	if mensajeVec[0] == "V" {

		fmt.Println("INICIO EXITOSO")

		fmt.Println("INGRESE EL TIEMPO DE REPORTE:")
		tiempoReporte := resources.EnviaRespuesta(&socketC)

		fmt.Println("Client# Se ha establecido la conexion [", socketC.RemoteAddr(), tiempoReporte, "]")

		time.Sleep(5 * time.Second)
		var wg sync.WaitGroup //Sincronizar la ejecucion de goroutines

		wg.Add(2) //Espera 2 goroutines

		//Interfaz de usuario
		app := tview.NewApplication()

		outputTable := tview.NewTable(). //Se utiliza para mostrar datos de salida en la interfaz de usuario.
							SetBorders(false).
							ScrollToEnd()

		inputField := tview.NewInputField(). //Se utiliza para que el usuario pueda ingresar comandos en la interfaz de usuario.
							SetLabel("Ingrese el comando: ").
							SetFieldWidth(0).
							SetLabelColor(tcell.ColorBlue)

		lowerBox := tview.NewTextView(). //Se utiliza para mostrar mensajes adicionales en la interfaz de usuario.
							SetLabelWidth(0).
							SetDynamicColors(true).
							SetRegions(true).
							SetChangedFunc(func() {
				app.Draw()
			})

		go recibeReporte(socketC, &wg, outputTable, app, lowerBox) //Goroutine que se encarga de recibir y mostrar datos del servidor en la tabla de salida.
		go enviaMensaje(socketC, &wg, inputField)                  //Goroutine que se encarga de enviar los comandos ingresados por el usuario al servidor.

		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			// Cambia el enfoque a la tabla cuando se presiona la tecla Tab
			if event.Key() == tcell.KeyTab {
				app.SetFocus(outputTable)
			}
			// Cambia el enfoque al campo de entrada cuando se presiona la tecla Esc
			if event.Key() == tcell.KeyEsc {
				app.SetFocus(inputField)
			}
			return event
		})

		flex := tview.NewFlex(). //Layout flexible ( Se agregan la tabla de salida, el campo de entrada y el cuadro de texto al layout.)
						SetDirection(tview.FlexRow).
						AddItem(outputTable, 0, 1, false).
						AddItem(inputField, 1, 1, true).
						AddItem(lowerBox, 3, 1, false)

		if err := app.SetRoot(flex, true).SetFocus(inputField).Run(); err != nil { //Lanza excepciones
			panic(err)
		}

		wg.Wait() //Espera que todas las goroutines terminen

		socketC.Close() //Se cierra la conexion con el servidor

		fmt.Println("===========================================")
		fmt.Println("|| GRACIAS POR USAR CLIENTE OPERATIVOS J&J ||")
		fmt.Println("===========================================")

	} else {
		socketC.Close()
		fmt.Println("numero de intentos acabados")
	}

}

func recibeReporte(socketC net.Conn, wg *sync.WaitGroup, outputTable *tview.Table, app *tview.Application, lowerBox *tview.TextView) {
	for {

		m := bufio.NewReader(socketC) //Leer datos de la conexion de red
		for {

			line, err := m.ReadString('\n') //Recibe un mensaje, en este caso string

			if err != nil {
				break
			}

			if strings.HasPrefix(line, "Memoria:") {
				app.QueueUpdateDraw(func() {
					lowerBox.Clear()
					fmt.Fprintf(lowerBox, line)
				})
				continue
			}
			if line == "" {
				continue
			} //Se hace para verificar que la linea no este vacia y asi no la imprima
			row := outputTable.GetRowCount()
			outputTable.SetCell(row, 0, tview.NewTableCell(line))
			outputTable.ScrollToEnd()
		}
		wg.Done()  //Se marca como completa la goroutine
		app.Stop() //Cuando todas las goroutines se hallan detenido
		break

	}
}

func enviaMensaje(socketC net.Conn, wg *sync.WaitGroup, inputField *tview.InputField) {
	//Funcion de acción que se activa cuando se completa la entrada en el campo de entrada
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter { //Verifica si la tecla es ENTER
			comando = inputField.GetText() //Obtiene el texto ingresado por el usuario en el campo de entrada
			inputField.SetText("")         //Limpia el campo de entrada

			if comando != "" { //Verifica que no este vacio el mensaje;
				enviar := bufio.NewWriter(socketC) //Permitirá escribir datos en el socket de manera eficiente.
				enviar.WriteString(comando + "\n") //Prepara el comando para ser enviado al servidor.
				enviar.Flush()                     //Envia el comando al servidor a través del socket.

				if comando == "bye" {
					fmt.Println("Se entro en el cierre")
					wg.Done() //Se marca como completa la goroutine
					return
				}

				time.Sleep(10 * time.Millisecond)
			}

		}
	})
}
