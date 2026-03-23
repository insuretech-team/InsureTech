-- =============================================================================
-- Insurance Engine - PostgreSQL Database Setup Script
-- Aligned with Technical Contracts & Modular Monolith Architecture
-- =============================================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- -----------------------------------------------------------------------------
-- 1. Schemas
-- -----------------------------------------------------------------------------
CREATE SCHEMA IF NOT EXISTS insurance_schema;
CREATE SCHEMA IF NOT EXISTS partners;

-- -----------------------------------------------------------------------------
-- 2. Sequences
-- -----------------------------------------------------------------------------
CREATE SEQUENCE IF NOT EXISTS insurance_schema.policy_number_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE IF NOT EXISTS insurance_schema.endorsement_number_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE IF NOT EXISTS insurance_schema.quote_number_seq START WITH 1 INCREMENT BY 1;

-- -----------------------------------------------------------------------------
-- 3. Partners Module (partners schema)
-- -----------------------------------------------------------------------------
CREATE TABLE partners.partners (
    partner_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    trade_license VARCHAR(100) NOT NULL,
    tin_number VARCHAR(50),
    bank_account VARCHAR(100) NOT NULL,
    bank_name VARCHAR(255),
    bank_branch VARCHAR(255),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    acquisition_commission_rate DOUBLE PRECISION NOT NULL,
    renewal_commission_rate DOUBLE PRECISION NOT NULL,
    claims_assistance_rate DOUBLE PRECISION,
    onboarded_at TIMESTAMPTZ,
    focal_person_id UUID,
    commission JSONB,
    benefits JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE partners.agents (
    agent_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    partner_id UUID NOT NULL REFERENCES partners.partners(partner_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    status VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- -----------------------------------------------------------------------------
-- 4. Products Module (insurance_schema)
-- -----------------------------------------------------------------------------
CREATE TABLE insurance_schema.products (
    product_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_code VARCHAR(50) NOT NULL UNIQUE,
    product_name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    base_premium BIGINT NOT NULL,
    base_premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    min_sum_insured BIGINT NOT NULL,
    min_sum_insured_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    max_sum_insured BIGINT NOT NULL,
    max_sum_insured_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    min_tenure_months INT,
    max_tenure_months INT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE insurance_schema.product_plans (
    plan_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES insurance_schema.products(product_id) ON DELETE CASCADE,
    plan_name VARCHAR(255) NOT NULL,
    premium_amount BIGINT NOT NULL,
    premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    min_sum_insured BIGINT NOT NULL,
    min_sum_insured_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    max_sum_insured BIGINT NOT NULL,
    max_sum_insured_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    plan_description TEXT,
    attributes JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.product_riders (
    rider_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES insurance_schema.products(product_id) ON DELETE CASCADE,
    rider_name VARCHAR(255) NOT NULL,
    description TEXT,
    premium BIGINT NOT NULL,
    premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    coverage BIGINT NOT NULL,
    coverage_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    is_mandatory BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.pricing_configs (
    pricing_config_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL UNIQUE REFERENCES insurance_schema.products(product_id) ON DELETE CASCADE,
    rules JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.risk_assessment_questions (
    risk_assessment_question_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES insurance_schema.products(product_id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    question_text_bn TEXT,
    options_json JSONB DEFAULT '[]',
    weight INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- -----------------------------------------------------------------------------
-- 5. Beneficiaries Module (insurance_schema)
-- -----------------------------------------------------------------------------
CREATE TABLE insurance_schema.beneficiaries (
    beneficiary_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL, -- INDIVIDUAL, BUSINESS
    code VARCHAR(20) NOT NULL UNIQUE,
    status JSONB NOT NULL,
    kyc_status JSONB NOT NULL,
    kyc_completed_at TIMESTAMPTZ,
    risk_score VARCHAR(20),
    partner_id UUID,
    referral_code VARCHAR(50),
    referred_by UUID,
    audit_info JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE insurance_schema.individual_beneficiaries (
    individual_beneficiary_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    beneficiary_id UUID NOT NULL REFERENCES insurance_schema.beneficiaries(beneficiary_id) ON DELETE CASCADE,
    full_name VARCHAR(255) NOT NULL,
    full_name_bn VARCHAR(255),
    date_of_birth DATE NOT NULL,
    gender VARCHAR(20) NOT NULL,
    nid_number VARCHAR(50),
    passport_number VARCHAR(50),
    birth_certificate_number VARCHAR(50),
    tin_number VARCHAR(50),
    marital_status VARCHAR(50),
    occupation VARCHAR(100),
    nominee_name VARCHAR(255),
    nominee_relationship VARCHAR(100),
    contact_info JSONB,
    permanent_address JSONB,
    present_address JSONB,
    audit_info JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE insurance_schema.business_beneficiaries (
    business_beneficiary_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    beneficiary_id UUID NOT NULL REFERENCES insurance_schema.beneficiaries(beneficiary_id) ON DELETE CASCADE,
    business_name VARCHAR(255) NOT NULL,
    business_name_bn VARCHAR(255),
    trade_license_number VARCHAR(100) NOT NULL,
    trade_license_issue_date DATE,
    trade_license_expiry_date DATE,
    tin_number VARCHAR(50) NOT NULL,
    bin_number VARCHAR(50),
    registration_number VARCHAR(100),
    tax_id VARCHAR(100),
    business_type VARCHAR(100) NOT NULL,
    industry_sector VARCHAR(100),
    employee_count INT,
    incorporation_date DATE,
    focal_person_name VARCHAR(255) NOT NULL,
    focal_person_designation VARCHAR(100),
    focal_person_nid VARCHAR(50),
    focal_person_contact JSONB NOT NULL,
    contact_info JSONB NOT NULL, -- Contract maps registered_address -> contact_info column
    registered_address JSONB NOT NULL, -- Contract maps registered_address -> registered_address (implied)
    business_address JSONB NOT NULL,
    active_policies_count INT NOT NULL DEFAULT 0,
    pending_actions_count INT NOT NULL DEFAULT 0,
    total_employees_covered INT NOT NULL DEFAULT 0,
    total_premium_amount BIGINT NOT NULL DEFAULT 0,
    total_premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    audit_info JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- -----------------------------------------------------------------------------
-- 6. Underwriting Module (insurance_schema)
-- -----------------------------------------------------------------------------
CREATE TABLE insurance_schema.quotes (
    quote_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quote_number VARCHAR(50) NOT NULL UNIQUE,
    beneficiary_id UUID NOT NULL,
    insurer_product_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    sum_assured_amount BIGINT,
    sum_assured_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    term_years INT NOT NULL,
    premium_payment_mode VARCHAR(50) NOT NULL,
    base_premium_amount BIGINT,
    base_premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    rider_premium_amount BIGINT,
    rider_premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    tax_amount BIGINT,
    tax_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    total_premium_amount BIGINT,
    total_premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    premium_calculation JSONB,
    selected_riders JSONB,
    applicant_age INT NOT NULL,
    applicant_occupation VARCHAR(100),
    smoker BOOLEAN,
    valid_until TIMESTAMPTZ NOT NULL,
    converted_policy_id UUID,
    converted_at TIMESTAMPTZ,
    audit_info JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE insurance_schema.health_declarations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quote_id UUID NOT NULL UNIQUE REFERENCES insurance_schema.quotes(quote_id) ON DELETE CASCADE,
    weight_kg DECIMAL(5,2),
    bmi DECIMAL(5,2),
    pre_existing_conditions JSONB,
    family_history JSONB,
    medical_exam_results JSONB,
    medical_documents JSONB,
    audit_info JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.underwriting_decisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quote_id UUID NOT NULL UNIQUE REFERENCES insurance_schema.quotes(quote_id) ON DELETE CASCADE,
    decision VARCHAR(50) NOT NULL,
    method VARCHAR(50) NOT NULL,
    risk_level VARCHAR(50),
    risk_score DECIMAL(5,2),
    adjusted_premium_amount BIGINT,
    adjusted_premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    risk_factors JSONB,
    conditions JSONB,
    audit_info JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- -----------------------------------------------------------------------------
-- 7. Policy Module (insurance_schema)
-- -----------------------------------------------------------------------------
CREATE TABLE insurance_schema.policies (
    policy_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_number VARCHAR(50) NOT NULL UNIQUE,
    product_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    partner_id UUID,
    agent_id UUID,
    quote_id UUID,
    underwriting_decision_id UUID,
    status VARCHAR(50) NOT NULL,
    premium_amount BIGINT NOT NULL,
    premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    sum_insured_amount BIGINT NOT NULL,
    sum_insured_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    tenure_months INT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    issued_at TIMESTAMPTZ,
    policy_document_url TEXT,
    payment_frequency VARCHAR(50),
    vat_tax_amount BIGINT,
    vat_tax_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    service_fee_amount BIGINT,
    service_fee_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    total_payable_amount BIGINT,
    total_payable_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    payment_gateway_reference VARCHAR(100),
    receipt_number VARCHAR(100),
    occupation_risk_class VARCHAR(50),
    has_existing_policies BOOLEAN,
    claims_history_summary TEXT,
    provider_name VARCHAR(255),
    enrollment_start_date TIMESTAMPTZ,
    enrollment_end_date TIMESTAMPTZ,
    underwriting_data JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE insurance_schema.policy_nominees (
    nominee_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES insurance_schema.policies(policy_id) ON DELETE CASCADE,
    full_name VARCHAR(255) NOT NULL,
    relationship VARCHAR(100) NOT NULL,
    share_percentage DOUBLE PRECISION NOT NULL,
    date_of_birth TIMESTAMPTZ NOT NULL,
    nid_number VARCHAR(50),
    phone_number VARCHAR(50),
    nominee_dob_text VARCHAR(100),
    nominee_share_percent DOUBLE PRECISION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.policy_riders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES insurance_schema.policies(policy_id) ON DELETE CASCADE,
    rider_name VARCHAR(255) NOT NULL,
    premium_amount BIGINT NOT NULL,
    premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    coverage_amount BIGINT NOT NULL,
    coverage_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.endorsements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES insurance_schema.policies(policy_id) ON DELETE CASCADE,
    endorsement_number VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    reason VARCHAR(500),
    changes JSONB,
    audit_info JSONB,
    premium_adjustment_amount BIGINT,
    premium_adjustment_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- -----------------------------------------------------------------------------
-- 8. Claims Module (insurance_schema)
-- -----------------------------------------------------------------------------
CREATE TABLE insurance_schema.claims (
    claim_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_number VARCHAR(50) NOT NULL UNIQUE,
    policy_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    claimed_amount BIGINT NOT NULL,
    claimed_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    approved_amount BIGINT,
    approved_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    settled_amount BIGINT,
    settled_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    incident_date TIMESTAMPTZ NOT NULL,
    incident_description TEXT NOT NULL,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_at TIMESTAMPTZ,
    settled_at TIMESTAMPTZ,
    rejection_reason TEXT,
    place_of_incident VARCHAR(255),
    bank_details_for_payout TEXT,
    appeal_option_available BOOLEAN NOT NULL DEFAULT FALSE,
    in_app_messages JSONB,
    processing_type VARCHAR(50) NOT NULL,
    deductible_amount BIGINT,
    deductible_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    co_pay_amount BIGINT,
    co_pay_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    processor_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE insurance_schema.claim_approvals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES insurance_schema.claims(claim_id) ON DELETE CASCADE,
    approver_id UUID NOT NULL,
    approver_role VARCHAR(100),
    approval_level INT NOT NULL,
    decision VARCHAR(50) NOT NULL,
    approved_amount BIGINT,
    approved_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    notes TEXT,
    approved_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.claim_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES insurance_schema.claims(claim_id) ON DELETE CASCADE,
    document_type VARCHAR(100) NOT NULL,
    file_url TEXT NOT NULL,
    file_hash VARCHAR(64),
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.fraud_checks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL UNIQUE REFERENCES insurance_schema.claims(claim_id) ON DELETE CASCADE,
    fraud_score DECIMAL(5,2),
    flagged BOOLEAN NOT NULL DEFAULT FALSE,
    findings JSONB,
    checked_rules JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- -----------------------------------------------------------------------------
-- 9. Indexes
-- -----------------------------------------------------------------------------

-- Policies
CREATE INDEX idx_policies_customer_id ON insurance_schema.policies(customer_id);
CREATE INDEX idx_policies_product_id ON insurance_schema.policies(product_id);
CREATE INDEX idx_policies_status ON insurance_schema.policies(status);

-- Claims
CREATE INDEX idx_claims_policy_id ON insurance_schema.claims(policy_id);
CREATE INDEX idx_claims_customer_id ON insurance_schema.claims(customer_id);
CREATE INDEX idx_claims_status ON insurance_schema.claims(status);
CREATE INDEX idx_claims_incident_date ON insurance_schema.claims(incident_date);

-- Nominees
CREATE INDEX idx_policy_nominees_policy_id ON insurance_schema.policy_nominees(policy_id);
CREATE INDEX idx_policy_nominees_nid ON insurance_schema.policy_nominees(nid_number);

-- Products
CREATE INDEX idx_products_category ON insurance_schema.products(category);
CREATE INDEX idx_products_status ON insurance_schema.products(status);

-- Risk Assessment Questions
CREATE INDEX idx_risk_assessment_questions_product_id ON insurance_schema.risk_assessment_questions(product_id);

-- Quotes
CREATE INDEX idx_quotes_beneficiary_id ON insurance_schema.quotes(beneficiary_id);
CREATE INDEX idx_quotes_product_id ON insurance_schema.quotes(insurer_product_id);

-- Beneficiaries
CREATE INDEX idx_beneficiaries_type ON insurance_schema.beneficiaries(type);
CREATE INDEX idx_beneficiaries_status ON insurance_schema.beneficiaries(status);

-- Partners
CREATE INDEX idx_agents_partner_id ON partners.agents(partner_id);

-- -----------------------------------------------------------------------------
-- 10. Helper Functions (Optional)
-- -----------------------------------------------------------------------------
-- Example: Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to major tables
CREATE TRIGGER update_policies_updated_at BEFORE UPDATE ON insurance_schema.policies FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_claims_updated_at BEFORE UPDATE ON insurance_schema.claims FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON insurance_schema.products FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
