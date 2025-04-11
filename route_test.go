// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routing

import (
	"bytes"
	"testing"
)

type mockStore struct {
	*store
	data map[string]any
}

func newMockStore() *mockStore {
	return &mockStore{newStore(), make(map[string]any)}
}

func (s *mockStore) Add(key string, data any) int {
	for _, handler := range data.([]Handler) {
		handler(nil)
	}
	return s.store.Add(key, data)
}

func TestRouteNew(t *testing.T) {
	router := New()
	group := newRouteGroup("/admin", router, nil)

	r1 := newRoute("/users", group)
	if r1.name != "/admin/users" {
		t.Errorf("route.name = %s, want %s", r1.name, "/admin/users")
	}
	if r1.path != "/admin/users" {
		t.Errorf("route.path = %s, want %s", r1.path, "/admin/users")
	}
	if r1.template != "/admin/users" {
		t.Errorf("route.template = %s, want %s", r1.template, "/admin/users")
	}
	_, exists := router.routes[r1.name]
	if !exists {
		t.Error("router.routes[name] does not exist")
	}

	r2 := newRoute("/users/<id:\\d+>/*", group)
	if r2.name != "/admin/users/<id:\\d+>/*" {
		t.Errorf("route.name = %s, want %s", r2.name, "/admin/users/<id:\\d+>/*")
	}
	if r2.path != "/admin/users/<id:\\d+>/<:.*>" {
		t.Errorf("route.path = %s, want %s", r2.path, "/admin/users/<id:\\d+>/<:.*>")
	}
	if r2.template != "/admin/users/<id>/<>" {
		t.Errorf("route.template = %s, want %s", r2.template, "/admin/users/<id>/<>")
	}
	_, exists = router.routes[r2.name]
	if !exists {
		t.Error("router.routes[name] does not exist")
	}
}

func TestRouteName(t *testing.T) {
	router := New()
	group := newRouteGroup("/admin", router, nil)

	r1 := newRoute("/users", group)
	if r1.name != "/admin/users" {
		t.Errorf("route.name = %s, want %s", r1.name, "/admin/users")
	}
	r1.Name("user")
	if r1.name != "user" {
		t.Errorf("route.name = %s, want %s", r1.name, "user")
	}
	_, exists := router.routes[r1.name]
	if !exists {
		t.Error("router.routes[name] does not exist")
	}
}

func TestRouteURL(t *testing.T) {
	router := New()
	group := newRouteGroup("/admin", router, nil)
	r := newRoute("/users/<id:\\d+>/<action>/*", group)
	if got := r.URL("id", 123, "action", "address"); got != "/admin/users/123/address/<>" {
		t.Errorf("Route.URL@1 = %s, want %s", got, "/admin/users/123/address/<>")
	}
	if got := r.URL("id", 123); got != "/admin/users/123/<action>/<>" {
		t.Errorf("Route.URL@2 = %s, want %s", got, "/admin/users/123/<action>/<>")
	}
	if got := r.URL("id", 123, "action"); got != "/admin/users/123//<>" {
		t.Errorf("Route.URL@3 = %s, want %s", got, "/admin/users/123//<>")
	}
	if got := r.URL("id", 123, "action", "profile", ""); got != "/admin/users/123/profile/" {
		t.Errorf("Route.URL@4 = %s, want %s", got, "/admin/users/123/profile/")
	}
	if got := r.URL("id", 123, "action", "profile", "", "xyz/abc"); got != "/admin/users/123/profile/xyz%2Fabc" {
		t.Errorf("Route.URL@5 = %s, want %s", got, "/admin/users/123/profile/xyz%2Fabc")
	}
	if got := r.URL("id", 123, "action", "a,<>?#"); got != "/admin/users/123/a%2C%3C%3E%3F%23/<>" {
		t.Errorf("Route.URL@6 = %s, want %s", got, "/admin/users/123/a%2C%3C%3E%3F%23/<>")
	}
}

func newHandler(tag string, buf *bytes.Buffer) Handler {
	return func(*Context) error {
		buf.WriteString(tag)
		return nil
	}
}

func TestRouteAdd(t *testing.T) {
	store := newMockStore()
	router := New()
	router.stores["GET"] = store
	if store.count != 0 {
		t.Errorf("router.stores[GET].count = %d, want %d", store.count, 0)
	}

	var buf bytes.Buffer

	group := newRouteGroup("/admin", router, []Handler{newHandler("1.", &buf), newHandler("2.", &buf)})
	newRoute("/users", group).Get(newHandler("3.", &buf), newHandler("4.", &buf))
	if buf.String() != "1.2.3.4." {
		t.Errorf("buf@1 = %s, want %s", buf.String(), "1.2.3.4.")
	}

	buf.Reset()
	group = newRouteGroup("/admin", router, []Handler{})
	newRoute("/users", group).Get(newHandler("3.", &buf), newHandler("4.", &buf))
	if buf.String() != "3.4." {
		t.Errorf("buf@2 = %s, want %s", buf.String(), "3.4.")
	}

	buf.Reset()
	group = newRouteGroup("/admin", router, []Handler{newHandler("1.", &buf), newHandler("2.", &buf)})
	newRoute("/users", group).Get()
	if buf.String() != "1.2." {
		t.Errorf("buf@3 = %s, want %s", buf.String(), "1.2.")
	}
}

func TestRouteMethods(t *testing.T) {
	router := New()
	for _, method := range Methods {
		store := newMockStore()
		router.stores[method] = store
		if store.count != 0 {
			t.Errorf("router.stores[%s].count = %d, want %d", method, store.count, 0)
		}
	}
	group := newRouteGroup("/admin", router, nil)

	newRoute("/users", group).Get()
	if router.stores["GET"].(*mockStore).count != 1 {
		t.Errorf("router.stores[GET].count = %d, want %d", router.stores["GET"].(*mockStore).count, 1)
	}
	newRoute("/users", group).Post()
	if router.stores["POST"].(*mockStore).count != 1 {
		t.Errorf("router.stores[POST].count = %d, want %d", router.stores["POST"].(*mockStore).count, 1)
	}
	newRoute("/users", group).Patch()
	if router.stores["PATCH"].(*mockStore).count != 1 {
		t.Errorf("router.stores[PATCH].count = %d, want %d", router.stores["PATCH"].(*mockStore).count, 1)
	}
	newRoute("/users", group).Put()
	if router.stores["PUT"].(*mockStore).count != 1 {
		t.Errorf("router.stores[PUT].count = %d, want %d", router.stores["PUT"].(*mockStore).count, 1)
	}
	newRoute("/users", group).Delete()
	if router.stores["DELETE"].(*mockStore).count != 1 {
		t.Errorf("router.stores[DELETE].count = %d, want %d", router.stores["DELETE"].(*mockStore).count, 1)
	}
	newRoute("/users", group).Connect()
	if router.stores["CONNECT"].(*mockStore).count != 1 {
		t.Errorf("router.stores[CONNECT].count = %d, want %d", router.stores["CONNECT"].(*mockStore).count, 1)
	}
	newRoute("/users", group).Head()
	if router.stores["HEAD"].(*mockStore).count != 1 {
		t.Errorf("router.stores[HEAD].count = %d, want %d", router.stores["HEAD"].(*mockStore).count, 1)
	}
	newRoute("/users", group).Options()
	if router.stores["OPTIONS"].(*mockStore).count != 1 {
		t.Errorf("router.stores[OPTIONS].count = %d, want %d", router.stores["OPTIONS"].(*mockStore).count, 1)
	}
	newRoute("/users", group).Trace()
	if router.stores["TRACE"].(*mockStore).count != 1 {
		t.Errorf("router.stores[TRACE].count = %d, want %d", router.stores["TRACE"].(*mockStore).count, 1)
	}

	newRoute("/posts", group).To("GET,POST")
	if router.stores["GET"].(*mockStore).count != 2 {
		t.Errorf("router.stores[GET].count = %d, want %d", router.stores["GET"].(*mockStore).count, 2)
	}
	if router.stores["POST"].(*mockStore).count != 2 {
		t.Errorf("router.stores[POST].count = %d, want %d", router.stores["POST"].(*mockStore).count, 2)
	}
	if router.stores["PUT"].(*mockStore).count != 1 {
		t.Errorf("router.stores[PUT].count = %d, want %d", router.stores["PUT"].(*mockStore).count, 1)
	}
}

func TestBuildURLTemplate(t *testing.T) {
	tests := []struct {
		path, expected string
	}{
		{"", ""},
		{"/users", "/users"},
		{"<id>", "<id>"},
		{"<id", "<id"},
		{"/users/<id>", "/users/<id>"},
		{"/users/<id:\\d+>", "/users/<id>"},
		{"/users/<:\\d+>", "/users/<>"},
		{"/users/<id>/xyz", "/users/<id>/xyz"},
		{"/users/<id:\\d+>/xyz", "/users/<id>/xyz"},
		{"/users/<id:\\d+>/<test>", "/users/<id>/<test>"},
		{"/users/<id:\\d+>/<test>/", "/users/<id>/<test>/"},
		{"/users/<id:\\d+><test>", "/users/<id><test>"},
		{"/users/<id:\\d+><test>/", "/users/<id><test>/"},
	}
	for _, test := range tests {
		actual := buildURLTemplate(test.path)
		if actual != test.expected {
			t.Errorf("buildURLTemplate(%s) = %s, want %s", test.path, actual, test.expected)
		}
	}
}
