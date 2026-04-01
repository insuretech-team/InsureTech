# Seed Data Validation Report

## Summary
Analysis of all SQL seed files and migration constraints for InsureTech database. This report identifies faulty seed data based on database constraints.

---

## Database Constraints Identified

### From Migrations:
1. **user_type** (authn_schema.users):
   - Valid values: `B2C_CUSTOMER`, `AGENT`, `BUSINESS_BENEFICIARY`, `SYSTEM_USER`, `B2B_ORG_ADMIN`
   - Source: `20250130_052_enhance_users_email_auth.up.sql` (lines 15-16)

2. **Email Required Constraint** (authn_schema.users):
   - Constraint: `chk_users_email_required_for_business`
   - Rule: BUSINESS_BENEFICIARY and SYSTEM_USER **MUST** have non-null, non-empty email
   - Source: `20250130_052_enhance_users_email_auth.up.sql` (lines 50-53)

3. **Email Format** (expected regex):
   - Pattern: `^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$`
   - Type: Standard RFC-compliant email validation

4. **Mobile Number Format** (expected regex):
   - Pattern: `^\+880[1][0-9]{9}$`
   - Format: +8801XXXXXXXXX (14 characters total)
   - Type: Bangladesh phone number validation

---

## Seed Files Analysis

### File 1: `20260318_001_seed_b2c_customer_default_role.sql`
**Location:** `E:\Projects\InsureTech\backend\inscore\db\seeds\authz_schema\`

**Content:** Seeds B2C customer default role and auto-assignment policies.
- Inserts into: `authz_schema.roles`, `authz_schema.casbin_rules`, `authz_schema.user_roles`
- No direct user/email/mobile data seeding
- **Status:** ✅ NO ISSUES

---

### File 2: `20260227_001_seed_portal_employees.sql`
**Location:** `E:\Projects\InsureTech\backend\inscore\db\seeds\b2b_schema\`

**Content Summary:**
- Seeds default B2B organisation
- Seeds departments (8 departments)
- Seeds 13 employee records (including 1 admin + 12 regular employees)
- Links admin user to organisation and roles

**Issues Found:**

#### Organisation Record (Lines 3-37):
```sql
contact_email: 'b2b-admin@lifeplus.local'
contact_phone: '+8801000000000'
```
- ⚠️ **ISSUE 1 - Invalid Mobile Number:**
  - Value: `+8801000000000`
  - Expected: `^\+880[1][0-9]{9}$` (14 chars total, 9 digits after +8801)
  - Found: 15 characters total (11 digits after +880)
  - **Status:** ❌ FAULTY - Too many digits (should be +8801XXXXXXXXX)

#### Employee Records (Lines 214-239):
All 12 seeded employees have:
- **MISSING fields:**
  - ❌ No `email` provided
  - ❌ No `mobile_number` provided
  - ✅ `user_id` field is NULL (acceptable)

**Detailed Employee Issues:**

| Employee ID | Name | Issue |
|---|---|---|
| 66666666-6666-6666-6666-666666666001 | John Doe | Missing email, mobile_number |
| 66666666-6666-6666-6666-666666666002 | Jane Smith | Missing email, mobile_number |
| 66666666-6666-6666-6666-666666666003 | Bob Johnson | Missing email, mobile_number |
| 66666666-6666-6666-6666-666666666004 | Alice Williams | Missing email, mobile_number |
| 66666666-6666-6666-6666-666666666005 | Rafiul Karim | Missing email, mobile_number |
| 66666666-6666-6666-6666-666666666006 | Nusrat Jahan | Missing email, mobile_number |
| 66666666-6666-6666-6666-666666666007 | Sabbir Hossain | Missing email, mobile_number |
| 66666666-6666-6666-6666-666666666008 | Mahi Rahman | Missing email, mobile_number |
| 66666666-6666-6666-6666-666666666009 | Tanvir Ahmed | Missing email, mobile_number |
| 66666666-6666-6666-6666-666666666010 | Farzana Akter | Missing email, mobile_number |
| 66666666-6666-6666-6666-666666666011 | Imran Kabir | Missing email, mobile_number |
| 66666666-6666-6666-6666-666666666012 | Sharmin Nahar | Missing email, mobile_number |

⚠️ **ISSUE 2 - Missing Required Fields:**
- All 12 employee records lack `email` and `mobile_number` data
- These may or may not be required fields depending on your business logic
- If contact information is expected, this represents incomplete seed data

---

### File 3: `20260301_002_seed_purchase_orders.sql`
**Location:** `E:\Projects\InsureTech\backend\inscore\db\seeds\b2b_schema\`

**Content:** Seeds 3 purchase order records.
- Inserts into: `b2b_schema.purchase_orders`
- No email or mobile number fields involved
- All required FK relationships reference valid seeded organisations/departments
- **Status:** ✅ NO ISSUES

---

## Critical Issues Summary

### 🔴 CRITICAL ISSUES:

1. **Invalid Mobile Number in Organisation Seed**
   - File: `20260227_001_seed_portal_employees.sql` (Line 22)
   - Value: `+8801000000000` (15 chars)
   - Expected: `+8801XXXXXXXXX` (14 chars, format: ^\+880[1][0-9]{9}$)
   - Impact: Will fail validation if mobile number constraint is enforced

### ⚠️ WARNINGS:

2. **Missing Email/Mobile in Employee Seeds**
   - File: `20260227_001_seed_portal_employees.sql` (Lines 214-239)
   - 12 employee records lack `email` and `mobile_number` fields
   - Severity: Depends on business requirements
   - If contact info is optional: Not an issue
   - If contact info is required: Major data quality issue

---

## Recommended Fixes

### Fix 1: Correct Organisation Mobile Number
```sql
-- In 20260227_001_seed_portal_employees.sql, line 22:
-- FROM:
contact_phone: '+8801000000000',
-- TO:
contact_phone: '+8801912345678',  -- or any valid +8801XXXXXXXXX format
```

### Fix 2: Add Employee Contact Information (if required)
Add email and mobile_number to each employee INSERT block:
```sql
-- Example for first employee (John Doe):
INSERT INTO b2b_schema.employees (
  ...existing fields...,
  email,
  mobile_number
)
VALUES (
  ...existing values...,
  'john.doe@example.com',
  '+8801912345670'
);
```

---

## Migration Constraints Review

### Reviewed Migration Files:
- ✅ `20250129_051_enhance_users.up.sql` - Indexes and triggers only
- ✅ `20250129_050_enhance_user_profiles.up.sql` - Indexes and triggers only
- ✅ `20250129_003_enhance_user_roles.up.sql` - Indexes only
- ✅ `20260301_004_enhance_employees.up.sql` - FK constraints, indexes, triggers, RLS
- ✅ `20250129_007_enhance_business_beneficiaries.up.sql` - Indexes and triggers only
- ✅ `20250130_052_enhance_users_email_auth.up.sql` - **KEY CONSTRAINTS FOUND**
  - Email required for BUSINESS_BENEFICIARY and SYSTEM_USER
  - User type validation
  - Email verification status tracking

### Constraint Validation:
- User types are validated in code (seed doesn't create users directly)
- Mobile number regex not found in migrations (likely app-level validation)
- Email regex not found in migrations (likely app-level validation)
- Email requirement is enforced at DB level for specific user types ✅

---

## Files Content Summary

### All SQL Files Reviewed:

**Seed Files (3 total):**
1. `20260318_001_seed_b2c_customer_default_role.sql` - 193 lines
2. `20260227_001_seed_portal_employees.sql` - 255 lines
3. `20260301_002_seed_purchase_orders.sql` - 85 lines

**Migration Files Analyzed (Sample):**
- Multiple enhancement files for various schemas
- No direct table creation DDL found (proto-generated)
- All provide indexes, triggers, constraints, and RLS policies

---

## Conclusion

**Total Issues Found: 2**
- **Critical:** 1 (Invalid mobile number format in organisation seed)
- **Warnings:** 1 (Missing contact data in 12 employee records - severity depends on requirements)

All other seed data and migration files are properly structured with no constraint violations detected.
