# endpoint-resolver

 
The endpoint resolver is a library developed and used internally at Detectify to ensure an endpoint is reachable before 
starting a scan on it.

# Installation

```
go get github.com/detectify/endpoint-resolver
```

# `applicationscanning` package

The `applicationscanning` package implements the `Resolver` interface and provides a full implementation of our
endpoint resolution process. This is precisely the resolver used internally by Detectify in the 
[Application Scanning](https://detectify.com/product/application-scanning) product.

It includes:
- native DNS resolver check
- third-party DNS resolver check
- open ports check
- HTTP request check
- User Agent check

# `opentelemetry` package

This package provides a tracing wrapper on the Resolver using OpenTelemetry. 

# Example

### Example without tracing

```go
package main

import (
    "fmt"
    endpointresolver "github.com/detectify/endpoint-resolver"
    "github.com/detectify/endpoint-resolver/applicationscanning"
)

func main() {
    awsDNS := os.GetEnv("EXTERNAL_DNS")
    resolver := applicationscanning.NewResolver([]string{awsDNS})

    conf := endpointresolver.ResolveConf{
        Endpoint:      "example.com",
        Ports:         []int{80, 8080, 8081, 443},
        UserAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
    }

    urls, err := resolver.Resolve(ctx, conf)
	if err != nil {
	    panic(err)	
    }
	
    fmt.Println("urls: %v", urls)
}
```

### Example with tracing

```go
package main 

import (
    "fmt"
    endpointresolver "github.com/detectify/endpoint-resolver"
)

func main() {
    awsDNS := os.GetEnv("EXTERNAL_DNS")
    resolver := endpointresolver.NewApplicationScanningResolverWithTracing([]string{awsDNS})

    conf := endpointresolver.ResolveConf{
        Endpoint:      "example.com",
        Ports:         []int{80, 8080, 8081, 443},
        UserAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
    }

    urls, err := resolver.Resolve(ctx, conf)
	if err != nil {
	    panic(err)	
    }
	
    fmt.Println("urls: %v", urls)
}
```

# Contributing

Please feel free to submit issues, fork the repository and send pull requests. In addition to fixes, new features are 
also welcome if you feel they are within the scope of the package. 

Feel free to reach out and discuss if you have any questions.

# License

This project is published under the terms of the MIT license.


