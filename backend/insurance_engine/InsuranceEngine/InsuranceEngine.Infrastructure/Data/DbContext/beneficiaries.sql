CREATE TABLE beneficiaries (
    beneficiary_id UUID PRIMARY KEY,
    policy_id UUID NOT NULL,
    type VARCHAR(32) NOT NULL,
    status VARCHAR(32) NOT NULL,
    gender VARCHAR(32),
    email VARCHAR(255),
    phone VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE beneficiary_individuals (
    beneficiary_id UUID PRIMARY KEY REFERENCES beneficiaries(beneficiary_id),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    date_of_birth DATE,
    national_id_number VARCHAR(100)
);

CREATE TABLE beneficiary_businesses (
    beneficiary_id UUID PRIMARY KEY REFERENCES beneficiaries(beneficiary_id),
    business_name VARCHAR(255),
    business_registration_number VARCHAR(100),
    contact_person_name VARCHAR(255),
    tax_identification_number VARCHAR(100)
);
