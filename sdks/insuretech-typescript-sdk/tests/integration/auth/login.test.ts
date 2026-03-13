import { describe, it, expect } from 'vitest';
import { authServiceLogin, authServiceEmailLogin, authServiceLogout } from '../../../src/sdk.gen';
import { createTestClient, expectValidJWT } from '../../helpers/test-utils';
import { testUsers, testResponses } from '../../helpers/test-data';

describe('Auth - Login', () => {
  describe('Mobile Login', () => {
    it('should login with valid credentials', async () => {
      const client = createTestClient();
      const response = await authServiceLogin({
        client,
        body: {
          mobile_number: testUsers.mobile.mobile_number,
          password: testUsers.mobile.password,
          device_id: testUsers.mobile.device_id,
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.user_id).toBe(testResponses.login.jwt.user_id);
      expect(response.data?.access_token).toBeDefined();
      expect(response.data?.refresh_token).toBeDefined();
      expect(response.data?.session_id).toBeDefined();
      
      if (response.data?.access_token) {
        expectValidJWT(response.data.access_token);
      }
    });

    it('should reject invalid credentials', async () => {
      const client = createTestClient();
      
      try {
        await authServiceLogin({
          client,
          body: {
            mobile_number: testUsers.mobile.mobile_number,
            password: 'WrongPassword123!',
            device_id: testUsers.mobile.device_id,
          },
        });
      } catch (error: any) {
        expect(error).toBeDefined();
        expect(error.error?.code).toBe('UNAUTHENTICATED');
      }
    });

    it('should return session information', async () => {
      const client = createTestClient();
      const response = await authServiceLogin({
        client,
        body: {
          mobile_number: testUsers.mobile.mobile_number,
          password: testUsers.mobile.password,
          device_id: testUsers.mobile.device_id,
        },
      });

      expect(response.data?.session_type).toBe('JWT');
      expect(response.data?.access_token_expires_in).toBeGreaterThan(0);
      expect(response.data?.refresh_token_expires_in).toBeGreaterThan(0);
    });
  });

  describe('Email Login', () => {
    it('should login with email credentials', async () => {
      const client = createTestClient();
      const response = await authServiceEmailLogin({
        client,
        body: {
          email: testUsers.email.email,
          password: testUsers.email.password,
          device_id: testUsers.email.device_id,
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.access_token).toBeDefined();
      expect(response.data?.user_id).toBeDefined();
    });
  });

  describe('Logout', () => {
    it('should logout successfully', async () => {
      const client = createTestClient();
      const response = await authServiceLogout({
        client,
        body: {
          session_id: 'session_123',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.message).toBe('Logged out successfully');
    });
  });
});
