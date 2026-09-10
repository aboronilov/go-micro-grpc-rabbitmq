// Restart wrapper used by Tilt live_update without pulling tiltdev/restart-helper.
package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	watchFile := flag.String("watch_file", "/tmp/.restart-proc", "file whose writes trigger a process restart")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("usage: tilt-restart-wrapper [--watch_file=path] command [args...]")
	}

	if err := ensureWatchFile(*watchFile); err != nil {
		log.Fatalf("watch file: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var cmd *exec.Cmd
	start := func() {
		cmd = exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		if err := cmd.Start(); err != nil {
			log.Fatalf("start %v: %v", args, err)
		}
	}

	stop := func() {
		if cmd == nil || cmd.Process == nil {
			return
		}
		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			_ = syscall.Kill(-pgid, syscall.SIGTERM)
		}
		done := make(chan struct{})
		go func() {
			_ = cmd.Wait()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			if err == nil {
				_ = syscall.Kill(-pgid, syscall.SIGKILL)
			}
			_ = cmd.Wait()
		}
	}

	start()
	stamp, _ := fileStamp(*watchFile)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case sig := <-sigs:
			if cmd != nil && cmd.Process != nil {
				_ = cmd.Process.Signal(sig)
				_ = cmd.Wait()
			}
			os.Exit(0)
		case <-ticker.C:
			next, err := fileStamp(*watchFile)
			if err != nil || next == stamp {
				continue
			}
			stamp = next
			stop()
			start()
		}
	}
}

func ensureWatchFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return err
	}
	return f.Close()
}

func fileStamp(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.ModTime().UnixNano() + info.Size(), nil
}
