// Package status contains status codes not unline the http
// statuses, but tailored toward the ddn ecosystem
package status

// Info statuses are used to convey that something has happened
// but has not finished yet. It is not a success, nor a failure.
//
// They can range from 1 to 99
const (
	Started    int = 1 // status.Started
	InProgress int = 2 // status.InProgress
)

// Success statuses are used to convey a successful result.
const (
	Success  int = 100 // status.Success
	Created  int = 101 // status.Created
	Accepted int = 102 // status.Accepted
	Update   int = 103 // status.Update
)

// Client errors are used to convey that something was
// wrong with a client request.
const (
	ClientError int = 200 // status.ClientError
	NotFound    int = 201 // status.NotFound
)

// Server errors are used to convey that something went wrong
// on the server.
const (
	ServerError int = 300 // status.ServerError
)
