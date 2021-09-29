# Wire-Jacket: IoC Container of google/wire
[![GoDoc][doc-img]][doc] [![Github release][release-img]][release] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][report-card-img]][report-card]

Jacket of google/wire: advanced DI approach wrapping google/wire.

- google/wire : https://github.com/google/wire

A jacket is an outer sheath that groups wires and protects the core 
from external issues.

wire-jacket wraps google/wire and provides advanced 
DI(Dependency Injection) experience.

![image](https://user-images.githubusercontent.com/24886864/132741198-7a92ef0a-7d59-4f3a-933c-fd5e830a31a4.png)



# Installation
Install Wire-Jacket by running:
```
go get github.com/bang9211/wire-jacket
```
and ensuring that $GOPATH/bin is added to your $PATH.

# Example
Wire-Jacket example of ossicones.
In this example, ossicones is simple blockchain package.
It consisted of only 3 components: Config, Database, 
OssiconesBlockchain.

Simple two Interface, Implement Definition.
```
type Database interface {
    Connect() bool
    Close() error   //necessary for Wire Jacket
}

type MySQL struct {
    cfg config.Config
}

func (m *MySQL) Connect() error {
    address := m.cfg.GetString("address", "localhost:3306")
    ...
    return nil
}

func (m *MySQL) Close() error {
    return nil
}
```

```
type Blockchain interface {
    Init() error
    Close() error   //necessary for Wire Jacket
}

type Ossicones struct {
    db *Database
}

func (o *Ossicones) Init() error {
    err := o.db.Connect()
    ...
    return nil
}

func (o *Ossicones) Close() error {
    return nil
}
```
Database depends on config.Config. Blockchain depends on Database.

You can use config.Config. It is contained in Wire-Jacket by default.

The pair of interface and implement called module in Wire-Jacket.

Wire-Jacket supports to close modules gracefully. So the modules 
are closable, have to implment Close().

### 1. Create wire.go with two injectors.
```
func InjectMySQL(cfg config.Config) (Blockchain, error) {
	wire.Build(func(cfg config.Config) Database {
        return &MySQL{cfg : cfg}
    })
	return nil, nil
}

func InjectOssicones(db Database) (Blockchain, error) {
	wire.Build(func(db Database) Blockchain {
        return &Ossicones{db : db}
    })
	return nil, nil
}

var injectors map[string]interface{} = {"mysql" : InjectMySQL}
var eagerInjectors map[string]interface{} = {"ossicones" : InjectOssicones}
```

### 2. Generate wire_gen.go using wire.
```
wire
```
wire is compile-time DI. It can verify validity of DI in compile-time.

### 3. Create wirejacket, Set injectors, Call DoWire().
```
// service name is test
wj := wirejacket.NewWithInjectors("test", injectors, eagerInjectors)

// you can read modules explicitly using SetActivatingModules().
// or implicitly in activating_modules of {service name}.conf by default.
wj.SetActivatingModules([]string{"mysql", "ossicones"})

// inject eager injectors.
wj.DoWire()
```
Except for eager injectors, injectors are loaded lazy by default.
If there is no any eager injector, all the injectors are called by 
Wire Jacket. 

Or If you call wj.GetModule() to get the module you need, 
all the dependencies of the module injected automatically.
you don't need to call DoWire() in this case.

#

## Features
- IoC Container based environment variables using viper config
- Lazy Loading, Eager Loading

#

## Why wire-jacket needs?
google/wire works statically because it performs DI at compile-time.
This approach is great for debugging, but it has some drawbacks.

### 1. Absence of IoC (Inversion of Control) Container

Most go DI libraries, including google/wire, do not have an IoC 
container. IoC Container makes it easy to version up and replace 
modules. You can also make a Plan B and keep it. DB Skip mode that 
does not use DB or emergency processing mode that does not actually 
connect with other nodes can be applied by changing the activation 
module and restarting.

For example, if you use the MySQL DB implementation in your app and 
want to replace the implementation with MongoDB, you don't need to 
change the code, just change the string from MySQL to MongoDB in 
config and restart the app.

### 2. Graceful handling of modules
You can Gracefully write functions that need to call functions on 
every singleton instance like Close , Reload , etc.

### 3. Interface binding Implement
Wirejacket uses the approach of injecting implement into the 
interface. This approach allows the team leader or designer to 
define the interface required for the project and effectively divide 
the work, and each implementation can be easily replaced with a 
plug-in method using the config file. You just need to change the 
module (implement) name in the configuration.


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
