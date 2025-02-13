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
	keybd_event            = user32.NewProc("keybd_event")
)

const (
	VK_MENU = 0x12 // Código da tecla Alt
	VK_TAB  = 0x09 // Código da tecla Tab
	KEYEVENTF_KEYUP = 0x0002 // Indica que a tecla foi solta
)

// LASTINPUTINFO struct usada pela API do Windows
type LASTINPUTINFO struct {
	CbSize uint32
	DwTime uint32
}

// getIdleTime retorna o tempo de inatividade do usuário
func getIdleTime() time.Duration {
	var lastInput LASTINPUTINFO
	lastInput.CbSize = uint32(unsafe.Sizeof(lastInput))

	if ret, _, _ := getLastInputInfo.Call(uintptr(unsafe.Pointer(&lastInput))); ret == 0 {
		return 0
	}

	tickCount, _, _ := getTickCount64.Call()
	currentTime := uint32(tickCount)

	idleTimeMs := currentTime - lastInput.DwTime
	return time.Duration(idleTimeMs) * time.Millisecond
}

// simulateAltTab executa a combinação Alt+Tab
func simulateAltTab() {
	fmt.Println("Executando Alt+Tab para resetar o tempo de inatividade...")

	keybd_event.Call(uintptr(VK_MENU), 0, 0, 0)
	time.Sleep(50 * time.Millisecond)

	keybd_event.Call(uintptr(VK_TAB), 0, 0, 0)
	time.Sleep(50 * time.Millisecond)

	keybd_event.Call(uintptr(VK_TAB), 0, KEYEVENTF_KEYUP, 0)
	time.Sleep(50 * time.Millisecond)

	keybd_event.Call(uintptr(VK_MENU), 0, KEYEVENTF_KEYUP, 0)
}

func main() {
	teamsIdleLimit := 300 * time.Second // Tempo para o Teams marcar como ausente (5 minutos)
	altTabTrigger := 240 * time.Second  // Tempo para ativar Alt+Tab (4 minutos)

	for {
		idleTime := getIdleTime()
		timeLeft := teamsIdleLimit - idleTime

		if timeLeft < 0 {
			timeLeft = 0
		}

		if idleTime >= altTabTrigger {
			simulateAltTab()
			time.Sleep(2 * time.Second)
			fmt.Printf("1")
		}

		time.Sleep(1 * time.Second)
	}
}
