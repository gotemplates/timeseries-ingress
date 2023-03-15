# timeseries-ingress

[Timeseries-ingress][timeseriespkg] is a service that provides a REST interface for timeseries data stored in TimescaleDB. The service itself is only composed of
the Http handlers for ingress traffic, a host package for initialization, and a resource package that contains host configuration. The application
functionality is contained in the gotemplates/timeseries module.



## handler

Http handlers for health and access log traffic.

## host

Service host initialization.


## resource

Host embedded configuration resources.

[timeseriespkg]: <https://pkg.go.dev/github.com/gotemplates/timeseries-ingress/pkg>


