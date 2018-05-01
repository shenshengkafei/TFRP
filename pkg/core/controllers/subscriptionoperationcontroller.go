//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package controllers

import (
	"net/http"

	restful "github.com/emicklei/go-restful"
)

// GetSubscriptionOperationController returns a 200
func GetSubscriptionOperationController(request *restful.Request, response *restful.Response) {
	response.WriteHeader(http.StatusOK)
}

// PutSubscriptionOperationController returns a 200
func PutSubscriptionOperationController(request *restful.Request, response *restful.Response) {
	response.WriteHeader(http.StatusOK)
}
