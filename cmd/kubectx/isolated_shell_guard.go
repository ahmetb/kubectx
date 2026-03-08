package main

import (
	"fmt"
	"os"

	"github.com/ahmetb/kubectx/internal/env"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
)

func checkIsolatedMode() error {
	if os.Getenv(env.EnvIsolatedShell) != "1" {
		return nil
	}

	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return fmt.Errorf("you are in a locked single-context shell, use 'exit' to leave")
	}

	cur := kc.GetCurrentContext()
	return fmt.Errorf("you are in a locked single-context shell (\"%s\"), use 'exit' to leave", cur)
}
