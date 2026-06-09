package resources

import (
	"bufio"       // Importa el paquete bufio para leer y escribir datos en buffers.
	"fmt"         // Importa el paquete fmt para imprimir mensajes.
	"net"         // Importa el paquete net para operaciones de red.
	"os"          // Importa el paquete os para operaciones de archivo y sistema operativo.
	"os/exec"     // Importa el paquete os/exec para ejecutar comandos del sistema.
	"strconv"     // Importa el paquete strconv para conversiones entre cadenas y enteros.
	"strings"     // Importa el paquete strings para operaciones con cadenas.
	"sync"        // Importa el paquete sync para sincronización entre goroutines.
	"time"        // Importa el paquete time para manejo de tiempo y duración.
)

// Una función sin argumentos devuelve un string, que es el puerto.
func LeerPuerto() (numPuerto string) {
	var puerto string // Declara una variable para almacenar el puerto.

	// Da la ubicación del archivo y lo abre.
	file, _ := os.Open("./resources/client.conf")

	// Se usa para leer el archivo línea por línea.
	scanner := bufio.NewScanner(file)

	// Hace un ciclo buscando línea por línea cuando la palabra sea igual a "puerto".
	for scanner.Scan() {
		linea := scanner.Text() // Lee la línea actual.
		if linea == "puerto" {
			scanner.Scan()         // Salta una línea abajo para extraer el puerto.
			puerto = scanner.Text() // Guarda el puerto en la variable.
			break
		}
	}
	return puerto // Retorna el puerto.
}

// Función que retorna la IP de tipo string.
func LeerIP() (numip string) {
	var ip string // Declara una variable para almacenar la IP.

	file, _ := os.Open("./resources/client.conf") // Abre el archivo de configuración.

	scanner := bufio.NewScanner(file) // Crea un escáner para leer el archivo línea por línea.

	for scanner.Scan() { // Lee cada línea del archivo.
		linea := scanner.Text() // Obtiene la línea actual.
		if linea == "ip" {       // Si la línea es "ip":
			scanner.Scan()  // Salta a la siguiente línea.
			ip = scanner.Text() // Guarda la IP en la variable.
			break
		}
	}
	return ip // Retorna la IP.
}

// Función que lee y retorna el tiempo de reporte.
func LeerTiempoReporte() (tiempo int) {
	var time int // Declara una variable para almacenar el tiempo.

	file, _ := os.Open("./resources/client.conf") // Abre el archivo de configuración.

	scanner := bufio.NewScanner(file) // Crea un escáner para leer el archivo línea por línea.

	for scanner.Scan() { // Lee cada línea del archivo.
		linea := scanner.Text() // Obtiene la línea actual.
		if linea == "reporte" { // Si la línea es "reporte":
			scanner.Scan()              // Salta a la siguiente línea.
			time, _ = strconv.Atoi(scanner.Text()) // Convierte la línea a un entero y la guarda en la variable.
			break
		}
	}
	return time // Retorna el tiempo.
}

// Función que lee y retorna el número de intentos permitidos.
func LeerIntentos() (attemps int) {
	var intentos int // Declara una variable para almacenar el número de intentos.

	file, _ := os.Open("./resources/client.conf") // Abre el archivo de configuración.

	scanner := bufio.NewScanner(file) // Crea un escáner para leer el archivo línea por línea.

	for scanner.Scan() { // Lee cada línea del archivo.
		linea := scanner.Text() // Obtiene la línea actual.
		if linea == "intentos" { // Si la línea es "intentos":
			scanner.Scan()                  // Salta a la siguiente línea.
			intentos, _ = strconv.Atoi(scanner.Text()) // Convierte la línea a un entero y la guarda en la variable.
			break
		}
	}
	return intentos // Retorna el número de intentos.
}

// Función que lee y retorna una lista de usuarios desde un archivo.
func LeerUserBSD() (vecUser [3]string) {
	var usersBd [3]string // Declara un arreglo para almacenar los usuarios.
	i := 0 // Inicializa un índice para el arreglo.

	file, _ := os.Open("./resources/usersBsd.txt") // Abre el archivo de usuarios.

	scanner := bufio.NewScanner(file) // Crea un escáner para leer el archivo línea por línea.

	scanner.Scan() // Lee la primera línea (probablemente un encabezado o vacío).
	for scanner.Scan() { // Lee cada línea del archivo.
		linea := scanner.Text() // Obtiene la línea actual.
		usersBd[i] = linea // Guarda la línea en el arreglo de usuarios.
		i++ // Incrementa el índice.
	}
	return usersBd // Retorna el arreglo de usuarios.
}

// Función que lee datos desde una conexión y los retorna como string.
func LeerDatos(socketS *net.Conn) (user string) {
	m, _ := bufio.NewReader(*socketS).ReadString('\n') // Lee una línea desde la conexión.
	return m // Retorna la línea leída.
}

// Función que envía una respuesta a través de una conexión.
func EnviarRespuesta(socketS *net.Conn, respuesta string) {
	env := bufio.NewWriter(*socketS) // Crea un escritor de búfer para la conexión.
	env.WriteString(respuesta + "\n") // Escribe la respuesta en el búfer.
	env.Flush() // Envía el contenido del búfer.
}

// Función que recibe mensajes desde una conexión y ejecuta comandos.
func RecibeMensaje(socketS *net.Conn, wg *sync.WaitGroup) {
	for {
		m, _ := bufio.NewReader(*socketS).ReadString('\n') // Lee una línea desde la conexión.
		fmt.Println("Server# Comando recibido [", m, "]") // Imprime el comando recibido.
		fmt.Println("Ejecutando el comando", m) // Imprime el mensaje de ejecución del comando.

		if m == "bye\n" { // Si el comando es "bye":
			wg.Done() // Indica que la goroutine ha terminado.
			break // Sale del bucle.
		}

		datoIn := strings.TrimRight(m, "\r\n") // Elimina los caracteres de nueva línea y retorno de carro.
		array_datosIn := strings.Fields(datoIn) // Divide la línea en palabras.
		shell := exec.Command(array_datosIn[0], array_datosIn[1:]...) // Crea un comando para ejecutar.
		stdout, _ := shell.Output() // Ejecuta el comando y obtiene la salida.

		command(string(stdout)) // Escribe la salida del comando en un archivo.
		comando := leerComando() // Lee el comando desde el archivo.

		if comando == "" { // Si el comando no existe:
			EnviarRespuesta(socketS, "~~~~~~COMANDO NO ENCONTRADO~~~~~~") // Envía un mensaje de error.
			fmt.Println("Shelloper# \n", string(stdout)) // Imprime la salida del comando.
			continue // Continúa con la siguiente iteración del bucle.
		}

		EnviarRespuesta(socketS, comando) // Envía el comando leído como respuesta.
		fmt.Println("Shelloper# \n", string(stdout)) // Imprime la salida del comando.
	}
}

// Función que lee el contenido de un archivo y lo retorna como string.
func leerComando() (comando string) {
	mensaje := "" // Inicializa una cadena para almacenar el mensaje.

	file, _ := os.Open("comando.txt") // Abre el archivo de comandos.

	scanner := bufio.NewScanner(file) // Crea un escáner para leer el archivo línea por línea.

	for scanner.Scan() { // Lee cada línea del archivo.
		linea := scanner.Text() // Obtiene la línea actual.
		mensaje = mensaje + linea + "\n" // Añade la línea al mensaje.
	}
	return mensaje // Retorna el mensaje completo.
}

// Función que escribe un comando en un archivo.
func command(out string) {
	escribir := os.WriteFile("comando.txt", []byte(out), 0644) // Escribe la salida en el archivo comando.txt.
	if escribir != nil { // Si ocurre un error al escribir:
		// Manejo de error, si es necesario.
	}
}

// Función que lee y calcula el porcentaje de uso de memoria.
func LeerMemoria() (memoria string) {
	var mensaje string // Inicializa una cadena para almacenar el mensaje.

	file, _ := os.Open("memoria.txt") // Abre el archivo de memoria.

	scanner := bufio.NewScanner(file) // Crea un escáner para leer el archivo línea por línea.

	for scanner.Scan() { // Lee cada línea del archivo.
		linea := scanner.Text() // Obtiene la línea actual.
		mensaje = mensaje + linea + "\n" // Añade la línea al mensaje.
	}
	lineas := strings.Split(mensaje, "\n") // Divide el mensaje en líneas.

	var memoriaLinea string // Declara una variable para la línea de memoria.

	for _, i := range lineas { // Recorre cada línea.
		if strings.Contains(i, "Mem:") { // Si la línea contiene "Mem:":
			memoriaLinea = i // Guarda la línea en la variable.
			break
		}
	}

	campos := strings.Fields(memoriaLinea) // Divide la línea en campos.
	memTotal := campos[1] // Obtiene el total de memoria.
	memUsada := campos[2] // Obtiene la memoria usada.
	memTotalInt, _ := strconv.Atoi(memTotal) // Convierte el total de memoria a entero.
	memUsadaInt, _ := strconv.Atoi(memUsada) // Convierte la memoria usada a entero.

	return strconv.Itoa((memUsadaInt * 100) / memTotalInt) // Calcula y retorna el porcentaje de memoria usada.
}

// Función que escribe el uso de memoria en un archivo.
func Mem(out string) {
	escribir := os.WriteFile("memoria.txt", []byte(out), 0644) // Escribe la salida en el archivo memoria.txt.
	if escribir != nil { // Si ocurre un error al escribir:
		// Manejo de error, si es necesario.
	}
}

// Función que lee y calcula el porcentaje de uso del procesador.
func LeeProcesador() (vecMensaje string) {
	var mensaje string // Inicializa una cadena para almacenar el mensaje.

	file, _ := os.Open("cpu.txt") // Abre el archivo de CPU.

	scanner := bufio.NewScanner(file) // Crea un escáner para leer el archivo línea por línea.

	scanner.Scan() // Lee la primera línea.
	mensaje = scanner.Text() // Guarda la línea en la variable.

	campos := strings.Fields(mensaje) // Divide la línea en campos.

	// Convierte los campos a enteros.
	user, _ := strconv.Atoi(campos[1])
	nice, _ := strconv.Atoi(campos[2])
	system, _ := strconv.Atoi(campos[3])
	idle, _ := strconv.Atoi(campos[4])
	iowat, _ := strconv.Atoi(campos[5])

	tiempoActivoCpu := user + nice + system // Calcula el tiempo activo del CPU.
	tiempoTotalCpu := user + nice + system + idle + iowat // Calcula el tiempo total del CPU.
	porcentajeCpu := (float64(tiempoActivoCpu) / float64(tiempoTotalCpu)) * 100 // Calcula el porcentaje de uso del CPU.
	stringCpu := fmt.Sprintf("%.2f", porcentajeCpu) // Formatea el porcentaje como una cadena.

	return stringCpu // Retorna el porcentaje de uso del CPU.
}

// Función que escribe el uso del CPU en un archivo.
func cpu(out string) {
	escribir := os.WriteFile("cpu.txt", []byte(out), 0644) // Escribe la salida en el archivo cpu.txt.
	if escribir != nil { // Si ocurre un error al escribir:
		// Manejo de error, si es necesario.
	}
}

// Función que escribe el uso del disco en un archivo.
func Disk(out string) {
	escribir := os.WriteFile("disco.txt", []byte(out), 0644) // Escribe la salida en el archivo disco.txt.
	if escribir != nil { // Si ocurre un error al escribir:
		// Manejo de error, si es necesario.
	}
}

// Función que lee y retorna el porcentaje de uso del disco.
func LeerDisco() (vecMensaje string) {
	var mensaje string // Inicializa una cadena para almacenar el mensaje.

	file, _ := os.Open("disco.txt") // Abre el archivo de disco.

	scanner := bufio.NewScanner(file) // Crea un escáner para leer el archivo línea por línea.

	for scanner.Scan() { // Lee cada línea del archivo.
		mensaje = scanner.Text() // Guarda la línea en la variable.
	}

	campos := strings.Fields(mensaje) // Divide la línea en campos.

	porcentajeDisk := campos[4] // Obtiene el porcentaje de uso del disco.

	return porcentajeDisk // Retorna el porcentaje de uso del disco.
}

// Función que envía un reporte periódicamente a través de una conexión.
func EnviaReporte(socketS *net.Conn, n int) {
	for {
		time.Sleep(time.Duration(n) * time.Second) // Espera durante n segundos.

		env := bufio.NewWriter(*socketS) // Crea un escritor de búfer para la conexión.

		// Prepara y ejecuta el comando para obtener el uso de memoria.
		datoInmem := strings.TrimRight("free -m", "\r\n")
		array_datosInmem := strings.Fields(datoInmem)
		shellmem := exec.Command(array_datosInmem[0], array_datosInmem[1:]...)
		stdoutmem, _ := shellmem.Output()

		// Prepara y ejecuta el comando para obtener el uso del CPU.
		datoInCpu := strings.TrimRight("cat /proc/stat", "\r\n")
		array_datosInCpu := strings.Fields(datoInCpu)
		shellCpu := exec.Command(array_datosInCpu[0], array_datosInCpu[1:]...)
		stdoutCpu, _ := shellCpu.Output()

		// Prepara y ejecuta el comando para obtener el uso del disco.
		datoInDisk := strings.TrimRight("df -m --total", "\r\n")
		array_datosInDisk := strings.Fields(datoInDisk)
		shellDisk := exec.Command(array_datosInDisk[0], array_datosInDisk[1:]...)
		stdoutDisk, _ := shellDisk.Output()

		// Guarda la salida de los comandos en archivos.
		Mem(string(stdoutmem))
		cpu(string(stdoutCpu))
		Disk(string(stdoutDisk))

		// Lee los datos de los archivos y calcula los porcentajes.
		memoria := LeerMemoria()
		cpu := LeeProcesador()
		disco := LeerDisco()

		// Prepara el mensaje con los porcentajes.
		envString := "Memoria: " + memoria + "%, " + " CPU: " + cpu + "%, " + " Disco:" + disco

		// Envía el mensaje a través de la conexión.
		env.WriteString(envString + "\n")
		env.Flush()
	}
}
