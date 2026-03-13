import { describe, it, expect } from 'vitest';
import {
  policyServiceCreatePolicy,
  policyServiceGetPolicy,
  policyServiceUpdatePolicy,
  policyServiceListUserPolicies,
  policyServiceCancelPolicy,
  policyServiceRenewPolicy,
  policyServiceIssuePolicy,
} from '../../../src/sdk.gen';
import { createTestClient } from '../../helpers/test-utils';
import { testResponses } from '../../helpers/test-data';

describe('Policy Service', () => {
  describe('Create Policy', () => {
    it('should create new policy', async () => {
      const client = createTestClient();
      const response = await policyServiceCreatePolicy({
        client,
        body: {
          user_id: 'user_123',
          product_id: 'prod_456',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.policy_id).toBe(testResponses.policy.created.policy_id);
      expect(response.data?.policy_number).toBeDefined();
      expect(response.data?.status).toBe('DRAFT');
    });
  });

  describe('Get Policy', () => {
    it('should get policy details', async () => {
      const client = createTestClient();
      const response = await policyServiceGetPolicy({
        client,
        path: {
          policy_id: 'pol_123',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.policy_id).toBe(testResponses.policy.details.policy_id);
      expect(response.data?.status).toBe('ACTIVE');
      expect(response.data?.premium_amount).toBeDefined();
      expect(response.data?.coverage_amount).toBeDefined();
    });
  });

  describe('Update Policy', () => {
    it('should update policy', async () => {
      const client = createTestClient();
      const response = await policyServiceUpdatePolicy({
        client,
        path: {
          policy_id: 'pol_123',
        },
        body: {
          premium_amount: 6000,
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.policy_id).toBeDefined();
    });
  });

  describe('List User Policies', () => {
    it('should list all user policies', async () => {
      const client = createTestClient();
      const response = await policyServiceListUserPolicies({
        client,
        query: {
          user_id: 'user_123',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.policies).toBeDefined();
      expect(Array.isArray(response.data?.policies)).toBe(true);
    });
  });

  describe('Cancel Policy', () => {
    it('should cancel policy', async () => {
      const client = createTestClient();
      const response = await policyServiceCancelPolicy({
        client,
        path: {
          policy_id: 'pol_123',
        },
        body: {
          reason: 'Customer request',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.message).toBe('Policy cancelled');
    });
  });

  describe('Renew Policy', () => {
    it('should renew policy', async () => {
      const client = createTestClient();
      const response = await policyServiceRenewPolicy({
        client,
        path: {
          policy_id: 'pol_123',
        },
        body: {},
      });

      expect(response.data).toBeDefined();
      expect(response.data?.policy_id).toBeDefined();
    });
  });

  describe('Issue Policy', () => {
    it('should issue policy', async () => {
      const client = createTestClient();
      const response = await policyServiceIssuePolicy({
        client,
        path: {
          policy_id: 'pol_123',
        },
        body: {},
      });

      expect(response.data).toBeDefined();
      expect(response.data?.issued).toBe(true);
    });
  });
});
