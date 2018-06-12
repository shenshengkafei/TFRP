//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package entities

// AsyncOperationResult is the async operation result
type AsyncOperationResult struct {
	Status string `json:",omitempty"`
	Error  *ExtendedErrorInfo
}

// ExtendedErrorInfo is the extended error info
type ExtendedErrorInfo struct {
	Code    string `json:",omitempty"`
	Message string `json:",omitempty"`
}
