# Internal Go Controllers

This repo is a experiment after:

- Reading Dave Cheney's article: ["Don’t force allocations on the callers of your API"](https://dave.cheney.net/2019/09/05/dont-force-allocations-on-the-callers-of-your-api). 
- Reading Mat Ryer's article: [How I write Go HTTP services after seven years](https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831)
- Using [gqlgen](https://gqlgen.com)

The repo declares a custom handler type that is used internally in Go, supports middleware, and is internal to the server process but is only exposed through GraphQL. 

```
HTTP Client -> Middleware -> GraphQL server -> Resolver -> Internal Middleware -> Method
```