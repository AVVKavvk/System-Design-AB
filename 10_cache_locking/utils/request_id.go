package utils

import (
	"bufio"
	"context"
	"crypto/rand"
	"io"
	"sync"

	"github.com/labstack/echo/v4"
)

// GetRequestID retrieves the request ID from the context
func GetRequestIDFromContext(c context.Context) string {
	requestID, ok := c.Value("request-id").(string)

	if !ok {
		var len uint8
		len = 32
		requestID = RandomString(len)
	}

	return requestID
}

// GetRequestCtx extracts the request ID from the response header of the Echo context and returns a new context with the request ID.
func GetRequestContextAndIdFromEchoContext(ctx echo.Context) (context.Context, string) {
	requestID := ctx.Response().Header().Get(echo.HeaderXRequestID)

	newContext := context.WithValue(context.Background(), "request-id", requestID)

	return newContext, requestID
}

// https://tip.golang.org/doc/go1.19#:~:text=Read%20no%20longer%20buffers%20random%20data%20obtained%20from%20the%20operating%20system%20between%20calls
var randomReaderPool = sync.Pool{New: func() interface{} {
	return bufio.NewReader(rand.Reader)
}}

const randomStringCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const randomStringCharsetLen = 52 // len(randomStringCharset)
const number = 256
const randomStringMaxByte = 255 - (number % randomStringCharsetLen)

func RandomString(length uint8) string {
	reader := randomReaderPool.Get().(*bufio.Reader)

	defer randomReaderPool.Put(reader)

	divident := 4
	b := make([]byte, length)
	r := make([]byte, length+(length/uint8(divident)))

	var i uint8 = 0

	for {
		_, err := io.ReadFull(reader, r)
		if err != nil {
			panic("unexpected error happened when reading from bufio.NewReader(crypto/rand.Reader)")
		}

		for _, rb := range r {
			// loop over range
			if rb > randomStringMaxByte {
				// Skip this number to avoid bias.
				continue
			}

			b[i] = randomStringCharset[rb%randomStringCharsetLen]
			i++

			if i == length {
				return string(b)
			}
		}
	}
}
