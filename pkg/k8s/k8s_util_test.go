package k8s

import (
	"testing"
)

func TestNewK8sClient(t *testing.T) {
	_, err := NewK8sClient()
	if err != nil {
		t.Errorf("Error creating k8sclient %v.", err)
	}
}
