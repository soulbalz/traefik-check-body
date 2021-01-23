# Check Body Request

This check body request plugin for [Traefik](https://github.com/traefik/traefik) which checks the incoming request for specific body and their values to be present and matching the configuration. If the request does not validate against the configured body, the middleware will return a 400 Bad Request status code.

## Limitations
#

 - JSON request are not working and it's will be make your traefik crash (We are fixing)

#
## Configuration

### Static

```yaml
pilot:
  token: xxxx

experimental:
  plugins:
    check-body:
      modulename: github.com/soulbalz/checkbodyplugin
      version: v1.0.5
```

### Dynamic configuration

```yaml
http:
  routers:
    my-router:
      rule: Path(`/whoami`)
      service: service-whoami
      entryPoints:
        - http
      middlewares:
        - check-body

  services:
   service-whoami:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
  
  middlewares:
    check-body:
      plugin:
        check-body:
          response:
            code: '1234'
            message: 'test'
            status: 401
            raw: |
              {
                "error": "example error"
              }
          body:
            - name: "BODY_1"
              matchtype: one
              values: 
                - "VALUE_1"
                - "VALUE_99"
            - name: "BODY_2"
              matchtype: one
              values: 
                - "VALUE_2"
            - name: "BODY_3"
              matchtype: one
              values: 
                - "VALUE_3"
              required: false
            - name: "BODY_4"
              matchtype: all
              values: 
                - "LUE_4"
                - "VALUE_5"
              contains: true
              required: true
```

#
## Configuration documentation

Supported configurations per body

| Setting   | Allowed values    | Required    | Description |
| :--       | :--               | :--         | :--         |
| name      | string            | Yes        | Name of the request body |
| matchtype | one, all          | Yes        | Match on all values or one of the values specified. The value 'all' is only allowed in combination with the 'contains' setting.|
| values    | []string          | Yes        | A list of allowed values which are matched against the request header value|
| contains  | boolean           | No        | If set to true (default false), the request is allowed if the request body value contains the value specified in the configuration |
| required  | boolean           | No        | If set to false (default true), the request is allowed if the body is absent or the value is empty|

#

response

| Setting   | Allowed values    | Required  | Description |
| :--       | :--               | :--       | :--         |
| code      | string            | No        | Response error message code (default: 400) |
| message | string              | No        | Response error message (default: "Invalid Request.") |
| status    | integer           | No        | Response error status code (default: 400) |
| raw  | string                 | No        | Custom your response message whatever you want |

#