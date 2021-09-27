# 🧰 wire-jacket 
[![GoDoc][doc-img]][doc] [![Github release][release-img]][release] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][report-card-img]][report-card]

Jacket of google/wire: advanced DI approach wrapping google/wire

![image](https://user-images.githubusercontent.com/24886864/132741198-7a92ef0a-7d59-4f3a-933c-fd5e830a31a4.png)

A jacket is an outer sheath that protects the wire core from mechanical, moisture and chemical issues.

wire-jacket wrapping google/wire and provides advanced DI(Dependency Injection).


Features
- IoC Container based environment variables using viper
- Lazy Loading
- Eager Loading


## Why wire-jacket needs?
google/wire works statically because it performs DI at compile-time.
This approach is great for debugging, but it has some drawbacks.

### 1. Absence of IoC (Inversion of Control) Container

Most go DI libraries, including google/wire, do not have an IoC container. IoC Container makes it easy to version up and replace modules. You can also make a Plan B and keep it. DB Skip mode that does not use DB or emergency processing mode that does not actually connect with other nodes can be applied by changing the activation module and restarting.

For example, if you use the MySQL DB implementation in your app and want to replace the implementation with MongoDB, you don't need to change the code, just change the string from MySQL to MongoDB in the IoC Container.

### 2. Graceful handling of modules
You can Gracefully write functions that need to call functions on every singleton instance like Close , Reload , etc.


[doc-img]: http://img.shields.io/badge/GoDoc-Reference-blue.svg
[doc]: https://pkg.go.dev/github.com/bang9211/wire-jacket

[release-img]: https://img.shields.io/github/release/bang9211/wire-jacket.svg
[release]: https://github.com/bang9211/wire-jacket/releases

[ci-img]: https://github.com/bang9211/wire-jacket/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/bang9211/wire-jacket/actions/workflows/go.yml

[cov-img]: https://codecov.io/gh/bang9211/wire-jacket/branch/main/graph/badge.svg
[cov]: https://codecov.io/gh/bang9211/wire-jacket/branch/main

[report-card-img]: https://goreportcard.com/badge/github.com/bang9211/wire-jacket
[report-card]: https://goreportcard.com/report/github.com/bang9211/wire-jacket

[release-policy]: https://golang.org/doc/devel/release.html#policy
