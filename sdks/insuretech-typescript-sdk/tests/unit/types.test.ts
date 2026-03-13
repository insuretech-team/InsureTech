import { describe, it, expect } from 'vitest';
import type {
  AuthnRegisterRequest,
  AuthnLoginRequest,
  AuthnSendOtpRequest,
  AuthnVerifyOtpRequest,
  PolicyCreateRequest,
  ClaimSubmitRequest,
} from '../../src/types.gen';

describe('Types', () => {
  describe('Auth Types', () => {
    it('should validate AuthnRegisterRequest structure', () => {
      const request: AuthnRegisterRequest = {
        mobile_number: '+8801712345678',
        password: 'SecurePass123!',
        device_id: 'device_123',
        device_type: 'DEVICE_TYPE_MOBILE',
        device_name: 'iPhone 14',
      };
      expect(request.mobile_number).toBeDefined();
      expect(request.password).toBeDefined();
      expect(request.device_id).toBeDefined();
    });

    it('should validate AuthnLoginRequest structure', () => {
      const request: AuthnLoginRequest = {
        mobile_number: '+8801712345678',
        password: 'SecurePass123!',
        device_id: 'device_123',
      };
      expect(request.mobile_number).toBeDefined();
      expect(request.password).toBeDefined();
    });

    it('should validate AuthnSendOtpRequest structure', () => {
      const request: AuthnSendOtpRequest = {
        mobile_number: '+8801712345678',
        purpose: 'REGISTRATION',
      };
      expect(request.mobile_number).toBeDefined();
      expect(request.purpose).toBeDefined();
    });

    it('should validate AuthnVerifyOtpRequest structure', () => {
      const request: AuthnVerifyOtpRequest = {
        otp_id: 'otp_123',
        otp_code: '123456',
        mobile_number: '+8801712345678',
      };
      expect(request.otp_id).toBeDefined();
      expect(request.otp_code).toBeDefined();
    });
  });

  describe('Policy Types', () => {
    it('should validate PolicyCreateRequest structure', () => {
      const request: Partial<PolicyCreateRequest> = {
        product_id: 'prod_123',
        user_id: 'user_123',
      };
      expect(request.product_id).toBeDefined();
      expect(request.user_id).toBeDefined();
    });
  });

  describe('Claim Types', () => {
    it('should validate ClaimSubmitRequest structure', () => {
      const request: Partial<ClaimSubmitRequest> = {
        policy_id: 'pol_123',
        claim_type: 'HEALTH',
      };
      expect(request.policy_id).toBeDefined();
      expect(request.claim_type).toBeDefined();
    });
  });
});
