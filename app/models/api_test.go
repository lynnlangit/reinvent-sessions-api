package models

import "testing"

func TestListAPI(t *testing.T) {
	var actual, expected interface{}
	apis := ListAPI()

	actual = len(apis)
	expected = 1
	if actual != expected {
		t.Errorf("Unexpected API count. Expected %v, but got %v", expected, actual)
		return
	}
	actual = apis[0].Name
	expected = "/reinvent/sessions"
	if actual != expected {
		t.Errorf("Unexpected name. Expected %v, but got %v", expected, actual)
		return
	}
	actual = len(apis[0].Parameters)
	expected = 3
	if actual != expected {
		t.Errorf("Unexpected parameters count. Expected %v, but got %v", expected, actual)
		return
	}
	actual = apis[0].Parameters[0].Necessary
	expected = false
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}
