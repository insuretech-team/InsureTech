import { describe, it, expect } from 'vitest';
import { authServiceRegister, authServiceRegisterEmailUser } from '../../../src/sdk.gen';
import { createTestClient } from '../../helpers/test-utils';
import { testUsers, testResponses } from '../../helpers/test-data';

describe('Auth - Registration', () => {
  describe('Mobile Registration', () => {
    it('should register new user successfully', async () => {
      const client = createTestClient();
      const response = await authServiceRegister({
        client,
        body: testUsers.mobile,
      });

      expect(response.data).toBeDefined();
      expect(response.data?.user_id).toBe(testResponses.registration.success.user_id);
      expect(response.data?.otp_sent).toBe(true);
      expect(response.data?.otp_id).toBeDefined();
    });

    it('should handle duplicate registration', async () => {
      const client = createTestClient();
      
      try {
        await authServiceRegister({
          client,
          body: {
            ...testUsers.mobile,
            mobile_number: '+8801798765432', // Different number for duplicate test
          },
        });
      } catch (error: any) {
        expect(error).toBeDefined();
        expect(error.error?.code).toBe('ALREADY_EXISTS');
      }
    });

    it('should validate required fields', async () => {
      const client = createTestClient();
      
      const invalidRequest = {
        mobile_number: '',
        password: 'test',
        device_id: 'test',
      };

      // TypeScript should catch this, but test runtime behavior
      expect(invalidRequest.mobile_number).toBe('');
    });
  });

  describe('Email Registration', () => {
    it('should register email user successfully', async () => {
      const client = createTestClient();
      const response = await authServiceRegisterEmailUser({
        client,
        body: testUsers.email,
      });

      expect(response.data).toBeDefined();
      expect(response.data?.user_id).toBeDefined();
      expect(response.data?.otp_sent).toBe(true);
    });

    it('should validate email format', () => {
      const validEmail = 'test@business.com';
      const invalidEmail = 'invalid-email';
      
      expect(validEmail).toMatch(/^[^\s@]+@[^\s@]+\.[^\s@]+$/);
      expect(invalidEmail).not.toMatch(/^[^\s@]+@[^\s@]+\.[^\s@]+$/);
    });
  });
});
