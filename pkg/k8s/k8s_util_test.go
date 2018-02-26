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

func TestInitialCheck(t *testing.T) {
	err := CheckK8s()
	if err != nil {
		t.Errorf("Error initial checking cluster %v.", err)
	}
}

func TestCreatePV(t *testing.T) {
	err := AddPV()
	if err != nil {
		t.Errorf("Error creating PV %v.", err)
	}
}

func TestDestroyPV(t *testing.T) {
	err := DestroyPV()
	if err != nil {
		t.Errorf("Error destroy PV %v.", err)
	}
}
