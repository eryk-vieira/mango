package internal

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/eryk-vieira/mango/internal/types"
)

func Run(settings *types.Settings) {
	isOpen := raw_connect("localhost", settings.Server.Port)

	if isOpen {
		log.Fatalf("Port %s already in use", settings.Server.Port)

		return
	}

	fmt.Println(fmt.Sprintf("Runnning server on port: %s", settings.Server.Port))

	cmd := exec.Command("./server")
	cmd.Dir = ".dist/"
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	cmd.Run()
}

func raw_connect(host string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)

	if err != nil {
		return false
	}

	if conn != nil {
		conn.Close()

		return true
	}

	return false
}
