# authz

authz is an low latency HTTP authorization service. It provides an HTTP endpoint for serving forward authorization requests by web servers or HTTP proxies. Forward authorization is the function of authorizing a client request based on the status of a subsequent request. En example of forward authorization is implemented by NGINX's [`auth_request`](http://nginx.org/en/docs/http/ngx_http_auth_request_module.html) directive. 

## Authorization rules

Authorization of a client request is calculated based on a set of rules defined authz's configuration file. Each rule defines a `URL` and an access control list or `ACL`. 

### URL

A `URL` is defined as the combination of host and URI of the client request. A `URL` can have a wildcard `*` at the end denoting that this rule applies to any URL matching the expression `url.*`. The `URL` definition `www.site.com/pub/*` for example would match the following client request urls:
- `www.site.com/pub`
- `www.site.com/pub/user`
- `www.site.com/pub/doc/en/`

but would not match 
- `www.site.com/pub`

### ACL

The `ACL` is a list of access control statements consisting of a subject and a list of HTTP verbs the subject is authorized to use. A subject can be a specific user or a group of subjects. 

### Subject groups

Subject groups are defined in authz's configuration file. Each subject group is a list of one or more subjects. Each subject group name must be unique accross all subject and subject group names. 

### Features

- Low latency response (< 1 ms)
- Supported authentication schemes: HTTP Basic
- URL wildcards
- ACLs based on subject or subject group

