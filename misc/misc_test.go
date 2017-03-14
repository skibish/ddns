package misc

import "testing"

func TestSuccess(t *testing.T) {
	if Success(199) {
		t.Error("Should be error on 199")
		return
	}

	if Success(300) {
		t.Error("Should be error on 300")
		return
	}

	if !Success(222) {
		t.Error("Should be success on 222")
		return
	}

}
