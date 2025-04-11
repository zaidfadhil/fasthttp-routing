// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routing

import (
	"bytes"
	"testing"
)

func TestRouteGroupTo(t *testing.T) {
	router := New()
	for _, method := range Methods {
		store := newMockStore()
		router.stores[method] = store
	}
	group := newRouteGroup("/admin", router, nil)

	group.Any("/users")
	for _, method := range Methods {
		if router.stores[method].(*mockStore).count != 1 {
			t.Errorf("router.stores[%s].count@1 = %d, want %d", method, router.stores[method].(*mockStore).count, 1)
		}
	}

	group.To("GET", "/articles")
	if router.stores["GET"].(*mockStore).count != 2 {
		t.Errorf("router.stores[GET].count@2 = %d, want %d", router.stores["GET"].(*mockStore).count, 2)
	}
	if router.stores["POST"].(*mockStore).count != 1 {
		t.Errorf("router.stores[POST].count@2 = %d, want %d", router.stores["POST"].(*mockStore).count, 1)
	}

	group.To("GET,POST", "/comments")
	if router.stores["GET"].(*mockStore).count != 3 {
		t.Errorf("router.stores[GET].count@3 = %d, want %d", router.stores["GET"].(*mockStore).count, 3)
	}
	if router.stores["POST"].(*mockStore).count != 2 {
		t.Errorf("router.stores[POST].count@3 = %d, want %d", router.stores["POST"].(*mockStore).count, 2)
	}
}

func TestRouteGroupMethods(t *testing.T) {
	router := New()
	for _, method := range Methods {
		store := newMockStore()
		router.stores[method] = store
		if store.count != 0 {
			t.Errorf("router.stores[%s].count = %d, want %d", method, store.count, 0)
		}
	}
	group := newRouteGroup("/admin", router, nil)

	group.Get("/users")
	if router.stores["GET"].(*mockStore).count != 1 {
		t.Errorf("router.stores[GET].count = %d, want %d", router.stores["GET"].(*mockStore).count, 1)
	}
	group.Post("/users")
	if router.stores["POST"].(*mockStore).count != 1 {
		t.Errorf("router.stores[POST].count = %d, want %d", router.stores["POST"].(*mockStore).count, 1)
	}
	group.Patch("/users")
	if router.stores["PATCH"].(*mockStore).count != 1 {
		t.Errorf("router.stores[PATCH].count = %d, want %d", router.stores["PATCH"].(*mockStore).count, 1)
	}
	group.Put("/users")
	if router.stores["PUT"].(*mockStore).count != 1 {
		t.Errorf("router.stores[PUT].count = %d, want %d", router.stores["PUT"].(*mockStore).count, 1)
	}
	group.Delete("/users")
	if router.stores["DELETE"].(*mockStore).count != 1 {
		t.Errorf("router.stores[DELETE].count = %d, want %d", router.stores["DELETE"].(*mockStore).count, 1)
	}
	group.Connect("/users")
	if router.stores["CONNECT"].(*mockStore).count != 1 {
		t.Errorf("router.stores[CONNECT].count = %d, want %d", router.stores["CONNECT"].(*mockStore).count, 1)
	}
	group.Head("/users")
	if router.stores["HEAD"].(*mockStore).count != 1 {
		t.Errorf("router.stores[HEAD].count = %d, want %d", router.stores["HEAD"].(*mockStore).count, 1)
	}
	group.Options("/users")
	if router.stores["OPTIONS"].(*mockStore).count != 1 {
		t.Errorf("router.stores[OPTIONS].count = %d, want %d", router.stores["OPTIONS"].(*mockStore).count, 1)
	}
	group.Trace("/users")
	if router.stores["TRACE"].(*mockStore).count != 1 {
		t.Errorf("router.stores[TRACE].count = %d, want %d", router.stores["TRACE"].(*mockStore).count, 1)
	}
}

func TestRouteGroupGroup(t *testing.T) {
	group := newRouteGroup("/admin", New(), nil)
	g1 := group.Group("/users")
	if g1.prefix != "/admin/users" {
		t.Errorf("g1.prefix = %s, want %s", g1.prefix, "/admin/users")
	}
	if len(g1.handlers) != 0 {
		t.Errorf("len(g1.handlers) = %d, want %d", len(g1.handlers), 0)
	}
	var buf bytes.Buffer
	g2 := group.Group("", newHandler("1", &buf), newHandler("2", &buf))
	if g2.prefix != "/admin" {
		t.Errorf("g2.prefix = %s, want %s", g2.prefix, "/admin")
	}
	if len(g2.handlers) != 2 {
		t.Errorf("len(g2.handlers) = %d, want %d", len(g2.handlers), 2)
	}

	group2 := newRouteGroup("/admin", New(), []Handler{newHandler("1", &buf), newHandler("2", &buf)})
	g3 := group2.Group("/users")
	if g3.prefix != "/admin/users" {
		t.Errorf("g3.prefix = %s, want %s", g3.prefix, "/admin/users")
	}
	if len(g3.handlers) != 2 {
		t.Errorf("len(g3.handlers) = %d, want %d", len(g3.handlers), 2)
	}
	g4 := group2.Group("", newHandler("3", &buf))
	if g4.prefix != "/admin" {
		t.Errorf("g4.prefix = %s, want %s", g4.prefix, "/admin")
	}
	if len(g4.handlers) != 1 {
		t.Errorf("len(g4.handlers) = %d, want %d", len(g4.handlers), 1)
	}
}

func TestRouteGroupUse(t *testing.T) {
	var buf bytes.Buffer
	group := newRouteGroup("/admin", New(), nil)
	group.Use(newHandler("1", &buf), newHandler("2", &buf))
	if len(group.handlers) != 2 {
		t.Errorf("len(group.handlers) = %d, want %d", len(group.handlers), 2)
	}

	group2 := newRouteGroup("/admin", New(), []Handler{newHandler("1", &buf), newHandler("2", &buf)})
	group2.Use(newHandler("3", &buf))
	if len(group2.handlers) != 3 {
		t.Errorf("len(group2.handlers) = %d, want %d", len(group2.handlers), 3)
	}
}
