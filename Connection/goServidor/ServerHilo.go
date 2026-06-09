package main

import (
	"crypto/sha256"        // Importa el paquete sha256 para el hashing de contraseñas.
	"fmt"                  // este paquete fmt es para la salida de texto.
	"net"                  // el paquete net para operaciones de red.
	"paquetes/resources"   // el paquete resources que contiene funciones auxiliares.
	"strconv"              // el paquete strconv para la conversión de tipos.
	"strings"              // el paquete strings para operaciones con cadenas.
	"sync"
)

func main() {

	fmt.Println("===================================")
	fmt.Println("|| SERVIDOR OPERTATIVOS J&J 2024.1 ||")
	fmt.Println("===================================")

	puerto := resources.LeerPuerto()   // Lee el puerto desde la configuración.
	tcpAddress, _ := net.ResolveTCPAddr("tcp4", ":"+puerto)   // Resuelve el puerto en una dirección TCP.

	socketServer, _ := net.ListenTCP("tcp", tcpAddress)   // Crea un listener TCP que escucha en la dirección y puerto especificados.

	socketS, _ := socketServer.Accept()   // Acepta una conexión entrante.

	socketString := socketS.RemoteAddr().String()   // Obtiene la dirección remota del cliente como cadena.

	ipremota := strings.Split(socketString, ":")   // Divide la dirección en IP y puerto.

	ip := resources.LeerIP()   // Lee la IP permitida desde la configuración.

	if ip != ipremota[0] {   // Verifica si la IP del cliente coincide con la permitida.
		fmt.Println("IP NO PERMITIDA")   // Si no coincide, muestra un mensaje y termina el programa.
		return
	}

	fmt.Println("Server# Se ha conectado [", socketS.RemoteAddr(), "]")   // Imprime un mensaje indicando que un cliente se ha conectado.
	usersBD := resources.LeerUserBSD()   // Lee los usuarios y contraseñas permitidos desde el resources.
	inicio := false   // Variable para indicar si el inicio de sesión fue exitoso.
	intentos := resources.LeerIntentos()   // Lee el número de intentos permitidos desde la configuración.

	for i := intentos; i > 0; i-- {

		usuario := resources.LeerDatos(&socketS)   // Lee el nombre de usuario del cliente.
		usuarioCam := strings.Fields(usuario)   // Divide la cadena de usuario en campos.
		constrasena := resources.LeerDatos(&socketS)   // Lee la contraseña del cliente.
		constrasenaCam := strings.Fields(constrasena)   // Divide la cadena de contraseña en campos.

		huser := sha256.Sum256([]byte(constrasenaCam[0]))   // Calcula el hash SHA-256 de la contraseña.
		hxuser := fmt.Sprintf("%x", huser)   // Convierte el hash a una cadena hexadecimal.

		for i := 0; i <= 2; i++ {   // Bucle para comparar el usuario y contraseña hash con los almacenados.
			if (usuarioCam[0] + ":" + hxuser) == usersBD[i] {
				inicio = true   // Si hay coincidencia, marca el inicio como exitoso.
			}
		}

		if inicio {
			resources.EnviarRespuesta(&socketS, "V:0")   // Envía una respuesta de éxito al cliente.
			break   // Sale del bucle si el inicio es exitoso.
		}
		resources.EnviarRespuesta(&socketS, "F:"+strconv.Itoa(i-1))   // Envía una respuesta de fallo al cliente con los intentos restantes.
	}
	if inicio {

		timeR := resources.LeerDatos(&socketS)   // Lee el tiempo de reporte desde el cliente.

		timeVec := strings.Fields(timeR)   // Divide la cadena de tiempo en campos.

		time, _ := strconv.Atoi(timeVec[0])   // Convierte el tiempo a un entero.

		var wg sync.WaitGroup   // Crea un grupo de espera para sincronización de goroutines.

		wg.Add(1)   // Añade una goroutine al grupo de espera.

		go resources.EnviaReporte(&socketS, time)   // Inicia una goroutine para enviar reportes periódicos.

		go resources.RecibeMensaje(&socketS, &wg)   // Inicia una goroutine para recibir mensajes del cliente.

		wg.Wait()   // Espera a que las goroutines terminen.

		socketS.Close()   // Cierra la conexión con el cliente.

	} else {
		fmt.Println("numero de intentos acabados")   // Imprime un mensaje si se agotaron los intentos de inicio de sesión.
	}

	fmt.Println("============================================")
	fmt.Println("|| GRACIAS POR USAR SERVIDOR OPERATIVOS J&J ||")
	fmt.Println("============================================")

}
