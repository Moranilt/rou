package rou

import "testing"

func TestRouterParams(t *testing.T) {
	t.Run("Set action", func(t *testing.T) {
		routerParams := &routerParams{value: make(map[string]string)}

		want := "Melony"
		routerParams.Set("user", "Melony")
		got := routerParams.value["user"]
		if got != want {
			t.Errorf("value of key is not equal. Got - %s, want %s", got, want)
		}
	})

	t.Run("Get action", func(t *testing.T) {
		routerParams := &routerParams{value: make(map[string]string)}
		want := "Melony"
		routerParams.Set("user", want)
		got := routerParams.Get("user")
		if routerParams.Get("user") != "Melony" {
			t.Errorf("value of key is not equal. Got - %s, want %s", got, want)
		}
	})

	t.Run("Has action", func(t *testing.T) {
		routerParams := &routerParams{value: make(map[string]string)}
		want := "Melony"
		routerParams.Set("user", want)

		if !routerParams.Has("user") {
			t.Errorf("Got - %t, want %t", false, true)
		}

		if routerParams.Has("bad") {
			t.Errorf("Got - %t, want %t", true, false)
		}
	})

	t.Run("Delete action", func(t *testing.T) {
		routerParams := &routerParams{value: make(map[string]string)}
		want := "Melony"
		routerParams.Set("user", want)
		routerParams.Set("age", "20")
		routerParams.Delete("user")
		if _, ok := routerParams.value["user"]; ok {
			t.Error("Key was not deleted from map")
		}
		if _, ok := routerParams.value["age"]; !ok {
			t.Error("The wrong key 'age' has been removed from map")
		}
	})

}
