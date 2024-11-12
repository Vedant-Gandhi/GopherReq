package gopherreq

/*
This Reader allows to read the contents of the body.

If the body is sent as a whole then the whole body is sent back in one go.

If the body is sent in chunked form it returns every chunk for each read and for the last chunk it returns an empty data.
*/
type ResponseBodyReader interface {
	// This method reads from the input pipeline.
	Read() ([]byte, error)
}
