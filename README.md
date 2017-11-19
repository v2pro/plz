# plz

Make application can be externally managed.

![plz](https://docs.google.com/drawings/d/e/2PACX-1vTkDCgDnGucsSPs1FgCcp40fA8JKAzmMTdfNAQQkHuIhsD-ivfkqBss0F75z0tURdHLaMrnvEAObK2e/pub?w=496&h=217)

Connections:

* countlog: provide observability via push or pull interface between application and plz
* mgmt: provide management callback that can be triggered between application and plz
* pnp: plug into local agent, and publish the plz interfaces to the public

Stateful roles:

* counselor: hold config and provide proper answer for inquiry
* witch: web interface for administrator