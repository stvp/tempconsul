package tempconsul

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var (
	StartupSuccess = "cluster leadership acquired"
	StartupTimeout = 5 * time.Second
)

type Server struct {
	cmd *exec.Cmd
}

func (s *Server) Start() error {
	if s.cmd != nil {
		return fmt.Errorf("consul has already been started")
	}

	// Build Cmd
	s.cmd = exec.Command("consul", "agent", "-dev", "-bind", "127.0.0.1")
	serverStdout, err := s.cmd.StdoutPipe()
	if err != nil {
		s.cmd = nil
		return err
	}

	// Start up the server and wait for success log line
	if err = s.cmd.Start(); err != nil {
		return err
	}
	if err = s.waitForSuccessfulStartup(serverStdout); err != nil {
		s.Term()
		return err
	}

	return nil
}

func (s *Server) Interrupt() (err error) {
	if s.cmd == nil {
		return fmt.Errorf("consul is not running")
	}

	s.cmd.Process.Signal(syscall.SIGINT)
	_, err = s.cmd.Process.Wait()
	if err != nil {
		return err
	}

	s.cmd = nil
	return nil
}

func (s *Server) Term() (err error) {
	if s.cmd == nil {
		return fmt.Errorf("consul is not running")
	}

	s.cmd.Process.Signal(syscall.SIGTERM)
	_, err = s.cmd.Process.Wait()
	if err != nil {
		return err
	}

	s.cmd = nil
	return nil
}

// TODO This is gross and weird. Should fix it someday...
func (s *Server) waitForSuccessfulStartup(r io.ReadCloser) error {
	scanner := bufio.NewScanner(r)
	line := ""

	success := make(chan bool, 1)
	failure := make(chan bool, 1)
	stopWaiting := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-stopWaiting: // Timeout
				return
			default:
				if scanner.Scan() {
					line = scanner.Text()
					if strings.Contains(line, StartupSuccess) {
						success <- true
						return
					}
				} else {
					failure <- true
					return
				}
			}
		}
	}()

	select {
	case <-success:
		return nil
	case <-failure:
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("couldn't read consul's stdout: %s", err.Error())
		} else {
			return fmt.Errorf("consul failed to start up: %s", line)
		}
	case <-time.After(StartupTimeout):
		stopWaiting <- true
		return fmt.Errorf("timed-out waiting for consul to start up successfully: %s", line)
	}

}
