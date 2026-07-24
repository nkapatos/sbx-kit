package sbxutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Runner executes sbx subcommands. Tests can substitute a fake.
type Runner struct {
	// LookPath defaults to exec.LookPath("sbx").
	LookPath func() (string, error)
	// Command builds a command; defaults to exec.Command.
	Command func(name string, args ...string) *exec.Cmd
}

func Default() *Runner {
	return &Runner{
		LookPath: func() (string, error) { return exec.LookPath("sbx") },
		Command:  exec.Command,
	}
}

func (r *Runner) require() (string, error) {
	if r.LookPath == nil {
		r.LookPath = func() (string, error) { return exec.LookPath("sbx") }
	}
	if r.Command == nil {
		r.Command = exec.Command
	}
	return r.LookPath()
}

// RunInteractive attaches stdin/stdout/stderr (for sbx run).
func (r *Runner) RunInteractive(args ...string) error {
	bin, err := r.require()
	if err != nil {
		return fmt.Errorf("sbx not found on PATH")
	}
	fmt.Printf("==> sbx %s\n", strings.Join(args, " "))
	cmd := r.Command(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

// RunEnv is like RunInteractive but with extra env.
func (r *Runner) RunEnv(env []string, args ...string) error {
	bin, err := r.require()
	if err != nil {
		return fmt.Errorf("sbx not found on PATH")
	}
	fmt.Printf("==> sbx %s\n", strings.Join(args, " "))
	cmd := r.Command(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	return cmd.Run()
}

// Output runs sbx and returns combined stdout (trimmed).
func (r *Runner) Output(args ...string) (string, error) {
	bin, err := r.require()
	if err != nil {
		return "", fmt.Errorf("sbx not found on PATH")
	}
	cmd := r.Command(bin, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("sbx %s: %s", strings.Join(args, " "), msg)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// Exec streams a command inside the sandbox to the terminal.
func (r *Runner) Exec(sandbox string, command ...string) error {
	args := append([]string{"exec", sandbox, "--"}, command...)
	return r.RunInteractive(args...)
}

// ExecVisible is an alias for Exec.
func (r *Runner) ExecVisible(sandbox string, command ...string) error {
	return r.Exec(sandbox, command...)
}

// Cp copies between host and sandbox. One side must be "name:path".
func (r *Runner) Cp(src, dst string) error {
	return r.RunInteractive("cp", src, dst)
}

// Rm removes a sandbox.
func (r *Runner) Rm(name string, force bool) error {
	args := []string{"rm", name}
	if force {
		args = []string{"rm", "--force", name}
	}
	return r.RunInteractive(args...)
}

// Sandbox is one row from sbx ls.
type Sandbox struct {
	Name      string
	Agent     string
	Status    string
	Workspace string
}

// Ls parses `sbx ls` into rows (best-effort whitespace split).
func (r *Runner) Ls() ([]Sandbox, error) {
	out, err := r.Output("ls")
	if err != nil {
		return nil, err
	}
	var rows []Sandbox
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(strings.ToUpper(line), "SANDBOX") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		s := Sandbox{Name: fields[0]}
		if len(fields) > 1 {
			s.Agent = fields[1]
		}
		if len(fields) > 2 {
			s.Status = fields[2]
		}
		if len(fields) > 3 {
			// PORTS may be "-" or mappings; workspace is last field typically
			s.Workspace = fields[len(fields)-1]
		}
		rows = append(rows, s)
	}
	return rows, nil
}

// Exists reports whether a sandbox name appears in sbx ls.
func (r *Runner) Exists(name string) (bool, error) {
	rows, err := r.Ls()
	if err != nil {
		return false, err
	}
	for _, s := range rows {
		if s.Name == name {
			return true, nil
		}
	}
	return false, nil
}
