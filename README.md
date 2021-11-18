# gomock

> This program creates mocks based on Go interface

> for example

> Running 

```sh

 gomock -s $GOROOT/src/io/io.go -n ReadCloser

```


> The Mock Created Is 

```go

type ReadCloserMock struct {
        CallbackRead func(p []byte) (n int, err error)
        CallbackClose func() (error)
}

func (m *ReadCloserMock) Read(p []byte) (n int, err error) {
        if m.CallbackRead != nil {
                return m.CallbackRead(p)
        }
        return int, nil
}

func (m *ReadCloserMock) Close() (error) {
        if m.CallbackClose != nil {
                return m.CallbackClose()
        }
        return nil
}

```