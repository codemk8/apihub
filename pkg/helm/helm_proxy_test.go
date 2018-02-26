package helm

import "testing"

func TestHelmList(t *testing.T) {
	err := ListRelease()
	if err != nil {
		t.Errorf("Error creating k8sclient %v.", err)
	}
}
