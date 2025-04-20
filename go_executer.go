package go_executer

import (
	"fmt"
	"github.com/99-m4n/go_logger"
	"os/exec"
	"os"
	"io"
	"bufio"
	"context"
	"time"
)

func ExecSh(file_path string) {
	go_logger.Info(fmt.Sprintf("File path: %s\n", file_path))

	file_content, err := os.ReadFile(file_path)
	if err != nil {go_logger.Error(fmt.Sprintf("%s", err))}

	go_logger.Info(fmt.Sprintf("File content:\n\n%s\n\n", file_content))

	reader, writer := io.Pipe()
	cmdCtx, cmdDone := context.WithCancel(context.Background())
	scannerStopped := make(chan struct{})

	go func() {
		defer close(scannerStopped)

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			t := time.Now().Format(time.RFC3339)
			fmt.Printf("%s -- %s\n", t, scanner.Text())
		}
	}()

	cmd := exec.Command("/bin/sh", "-c", string(file_content))
	cmd.Stdout = writer
	_ = cmd.Start()
	go func() {
		_ = cmd.Wait()
		cmdDone()
		writer.Close()
	}()
	<-cmdCtx.Done()

	<-scannerStopped

	go_logger.Info("Execution finished")
}

