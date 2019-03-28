package yubioath

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Command interface {
	CombinedOutput() ([]byte, error)
}

type Commander func(string, ...string) ([]byte, error)

func WaitForKeys(ctx context.Context) (Keys, error) {
	_, err := exec.Command("which", "ykman").CombinedOutput()
	if err != nil {
		return nil, errors.Errorf("ykman seems to be not installed")
	}
	fn := func() (Keys, bool, error) {
		c := exec.Command("ykman", "oath", "code")
		b, err := c.CombinedOutput()
		if err != nil {
			if c.ProcessState.ExitCode() == 2 {
				return nil, false, nil
			}
		}
		return parseOutput(b), true, nil
	}

	keys, ok, err := fn()
	if err != nil {
		return nil, errors.WithStack(err)
	} else if ok {
		return keys, nil
	}
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-t.C:
			keys, found, err := fn()
			if err != nil {
				return nil, err
			} else if found {
				return keys, nil
			}
		}
	}

}

type Keys map[string]string

func (k Keys) Lookup(key string) (string, bool) {
	v, ok := k[key]
	return v, ok
}

func parseOutput(in []byte) Keys {
	m := Keys{}
	for _, line := range strings.Split(strings.TrimSpace(string(in)), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name := strings.Join(fields[0:len(fields)-1], " ")
		m[name] = fields[len(fields)-1]
	}
	return m
}
