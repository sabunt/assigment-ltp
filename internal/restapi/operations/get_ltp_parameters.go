// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewGetLtpParams creates a new GetLtpParams object
//
// There are no default values defined in the spec.
func NewGetLtpParams() GetLtpParams {

	return GetLtpParams{}
}

// GetLtpParams contains all the bound params for the get ltp operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetLtp
type GetLtpParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*List of currency pairs (e.g., BTCUSD, BTCCHF, BTCEUR)
	  Required: true
	  In: query
	*/
	Pair []string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetLtpParams() beforehand.
func (o *GetLtpParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qPair, qhkPair, _ := qs.GetOK("pair")
	if err := o.bindPair(qPair, qhkPair, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindPair binds and validates array parameter Pair from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetLtpParams) bindPair(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("pair", "query", rawData)
	}
	var qvPair string
	if len(rawData) > 0 {
		qvPair = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	pairIC := swag.SplitByFormat(qvPair, "")
	if len(pairIC) == 0 {
		return errors.Required("pair", "query", pairIC)
	}

	var pairIR []string
	for _, pairIV := range pairIC {
		pairI := pairIV

		pairIR = append(pairIR, pairI)
	}

	o.Pair = pairIR

	return nil
}
