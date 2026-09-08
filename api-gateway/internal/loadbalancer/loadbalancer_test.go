package loadbalancer

import "testing"

func TestRoundRobin(t *testing.T) {
	lb := New([]string{
		"http://orders-1:8080",
		"http://orders-2:8080",
		"http://orders-3:8080",
	})

	expected := []string{
		"http://orders-1:8080",
		"http://orders-2:8080",
		"http://orders-3:8080",
		"http://orders-1:8080",
		"http://orders-2:8080",
		"http://orders-3:8080",
	}

	for _, expectedURL := range expected {
		backend, ok := lb.Next()

		if !ok {
			t.Fatal("expected a healthy backend")
		}

		if backend.URL != expectedURL {
			t.Errorf(
				"expected %s, got %s",
				expectedURL,
				backend.URL,
			)
		}
	}
}
func TestUnhealthyBackendIsSkipped(t *testing.T) {
	lb := New([]string{
		"http://orders-1:8080",
		"http://orders-2:8080",
		"http://orders-3:8080",
	})

	backend, ok := lb.Next()
	if !ok {
		t.Fatal("expected a backend")
	}

	if backend.URL != "http://orders-1:8080" {
		t.Fatalf("expected orders-1, got %s", backend.URL)
	}

	lb.MarkUnhealthy(backend)

	expected := []string{
		"http://orders-2:8080",
		"http://orders-3:8080",
		"http://orders-2:8080",
		"http://orders-3:8080",
	}

	for _, expectedURL := range expected {
		backend, ok := lb.Next()

		if !ok {
			t.Fatal("expected a healthy backend")
		}

		if backend.URL != expectedURL {
			t.Errorf(
				"expected %s, got %s",
				expectedURL,
				backend.URL,
			)
		}
	}
}
func TestNoHealthyBackends(t *testing.T) {
	lb := New([]string{
		"http://orders-1:8080",
		"http://orders-2:8080",
	})

	backend, _ := lb.Next()
	lb.MarkUnhealthy(backend)

	backend, _ = lb.Next()
	lb.MarkUnhealthy(backend)

	_, ok := lb.Next()

	if ok {
		t.Fatal("expected no healthy backends")
	}
}