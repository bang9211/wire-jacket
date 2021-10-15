# Wire-Jacket: IoC Container of google/wire for cloud-native
[![GoDoc][doc-img]][doc] [![Github release][release-img]][release] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][report-card-img]][report-card]

<img src="https://user-images.githubusercontent.com/24886864/135562733-0c9cb4d5-ece2-428c-b8e3-0a155b675ace.png" width="560" height="200"/>

Jacket of google/wire: advanced DI approach wrapping google/wire for cloud.

- google/wire : https://github.com/google/wire

<img src="https://user-images.githubusercontent.com/24886864/136348668-1b127103-fbc7-482f-a1b5-be2ee63f4875.png" width="220" height="220"/> <img src="https://user-images.githubusercontent.com/24886864/132741198-7a92ef0a-7d59-4f3a-933c-fd5e830a31a4.png" width="255" height="219"/>


A jacket is an outer sheath that groups wires and protects the core 
from external issues.

wire-jacket wraps google/wire and provides advanced 
DI(Dependency Injection) experience in cloud-native.

# Installation
Install Wire-Jacket by running:
```
go get github.com/bang9211/wire-jacket
```
and ensuring that $GOPATH/bin is added to your $PATH.

# Example
Wire-Jacket example of ossicones.
In this example, ossicones is simple blockchain package.
It consisted of only 3 components: Config, Database, Blockchain.

Define simple two Interface, Implement.
```go
type Database interface {
    Connect() bool
    Close() error   //necessary for Wire Jacket
}

type MySQL struct {
    cfg config.Config
}

func NewMySQL(cfg config.Config) Database {
    return &MySQL{cfg : cfg}
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

```go
type Blockchain interface {
    Init() error
    Close() error   //necessary for Wire Jacket
}

type Ossicones struct {
    db *Database
}

func NewOssicones(db Database) Blockchain {
    return &Ossicones{db : db}
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
Suppose there is MongoDB that implements Database 
Interface like MySQL. And Wire-Jacket has ViperJacket by default.

Then, there are 3 `Interfaces` and 4 `Implements`.

- Database `Interface` - MySQL, MongoDB `Implement`
- Blockchain `Interface` - Ossicones `Implement`
- (Default) Config `Interface` - ViperJacket `Implement`

Database depends on config.Config. Blockchain depends on Database.

The pair of interface and implement called module in Wire-Jacket.

Wire-Jacket helps to replace implement of interface easy way. 
And close modules gracefully. So the modules are closable, 
have to implment Close().

### 1. Create wire.go with injectors.
```go
package wire

func InjectMySQL(cfg config.Config) (Database, error) {
	wire.Build(NewMySQL)
	return nil, nil
}

func InjectMongoDB(cfg config.Config) (Databsae, error) {
    wire.Build(NewMongoDB)
    return nil, nil
}

func InjectOssicones(db Database) (Blockchain, error) {
	wire.Build(NewOssicones)
	return nil, nil
}

// key will use in app.conf
var Injectors map[string]interface{} = {"mysql" : InjectMySQL}
var EagerInjectors map[string]interface{} = {"ossicones" : InjectOssicones}
```

### 2. Generate wire_gen.go using wire.
```
wire wire.go
```
wire is compile-time DI. It can verify validity of DI in 
compile-time.

### 3. Create app.conf(default way)
```
# Specify module to use.
modules=mysql ossicones

# And you can add any config to use.
db_host=localhost
db_port=3306
```
Choose modules to use mysql, ossicones.

Database binds to MySQL, Blockchain binds to Ossicones.

### 4. Create wirejacket, Set injectors, Call DoWire().
```go
wj := wirejacket.New().
    SetInjectors(wire.Injectors).
    SetEagerInjectors(wire.EagerInjectors)

// inject eager injectors.
wj.DoWire()
```
Except for eager injectors, injectors are loaded lazy by default.

Or If you call wj.GetModule() to get the module you need, 
all the dependencies of the module will be injected automatically.
you don't need to call DoWire() in this case. It is not necessary 
to call DoWire().

Assume that there is `mongodb` like `mysql` as the implementation of Database.

If you want to change implement of Database to `mongodb`, 
Just modify modules in config and restart app.
```
modules=mongodb ossicones
```
#

## Features
- google/wire based IoC Container
- Environment-variable based setting for cloud
- Lazy Loading, Eager Loading

#

## Why Wire-Jacket needs?
google/wire works statically because it performs DI at compile-time.
This approach is great for pre-debugging of DI.

It is difficult to debug problems that occur at runtime like dependency 
injection in a cloud environment.
So It can be usefully used in a cloud environment.

Wire-Jacket wraps google/wire in order to use google/wire appropriately 
for the cloud environment.

### 1. IoC (Inversion of Control) Container

google/wire just provides injecting and binding features, do not 
have an IoC container. IoC Container makes it easy to version up 
and replace modules. You can also make a Plan B and keep it. 
For example, DB Skip mode that does not use DB or emergency 
processing mode that does not actually connect with other nodes 
can be applied by replacing module and restarting.

And it helps not to use singleton as a global variable.
Automatically maintains and recycles one singleton object.

### 2. Binding based environment-variable for cloud.

The Twelve-Factors recommends use environment variables for
configuration. Because it's good in container, cloud environments.
However, it is not efficient to express all configs as environment
variables.

So, we adopted viper, which integrates most config formats with
environment variables in go.
By using viper, you can use various config file formats without
worrying about conversion to environment variables even if it is
not in env format like `json`, `yaml`, `toml`, etc.

In Wire-Jacket, modules to use set in config file.
For example, if you use the MySQL DB implementation in your app and 
want to replace the implementation with MongoDB, you don't need to 
change the code.

just change MySQL to MongoDB in config and restart the app.

### 3. Modulization, Loose-Coupling.
Wirejacket uses the simple approach of injecting implement into the 
interface. This approach simplifies and allows the team leader or 
designer to define the interface required for the project and 
effectively divide the work. the each implementation can be integrated 
easily and replaced as a plug-in using the config file. 

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
