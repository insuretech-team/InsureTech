import fs from 'fs';
import path from 'path';

function walk(dir) {
    let results = [];
    const list = fs.readdirSync(dir);
    list.forEach(file => {
        file = path.resolve(dir, file);
        const stat = fs.statSync(file);
        if (stat && stat.isDirectory()) {
            results = results.concat(walk(file));
        } else if (file.match(/\.(svelte|ts)$/)) {
            results.push(file);
        }
    });
    return results;
}

const files = walk('./src');

// List of exact words to replace
const repl = {
    'claimId': 'claim_id',
    'claimNumber': 'claim_number',
    'policyId': 'policy_id',
    'customerId': 'customer_id',
    'claimedAmount': 'claimed_amount',
    'approvedAmount': 'approved_amount',
    'settledAmount': 'settled_amount',
    'incidentDate': 'incident_date',
    'incidentDescription': 'incident_description',
    'submittedAt': 'submitted_at',
    'approvedAt': 'approved_at',
    'settledAt': 'settled_at',
    'createdAt': 'created_at',
    'updatedAt': 'updated_at',
    'documentId': 'document_id',
    'documentType': 'document_type',
    'fileUrl': 'file_url',
    'fileHash': 'file_hash',
    'uploadedAt': 'uploaded_at',
    'verifiedBy': 'verified_by',
    'approvalId': 'approval_id',
    'approverId': 'approver_id',
    'approverRole': 'approver_role',
    'approvalLevel': 'approval_level',
    'rejectionReason': 'rejection_reason',
    'fraudCheck': 'fraud_check',
    'fraudCheckId': 'fraud_check_id',
    'fraudScore': 'fraud_score',
    'riskFactors': 'risk_factors',
    'reviewedBy': 'reviewed_by',
    'reviewedAt': 'reviewed_at',
    'productId': 'product_id',
    'productName': 'product_name',
    'productCode': 'product_code',
    'basePremium': 'base_premium',
    'minSumInsured': 'min_sum_insured',
    'maxSumInsured': 'max_sum_insured',
    'minTenureMonths': 'min_tenure_months',
    'maxTenureMonths': 'max_tenure_months',
    'availableRiders': 'available_riders',
    'riderName': 'rider_name',
    'premiumAmount': 'premium_amount',
    'coverageAmount': 'coverage_amount',
    'isMandatory': 'is_mandatory',
    'createdBy': 'created_by'
};

files.forEach(file => {
    let content = fs.readFileSync(file, 'utf8');
    let original = content;

    for (const [camel, snake] of Object.entries(repl)) {
        // We want to replace properties safely. 
        // Example: c.claimId -> c.claim_id
        // Example: claimId: 'clm' -> claim_id: 'clm'
        // Example: { claimId } -> { claim_id } 
        // We use a regex that matches word boundaries.
        const regex = new RegExp(`\\b${camel}\\b`, 'g');
        content = content.replace(regex, snake);
    }

    if (content !== original) {
        fs.writeFileSync(file, content);
        console.log('Fixed cases in:', file);
    }
});
