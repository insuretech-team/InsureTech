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
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    status VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE partners.agents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    partner_id UUID NOT NULL REFERENCES partners.partners(id) ON DELETE CASCADE,
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
    category VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    base_premium BIGINT NOT NULL,
    base_premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    min_sum_insured BIGINT NOT NULL,
    min_sum_insured_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    max_sum_insured BIGINT NOT NULL,
    max_sum_insured_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    exclusions TEXT[],
    product_attributes JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.product_plans (
    plan_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES insurance_schema.products(product_id) ON DELETE CASCADE,
    plan_name VARCHAR(255) NOT NULL,
    premium_amount BIGINT NOT NULL,
    premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    min_sum_insured_amount BIGINT NOT NULL,
    min_sum_insured_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    max_sum_insured_amount BIGINT NOT NULL,
    max_sum_insured_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    attributes JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.product_riders (
    rider_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES insurance_schema.products(product_id) ON DELETE CASCADE,
    rider_name VARCHAR(255) NOT NULL,
    premium_amount BIGINT NOT NULL,
    premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    coverage_amount BIGINT NOT NULL,
    coverage_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
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

-- -----------------------------------------------------------------------------
-- 5. Beneficiaries Module (insurance_schema)
-- -----------------------------------------------------------------------------
CREATE TABLE insurance_schema.beneficiaries (
    beneficiary_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    kyc_status VARCHAR(50) NOT NULL,
    referral_code VARCHAR(20),
    referred_by UUID,
    audit_info JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.individual_beneficiaries (
    beneficiary_id UUID PRIMARY KEY REFERENCES insurance_schema.beneficiaries(beneficiary_id) ON DELETE CASCADE,
    gender VARCHAR(20),
    marital_status VARCHAR(20),
    contact_info JSONB,
    permanent_address JSONB,
    present_address JSONB,
    audit_info JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.business_beneficiaries (
    beneficiary_id UUID PRIMARY KEY REFERENCES insurance_schema.beneficiaries(beneficiary_id) ON DELETE CASCADE,
    business_type VARCHAR(50),
    contact_info JSONB,
    registered_address JSONB,
    business_address JSONB,
    audit_info JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- -----------------------------------------------------------------------------
-- 6. Underwriting Module (insurance_schema)
-- -----------------------------------------------------------------------------
CREATE TABLE insurance_schema.quotes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quote_number VARCHAR(50) NOT NULL UNIQUE,
    beneficiary_id UUID NOT NULL,
    insurer_product_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    sum_assured_amount BIGINT NOT NULL,
    sum_assured_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    base_premium_amount BIGINT NOT NULL,
    rider_premium_amount BIGINT,
    tax_amount BIGINT,
    total_premium_amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    premium_calculation JSONB,
    selected_riders JSONB,
    audit_info JSONB,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.health_declarations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quote_id UUID NOT NULL UNIQUE REFERENCES insurance_schema.quotes(id) ON DELETE CASCADE,
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
    quote_id UUID NOT NULL UNIQUE REFERENCES insurance_schema.quotes(id) ON DELETE CASCADE,
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
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_number VARCHAR(50) NOT NULL UNIQUE,
    product_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    partner_id UUID,
    agent_id UUID,
    status VARCHAR(50) NOT NULL,
    premium_amount BIGINT NOT NULL,
    premium_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    sum_insured_amount BIGINT NOT NULL,
    sum_insured_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    vat_tax_amount BIGINT,
    service_fee_amount BIGINT,
    total_payable_amount BIGINT,
    tenure_months INT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    issued_at TIMESTAMPTZ,
    proposer_details JSONB,
    underwriting_data JSONB,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.policy_nominees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES insurance_schema.policies(id) ON DELETE CASCADE,
    beneficiary_id UUID,
    full_name VARCHAR(200) NOT NULL,
    relationship VARCHAR(50) NOT NULL,
    date_of_birth DATE,
    nominee_dob_text VARCHAR(50),
    nid_number VARCHAR(20),
    phone_number VARCHAR(20),
    share_percentage DECIMAL(5,2) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.policy_riders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES insurance_schema.policies(id) ON DELETE CASCADE,
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
    policy_id UUID NOT NULL REFERENCES insurance_schema.policies(id) ON DELETE CASCADE,
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
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_number VARCHAR(50) NOT NULL UNIQUE,
    policy_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    processing_type VARCHAR(50) NOT NULL,
    claimed_amount BIGINT NOT NULL,
    claimed_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    approved_amount BIGINT,
    approved_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    settled_amount BIGINT,
    settled_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    incident_date TIMESTAMPTZ NOT NULL,
    incident_description TEXT NOT NULL,
    place_of_incident VARCHAR(255),
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_at TIMESTAMPTZ,
    settled_at TIMESTAMPTZ,
    rejection_reason TEXT,
    deductible_amount BIGINT,
    deductible_currency VARCHAR(3) NOT NULL DEFAULT 'BDT',
    co_pay_percentage DECIMAL(5,2),
    bank_details_for_payout TEXT,
    appeal_option_available BOOLEAN NOT NULL DEFAULT FALSE,
    in_app_messages JSONB,
    processor_notes TEXT,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE insurance_schema.claim_approvals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES insurance_schema.claims(id) ON DELETE CASCADE,
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
    claim_id UUID NOT NULL REFERENCES insurance_schema.claims(id) ON DELETE CASCADE,
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
    claim_id UUID NOT NULL UNIQUE REFERENCES insurance_schema.claims(id) ON DELETE CASCADE,
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
