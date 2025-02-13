package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                 = syscall.NewLazyDLL("user32.dll")
	getLastInputInfo       = user32.NewProc("GetLastInputInfo")
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	getTickCount64         = kernel32.NewProc("GetTickCount64")
)

// LASTINPUTINFO struct usada pela API do Windows
type LASTINPUTINFO struct {
	CbSize uint32
	DwTime uint32
}

func getIdleTime() time.Duration {
	var lastInput LASTINPUTINFO
	lastInput.CbSize = uint32(unsafe.Sizeof(lastInput))

	// Obtém o tempo da última entrada do usuário
	if ret, _, _ := getLastInputInfo.Call(uintptr(unsafe.Pointer(&lastInput))); ret == 0 {
		return 0
	}

	// Obtém o tempo atual do sistema
	tickCount, _, _ := getTickCount64.Call()
	currentTime := uint32(tickCount)

	// Calcula o tempo ocioso
	idleTimeMs := currentTime - lastInput.DwTime
	return time.Duration(idleTimeMs) * time.Millisecond
}

func main() {
	teamsIdleLimit := 300 * time.Second // Exemplo: O Teams marca como "Ausente" após 5 minutos (ajuste conforme necessário)

	for {
		idleTime := getIdleTime()
		timeLeft := teamsIdleLimit - idleTime

		if timeLeft < 0 {
			timeLeft = 0
		}

		fmt.Printf("Tempo ocioso: %v | Tempo restante para 'Ausente' no Teams: %v\n", idleTime, timeLeft)
		time.Sleep(1 * time.Second)
	}
}
