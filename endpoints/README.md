endpoints package
=================

* `common.go`
    * exists to provide generic functionality to generate server side HTML from html/template(s)
    * current best practice is to use template with minification
    * which is applying minification on load of the files and then parsing them as templates
* `static.go`
    * exists to serve static files from disk
* `login.go`
    * exists to provide login, logout functionality


Adding a new endpoint
=====================

* Create a new file under endpoints and code in the handlers.
* Most endpoints will be rendering server side templates
    * So, create the template under web/templates
    * The logic should be to load up the required template
    * And then render as output from the endpoint
* HTML might need some JS as well - which is served statically out of web/static