package yubiauth

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Command interface {
	CombinedOutput() ([]byte, error)
}

type Commander func(string, ...string) ([]byte, error)

func defaultCommander(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).CombinedOutput()
}

func WaitForKeys(ctx context.Context) (Keys, error) {
	return waitForKeysWithCommander(ctx, defaultCommander)
}

func waitForKeysWithCommander(ctx context.Context, cmd Commander) (Keys, error) {
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case _ = <-t.C:
			b, err := cmd("yubioath")
			if err != nil {
				if bytes.Contains(b, msgYubiKeyNotFound) {
					continue
				}
				return nil, fmt.Errorf("%s\n%s", b, err)
			}
			return parseOutput(b), nil
		}
	}
}

var ErrTimeoutWaitingForKeys = fmt.Errorf("timeout waiting for keys")

var msgYubiKeyNotFound = []byte("No YubiKey found!")

func isExecutableNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "executable file not found")
}

type Keys map[string]int

func (k Keys) Lookup(key string) (int, bool) {
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
		value, err := strconv.Atoi(fields[len(fields)-1])
		if err != nil {
			continue
		}
		m[name] = value
	}
	return m
}