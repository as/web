// Package main is an example of how implement a web service without verbose
// noise (such as, context.Context in every method call), see issue golang/go#22602
// for more details about the latter.
//
// This is mainly here to serve as an easily-accessible template now that the
// underlying ideas have been shared.
//
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"path"
	"strings"
)

func main() {
	log.Println(http.ListenAndServe(":80", Srv{}))
}

// Srv implements an http server. It stores request-scoped
// data. See ServeHTTP for details.
type Srv struct {
	file, path string
	ctx        context.Context
	err        error
}

// ServeHTTP is the entry point to the server. Note the lack of a pointer reciever.
func (s Srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.path = r.URL.Path
	s.ctx = r.Context()
	switch s.chop() {
	case "":
	case "pub":
		s.Pub(w, r)
	case "priv":
		s.auth(w, r, http.HandlerFunc(s.Priv))
	}
}

func (s *Srv) Pub(w http.ResponseWriter, r *http.Request)  { fmt.Fprintln(w, "public") }
func (s *Srv) Priv(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "private") }

func (s *Srv) auth(w http.ResponseWriter, r *http.Request, next http.Handler) {
	if rand.Intn(2) != 0 { // obviously not a good authentication scheme here
		http.Error(w, "heads", 401)
		return
	}
	next.ServeHTTP(w, r)
}

func (s *Srv) chop() string {
	s.file, s.path = chop(s.path)
	return s.file
}

// chop cleans the absolute path and cuts it at the next slash,
// returning file in relative form and the next path in absolute
// form. File never contains a slash, and next always contains
// one or more.
//
// The final step in the path returns the pair ("", "/")
//
func chop(p string) (file, next string) {
	p = path.Clean(p)[1:]
	if n := strings.Index(p, "/"); n >= 0 {
		return p[:n], p[n:]
	}
	return p, "/"
}
