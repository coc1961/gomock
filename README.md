# gomock

> This program creates mocks based on Go interface

> for example

> Running 

```sh

 gomock -s $GOROOT/src/io/io.go -n ReadCloser

```


> The Mock Created Is 

```go

// Interface compatible with ReadCloser that contains
// the Mock function to access the Mock instance
type ReadCloserMockInterface interface {
	io.ReadCloser
	Mock() *ReadCloserMock
}

// function to create the mock
func NewReadCloserMock() ReadCloserMockInterface {
	return &ReadCloserMock{}
}

// function to access the mock instance
// for example
// 	var myVar ReadCloser
// 	mock := NewReadCloserMock()
// 	mock.Mock().Callbackxxx = func(...)...{} // Modifies the default behavior of the mock function
// 	myVar = mock // Ok! compatible interface
func (m *ReadCloserMock) Mock() *ReadCloserMock {
	return m
}

// Mock for ReadCloser interface
type ReadCloserMock struct {
	CallbackRead  func(p []byte) (n int, err error)
	CallbackClose func() (retVar0 error)
}

// Read function
func (m *ReadCloserMock) Read(p []byte) (n int, err error) {
	if m.CallbackRead != nil {
		return m.CallbackRead(p)
	}
	return n, err
}

// Close function
func (m *ReadCloserMock) Close() (retVar0 error) {
	if m.CallbackClose != nil {
		return m.CallbackClose()
	}
	return retVar0
}

```