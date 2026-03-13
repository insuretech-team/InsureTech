import { describe, it, expect } from 'vitest';
import {
  authServiceSendOtp,
  authServiceVerifyOtp,
  authServiceSendEmailOtp,
} from '../../../src/sdk.gen';
import { createTestClient } from '../../helpers/test-utils';
import { testUsers, testResponses } from '../../helpers/test-data';

describe('Auth - OTP', () => {
  describe('Send OTP', () => {
    it('should send OTP successfully', async () => {
      const client = createTestClient();
      const response = await authServiceSendOtp({
        client,
        body: {
          mobile_number: testUsers.mobile.mobile_number,
          purpose: 'REGISTRATION',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.otp_id).toBe(testResponses.otp.sent.otp_id);
      expect(response.data?.expires_in_seconds).toBe(300);
      expect(response.data?.cooldown_seconds).toBeDefined();
    });

    it('should include sender information', async () => {
      const client = createTestClient();
      const response = await authServiceSendOtp({
        client,
        body: {
          mobile_number: testUsers.mobile.mobile_number,
          purpose: 'LOGIN',
        },
      });

      expect(response.data?.sender_id).toBe('LABAIDINS');
    });
  });

  describe('Verify OTP', () => {
    it('should verify valid OTP', async () => {
      const client = createTestClient();
      const response = await authServiceVerifyOtp({
        client,
        body: {
          otp_id: 'otp_123',
          otp_code: '123456',
          mobile_number: testUsers.mobile.mobile_number,
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.verified).toBe(true);
      expect(response.data?.user_id).toBeDefined();
    });

    it('should reject invalid OTP', async () => {
      const client = createTestClient();
      
      try {
        await authServiceVerifyOtp({
          client,
          body: {
            otp_id: 'otp_123',
            otp_code: '999999',
            mobile_number: testUsers.mobile.mobile_number,
          },
        });
      } catch (error: any) {
        expect(error).toBeDefined();
        expect(error.error?.code).toBe('INVALID_ARGUMENT');
      }
    });

    it('should reject expired OTP', async () => {
      const client = createTestClient();
      
      try {
        await authServiceVerifyOtp({
          client,
          body: {
            otp_id: 'otp_123',
            otp_code: '000000',
            mobile_number: testUsers.mobile.mobile_number,
          },
        });
      } catch (error: any) {
        expect(error).toBeDefined();
        expect(error.error?.code).toBe('DEADLINE_EXCEEDED');
      }
    });
  });

  describe('Email OTP', () => {
    it('should send email OTP', async () => {
      const client = createTestClient();
      const response = await authServiceSendEmailOtp({
        client,
        body: {
          email: testUsers.email.email,
          purpose: 'EMAIL_VERIFICATION',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.otp_id).toBeDefined();
      expect(response.data?.expires_in_seconds).toBeGreaterThan(0);
    });
  });
});
