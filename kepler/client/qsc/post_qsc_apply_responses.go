// Code generated by go-swagger; DO NOT EDIT.

package qsc

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// PostQscApplyReader is a Reader for the PostQscApply structure.
type PostQscApplyReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostQscApplyReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPostQscApplyOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPostQscApplyOK creates a PostQscApplyOK with default headers values
func NewPostQscApplyOK() *PostQscApplyOK {
	return &PostQscApplyOK{}
}

/*PostQscApplyOK handles this case with default header values.

OK
*/
type PostQscApplyOK struct {
	Payload int64
}

func (o *PostQscApplyOK) Error() string {
	return fmt.Sprintf("[POST /qsc/apply][%d] postQscApplyOK  %+v", 200, o.Payload)
}

func (o *PostQscApplyOK) GetPayload() int64 {
	return o.Payload
}

func (o *PostQscApplyOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
