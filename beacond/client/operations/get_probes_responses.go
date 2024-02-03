// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"beacon/beacond/models"
)

// GetProbesReader is a Reader for the GetProbes structure.
type GetProbesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetProbesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetProbesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /probes] GetProbes", response, response.Code())
	}
}

// NewGetProbesOK creates a GetProbesOK with default headers values
func NewGetProbesOK() *GetProbesOK {
	return &GetProbesOK{}
}

/*
GetProbesOK describes a response with status code 200, with default header values.

OK
*/
type GetProbesOK struct {
	Payload *models.ServerListProbesResponse
}

// IsSuccess returns true when this get probes o k response has a 2xx status code
func (o *GetProbesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get probes o k response has a 3xx status code
func (o *GetProbesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get probes o k response has a 4xx status code
func (o *GetProbesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get probes o k response has a 5xx status code
func (o *GetProbesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get probes o k response a status code equal to that given
func (o *GetProbesOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get probes o k response
func (o *GetProbesOK) Code() int {
	return 200
}

func (o *GetProbesOK) Error() string {
	return fmt.Sprintf("[GET /probes][%d] getProbesOK  %+v", 200, o.Payload)
}

func (o *GetProbesOK) String() string {
	return fmt.Sprintf("[GET /probes][%d] getProbesOK  %+v", 200, o.Payload)
}

func (o *GetProbesOK) GetPayload() *models.ServerListProbesResponse {
	return o.Payload
}

func (o *GetProbesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ServerListProbesResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
