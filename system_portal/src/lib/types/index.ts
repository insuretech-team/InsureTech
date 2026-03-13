// Re-export SDK classes mapped for SvelteKit usage 
import type {
    ClaimStatus as SDKClaimStatus,
    ClaimType as SDKClaimType,
    PolicyStatus as SDKPolicyStatus,
    ProductsProductStatus as SDKProductStatus,
    ProductCategory as SDKProductCategory,
    ApprovalDecision as SDKApprovalDecision,
    Nominee as SDKNominee,
    PolicyRider as SDKPolicyRider,
    ProductsRider as SDKProductsRider
} from '@lifeplus/insuretech-sdk';

export * from '@lifeplus/insuretech-sdk';

export type Nominee = SDKNominee;
export type PolicyRider = SDKPolicyRider;
export type Rider = SDKProductsRider;

// Map Enum objects so existing components don't crash when comparing
export const ClaimStatus = {
    CLAIM_STATUS_SUBMITTED: 'CLAIM_STATUS_SUBMITTED' as SDKClaimStatus,
    CLAIM_STATUS_UNDER_REVIEW: 'CLAIM_STATUS_UNDER_REVIEW' as SDKClaimStatus,
    CLAIM_STATUS_PENDING_DOCUMENTS: 'CLAIM_STATUS_PENDING_DOCUMENTS' as SDKClaimStatus,
    CLAIM_STATUS_APPROVED: 'CLAIM_STATUS_APPROVED' as SDKClaimStatus,
    CLAIM_STATUS_REJECTED: 'CLAIM_STATUS_REJECTED' as SDKClaimStatus,
    CLAIM_STATUS_SETTLED: 'CLAIM_STATUS_SETTLED' as SDKClaimStatus,
    CLAIM_STATUS_DISPUTED: 'CLAIM_STATUS_DISPUTED' as SDKClaimStatus,
}

export const ClaimType = {
    CLAIM_TYPE_HEALTH_HOSPITALIZATION: 'CLAIM_TYPE_HEALTH_HOSPITALIZATION' as SDKClaimType,
    CLAIM_TYPE_HEALTH_SURGERY: 'CLAIM_TYPE_HEALTH_SURGERY' as SDKClaimType,
    CLAIM_TYPE_MOTOR_ACCIDENT: 'CLAIM_TYPE_MOTOR_ACCIDENT' as SDKClaimType,
    CLAIM_TYPE_MOTOR_THEFT: 'CLAIM_TYPE_MOTOR_THEFT' as SDKClaimType,
    CLAIM_TYPE_TRAVEL_MEDICAL: 'CLAIM_TYPE_TRAVEL_MEDICAL' as SDKClaimType,
    CLAIM_TYPE_TRAVEL_BAGGAGE_LOSS: 'CLAIM_TYPE_TRAVEL_BAGGAGE_LOSS' as SDKClaimType,
    CLAIM_TYPE_DEVICE_DAMAGE: 'CLAIM_TYPE_DEVICE_DAMAGE' as SDKClaimType,
    CLAIM_TYPE_DEVICE_THEFT: 'CLAIM_TYPE_DEVICE_THEFT' as SDKClaimType,
}

export const PolicyStatus = {
    POLICY_STATUS_ACTIVE: 'POLICY_STATUS_ACTIVE' as SDKPolicyStatus,
    POLICY_STATUS_LAPSED: 'POLICY_STATUS_LAPSED' as SDKPolicyStatus,
    POLICY_STATUS_CANCELLED: 'POLICY_STATUS_CANCELLED' as SDKPolicyStatus,
}

export const ProductStatus = {
    PRODUCT_STATUS_DRAFT: 'PRODUCT_STATUS_DRAFT' as SDKProductStatus,
    PRODUCT_STATUS_ACTIVE: 'PRODUCT_STATUS_ACTIVE' as SDKProductStatus,
    PRODUCT_STATUS_INACTIVE: 'PRODUCT_STATUS_INACTIVE' as SDKProductStatus,
    PRODUCT_STATUS_DISCONTINUED: 'PRODUCT_STATUS_DISCONTINUED' as SDKProductStatus,
}

export const ProductCategory = {
    PRODUCT_CATEGORY_MOTOR: 'PRODUCT_CATEGORY_MOTOR' as SDKProductCategory,
    PRODUCT_CATEGORY_HEALTH: 'PRODUCT_CATEGORY_HEALTH' as SDKProductCategory,
    PRODUCT_CATEGORY_TRAVEL: 'PRODUCT_CATEGORY_TRAVEL' as SDKProductCategory,
    PRODUCT_CATEGORY_HOME: 'PRODUCT_CATEGORY_HOME' as SDKProductCategory,
    PRODUCT_CATEGORY_DEVICE: 'PRODUCT_CATEGORY_DEVICE' as SDKProductCategory,
    PRODUCT_CATEGORY_AGRICULTURAL: 'PRODUCT_CATEGORY_AGRICULTURAL' as SDKProductCategory,
    PRODUCT_CATEGORY_LIFE: 'PRODUCT_CATEGORY_LIFE' as SDKProductCategory,
}

export const ApprovalDecision = {
    APPROVAL_DECISION_PENDING: 'APPROVAL_DECISION_PENDING' as SDKApprovalDecision,
    APPROVAL_DECISION_APPROVED: 'APPROVAL_DECISION_APPROVED' as SDKApprovalDecision,
    APPROVAL_DECISION_REJECTED: 'APPROVAL_DECISION_REJECTED' as SDKApprovalDecision,
    APPROVAL_DECISION_ESCALATED: 'APPROVAL_DECISION_ESCALATED' as SDKApprovalDecision,
}
