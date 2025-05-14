# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [internal/adapter/grpc/proto/user.proto](#internal_adapter_grpc_proto_user-proto)
    - [CreateUserRequest](#user-CreateUserRequest)
    - [CreateUserResponse](#user-CreateUserResponse)
    - [GetUserRequest](#user-GetUserRequest)
    - [GetUserResponse](#user-GetUserResponse)
    - [User](#user-User)
    - [User.ScopesEntry](#user-User-ScopesEntry)
    - [UserResponse](#user-UserResponse)
    - [UserResponse.ScopesEntry](#user-UserResponse-ScopesEntry)
  
    - [UserService](#user-UserService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="internal_adapter_grpc_proto_user-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## internal/adapter/grpc/proto/user.proto



<a name="user-CreateUserRequest"></a>

### CreateUserRequest
Request to create a user.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email | [string](#string) |  |  |
| password | [string](#string) |  |  |






<a name="user-CreateUserResponse"></a>

### CreateUserResponse
Response after creating a user.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [UserResponse](#user-UserResponse) |  |  |






<a name="user-GetUserRequest"></a>

### GetUserRequest
Request to get a user by ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="user-GetUserResponse"></a>

### GetUserResponse
Response containing the user data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [UserResponse](#user-UserResponse) |  |  |






<a name="user-User"></a>

### User
A user entity in the system.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | Unique user ID |
| name | [string](#string) |  | Full name |
| email | [string](#string) |  | Email address |
| password | [string](#string) |  | User password (should be hashed) |
| createdAt | [string](#string) |  | Timestamp of creation |
| updatedAt | [string](#string) |  | Timestamp of last update |
| scopes | [User.ScopesEntry](#user-User-ScopesEntry) | repeated | Permission scopes |






<a name="user-User-ScopesEntry"></a>

### User.ScopesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [int32](#int32) |  |  |






<a name="user-UserResponse"></a>

### UserResponse
User object returned to clients (without password).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| email | [string](#string) |  |  |
| createdAt | [string](#string) |  |  |
| updatedAt | [string](#string) |  |  |
| scopes | [UserResponse.ScopesEntry](#user-UserResponse-ScopesEntry) | repeated |  |






<a name="user-UserResponse-ScopesEntry"></a>

### UserResponse.ScopesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [int32](#int32) |  |  |





 

 

 


<a name="user-UserService"></a>

### UserService
Service for user management.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateUser | [CreateUserRequest](#user-CreateUserRequest) | [CreateUserResponse](#user-CreateUserResponse) | Create a new user. |
| GetUser | [GetUserRequest](#user-GetUserRequest) | [GetUserResponse](#user-GetUserResponse) | Retrieve a user by ID. |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

