package query

import "io"

func StreamCloseError() error {
	return io.EOF
}
