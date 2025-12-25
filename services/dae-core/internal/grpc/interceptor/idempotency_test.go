package interceptor

import "testing"

func TestIsWriteMethod(t *testing.T) {
	tests := map[string]bool{
		"CreateUser":           true,
		"GetUser":              false,
		"ListUsers":            false,
		"UpdateUser":           true,
		"AdminSetUserRoles":    true,
		"AdminSetUserDisabled": true,
		"DeleteOrder":          true,
		"CloseSheet":           true,
		"ReopenSheet":          true,
		"JoinSheet":            true,
		"LeaveSheet":           true,
		"StreamOrders":         false,
		"ListOrders":           false,
	}

	for name, want := range tests {
		got := isWriteMethod(name)
		if got != want {
			t.Fatalf("isWriteMethod(%s) = %v, want %v", name, got, want)
		}
	}
}
