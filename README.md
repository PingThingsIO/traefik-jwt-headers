# traefik-jwt-headers
Traefik middleware plugin which decodes a JWT token and forwards JWT claims as request headers. It can also rewrite values in tokens as needed.

## Installation
The plugin needs to be configured in the Traefik static configuration before it can be used.

### Installation on Kubernetes with Helm
The following snippet can be used as an example for the `values.yaml` file:
```values.yaml
experimental:
  plugins:
    enabled: true

additionalArguments:
- --experimental.plugins.traefik-jwt-headers.modulename=github.com/PingThingsIO/traefik-jwt-headers
- --experimental.plugins.traefik-jwt-headers.version=v0.0.1
```

### Installation via command line

```shell
traefik \
  --experimental.plugins.traefik-jwt-headers.moduleName=github.com/PingThingsIO/traefik-jwt-headers \
  --experimental.plugins.traefik-jwt-headers.version=v0.0.1
```

## Configuration

### Kubernetes

``` tab="File (Kubernetes)"
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: jwt-headers
spec:
  plugin:
    traefik-jwt-headers:
      claimsPrefix: attr
      headers:
        displayName: X-WEBAUTH-NAME
        email: X-WEBAUTH-EMAIL
        username: X-WEBAUTH-USERNAME
      unboxFirstElement: true
```

### Example with Value Rewrite

``` tab="File (Kubernetes)"
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: jwt-headers
spec:
  plugin:
    traefik-jwt-headers:
      claimsPrefix: attr
      headers:
        displayName: X-AUTH-NAME
        email: X-AUTH-EMAIL
        username: X-AUTH-USERNAME
      unboxFirstElement: true
      valueRewrite:
        username: #Claim to rewrite
          alice: charlie # Old Value -> New Value
```

## License

This software is released under the Apache 2.0 License
