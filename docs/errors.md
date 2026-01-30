# Errors

Errors differ based on request. If the request is plain query (GET request) it will return error with GRPC error code.
GRPC error codes are defined in grpc [doc](https://pkg.go.dev/google.golang.org/grpc/codes).

If the request is transactional it will error with DCL's error codes.

DCL error codes are separated by modules:

1. [Auth](../x/dclauth/types/errors.go)
2. [Compliance](../x/compliance/types/errors.go)
3. [PKI](../types/pki/errors.go)
4. [Model](../x/model/types/errors.go)
5. [Validator](../x/validator/types/errors.go)
6. [VendorInfo](../x/vendorinfo/types/errors.go)
7. [DCL Upgrade](../x/dclupgrade/types/errors.go)