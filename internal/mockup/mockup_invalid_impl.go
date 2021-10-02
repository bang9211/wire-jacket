package mockup

import (
	"fmt"
)

type TestInterface interface {
	Test() error
	Close() error
}

type TestImplement struct{}

func NewTestImplement() TestInterface {
	return &TestImplement{}
}

func (ti *TestImplement) Test() error {
	return nil
}

func (ti *TestImplement) Close() error {
	return fmt.Errorf("mockup error")
}
