# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [common.proto](#common-proto)
    - [Common](#wallet-v1-insys-Common)
  
    - [CreateWallet](#wallet-v1-insys-CreateWallet)
    - [CreateWallet.Req](#wallet-v1-insys-CreateWallet-Req)
    - [CreateWallet.Res](#wallet-v1-insys-CreateWallet-Res)
  
- [Scalar Value Types](#scalar-value-types)



<a name="common-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## **common.proto**
Internal System gRPC Services

Common task related messages.
This proto contains everything related to common task services
accepted payload, returned response and method.

---
</br>


<a name="wallet-v1-insys-Common"></a>

### **Common**
Service for handling common or basic task on wallet.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateWallet | [CreateWallet.Req](#wallet-v1-insys-CreateWallet-Req) | [CreateWallet.Res](#wallet-v1-insys-CreateWallet-Res) | Used to creating or inserting a new wallet record. |

 <!-- end services -->

---
</br>


<a name="wallet-v1-insys-CreateWallet"></a>

### **CreateWallet**
Represents the create wallet request and response standard message.






<a name="wallet-v1-insys-CreateWallet-Req"></a>

### **CreateWallet.Req**
Accepted request payload for create wallet service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| branch_id | [uint32](#uint32) |  | ID of branch. |
| member_id | [uint64](#uint64) |  | ID of member. |
| p_id | [string](#string) |  | pID of member. |
| currency | [string](#string) |  | Currency of member. |
| username | [string](#string) |  | Username of member. |






<a name="wallet-v1-insys-CreateWallet-Res"></a>

### **CreateWallet.Res**
Response message that will be returned from create wallet service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  | Status of the response. |
| code | [int32](#int32) |  | Code of the response. |
| error | [google.protobuf.StringValue](#google-protobuf-StringValue) |  | Error info. |
| data | [string](#string) |  | Data of the response. |





 <!-- end messages -->

---
</br>

 <!-- end enums -->

 <!-- end HasExtensions -->



---
</br>

## Scalar Value Types

| .proto Type | Notes | Go | PHP |
| ----------- | ----- | -- | --- |
| <a name="double" /> double |  | float64 | float|
| <a name="float" /> float |  | float32 | float|
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | integer|
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | integer/string|
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | integer|
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | integer/string|
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | integer|
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | integer/string|
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | integer|
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | integer/string|
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | integer|
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | integer/string|
| <a name="bool" /> bool |  | bool | boolean|
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | string|
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | []byte | string|
