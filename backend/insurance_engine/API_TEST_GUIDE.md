# Insurance Engine API Documentation

## Base URL
```
http://localhost:55877
```

## Product API Endpoints

### 1. List Products
**GET** `/v1/products`

**Query Parameters (Optional):**
| Parameter | Type | Description |
|-----------|------|-------------|
| category | string | Filter by category (e.g., MOTOR, LIFE, HEALTH) |
| status | string | Filter by status (e.g., ACTIVE, DRAFT, INACTIVE) |
| page | int | Page number (default: 1) |
| page_size | int | Items per page (default: 10) |

**Example:**
```bash
curl -s http://localhost:55877/v1/products
```

---

### 2. Get Product
**GET** `/v1/products/{product_id}`

**Example:**
```bash
curl -s http://localhost:55877/v1/products/91c65ff8-a834-413e-995c-a3f7677c4a1d
```

---

### 3. Create Product
**POST** `/v1/products`

**Request Body:**
```json
{
  "product": {
    "productCode": "LIFE-001",
    "productName": "Term Life Insurance",
    "category": "LIFE",
    "description": "Comprehensive term life insurance coverage",
    "basePremium": {
      "amount": 5000000,
      "currency": "BDT"
    },
    "minSumInsured": {
      "amount": 100000000,
      "currency": "BDT"
    },
    "maxSumInsured": {
      "amount": 5000000000,
      "currency": "BDT"
    },
    "minTenureMonths": 12,
    "maxTenureMonths": 60
  }
}
```

**Required Fields:**
| Field | Type | Description |
|-------|------|-------------|
| product.productCode | string | Unique product code |
| product.productName | string | Product name |
| product.category | string | Category (MOTOR, LIFE, HEALTH, TRAVEL, PROPERTY) |
| product.basePremium | object | Base premium with amount and currency |
| product.minSumInsured | object | Minimum sum insured |
| product.maxSumInsured | object | Maximum sum insured |
| product.minTenureMonths | int | Minimum tenure in months |
| product.maxTenureMonths | int | Maximum tenure in months |

**Test Command:**
```bash
curl -s -X POST http://localhost:55877/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "product": {
      "productCode": "LIFE-002",
      "productName": "Term Life Insurance",
      "category": "LIFE",
      "description": "Comprehensive term life insurance coverage",
      "basePremium": {"amount": 5000000, "currency": "BDT"},
      "minSumInsured": {"amount": 100000000, "currency": "BDT"},
      "maxSumInsured": {"amount": 5000000000, "currency": "BDT"},
      "minTenureMonths": 12,
      "maxTenureMonths": 60
    }
  }'
```

---

### 4. Update Product
**PATCH** `/v1/products/{product_id}`

**Request Body:**
```json
{
  "product_id": "91c65ff8-a834-413e-995c-a3f7677c4a1d",
  "product": {
    "productCode": "MOT-001",
    "productName": "Motor Insurance - Updated",
    "description": "Updated description",
    "basePremium": {"amount": 50000000, "currency": "BDT"},
    "minSumInsured": {"amount": 1000000000, "currency": "BDT"},
    "maxSumInsured": {"amount": 50000000000, "currency": "BDT"},
    "minTenureMonths": 12,
    "maxTenureMonths": 60
  }
}
```

**Test Command:**
```bash
curl -s -X PATCH http://localhost:55877/v1/products/91c65ff8-a834-413e-995c-a3f7677c4a1d \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "91c65ff8-a834-413e-995c-a3f7677c4a1d",
    "product": {
      "productCode": "MOT-001",
      "productName": "Motor Insurance - Updated",
      "description": "Updated description",
      "basePremium": {"amount": 50000000, "currency": "BDT"},
      "minSumInsured": {"amount": 1000000000, "currency": "BDT"},
      "maxSumInsured": {"amount": 50000000000, "currency": "BDT"},
      "minTenureMonths": 12,
      "maxTenureMonths": 60
    }
  }'
```

---

### 5. Activate Product
**POST** `/v1/products/{product_id}:activate`

**Request Body:** Empty `{}` required

**Test Command:**
```bash
curl -s -X POST http://localhost:55877/v1/products/91c65ff8-a834-413e-995c-a3f7677c4a1d:activate \
  -H "Content-Type: application/json" \
  -d '{}'
```

---

### 6. Deactivate Product
**POST** `/v1/products/{product_id}:deactivate`

**Request Body:** Empty `{}` required

**Test Command:**
```bash
curl -s -X POST http://localhost:55877/v1/products/91c65ff8-a834-413e-995c-a3f7677c4a1d:deactivate \
  -H "Content-Type: application/json" \
  -d '{}'
```

---

## Policy API Endpoints

### 1. List User Policies
**GET** `/v1/users/{customer_id}/policies`

**Query Parameters (Optional):**
| Parameter | Type | Description |
|-----------|------|-------------|
| status | string | Filter by policy status |
| page | int | Page number (default: 1) |
| page_size | int | Items per page (default: 10) |

**Valid User ID:**
```
00000000-0000-0000-0000-000000000001
```

**Test Command:**
```bash
curl -s "http://localhost:55877/v1/users/00000000-0000-0000-0000-000000000001/policies"
```

---

### 2. Get Policy
**GET** `/v1/policies/{policy_id}`

**Test Command:**
```bash
curl -s http://localhost:55877/v1/policies/cda3bf3a-ad2e-4b13-ba59-14c923d91c62
```

---

### 3. Create Policy
**POST** `/v1/policies`

**Request Body:**
```json
{
  "productId": "91c65ff8-a834-413e-995c-a3f7677c4a1d",
  "customerId": "00000000-0000-0000-0000-000000000001",
  "premiumAmount": {
    "amount": 5000000,
    "currency": "BDT"
  },
  "sumInsured": {
    "amount": 100000000,
    "currency": "BDT"
  },
  "tenureMonths": 12
}
```

**Required Fields:**
| Field | Type | Description |
|-------|------|-------------|
| productId | string (UUID) | Product ID from product list |
| customerId | string (UUID) | Customer/User ID |
| premiumAmount.amount | long | Premium amount in paisa (amount × 100) |
| premiumAmount.currency | string | Currency code (BDT) |
| sumInsured.amount | long | Sum insured in paisa |
| sumInsured.currency | string | Currency code (BDT) |
| tenureMonths | int | Policy tenure in months |

**Test Command:**
```bash
curl -s -X POST http://localhost:55877/v1/policies \
  -H "Content-Type: application/json" \
  -d '{
    "productId": "91c65ff8-a834-413e-995c-a3f7677c4a1d",
    "customerId": "00000000-0000-0000-0000-000000000001",
    "premiumAmount": {"amount": 5000000, "currency": "BDT"},
    "sumInsured": {"amount": 100000000, "currency": "BDT"},
    "tenureMonths": 12
  }'
```

---

### 4. Issue Policy
**POST** `/v1/policies/{policy_id}:issue`

**Request Body:**
```json
{
  "policyId": "cda3bf3a-ad2e-4b13-ba59-14c923d91c62",
  "quoteId": "00000000-0000-0000-0000-000000000000",
  "paymentId": "00000000-0000-0000-0000-000000000000"
}
```

**Required Fields:**
| Field | Type | Description |
|-------|------|-------------|
| policyId | string (UUID) | Policy ID to issue |
| quoteId | string (UUID) | Quote ID (can be placeholder) |
| paymentId | string (UUID) | Payment ID (can be placeholder) |

**Test Command:**
```bash
curl -s -X POST http://localhost:55877/v1/policies/cda3bf3a-ad2e-4b13-ba59-14c923d91c62:issue \
  -H "Content-Type: application/json" \
  -d '{
    "policyId": "cda3bf3a-ad2e-4b13-ba59-14c923d91c62",
    "quoteId": "00000000-0000-0000-0000-000000000000",
    "paymentId": "00000000-0000-0000-0000-000000000000"
  }'
```

---

## Beneficiary API Endpoints

### Valid User ID for Testing
```
00000000-0000-0000-0000-000000000001
```

### 1. List Beneficiaries
**GET** `/v1/beneficiaries`

**Query Parameters (Optional):**
| Parameter | Type | Description |
|-----------|------|-------------|
| type | string | Filter by type (INDIVIDUAL, BUSINESS) |
| status | string | Filter by status |
| page | int | Page number |
| pageSize | int | Items per page |

**Test Command:**
```bash
curl -s http://localhost:55877/v1/beneficiaries
```

---

### 2. Create Individual Beneficiary
**POST** `/v1/beneficiaries/individual`

**Request Body:**
```json
{
  "userId": "00000000-0000-0000-0000-000000000001",
  "fullName": "John Doe",
  "dateOfBirth": "1990-01-15",
  "gender": "Male",
  "nidNumber": "12345678901234567",
  "mobileNumber": "8801711123456",
  "email": "john@example.com"
}
```

**Required Fields:**
| Field | Type | Description |
|-------|------|-------------|
| userId | string (UUID) | User ID |
| fullName | string | Full name |
| dateOfBirth | string | Date of birth (YYYY-MM-DD) |
| gender | string | Gender (Male/Female/Other) |
| nidNumber | string | National ID number |
| mobileNumber | string | Mobile number |
| email | string | Email address |

**Test Command:**
```bash
curl -s -X POST http://localhost:55877/v1/beneficiaries/individual \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "00000000-0000-0000-0000-000000000001",
    "fullName": "John Doe",
    "dateOfBirth": "1990-01-15",
    "gender": "Male",
    "nidNumber": "12345678901234567",
    "mobileNumber": "8801711123456",
    "email": "john@example.com"
  }'
```

---

### 3. Create Business Beneficiary
**POST** `/v1/beneficiaries/business`

**Request Body:**
```json
{
  "userId": "00000000-0000-0000-0000-000000000001",
  "businessName": "Acme Corporation",
  "tradeLicenseNumber": "TL-123456",
  "tinNumber": "TIN-789012",
  "focalPersonName": "Jane Smith",
  "focalPersonMobile": "8801711987654"
}
```

**Required Fields:**
| Field | Type | Description |
|-------|------|-------------|
| userId | string (UUID) | User ID |
| businessName | string | Business name |
| tradeLicenseNumber | string | Trade license number |
| tinNumber | string | TIN number |
| focalPersonName | string | Contact person name |
| focalPersonMobile | string | Contact person mobile |

**Test Command:**
```bash
curl -s -X POST http://localhost:55877/v1/beneficiaries/business \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "00000000-0000-0000-0000-000000000001",
    "businessName": "Acme Corporation",
    "tradeLicenseNumber": "TL-123456",
    "tinNumber": "TIN-789012",
    "focalPersonName": "Jane Smith",
    "focalPersonMobile": "8801711987654"
  }'
```

---

## Quick Test Script

Copy and paste this script to test all endpoints:

```bash
# Base URL
BASE_URL="http://localhost:55877"

echo "=== Testing Product API ==="
echo "1. List Products:"
curl -s "$BASE_URL/v1/products"

echo -e "\n\n2. Get Product:"
curl -s "$BASE_URL/v1/products/91c65ff8-a834-413e-995c-a3f7677c4a1d"

echo -e "\n\n3. Create Product:"
curl -s -X POST "$BASE_URL/v1/products" \
  -H "Content-Type: application/json" \
  -d '{
    "product": {
      "productCode": "TEST-001",
      "productName": "Test Insurance",
      "category": "HEALTH",
      "description": "Test product",
      "basePremium": {"amount": 1000000, "currency": "BDT"},
      "minSumInsured": {"amount": 50000000, "currency": "BDT"},
      "maxSumInsured": {"amount": 500000000, "currency": "BDT"},
      "minTenureMonths": 6,
      "maxTenureMonths": 36
    }
  }'

echo -e "\n\n=== Testing Policy API ==="
echo "4. List User Policies:"
curl -s "$BASE_URL/v1/users/00000000-0000-0000-0000-000000000001/policies"

echo -e "\n\n5. Create Policy:"
curl -s -X POST "$BASE_URL/v1/policies" \
  -H "Content-Type: application/json" \
  -d '{
    "productId": "91c65ff8-a834-413e-995c-a3f7677c4a1d",
    "customerId": "00000000-0000-0000-0000-000000000001",
    "premiumAmount": {"amount": 3000000, "currency": "BDT"},
    "sumInsured": {"amount": 50000000, "currency": "BDT"},
    "tenureMonths": 12
  }'

echo -e "\n\n=== Testing Beneficiary API ==="
echo "6. List Beneficiaries:"
curl -s "$BASE_URL/v1/beneficiaries"

echo -e "\n\n7. Create Individual Beneficiary:"
curl -s -X POST "$BASE_URL/v1/beneficiaries/individual" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "00000000-0000-0000-0000-000000000001",
    "fullName": "Test User",
    "dateOfBirth": "1985-06-15",
    "gender": "Male",
    "nidNumber": "98765432109876543",
    "mobileNumber": "8801912345678",
    "email": "testuser@example.com"
  }'

echo -e "\n\n=== All Tests Complete ==="
```

---

## Common Issues & Solutions

### 1. UUID Format Error
**Error:** `operator does not exist: uuid = text`
**Solution:** Ensure the UUID is in valid format (e.g., `91c65ff8-a834-413e-995c-a3f7677c4a1d`)

### 2. Amount Value Issue
**Error:** Check constraint violation
**Solution:** Amount is in **paisa**, not taka. `5000 BDT = 500000` (multiply by 100)

### 3. Invalid Product ID
**Error:** Product not found or FK constraint error
**Solution:** Use a valid product ID from the product list endpoint

### 4. Invalid Customer ID
**Error:** Foreign key constraint violation
**Solution:** Use `00000000-0000-0000-0000-000000000001` for testing

### 5. Policy Number Format
**Error:** Check constraint violation on policy_number
**Solution:** Format must be `LBT-YYYY-XXXX-NNNNNN` (e.g., `LBT-2026-0001-123456`)

---

## Product Categories
| Category | Code |
|----------|------|
| Motor | MOTOR |
| Life | LIFE |
| Health | HEALTH |
| Travel | TRAVEL |
| Property | PROPERTY |

## Policy Status
| Status | Description |
|--------|-------------|
| DRAFT | Policy created but not issued |
| PENDING | Awaiting payment |
| ACTIVE | Policy is active |
| ISSUED | Policy has been issued |
| EXPIRED | Policy has expired |
| CANCELLED | Policy was cancelled |

## Beneficiary Types
| Type | Code |
|------|------|
| Individual | INDIVIDUAL |
| Business | BUSINESS |
