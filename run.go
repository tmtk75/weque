package weque

import (
	"context"
	"io"
	"log"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

var stdout io.Writer = nil // for test

func Run(env []string, wd, s string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, s, args...)
	cmd.Dir = wd
	cmd.Env = env

	cmd.Stdout = stdout
	//errbuf := bytes.NewBuffer([]byte{})  // https://github.com/golang/go/issues/23019
	//cmd.Stderr = bufio.NewWriter(errbuf)

	log.Printf("start: %v %v in %v", s, args, wd)
	now := time.Now()
	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("[cancel] %v %v %v", err, s, args)
		return ctx.Err()
	}

	if err != nil {
		log.Printf("[error] %v %v %v", time.Since(now), err, args)
		return errors.Wrapf(err, "failed to run")
	}

	log.Printf("%v %v %v in %v", time.Since(now), s, args, wd)
	return nil
}
