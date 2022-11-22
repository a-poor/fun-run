package funrun

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Manager struct {
	conf   *Conf
	cancel context.CancelFunc
	cmds   []*Command
	lock   sync.RWMutex
	wout   io.Writer
	werr   io.Writer
}

func NewManager(conf *Conf) *Manager {
	return &Manager{
		conf: conf,
		wout: &SyncWriter{Writer: os.Stdout},
		werr: &SyncWriter{Writer: os.Stderr},
	}
}

func (m *Manager) SetOutputs(wout, werr io.Writer) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.wout = &SyncWriter{Writer: wout}
	m.werr = &SyncWriter{Writer: werr}
}

func (m *Manager) createCmds() []*Command {
	m.lock.Lock()
	defer m.lock.Unlock()

	nw := m.conf.maxNameLength()

	cmds := make([]*Command, len(m.conf.Procs))
	for i, proc := range m.conf.Procs {
		// Create the command...
		cmd := NewCommand(proc)

		// Set the outputs...
		wout := NewPrefixWriter(
			proc.Name,
			"stdout",
			nw,
			i,
			m.wout,
		)
		werr := NewPrefixWriter(
			proc.Name,
			"stderr",
			nw,
			i,
			m.werr,
		)
		cmd.SetOutputs(wout, werr)

		// Store the command
		cmds[i] = cmd
	}

	return cmds
}

func (m *Manager) Shutdown() {
	m.Cancel()
}

func (m *Manager) Run(ctx context.Context) error {
	// Create the parent context
	ctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel

	// Check for interrupts
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		select {
		case <-sigs:
			fmt.Fprintln(m.wout, "Received signal, shutting down...")
			cancel()
		case <-ctx.Done():
		}
		done <- true
	}()

	// Create the commands
	m.cmds = m.createCmds()

	// Run the commands
	var wg sync.WaitGroup
	for _, cmd := range m.cmds {
		wg.Add(1)
		go func(cmd *Command) {
			defer wg.Done()
			cmd.Run(ctx)
		}(cmd)
	}

	// Wait for the commands to finish
	wg.Wait()

	// Wait for the context to finish
	cancel()
	<-done

	// Return the error
	return m.Error()
}

func (m *Manager) setCancel(c context.CancelFunc) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.cancel = c
}

func (m *Manager) Cancel() {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.cancel != nil {
		m.cancel()
	}
}

func (m *Manager) Error() error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	me := MapError{}
	for _, cmd := range m.cmds {
		if err := cmd.Error(); err != nil {
			me[cmd.Name()] = err
		}
	}

	if len(me) == 0 {
		return nil
	}
	return me
}
