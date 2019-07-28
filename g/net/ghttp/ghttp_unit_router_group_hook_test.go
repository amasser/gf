// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
	"time"
)

func Test_Router_Group_HOOK1(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	g := s.Group("/api")
	g.GET("/handler", func(r *ghttp.Request) {
		r.Response.Write("1")
	})
	g.ALL("/handler", func(r *ghttp.Request) {
		r.Response.Write("0")
	}, ghttp.HOOK_BEFORE_SERVE)
	g.ALL("/handler", func(r *ghttp.Request) {
		r.Response.Write("2")
	}, ghttp.HOOK_AFTER_SERVE)

	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		gtest.Assert(client.GetContent("/api/handler"), "012")
		gtest.Assert(client.PostContent("/api/handler"), "02")
		gtest.Assert(client.GetContent("/api/ThisDoesNotExist"), "Not Found")
	})
}

func Test_Router_Group_HOOK2(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	g := s.Group("/api")
	g.GET("/handler", func(r *ghttp.Request) {
		r.Response.Write("1")
	})
	g.GET("/*", func(r *ghttp.Request) {
		r.Response.Write("0")
	}, ghttp.HOOK_BEFORE_SERVE)
	g.GET("/*", func(r *ghttp.Request) {
		r.Response.Write("2")
	}, ghttp.HOOK_AFTER_SERVE)

	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		gtest.Assert(client.GetContent("/api/handler"), "012")
		gtest.Assert(client.PostContent("/api/handler"), "Not Found")
		gtest.Assert(client.GetContent("/api/ThisDoesNotExist"), "02")
		gtest.Assert(client.PostContent("/api/ThisDoesNotExist"), "Not Found")
	})
}

func Test_Router_Group_HOOK3(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.Group("/api").Bind([]g.Slice{
		{"ALL", "handler", func(r *ghttp.Request) {
			r.Response.Write("1")
		}},
		{"ALL", "/*", func(r *ghttp.Request) {
			r.Response.Write("0")
		}, ghttp.HOOK_BEFORE_SERVE},
		{"ALL", "/*", func(r *ghttp.Request) {
			r.Response.Write("2")
		}, ghttp.HOOK_AFTER_SERVE},
	})

	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		gtest.Assert(client.GetContent("/api/handler"), "012")
		gtest.Assert(client.PostContent("/api/handler"), "012")
		gtest.Assert(client.DeleteContent("/api/ThisDoesNotExist"), "02")
	})
}
