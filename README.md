# kubeclient
General kubernetes go client, can deploy any resource kind.

# Example

```go
package main

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	clientv1 "github.com/alandtsang/kubeclient/pkg/client/v1"
	"github.com/alandtsang/kubeclient/pkg/config"
)

func main() {
	ctx := context.Background()

	kubeConfig, err := config.DefaultKubeConfig()
	if err != nil {
		// ...
	}

	kubeClient, err := clientv1.NewKubeClient(kubeConfig)
	if err != nil {
		// ...
	}

	deploy := &appsv1.Deployment{
		// ...
	}

	if err = kubeClient.Apply(ctx, deploy); err != nil {
		// ...
    }
}
```