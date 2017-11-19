# plz

Make application can be externally managed.

![plz](https://docs.google.com/drawings/d/e/2PACX-1vTkDCgDnGucsSPs1FgCcp40fA8JKAzmMTdfNAQQkHuIhsD-ivfkqBss0F75z0tURdHLaMrnvEAObK2e/pub?w=496&h=217)

Connections:

* countlog: provide observability via push or pull interface between application and plz
* mgmt: provide management callback between application and plz
* pnp: plug into local agent, and publish the plz interfaces to the public without binding to tcp port

Stateful roles:

* counselor: hold config and provide proper answer for application inquiry
* witch: direct web interface for administrator to manage this process, mainly for debugging purpose

From application point of view:

* it use countlog to tell the world what is going on
* it use mgmt to allow its state being updated externally (like a backdoor)
* it use counselor to query config