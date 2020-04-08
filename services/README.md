services package
================

* `server.go`
    * Exists to capture best practices in starting up an http.Server
    * Traps ^C to provide graceful shutdown functionality
* `sessions.go`
    * Exists to provide services around sessions
    * Provides services around auth - and perhaps that is all that is required
* `auth.go`
    * Exists to wrap a given backing store to Register and Authenticate users