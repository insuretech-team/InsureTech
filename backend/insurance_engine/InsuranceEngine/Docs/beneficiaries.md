# Beneficiary API & Data Model

## Enums

- **BusinessType**
  - BUSINESS_TYPE_UNSPECIFIED
  - BUSINESS_TYPE_SOLE_PROPRIETORSHIP
  - BUSINESS_TYPE_PARTNERSHIP
  - BUSINESS_TYPE_PRIVATE_LIMITED
  - BUSINESS_TYPE_PUBLIC_LIMITED
  - BUSINESS_TYPE_NGO
  - BUSINESS_TYPE_GOVERNMENT

- **MaritalStatus**
  - MARITAL_STATUS_UNSPECIFIED
  - MARITAL_STATUS_SINGLE
  - MARITAL_STATUS_MARRIED
  - MARITAL_STATUS_DIVORCED
  - MARITAL_STATUS_WIDOWED

- **BeneficiaryType**
  - BENEFICIARY_TYPE_UNSPECIFIED
  - BENEFICIARY_TYPE_INDIVIDUAL
  - BENEFICIARY_TYPE_BUSINESS

- **BeneficiaryGender**
  - GENDER_UNSPECIFIED
  - GENDER_MALE
  - GENDER_FEMALE
  - GENDER_OTHER

## Core Tables & Owned Types

### AuditInfo (owned)
| Field | Type | Required | Description |
| --- | --- | --- | --- |
| created_at | string (datetime) | Yes | Date and time the record was created |
| updated_at | string (datetime) | Yes | Date and time the record was last updated |
| created_by | string | Yes | User who created the record |
| updated_by | string | Yes | User who last updated the record |
| deleted_at | string (datetime) | No | Date and time the record was soft-deleted |
| deleted_by | string | No | User who performed the soft delete |

### ContactInfo (owned)
| Field | Type | Required | Description |
| --- | --- | --- | --- |
| mobile_number | string | Yes | e.g., +880 1XXX XXXXXX |
| email | string | No | Optional email address |
| alternate_mobile | string | No | Optional alternate mobile number |
| landline | string | No | Optional landline number |

### Address (owned)
| Field | Type | Required | Description |
| --- | --- | --- | --- |
| address_line1 | string | Yes | Primary address line |
| address_line2 | string | No | Secondary address line |
| city | string | Yes | City of residence |
| district | string | Yes | District of residence |
| division | string | Yes | Division (region) |
| postal_code | string | Yes | Postal code |
| country | string | Yes | Default: Bangladesh |
| latitude | number | Yes | Latitude of the location |
| longitude | number | Yes | Longitude of the location |

## Entities & Relationships

### Beneficiary (primary customer record)
| Field | Type | Required | Description |
| --- | --- | --- | --- |
| beneficiary_id | string | Yes | Unique beneficiary identifier |
| user_id | string | Yes | User identifier |
| type | BeneficiaryType | Yes | Individual or business |
| code | string | Yes | Beneficiary code |
| status | BeneficiaryStatus | Yes | Status enum |
| kyc_status | string | Yes | KYC status code |
| kyc_completed_at | string | No | KYC completion timestamp |
| risk_score | string | No | Risk scoring label |
| referral_code | string | No | Referral code |
| referred_by | string | No | Referrer identifier |
| partner_id | string | No | Partner identifier |
| audit_info | AuditInfo | Yes | Audit metadata |

Relationships:
- One Beneficiary to one IndividualBeneficiary (for individuals)
- One Beneficiary to one BusinessBeneficiary (for businesses)

### IndividualBeneficiary (individual beneficiary details)
| Field | Type | Required | Description |
| --- | --- | --- | --- |
| id | string | Yes | Unique identifier |
| beneficiary_id | string | Yes | Link to Beneficiary |
| full_name | string | Yes | Full name |
| full_name_bn | string | No | Bangla full name |
| date_of_birth | string | Yes | Date of birth |
| gender | BeneficiaryGender | Yes | Gender enum |
| nid_number | string | No | National ID |
| passport_number | string | No | Passport number |
| birth_certificate_number | string | No | Birth certificate number |
| tin_number | string | No | Tax identification number |
| marital_status | MaritalStatus | No | Marital status |
| occupation | string | No | Occupation |
| contact_info | ContactInfo | No | Contact details |
| permanent_address | Address | No | Permanent address |
| present_address | Address | No | Present address |
| nominee_name | string | No | Nominee name |
| nominee_relationship | string | No | Nominee relationship |
| audit_info | AuditInfo | Yes | Audit metadata |

### BusinessBeneficiary (business entity details)
| Field | Type | Required | Description |
| --- | --- | --- | --- |
| id | string | Yes | Unique identifier |
| beneficiary_id | string | Yes | Link to Beneficiary |
| business_name | string | Yes | Business legal name |
| business_name_bn | string | No | Bangla business name |
| trade_license_number | string | Yes | Trade license number |
| trade_license_issue_date | string | No | License issue date |
| trade_license_expiry_date | string | No | License expiry date |
| tin_number | string | Yes | Tax identification number |
| bin_number | string | No | Business identification number |
| business_type | BusinessType | Yes | Business type enum |
| industry_sector | string | No | Industry sector |
| employee_count | integer | No | Employee count |
| incorporation_date | string | No | Incorporation date |
| contact_info | ContactInfo | No | Business contact info |
| registered_address | Address | No | Registered address |
| business_address | Address | No | Business address |
| focal_person_name | string | Yes | Focal person name |
| focal_person_designation | string | No | Focal person designation |
| focal_person_nid | string | No | Focal person NID |
| focal_person_contact | ContactInfo | No | Focal person contact |
| audit_info | AuditInfo | Yes | Audit metadata |

## Endpoints (v1)

- `GET /v1/beneficiaries`
  - Pagination: `page`, `page_size`
- `GET /v1/beneficiaries/{beneficiaryId}`
  - Returns beneficiary profile with individual or business details
- `POST /v1/beneficiaries/individual`
  - Creates an individual beneficiary using the full schema
- `POST /v1/beneficiaries/business`
  - Creates a business beneficiary using the full schema
- `PATCH /v1/beneficiaries/{beneficiaryId}`
  - Updates beneficiary base data and nested individual/business details
