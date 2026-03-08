# Beneficiary API DTOs

## Error model

All responses return a structured error object to standardize error handling.

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| code | string | Yes | Machine-readable error code (UPPER_SNAKE_CASE). |
| message | string | Yes | Human-readable error message. |
| details | object | No | Additional error details as key-value pairs. |
| field_violations | FieldViolation[] | No | Field-level validation errors. |
| retryable | boolean | Yes | Indicates whether retry is safe. |
| retry_after_seconds | integer | No | Suggested retry delay when retryable is true. |
| http_status_code | integer | Yes | HTTP status code equivalent. |
| error_id | string | Yes | Unique error instance ID. |
| documentation_url | string | No | Link to error documentation. |

### FieldViolation model

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| field | string | Yes | Field path (e.g., `applicant.date_of_birth`). |
| code | string | Yes | Field-level error code. |
| description | string | Yes | Human-readable description of the violation. |
| rejected_value | string | Yes | Invalid value (when safe to return). |

## Requests

### BeneficiariesListingRequest

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| type | string | Yes | Beneficiary type filter. |
| status | string | No | Status filter. |
| page | integer | Yes | Page number (>= 1). |
| page_size | integer | Yes | Page size (1-100). |

### BeneficiaryRetrievalRequest

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| beneficiary_id | string | Yes | Beneficiary identifier. |

### BeneficiaryUpdateRequest

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| beneficiary_id | string | Yes | Beneficiary identifier. |
| mobile_number | string | No | Updated mobile number. |
| email | string | Yes | Updated email. |
| address | string | No | Updated address. |

### BusinessBeneficiaryCreationRequest

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| user_id | string | Yes | User identifier. |
| business_name | string | No | Business name. |
| trade_license_number | string | No | Trade license number. |
| tin_number | string | No | TIN number. |
| focal_person_name | string | No | Focal person name. |
| focal_person_mobile | string | No | Focal person mobile. |
| partner_id | string | Yes | Partner identifier. |

### IndividualBeneficiaryCreationRequest

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| user_id | string | Yes | User identifier. |
| full_name | string | No | Full name. |
| date_of_birth | string | No | Date of birth. |
| gender | string | No | Gender. |
| nid_number | string | No | National ID number. |
| mobile_number | string | No | Mobile number. |
| email | string | Yes | Email address. |
| partner_id | string | Yes | Partner identifier. |

## Responses

### BeneficiariesListingResponse

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| beneficiaries | Beneficiary[] | Yes | Beneficiary list. |
| total_count | integer | Yes | Total results count. |
| error | Error | Yes | Structured error payload. |

### BeneficiaryRetrievalResponse

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| beneficiary | Beneficiary | Yes | Beneficiary summary. |
| individual_details | IndividualBeneficiary | Yes | Individual details. |
| business_details | BusinessBeneficiary | Yes | Business details. |
| error | Error | Yes | Structured error payload. |

### BeneficiaryUpdateResponse

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| message | string | Yes | Status message. |
| error | Error | Yes | Structured error payload. |

### BusinessBeneficiaryCreationResponse

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| beneficiary_id | string | Yes | Beneficiary identifier. |
| beneficiary_code | string | Yes | Beneficiary code. |
| message | string | Yes | Status message. |
| error | Error | Yes | Structured error payload. |

### IndividualBeneficiaryCreationResponse

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| beneficiary_id | string | Yes | Beneficiary identifier. |
| beneficiary_code | string | Yes | Beneficiary code. |
| message | string | Yes | Status message. |
| error | Error | Yes | Structured error payload. |
