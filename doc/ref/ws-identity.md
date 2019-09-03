# Request identification

Back to: [Reference](README.md) | [Web Services](ws-index.md)

---

It is common practise to generate an unique ID for every HTTP request received by a web service. This ID often takes the
form of a [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier) and can be referred to as a request ID or
a correlation ID, among other terms.

These IDs are often generated at the edge of a web architecture in the first server receiving a request from a user 
(a web page or a public API endpoint) then propagated down through calls to internal services. This allows a single
logical request to be traced down through your whole architecture.


Granitic provides an interface [httpserver.IdentifiedRequestContextBuilder](https://godoc.org/github.com/graniticio/granitic/facility/httpserver#IdentifiedRequestContextBuilder)
which allows you to define a component that:

 * Generates new IDs or recovers them from an inbound HTTP request
 * Store that ID in a new [context.Context](https://golang.org/pkg/context/) 
 * Provide a way of extracting an ID from an existing context
 
If you create a component that implements [httpserver.IdentifiedRequestContextBuilder](https://godoc.org/github.com/graniticio/granitic/facility/httpserver#IdentifiedRequestContextBuilder),
it will automatically be injected into your [HTTPServer](fac-http-server.md) using a [decorator](ioc-decorators.md).

The ID will automatically be made available to any [request instrumentation](ws-instrumentation.md) you have set up and,
by using the context key you have used to store the ID in the context, can be logged in [application](fac-logger.md) and
[access](fac-http-server.md) logging.

If your application needs to recover the ID, it should be given a reference to the same 
[httpserver.IdentifiedRequestContextBuilder](https://godoc.org/github.com/graniticio/granitic/facility/httpserver#IdentifiedRequestContextBuilder),
component and use its `ID` method.

---
**Next**: [Rule based validation](vld-index.md)

**Prev**: [Instrumentation](ws-instrumentation.md)