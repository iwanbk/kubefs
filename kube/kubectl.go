package kube

import (
	"context"
	"os/exec"
)

func (c *Client) kubeCliNs(ctx context.Context, ns string, args ...string) ([]byte, error) {
	args = append([]string{"-n", ns}, args...)
	return exec.CommandContext(ctx, "kubectl", args...).CombinedOutput()
}
