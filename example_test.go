package wirejacket_test

import (
	"fmt"
	"log"

	wirejacket "github.com/bang9211/wire-jacket"
	"github.com/bang9211/wire-jacket/internal/mockup"
)

// Example of default use case to use New().
// Wire-Jacket defaultly uses 'app.conf' for setting modules
// to activate. Or you can use the flag '--config {file_name}'.
func Example() {
	wj := wirejacket.New().
		SetEagerInjectors(mockup.EagerInjectors).
		SetInjectors(mockup.Injectors)

	if err := wj.DoWire(); err != nil {
		log.Fatal(err)
	}
}

// Example of second use case to use NewWithServiceName()
func Example_second() {
	fmt.Println("Example_second")
}

// Example of thrid use case without app.conf
func Example_thrid() {
	fmt.Println("Example_third")
}

func ExampleNew() {
	wirejacket.New()
}

func ExampleNewWithServiceName() {
	wirejacket.NewWithServiceName("test service")
}
