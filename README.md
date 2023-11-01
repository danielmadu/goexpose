GoExpose
===============

Project inspired by [Expose](https://github.com/beyondcode/expose)

With GoExpose is possible to create public URLs for local sites. You can receive Webhooks on your local eviorement and share your local projects with others.

## Reading the code

⚠️ GoExpose is my first golang project. I made them for studies purpose and is not complete, feel free to send PRs.

## Requirements

 - Go 1.21+

 ## How to install

You can compile running `go build .`

## Example

To create a server: `goexpose server`

To share your local site: `goexpose share --server=localhost:3000 http://localhost:8080`

