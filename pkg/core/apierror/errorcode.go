//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package apierror

// ErrorCategory indicates the kind of error
type ErrorCategory string

const (
	// ClientError is expected error
	ClientError ErrorCategory = "ClientError"

	// InternalError is system or internal error
	InternalError ErrorCategory = "InternalError"
)

// Common Azure Resource Provider API error code
type ErrorCode string

const (
	// From Microsoft.Azure.ResourceProvider.API.ErrorCode
	InvalidParameter                          ErrorCode = "InvalidParameter"
	BadRequest                                ErrorCode = "BadRequest"
	NotFound                                  ErrorCode = "NotFound"
	Conflict                                  ErrorCode = "Conflict"
	PreconditionFailed                        ErrorCode = "PreconditionFailed"
	OperationNotAllowed                       ErrorCode = "OperationNotAllowed"
	OperationPreempted                        ErrorCode = "OperationPreempted"
	PropertyChangeNotAllowed                  ErrorCode = "PropertyChangeNotAllowed"
	InternalOperationError                    ErrorCode = "InternalOperationError"
	InvalidSubscriptionStateTransition        ErrorCode = "InvalidSubscriptionStateTransition"
	UnregisterWithResourcesNotAllowed         ErrorCode = "UnregisterWithResourcesNotAllowed"
	InvalidParameterConflictingProperties     ErrorCode = "InvalidParameterConflictingProperties"
	SubscriptionNotRegistered                 ErrorCode = "SubscriptionNotRegistered"
	ConflictingUserInput                      ErrorCode = "ConflictingUserInput"
	ProvisioningInternalError                 ErrorCode = "ProvisioningInternalError"
	ProvisioningFailed                        ErrorCode = "ProvisioningFailed"
	NetworkingInternalOperationError          ErrorCode = "NetworkingInternalOperationError"
	QuotaExceeded                             ErrorCode = "QuotaExceeded"
	Unauthorized                              ErrorCode = "Unauthorized"
	ResourcesOverConstrained                  ErrorCode = "ResourcesOverConstrained"
	ControlPlaneProvisioningInternalError     ErrorCode = "ControlPlaneProvisioningInternalError"
	ControlPlaneProvisioningSyncError         ErrorCode = "ControlPlaneProvisioningSyncError"
	ControlPlaneInternalError                 ErrorCode = "ControlPlaneInternalError"
	ControlPlaneUnexpectedValue               ErrorCode = "ControlPlaneUnexpectedValue"
	ControlPlaneCloudProviderNotSet           ErrorCode = "ControlPlaneCloudProviderNotSet"
	ControlPlaneNotAvailable                  ErrorCode = "ControlPlaneNotAvailable"
	RegionOverConstrained                     ErrorCode = "RegionOverConstrained"
	DeleteResourceGroupInternalOperationError ErrorCode = "DeleteResourceGroupInternalOperationError"

	// From Microsoft.WindowsAzure.ContainerService.API.AcsErrorCode
	ScaleDownInternalError ErrorCode = "ScaleDownInternalError"

	// New
	PreconditionCheckTimeOut        ErrorCode = "PreconditionCheckTimeOut"
	UpgradeFailed                   ErrorCode = "UpgradeFailed"
	ScaleError                      ErrorCode = "ScaleError"
	CreateRoleAssignmentError       ErrorCode = "CreateRoleAssignmentError"
	ServicePrincipalNotFound        ErrorCode = "ServicePrincipalNotFound"
	ClusterResourceGroupNotFound    ErrorCode = "ClusterResourceGroupNotFound"
	KubeConfigError                 ErrorCode = "KubeConfigError"
	ServicePrincipalSecretTooLong   ErrorCode = "ServicePrincipalSecretTooLong"
	RegionNotSupported              ErrorCode = "RegionNotSupported"
	ControlPlaneAddOnsNotReady      ErrorCode = "ControlPlaneAddOnsNotReady"
	NodesNotFound                   ErrorCode = "NodesNotFound"
	NodesNotReady                   ErrorCode = "NodesNotReady"
	ControlPlaneProvisioningTimeout ErrorCode = "ControlPlaneProvisioningTimeout"

	// Error codes returned by HCP
	UnderlayNotFound         ErrorCode = "UnderlayNotFound"
	UnderlaysOverConstrained ErrorCode = "UnderlaysOverConstrained"
	UnexpectedUnderlayCount  ErrorCode = "UnexpectedUnderlayCount"
	ControlPlaneNotReady     ErrorCode = "ControlPlaneNotReady"
	TillerRecordNotFound     ErrorCode = "TillerRecordNotFound"
	RegionCreatesDisabled    ErrorCode = "RegionCreatesDisabled"
	RegionUpdatesDisabled    ErrorCode = "RegionUpdatesDisabled"

	// Geneva Action related error codes
	GenevaActionInternalError ErrorCode = "GenevaActionInternalError"

	// Error codes related to Addon
	AddonInvalid                        ErrorCode = "AddonInvalid"
	OmsAgentAddonConfigInvalid          ErrorCode = "OmsAgentAddonConfigInvalid"
	OmsAgentAddonConfigWorkspaceInvalid ErrorCode = "OmsAgentAddonConfigWorkspaceInvalid"
	GetLogAnalyticsWorkspaceError       ErrorCode = "GetLogAnalyticsWorkspaceError"
)

//Subcode strings
const (
	DeleteResourceGroupFailed = "DeleteResourceGroupFailed"
	DeleteEntityFailed        = "DeleteEntityFailed"
	GetEntityFailed           = "GetEntityFailed"
	OperationTimeout          = "OperationTimeout"
	PanicCaught               = "PanicCaught"
)
