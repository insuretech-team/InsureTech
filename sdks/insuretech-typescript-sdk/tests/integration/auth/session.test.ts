import { describe, it, expect } from 'vitest';
import {
  authServiceGetSession,
  authServiceListSessions,
  authServiceGetCurrentSession,
  authServiceRevokeSession,
  authServiceRevokeAllSessions,
} from '../../../src/sdk.gen';
import { createTestClient } from '../../helpers/test-utils';
import { testResponses } from '../../helpers/test-data';

describe('Auth - Session Management', () => {
  describe('Get Session', () => {
    it('should get session details', async () => {
      const client = createTestClient();
      const response = await authServiceGetSession({
        client,
        path: {
          session_id: 'session_123',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.session_id).toBe(testResponses.session.details.session_id);
      expect(response.data?.user_id).toBeDefined();
      expect(response.data?.device_type).toBeDefined();
    });

    it('should include session timestamps', async () => {
      const client = createTestClient();
      const response = await authServiceGetSession({
        client,
        path: {
          session_id: 'session_123',
        },
      });

      expect(response.data?.created_at).toBeDefined();
      expect(response.data?.expires_at).toBeDefined();
      expect(response.data?.last_activity_at).toBeDefined();
    });
  });

  describe('List Sessions', () => {
    it('should list all user sessions', async () => {
      const client = createTestClient();
      const response = await authServiceListSessions({
        client,
        path: {
          user_id: 'user_123',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.sessions).toBeDefined();
      expect(Array.isArray(response.data?.sessions)).toBe(true);
      expect(response.data?.sessions?.length).toBeGreaterThan(0);
    });

    it('should include device information', async () => {
      const client = createTestClient();
      const response = await authServiceListSessions({
        client,
        path: {
          user_id: 'user_123',
        },
      });

      const firstSession = response.data?.sessions?.[0];
      expect(firstSession?.device_name).toBeDefined();
      expect(firstSession?.device_type).toBeDefined();
    });
  });

  describe('Get Current Session', () => {
    it('should get current session', async () => {
      const client = createTestClient();
      const response = await authServiceGetCurrentSession({
        client,
      });

      expect(response.data).toBeDefined();
      expect(response.data?.session_id).toBeDefined();
      expect(response.data?.user_id).toBeDefined();
    });
  });

  describe('Revoke Session', () => {
    it('should revoke specific session', async () => {
      const client = createTestClient();
      const response = await authServiceRevokeSession({
        client,
        path: {
          session_id: 'session_123',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.session_revoked).toBe(true);
      expect(response.data?.message).toBe('Session revoked successfully');
    });
  });

  describe('Revoke All Sessions', () => {
    it('should revoke all user sessions', async () => {
      const client = createTestClient();
      const response = await authServiceRevokeAllSessions({
        client,
        body: {
          user_id: 'user_123',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.message).toBe('All sessions revoked');
    });
  });
});
