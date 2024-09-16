package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

func waitForErrChan(c chan error, t time.Duration) bool {
	select {
	case <-time.After(t):
		return false
	case <-c:
		return true
	}
}

func termChild(running bool, cmd *exec.Cmd, ch chan error,
	wait int, out io.Writer, wg *sync.WaitGroup) {
	if nil != wg {
		defer wg.Done()
	}
	if !running {
		return
	}

	var err error

	
	err = Terminate(cmd)

	if nil != err {
		_, e := fmt.Fprintln(out, "Error sending TERM signal: ", err)
		if nil != e {
			log.Println("Error writing log: ", err)
		}
	}

	if !waitForErrChan(ch, time.Duration(wait)*time.Second) {
		_, e := fmt.Fprintln(out, "Process is still running, sending kill signal")
		if nil != e {
			log.Println("Error writing log: ", err)
		}

		
		err = Kill(cmd)

		if nil != err {
			_, e := fmt.Fprintln(out, "Error sending KILL signal: ", err)
			if nil != e {
				log.Println("Error writing log: ", err)
			}
		}
	}
}


func Terminate(cmd *exec.Cmd) error{
	var err error = nil
	if runtime.GOOS == "windows" {
		// Windows-specific code
		if cmd.ProcessState == nil || cmd.ProcessState.Exited() {
			return nil // Process has already finished
		}
		handle := windows.Handle(cmd.Process.Pid)
		
		err = windows.TerminateProcess(handle, 1)
	} else {
		// Unix-specific code
		err = cmd.Process.Signal(syscall.SIGTERM)
	}
	return err
}

func Kill(cmd *exec.Cmd) error{
	var err error = nil
	if runtime.GOOS == "windows" {
		// Windows-specific code
		if cmd.ProcessState == nil || cmd.ProcessState.Exited() {
			return nil // Process has already finished
		}
		handle := windows.Handle(cmd.Process.Pid)
		
		err = windows.TerminateProcess(handle, 1)
	} else {
		// Unix-specific code
		err = cmd.Process.Kill()
	}
	return err
}