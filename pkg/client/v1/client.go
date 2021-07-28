package client

import (
	"context"
	"log"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// KubeClient is a client structure containing ClientSet and control-runtime Client.
type KubeClient struct {
	clientSet  *kubernetes.Clientset
	ctrlClient crclient.Client
}

// NewKubeClient return a KubeClient instance.
func NewKubeClient(config *rest.Config) (kubeClient *KubeClient, err error) {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("NewKubeClient new client set failed, %s", err)
		return nil, err
	}

	ctrlClient, err := crclient.New(config, crclient.Options{})
	if err != nil {
		log.Printf("NewKubeClient new controller client failed, %s", err)
		return nil, err
	}

	kubeClient = &KubeClient{
		clientSet:  clientSet,
		ctrlClient: ctrlClient,
	}
	return kubeClient, nil
}

// Apply creates when the object is not deployed, otherwise patches the object.
func (c *KubeClient) Apply(ctx context.Context, obj crclient.Object) (err error) {
	if err = c.ctrlClient.Create(ctx, obj); err == nil {
		return nil
	}

	if !apierrors.IsAlreadyExists(err) {
		log.Printf("Apply object failed, %s", err)
		return err
	}

	key := crclient.ObjectKeyFromObject(obj)

	gvk := obj.GetObjectKind().GroupVersionKind()
	oldObj := unstructured.Unstructured{}
	oldObj.SetGroupVersionKind(gvk)
	if err = c.ctrlClient.Get(ctx, key, &oldObj); err != nil {
		log.Printf("Apply get object failed, %s", err)
		return err
	}

	patch := crclient.MergeFrom(oldObj.DeepCopy())
	if err = c.ctrlClient.Patch(ctx, obj, patch); err != nil {
		log.Printf("Apply patch object failed, %s", err)
		return err
	}
	return nil
}

// Delete deletes the object,
func (c *KubeClient) Delete(ctx context.Context, obj crclient.Object) error {
	err := c.ctrlClient.Delete(ctx, obj, crclient.PropagationPolicy(metav1.DeletePropagationBackground))
	if err != nil {
		if !apierrors.IsNotFound(err) {
			log.Printf("Delete failed, %s", err)
			return err
		}
	}
	return nil
}

// Get gets the object.
func (c *KubeClient) Get(ctx context.Context, obj crclient.Object) (err error) {
	key := crclient.ObjectKeyFromObject(obj)
	err = c.ctrlClient.Get(ctx, key, obj)
	if apierrors.IsNotFound(err) {
		log.Printf("Get failed, %s", err)
		return err
	}
	return nil
}

// List lists the kind in namespace with labels and fields.
func (c *KubeClient) List(ctx context.Context, objList crclient.ObjectList, namespace string, labels, fields map[string]string) (err error) {
	listOpts := []crclient.ListOption{
		crclient.InNamespace(namespace),
		crclient.MatchingLabels(labels),
		crclient.MatchingFields(fields),
	}

	if err = c.ctrlClient.List(ctx, objList, listOpts...); err != nil {
		log.Printf("List failed, %s", err)
		return err
	}
	return nil
}
