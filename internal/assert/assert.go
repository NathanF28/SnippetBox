package assert

import "testing"

func Equal[T comparable](t *testing.T, actualValue, expectedValue T) {

	t.Helper()

	if actualValue != expectedValue {
		t.Errorf("got %v; want %v", actualValue, expectedValue)
	}
}
