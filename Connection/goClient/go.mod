//Define nombre del modulo
module paquete

//Version del go que se necesita para el modulo
go 1.20
//Lista las dependencias directas del proyecto y sus versiones.
require (
	github.com/gdamore/tcell/v2 v2.7.0 //Utilizada para manipulacion de terminales y eventos de entrada
	github.com/rivo/tview v0.0.0-20231206124440-5f078138442e //Construir interfaces de usuario en terminales en Go, con widgets
)
//Lista las dependencias indirectas del proyecto y sus versiones
require (
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.4.3 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/term v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
