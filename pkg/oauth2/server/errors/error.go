package errors

import "errors"

var (
	ErrInvalidRequest          = errors.New("invalid_request")
	ErrUnauthorizedClient      = errors.New("unauthorized_client")
	ErrAccessDenied            = errors.New("access_denied")
	ErrUnsupportedResponseType = errors.New("unsupported_response_type")
	ErrInvalidScope            = errors.New("invalid_scope")
	ErrServerError             = errors.New("server_error")
	ErrTemporarilyUnavailable  = errors.New("temporarily_unavailable")
	ErrInvalidClient           = errors.New("invalid_client")
	ErrInvalidGrant            = errors.New("invalid_grant")
	ErrUnsupportedGrantType    = errors.New("unsupported_grant_type")
	ErrInvalidMacAlgorithm   = errors.New("invalid_mac_algorithm")

	ErrInvalidRedirectURI   = errors.New("invalid redirect uri")
	ErrInvalidAuthorizeCode = errors.New("invalid authorize code")
	ErrInvalidAccessToken   = errors.New("invalid access token")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrExpiredAccessToken   = errors.New("expired access token")
	ErrExpiredRefreshToken  = errors.New("expired refresh token")
	ErrUnsupportedSignMethod = errors.New("unsupported sign method")
)

var Descriptions = map[error]string{
	ErrInvalidRequest:          "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed",
	ErrUnauthorizedClient:      "The client is not authorized to request an authorization code using this method",
	ErrAccessDenied:            "The resource owner or authorization server denied the request",
	ErrUnsupportedResponseType: "The authorization server does not support obtaining an authorization code using this method",
	ErrInvalidScope:            "The requested scope is invalid, unknown, or malformed",
	ErrServerError:             "The authorization server encountered an unexpected condition that prevented it from fulfilling the request",
	ErrTemporarilyUnavailable:  "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server",
	ErrInvalidClient:           "Client authentication failed",
	ErrInvalidGrant:            "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client",
	ErrUnsupportedGrantType:    "The authorization grant type is not supported by the authorization server",
	ErrInvalidMacAlgorithm: 	"The mac algorithm is no",
}

var StatusCodes = map[error]int{
	ErrInvalidRequest:          400,
	ErrUnauthorizedClient:      401,
	ErrAccessDenied:            403,
	ErrUnsupportedResponseType: 401,
	ErrInvalidScope:            400,
	ErrServerError:             500,
	ErrTemporarilyUnavailable:  503,
	ErrInvalidClient:           401,
	ErrInvalidGrant:            401,
	ErrUnsupportedGrantType:    401,
	ErrInvalidMacAlgorithm:     401,
}