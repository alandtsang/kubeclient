package main

import (
	"context"
	"fmt"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clientv1 "github.com/alandtsang/kubeclient/pkg/client/v1"
	"github.com/alandtsang/kubeclient/pkg/config"
)

const (
	namespace  = "default"
	deployName = "mydeploy"
)

func main() {
	ctx := context.Background()

	kubeConfig, err := config.DefaultKubeConfig()
	if err != nil {
		log.Fatal(err)
	}

	kubeClient, err := clientv1.NewKubeClient(kubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err = applyTest(ctx, kubeClient); err != nil {
		return
	}
	time.Sleep(10 * time.Second)

	if err = getTest(ctx, kubeClient); err != nil {
		return
	}

	if err = listTest(ctx, kubeClient); err != nil {
		return
	}

	if err = deleteTest(ctx, kubeClient); err != nil {
		return
	}
}

func buildDeployment() *appsv1.Deployment {
	var replicas int32 = 1

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      deployName,
			Labels: map[string]string{
				"app": deployName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deployName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      deployName,
					Labels: map[string]string{
						"app": deployName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.21.1",
						},
					},
				},
			},
		},
	}
}

func applyTest(ctx context.Context, client *clientv1.KubeClient) (err error) {
	newDeploy := buildDeployment()
	if err = client.Apply(ctx, newDeploy); err != nil {
		return err
	}
	return nil
}

func getTest(ctx context.Context, client *clientv1.KubeClient) (err error) {
	obj := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      deployName,
		},
	}
	if err = client.Get(ctx, obj); err != nil {
		return err
	}
	fmt.Printf("obj: %+v\n\n", obj)
	return nil
}

func listTest(ctx context.Context, client *clientv1.KubeClient) (err error) {
	deployList := &appsv1.DeploymentList{}
	labels := map[string]string{
		"app": deployName,
	}

	if err = client.List(ctx, deployList, namespace, labels, nil); err != nil {
		return err
	}
	for _, dep := range deployList.Items {
		fmt.Printf("deploy: %+v\n\n", dep)
	}

	return nil
}

func deleteTest(ctx context.Context, client *clientv1.KubeClient) (err error) {
	obj := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      deployName,
		},
	}

	if err = client.Delete(ctx, obj); err != nil {
		return err
	}
	return nil
}
