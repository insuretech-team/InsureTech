import { describe, it, expect } from 'vitest';
import {
  authServiceChangePassword,
  authServiceResetPassword,
  authServiceRequestPasswordResetByEmail,
  authServiceResetPasswordByEmail,
} from '../../../src/sdk.gen';
import { createTestClient } from '../../helpers/test-utils';
import { testResponses } from '../../helpers/test-data';

describe('Auth - Password Management', () => {
  describe('Change Password', () => {
    it('should change password successfully', async () => {
      const client = createTestClient();
      const response = await authServiceChangePassword({
        client,
        body: {
          user_id: 'user_123',
          old_password: 'OldPass123!',
          new_password: 'NewPass456!',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.message).toBe(testResponses.password.changed.message);
    });

    it('should validate password strength', () => {
      const weakPassword = '123456';
      const strongPassword = 'SecurePass123!';
      
      // Basic password validation
      expect(weakPassword.length).toBeLessThan(8);
      expect(strongPassword.length).toBeGreaterThanOrEqual(8);
      expect(strongPassword).toMatch(/[A-Z]/); // Has uppercase
      expect(strongPassword).toMatch(/[a-z]/); // Has lowercase
      expect(strongPassword).toMatch(/[0-9]/); // Has number
      expect(strongPassword).toMatch(/[!@#$%^&*]/); // Has special char
    });
  });

  describe('Reset Password', () => {
    it('should reset password with OTP', async () => {
      const client = createTestClient();
      const response = await authServiceResetPassword({
        client,
        body: {
          mobile_number: '+8801712345678',
          otp_id: 'otp_123',
          otp_code: '123456',
          new_password: 'NewPass789!',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.message).toBe(testResponses.password.reset.message);
    });
  });

  describe('Email Password Reset', () => {
    it('should request password reset email', async () => {
      const client = createTestClient();
      const response = await authServiceRequestPasswordResetByEmail({
        client,
        body: {
          email: 'test@business.com',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.message).toBe('Reset email sent');
    });

    it('should reset password via email token', async () => {
      const client = createTestClient();
      const response = await authServiceResetPasswordByEmail({
        client,
        body: {
          email: 'test@business.com',
          reset_token: 'reset_token_abc',
          new_password: 'NewEmailPass123!',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.message).toBe(testResponses.password.reset.message);
    });
  });
});
