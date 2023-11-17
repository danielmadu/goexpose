GoExpose
===============

Project inspired by [Expose](https://github.com/beyondcode/expose)

With GoExpose is possible to create public URLs for local sites. You can receive Webhooks on your local environment and share your local projects with others.

## Reading the code

⚠️ GoExpose is my first golang project. I made them for studies purpose and is not complete, feel free to send PRs.

## Requirements

 - Go 1.21+

 ## How to install

You can compile running `go build .`

## Example

To create a server: `goexpose server`

To share your local site: `goexpose share --server=http://localhost:3000 http://localhost:8080`

## SSL Suport

To add support to SSL connections use the flags `--certFile` and `--keyFile` in server command:

```bash
goexpose server --certFile=/etc/letsencrypt/live/example.com/fullchain.crt --keyFile=/etc/letsencrypt/live/example.com/privkey.key
```

When execute the client just change the server url protocol to `https`

## Basic Auth

To enable the Basic Auth just pass the `--basicAuth` flag with the user and password separated by `:` like the example bellow:

`goexpose share --basicAuth=test:password --server=http://localhost:3000 http://localhost:8080`