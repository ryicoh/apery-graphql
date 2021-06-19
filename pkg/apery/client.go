package apery

import (
	"context"
	"os/exec"
)

type (
	AperyClient interface {
		Evaluate(ctx context.Context, sfen string) error
	}

	aperyClient struct {
		bin string
	}
)

func NewAperyClient(bin string) AperyClient {
	return &aperyClient{bin}
}

func (a *aperyClient) Evaluate(ctx context.Context, sfen string) error {
	cmd := exec.CommandContext(ctx, a.bin)

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := cmd.Stdin.Read([]byte("isready")); err != nil {
		return err
	}

	return nil
}
