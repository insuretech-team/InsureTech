# OpenAPI Schema Descriptions

This directory contains markdown description files for OpenAPI schemas generated from proto definitions.

## Structure

```
descriptions/
├── dto/              # Data Transfer Objects (Requests/Responses)
│   └── insuretech/
│       ├── policy/services/v1/
│       │   ├── PolicyCreationRequest.md
│       │   └── PolicyCreationResponse.md
│       └── claims/services/v1/
│           └── ClaimEvaluationRequest.md
├── entity/           # Domain Entities
│   └── insuretech/
│       ├── policy/entity/v1/
│       │   └── Policy.md
│       └── claims/entity/v1/
│           └── Claim.md
└── event/            # Event Schemas
    └── insuretech/
        ├── policy/events/v1/
        │   └── PolicyCreatedEvent.md
        └── claims/events/v1/
            └── ClaimSubmittedEvent.md
```

## Description Template Format

Each description file follows this structure:

### Header
- **Overview**: Type, proto path, purpose
- **Proto Comment**: Original comment from proto file

### Description
- Detailed explanation of the schema
- Validation rules
- Business logic notes

### Fields
- Each field with type, requirement, and description
- Proto comments preserved
- Additional context added

### Usage Examples
- Request/Response examples in JSON
- Common use cases

### Notes
- Additional constraints
- Business rules
- Implementation notes

### Related
- Links to related schemas
- API endpoints using this schema
- Documentation references

## Writing Guidelines

### 1. Be Clear and Concise
- Use simple language
- Avoid jargon unless necessary
- Explain technical terms

### 2. Provide Context
- Explain why the schema exists
- Describe typical use cases
- Note any constraints or limitations

### 3. Include Examples
- Provide realistic JSON examples
- Show edge cases
- Include validation examples

### 4. Maintain Consistency
- Use consistent terminology
- Follow the template structure
- Keep similar schemas aligned

### 5. Update Regularly
- Mark last updated date
- Note changes in business logic
- Keep examples current

## Priority Schemas

### High Priority (Core Operations)
1. **Policy Management**
   - PolicyCreationRequest
   - PolicyModificationRequest
   - PolicyCancellationRequest
   - PolicyRenewalRequest

2. **Claims Processing**
   - ClaimSubmissionRequest
   - ClaimEvaluationRequest
   - ClaimSettlementRequest

3. **Partner Integration**
   - PartnerRegistrationRequest
   - PartnerVerificationRequest

4. **Payment Processing**
   - PaymentInitiationRequest
   - PaymentConfirmationRequest

### Medium Priority (Supporting Operations)
- User management schemas
- Document upload schemas
- Notification schemas
- Report generation schemas

### Low Priority (Internal/Admin)
- Audit schemas
- System configuration schemas
- Internal workflow schemas

## Description Loader Integration

The `description_loader.py` module loads these descriptions using a 3-tier priority:

1. **Markdown File**: Load from this directory (highest priority)
2. **Proto Comment**: Use comment from proto file
3. **Generated**: Auto-generate from field names

### Loading Logic
```python
# insuretech.policy.services.v1.PolicyCreationRequest
# Looks for: descriptions/dto/insuretech/policy/services/v1/PolicyCreationRequest.md
```

## Status

- **Generated**: 49 templates (25 DTOs, 12 Entities, 12 Events)
- **Completed**: 0 (pending manual completion)
- **In Progress**: 0

## Next Steps

1. Review generated templates
2. Fill in detailed descriptions (start with High Priority)
3. Add realistic examples
4. Add validation rules and business logic
5. Link to related documentation
6. Test descriptions in regenerated OpenAPI spec

---

**Created**: 2026-01-02  
**Generator**: description_template_generator.py  
**Total Templates**: 49
