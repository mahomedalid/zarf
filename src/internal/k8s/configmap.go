package k8s

import (
	"context"
	"fmt"

	"github.com/defenseunicorns/zarf/src/config"
	"github.com/defenseunicorns/zarf/src/internal/message"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReplaceConfigmap deletes and recreates a configmap
func ReplaceConfigmap(namespace, name string, labels map[string]string, data map[string][]byte) (*corev1.ConfigMap, error) {
	message.Debugf("k8s.ReplaceConfigmap(%s, %s, data)", namespace, name)

	if err := DeleteConfigmap(namespace, name); err != nil {
		return nil, err
	}

	return CreateConfigmap(namespace, name, labels, data)
}

// CreateConfigmap applys a configmap to the cluster
func CreateConfigmap(namespace, name string, labels map[string]string, data map[string][]byte) (*corev1.ConfigMap, error) {
	message.Debugf("k8s.CreateConfigmap(%s, %s, data)", namespace, name)
	clientset := getClientset()

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				// track the creation of this ns by zarf
				config.ZarfManagedByLabel: "zarf",
			},
		},
		BinaryData: data,
	}

	// Add any additional labels to the service
	for key, val := range labels {
		configMap.ObjectMeta.Labels[key] = val
	}

	createOptions := metav1.CreateOptions{}
	return clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, createOptions)
}

// DeleteConfigmap delets a confimap by name
func DeleteConfigmap(namespace, name string) error {
	message.Debugf("k8s.DeleteConfigmap(%s, %s)", namespace, name)
	clientSet := getClientset()

	namespaceConfigmap := clientSet.CoreV1().ConfigMaps(namespace)

	err := namespaceConfigmap.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("error deleting the configmap: %w", err)
	}

	return nil
}

// DeleteConfigMapsByLabel deletes a configmap by label(s)
func DeleteConfigMapsByLabel(namespace string, labels map[string]string) error {
	message.Debugf("k8s.DeleteConfigMapsByLabel(%s, %v)", namespace, labels)
	clientSet := getClientset()

	labelSelector, _ := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: labels,
	})
	metaOptions := metav1.DeleteOptions{}
	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector.String(),
	}

	return clientSet.CoreV1().ConfigMaps(namespace).DeleteCollection(context.TODO(), metaOptions, listOptions)
}
