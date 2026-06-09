package resources

import (
	"bufio"
	"net"
	"os"
)

func LeeDatos(socketS net.Conn) (user string) {

	m, _ := bufio.NewReader(socketS).ReadString('\n') //Permite leer datos de esa conexion, hasta que encuentre un salto de linea

	return m
}

func EnviaRespuesta(socketC *net.Conn) (mensaje string) {

	env := bufio.NewWriter(*socketC) //Permite escribir datos en la conexión de manera eficiente.

	lector := bufio.NewReader(os.Stdin) //Permite leer la entrada del usuario desde el teclado

	dato, _ := lector.ReadString('\n') //Hasta que haya un salto de linea

	env.WriteString(dato) //Se escribe en la conexión de red.

	env.Flush() //Asegurar que todos los datos escritos en el buffer se envien a traves de la conexion

	return dato
}
