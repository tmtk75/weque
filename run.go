package weque

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	KeyHandlersTimeout = "handlers.timeout"
)

var (
	Stdout io.Writer = nil // for test
	Stderr io.Writer = nil // for test
)

func init() {
	viper.BindEnv(KeyHandlersTimeout, "HANDLERS_TIMEOUT") // TODO: consier to reanme
}

func Run(env []string, wd, s string, args ...string) error {
	to := viper.GetDuration(KeyHandlersTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), to)
	defer cancel()

	cmd := exec.CommandContext(ctx, s, args...)
	cmd.Dir = wd
	cmd.Env = env

	cmd.Stdout = Stdout
	cmd.Stderr = Stderr
	//errbuf := bytes.NewBuffer([]byte{})  // https://github.com/golang/go/issues/23019
	//cmd.Stderr = bufio.NewWriter(errbuf)

	log.Printf("run: %v %v in %v", s, args, wd)
	now := time.Now()
	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("[cancel] %v %v %v", err, s, args)
		return errors.Wrapf(ctx.Err(), fmt.Sprintf("%v passed for %s", time.Since(now), s))
	}

	if err != nil {
		log.Printf("[error] %v %v %v", time.Since(now), err, args)
		return errors.Wrapf(err, "failed to run")
	}

	log.Printf("%v %v %v in %v", time.Since(now), s, args, wd)
	return nil
}
