package apery

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type (
	AperyClient interface {
		Evaluate(ctx context.Context, sfen string, moves []string, timeout time.Duration) (value int, bestmove string, pv []string, err error)
	}

	aperyClient struct {
		bin string
	}
)

func NewAperyClient(bin string) AperyClient {
	return &aperyClient{bin}
}

func (a *aperyClient) Evaluate(ctx context.Context, sfen string, moves []string, timeout time.Duration) (value int, bestmove string, pv []string, err error) {
	cmd := exec.CommandContext(ctx, a.bin)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return 0, "", nil, err
	}
	defer stdin.Close()

	if err := cmd.Start(); err != nil {
		return 0, "", nil, err
	}

	if err := a.isReady(stdin, &stdout); err != nil {
		return 0, "", nil, err
	}

	if err := a.setPosition(stdin, &stdout, sfen, moves); err != nil {
		return 0, "", nil, err
	}

	if err := a._go(stdin); err != nil {
		return 0, "", nil, err
	}

	time.Sleep(timeout)

	if err := a.stop(stdin); err != nil {
		return 0, "", nil, err
	}

	value, bestmove, pv, err = a.getResult(&stdout)
	if err != nil {
		return 0, "", nil, err
	}

	return value, bestmove, pv, nil
}

func (a *aperyClient) isReady(stdin io.Writer, stdout io.Reader) error {
	if _, err := io.WriteString(stdin, "isready\n"); err != nil {
		return err
	}

	res, err := a.waitResponse(stdout, 1000, 100*time.Millisecond)
	if err != nil || res != "readyok\n" {
		return fmt.Errorf("apery がisreadyに対して%d秒以内にreadyokを返しません", 10)
	}
	return nil
}

func (a *aperyClient) waitResponse(stdout io.Reader, attemptLimit int, interval time.Duration) (string, error) {
	for i := 0; i < attemptLimit; i++ {
		bytes, err := io.ReadAll(stdout)
		if err != nil {
			return "", err
		}

		if len(bytes) == 0 {
			time.Sleep(interval)
			continue
		}

		return string(bytes), nil
	}

	return "", errors.New("attempt limit exceeded")
}

func (a *aperyClient) setPosition(stdin io.Writer, stdout io.Reader, sfen string, moves []string) error {
	if _, err := io.WriteString(
		stdin, fmt.Sprintf("position sfen %s moves %s\n", sfen, strings.Join(moves, " "))); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	bytes, err := io.ReadAll(stdout)
	if err != nil {
		return err
	}

	if len(bytes) != 0 {
		return errors.New(string(bytes))
	}

	return nil
}

func (a *aperyClient) _go(stdin io.Writer) error {
	if _, err := io.WriteString(stdin, "go\n"); err != nil {
		return err
	}
	return nil
}

func (a *aperyClient) stop(stdin io.Writer) error {
	if _, err := io.WriteString(stdin, "stop\n"); err != nil {
		return err
	}
	return nil
}

func (a *aperyClient) getResult(stdout io.Reader) (value int, bestmove string, pv []string, err error) {
	res, err := a.waitResponse(stdout, 10, 100*time.Millisecond)
	logs := strings.TrimRight(string(res), "\n")
	lines := strings.Split(logs, "\n")
	if err != nil || !strings.Contains(logs, "bestmove") || len(lines) <= 2 {
		return 0, "", nil, fmt.Errorf("bestmoveが得られません")
	}

	bestmoveline := lines[len(lines)-1]
	bestmove = strings.Split(bestmoveline, " ")[1]

	lastInfoLine := lines[len(lines)-2]
	lastInfoLineParts := strings.Split(lastInfoLine, " ")
	pv = make([]string, 0, 10)

	for i, part := range lastInfoLineParts {
		if part == "cp" {
			value, err = strconv.Atoi(lastInfoLineParts[i+1])
			if err != nil {
				return 0, "", nil, err
			}
		}

		if part == "pv" {
			for _, p := range lastInfoLineParts[i+1:] {
				pv = append(pv, p)
			}
			break
		}
	}

	return value, bestmove, pv, nil
}
