package main

import (
	"fmt"
	"launch/config"
	"launch/terminal"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	cfg.Env.Set()

	names := []config.ProcName{}
	if len(os.Args) > 1 {
		for _, name := range os.Args[1] {
			names = append(names, config.ProcName(string(name)))
		}
	}

	terminal := terminal.NewTerminal()
	defer terminal.Close()

	processes := make([]*exec.Cmd, 0, len(names))
	for _, name := range names {
		process, ok := cfg.Processes[name]
		if ok {
			cmd, err := terminal.StartProcess(process)
			if err != nil {
				fmt.Printf("failed to start %q:\n\t%s %s\n%s\n", name, process.Path, strings.Join(process.Args, " "), err.Error())
			}

			processes = append(processes, cmd)
		}
	}

	done := terminal.ProcessInput()

	select {
	case <-done:
		kill(processes)
	case <-terminated(processes):
		kill(processes)
	case <-finished(processes):
	}
}

func kill(processes []*exec.Cmd) {
	for _, process := range processes {
		process.Process.Kill()
	}
}

func terminated(processes []*exec.Cmd) chan struct{} {
	c := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		c <- struct{}{}
	}()
	return c
}

func finished(processes []*exec.Cmd) chan struct{} {
	c := make(chan struct{})
	go func() {
		wait(processes)
		c <- struct{}{}
	}()
	return c
}

func wait(processes []*exec.Cmd) {
	var wg sync.WaitGroup

	if len(processes) > 0 {
		wg.Add(len(processes))
	} else {
		wg.Add(1) // wait infinitely
	}

	for _, process := range processes {
		go func(p *exec.Cmd) {
			p.Wait()
			wg.Done()
		}(process)
	}

	wg.Wait()
}
