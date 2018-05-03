//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package apierror

import restful "github.com/emicklei/go-restful"

// New creates an ErrorResponse
func New(errorCategory ErrorCategory, errorCode ErrorCode, message string) *ErrorResponse {
	return &ErrorResponse{
		Body: Error{
			Code:     errorCode,
			Message:  message,
			Category: errorCategory,
		},
	}
}

// WriteErrorToResponse writes an ErrorResponse
func WriteErrorToResponse(resp *restful.Response, httpStatus int, errorCategory ErrorCategory, errorCode ErrorCode, message string) {
	err := New(errorCategory, errorCode, message)
	resp.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	resp.WriteError(httpStatus, err)
}

// WriteErrorToResponseWitAPIError writes an ErrorResponse
func WriteErrorToResponseWitAPIError(resp *restful.Response, httpStatus int, errorResponse *ErrorResponse) {
	resp.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	resp.WriteError(httpStatus, errorResponse)
}
