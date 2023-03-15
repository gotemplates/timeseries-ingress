# timeseries-ingress

Timeseries-ingress is a service that provides a REST interface for timeseries data stored in TimescaleDB. The service itself is only composed of
the Http handlers for ingress traffic, a host package for initialization, and a resource package for host configuration. The application
functionality is contained in the gotemplates/timeseries module.



## handler

Http [handlers][handlerpkg] for health and access log ingress traffic.

## host

Service [host][hostpkg] initialization.


## resource

Host embedded configuration [resources][resourcepkg].


[handlerpkg]: <https://pkg.go.dev/github.com/gotemplates/timeseries-ingress/pkg/handler>
[hostpkg]: <https://pkg.go.dev/github.com/gotemplates/timeseries-ingress/pkg/host>
[resourcepkg]: <https://pkg.go.dev/github.com/gotemplates/timeseries-ingress/pkg/resource>
