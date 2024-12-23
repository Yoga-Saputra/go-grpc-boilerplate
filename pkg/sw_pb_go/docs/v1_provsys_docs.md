# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [syncp.proto](#syncp-proto)
    - [Syncp](#wallet-v1-provsys-Syncp)
  
    - [SyncCredit](#wallet-v1-provsys-SyncCredit)
    - [SyncCredit.Data](#wallet-v1-provsys-SyncCredit-Data)
    - [SyncCredit.Req](#wallet-v1-provsys-SyncCredit-Req)
    - [SyncCredit.Res](#wallet-v1-provsys-SyncCredit-Res)
  
- [Scalar Value Types](#scalar-value-types)



<a name="syncp-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## **syncp.proto**
Provider System gRPC Services

Sync Credit Provider related messages.
This proto contains everything related to synchronize credit provider services
accepted payload, returned response and method.

---
</br>


<a name="wallet-v1-provsys-Syncp"></a>

### **Syncp**
Service for handling syncing member's credit that called by provider.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Credit | [SyncCredit.Req](#wallet-v1-provsys-SyncCredit-Req) | [SyncCredit.Res](#wallet-v1-provsys-SyncCredit-Res) | Used to getting member credit. |

 <!-- end services -->

---
</br>


<a name="wallet-v1-provsys-SyncCredit"></a>

### **SyncCredit**
Represents the sync credit request and response standard message.






<a name="wallet-v1-provsys-SyncCredit-Data"></a>

### **SyncCredit.Data**
Response message that wil be returned in response "data".


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| p_id | [string](#string) |  | pId of member. |
| currency | [string](#string) |  | Currency of member. |
| credit | [double](#double) |  | Credit amount of member. |
| last_update | [int64](#int64) |  | Last update credit info. |
| is_locked | [bool](#bool) |  | Lock wallet status of member. |
| is_disabled | [bool](#bool) |  | Disable wallet status of member. |






<a name="wallet-v1-provsys-SyncCredit-Req"></a>

### **SyncCredit.Req**
Accepted request payload.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| p_id | [string](#string) |  | pId of member. |
| currency | [string](#string) | optional | Currency of member. |
| provider_code | [string](#string) |  | provider code, e.g pgs. |






<a name="wallet-v1-provsys-SyncCredit-Res"></a>

### **SyncCredit.Res**
Response message that will be returned by sync credit services.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  | Status of the response. |
| code | [int32](#int32) |  | Code of the response. |
| error | [google.protobuf.StringValue](#google-protobuf-StringValue) |  | Error info. |
| data | [SyncCredit.Data](#wallet-v1-provsys-SyncCredit-Data) |  | Data of the response. |





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
