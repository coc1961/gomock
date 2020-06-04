# gomock

> This program creates mocks based on Go interface

> for example

given the interface `TestInterface`

```go

type T1 struct {
}

type T1Interface interface {
}

type TestInterface interface {
	Func1(i int, pt1 t1.T1) (rt1 *t1.T1, xx t1.T1Interface, err error)
}

```

Create the mock 

```go

type TestInterfaceMock struct {
	CallbackFunc1 func(i int, pt1 t1.T1) (rt1 *t1.T1, xx t1.T1Interface, err error)
}

func (m *TestInterfaceMock) Func1(i int, pt1 t1.T1) (rt1 *t1.T1, xx t1.T1Interface, err error) {
	if m.CallbackFunc1 != nil {
		return m.CallbackFunc1(i, pt1)
	}
	return &t1.T1{}, t1.T1Interface, nil
}


```