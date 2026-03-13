import { describe, it, expect } from 'vitest';
import {
  claimServiceSubmitClaim,
  claimServiceGetClaim,
  claimServiceListUserClaims,
  claimServiceApproveClaim,
  claimServiceRejectClaim,
  claimServiceSettleClaim,
} from '../../../src/sdk.gen';
import { createTestClient } from '../../helpers/test-utils';

describe('Claim Service', () => {
  describe('Submit Claim', () => {
    it('should submit new claim', async () => {
      const client = createTestClient();
      const response = await claimServiceSubmitClaim({
        client,
        body: {
          policy_id: 'pol_123',
          claim_type: 'HEALTH',
          claim_amount: 50000,
          incident_date: '2024-01-01',
          description: 'Medical treatment',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.claim_id).toBe('claim_123');
      expect(response.data?.status).toBe('SUBMITTED');
    });
  });

  describe('Get Claim', () => {
    it('should get claim details', async () => {
      const client = createTestClient();
      const response = await claimServiceGetClaim({
        client,
        path: {
          claim_id: 'claim_123',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.claim_id).toBe('claim_123');
      expect(response.data?.status).toBe('APPROVED');
    });
  });

  describe('List User Claims', () => {
    it('should list all user claims', async () => {
      const client = createTestClient();
      const response = await claimServiceListUserClaims({
        client,
        query: {
          user_id: 'user_123',
        },
      });

      expect(response.data).toBeDefined();
      expect(Array.isArray(response.data?.claims)).toBe(true);
    });
  });

  describe('Approve Claim', () => {
    it('should approve claim', async () => {
      const client = createTestClient();
      const response = await claimServiceApproveClaim({
        client,
        path: {
          claim_id: 'claim_123',
        },
        body: {
          approved_amount: 45000,
          notes: 'Approved after review',
        },
      });

      expect(response.data).toBeDefined();
    });
  });

  describe('Reject Claim', () => {
    it('should reject claim', async () => {
      const client = createTestClient();
      const response = await claimServiceRejectClaim({
        client,
        path: {
          claim_id: 'claim_123',
        },
        body: {
          reason: 'Insufficient documentation',
        },
      });

      expect(response.data).toBeDefined();
    });
  });

  describe('Settle Claim', () => {
    it('should settle claim', async () => {
      const client = createTestClient();
      const response = await claimServiceSettleClaim({
        client,
        path: {
          claim_id: 'claim_123',
        },
        body: {
          settlement_amount: 45000,
          payment_method: 'BANK_TRANSFER',
        },
      });

      expect(response.data).toBeDefined();
    });
  });
});
