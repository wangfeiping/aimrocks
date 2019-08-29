// Code generated by go-swagger; DO NOT EDIT.

package qsc

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetQscApplyParams creates a new GetQscApplyParams object
// with the default values initialized.
func NewGetQscApplyParams() *GetQscApplyParams {
	var ()
	return &GetQscApplyParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetQscApplyParamsWithTimeout creates a new GetQscApplyParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetQscApplyParamsWithTimeout(timeout time.Duration) *GetQscApplyParams {
	var ()
	return &GetQscApplyParams{

		timeout: timeout,
	}
}

// NewGetQscApplyParamsWithContext creates a new GetQscApplyParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetQscApplyParamsWithContext(ctx context.Context) *GetQscApplyParams {
	var ()
	return &GetQscApplyParams{

		Context: ctx,
	}
}

// NewGetQscApplyParamsWithHTTPClient creates a new GetQscApplyParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetQscApplyParamsWithHTTPClient(client *http.Client) *GetQscApplyParams {
	var ()
	return &GetQscApplyParams{
		HTTPClient: client,
	}
}

/*GetQscApplyParams contains all the parameters to send to the API endpoint
for the get qsc apply operation typically these are written to a http.Request
*/
type GetQscApplyParams struct {

	/*Email
	  邮箱

	*/
	Email string
	/*Phone
	  手机号

	*/
	Phone string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get qsc apply params
func (o *GetQscApplyParams) WithTimeout(timeout time.Duration) *GetQscApplyParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get qsc apply params
func (o *GetQscApplyParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get qsc apply params
func (o *GetQscApplyParams) WithContext(ctx context.Context) *GetQscApplyParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get qsc apply params
func (o *GetQscApplyParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get qsc apply params
func (o *GetQscApplyParams) WithHTTPClient(client *http.Client) *GetQscApplyParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get qsc apply params
func (o *GetQscApplyParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithEmail adds the email to the get qsc apply params
func (o *GetQscApplyParams) WithEmail(email string) *GetQscApplyParams {
	o.SetEmail(email)
	return o
}

// SetEmail adds the email to the get qsc apply params
func (o *GetQscApplyParams) SetEmail(email string) {
	o.Email = email
}

// WithPhone adds the phone to the get qsc apply params
func (o *GetQscApplyParams) WithPhone(phone string) *GetQscApplyParams {
	o.SetPhone(phone)
	return o
}

// SetPhone adds the phone to the get qsc apply params
func (o *GetQscApplyParams) SetPhone(phone string) {
	o.Phone = phone
}

// WriteToRequest writes these params to a swagger request
func (o *GetQscApplyParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param email
	qrEmail := o.Email
	qEmail := qrEmail
	if qEmail != "" {
		if err := r.SetQueryParam("email", qEmail); err != nil {
			return err
		}
	}

	// query param phone
	qrPhone := o.Phone
	qPhone := qrPhone
	if qPhone != "" {
		if err := r.SetQueryParam("phone", qPhone); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
