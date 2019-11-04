# Gositemap

Simple program made in Go to create a sitemap of any website.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Installing

Just download the repository and compile it with `go build`.  

## Run 

The following flags are available:

- -o    Filename to output the sitemap
- -s    Number of maximum concurrent connections
- -q    If used, queries on URLs are not ignored

Example:

```
./gositemap https://example.com -o sitemap.xml -s 50
```

## License

This project is licensed under the Apache License, Version 2.0 - see http://www.apache.org/licenses/LICENSE-2.0 for more details
