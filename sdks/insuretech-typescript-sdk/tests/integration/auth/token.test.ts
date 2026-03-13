import { describe, it, expect } from 'vitest';
import {
  authServiceRefreshToken,
  authServiceValidateToken,
  authServiceValidateCsrf,
} from '../../../src/sdk.gen';
import { createTestClient, expectValidJWT } from '../../helpers/test-utils';
import { testResponses } from '../../helpers/test-data';

describe('Auth - Token Management', () => {
  describe('Refresh Token', () => {
    it('should refresh access token', async () => {
      const client = createTestClient();
      const response = await authServiceRefreshToken({
        client,
        body: {
          refresh_token: 'refresh_token_xyz',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.access_token).toBeDefined();
      expect(response.data?.refresh_token).toBeDefined();
      
      if (response.data?.access_token) {
        expectValidJWT(response.data.access_token);
      }
    });

    it('should return new expiration times', async () => {
      const client = createTestClient();
      const response = await authServiceRefreshToken({
        client,
        body: {
          refresh_token: 'refresh_token_xyz',
        },
      });

      expect(response.data?.access_token_expires_in).toBe(3600);
      expect(response.data?.refresh_token_expires_in).toBe(86400);
      expect(response.data?.session_expires_at).toBeDefined();
    });
  });

  describe('Validate Token', () => {
    it('should validate valid token', async () => {
      const client = createTestClient();
      const response = await authServiceValidateToken({
        client,
        body: {
          access_token: 'valid_token_abc',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.valid).toBe(true);
      expect(response.data?.user_id).toBeDefined();
      expect(response.data?.session_id).toBeDefined();
    });

    it('should include token metadata', async () => {
      const client = createTestClient();
      const response = await authServiceValidateToken({
        client,
        body: {
          access_token: 'valid_token_abc',
        },
      });

      expect(response.data?.expires_at).toBeDefined();
      expect(response.data?.session_type).toBe('JWT');
    });
  });

  describe('Validate CSRF', () => {
    it('should validate CSRF token', async () => {
      const client = createTestClient();
      const response = await authServiceValidateCsrf({
        client,
        body: {
          csrf_token: 'csrf_token_def',
          session_id: 'session_123',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.valid).toBe(true);
    });
  });
});
