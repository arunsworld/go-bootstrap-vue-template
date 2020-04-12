go-bootstrap-vue template
=========================

A starter template with:
* Go back-end serving templatized HTML (and possibly JS) code and JSON API
* Bootstrap for styling
* Vue for any dynamic behaviors


Do we need yet another framework?
=================================

No, we don't. Hence this is not a framework. This is a template with many things thought through. Take it, modify it and use it. This is a starter pack.


Design - Packages
=================
* Services
    * ServerStart - starts up the HTTP server on a port with timeouts and graceful shutdown taken care of
    * SessionStore - a wrapper around gorilla/sessions.Store that provides:
        * authMiddlerware
        * Dealing with UserID (Store, Get, Logout, NextURL)
* Endpoints
    * These are the web server equivalent of func main.
    * This is also equivalent to Controllers in the MVC layering architecture.
    * Leverages html/template to do server side rendering
* Web
    * Not a package but location for all asserts - templates and static assets
    * This should be deployed alongside the binary
* App (main.go)
    * Instantiates concrete implementations (for injection) of:
        * Auth Store - to support Auth
    * Creates single instances (for sharing) of:
        * Session store
        * Template files (minified)

Each package is individually documented in detail.

Getting started
===============
* Make a copy of this template (everything in here)
* Change the module name to match in `go.mod`
* Create all the domain entities - they should go in the root directory - and be fully testable and tested
* If depending on any external I/O (DB/FS)
    * Create an abstraction interface for the I/O
    * Implement the abstraction within a model/adapter directory
* Create endpoints (please review README there)
    * That depend on the abstract I/O interface and the domain models
    * To provide the requested service
    * Concrete implementation of the I/O interface should be injected from main

Deploying
=========
* `make`
* `docker push arunsworld/go-vue-app:latest`

Reference Resources
===================
* [Go Templates](https://golang.org/pkg/text/template/). Refer to Actions section for details.
* [Templates Cheatsheet](https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet).
* [Bootstrap](https://getbootstrap.com/docs/4.4/getting-started/introduction/).
* [vue.js Guide](https://vuejs.org/v2/guide/).
* [vuetify icons](https://materialdesignicons.com/cdn/2.0.46/).
* [vuetify](https://vuetifyjs.com/en/).
* [vuetify color](https://vuetifyjs.com/en/styles/colors/).