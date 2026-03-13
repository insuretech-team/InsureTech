import { describe, it, expect } from 'vitest';
import * as SDK from '../../src/index';

describe('Exports', () => {
  it('should export createInsureTechClient', () => {
    expect(SDK.createInsureTechClient).toBeDefined();
    expect(typeof SDK.createInsureTechClient).toBe('function');
  });

  it('should export InsureTechClientConfig type', () => {
    // Type check - if this compiles, the type is exported
    const config: SDK.InsureTechClientConfig = {
      apiKey: 'test',
    };
    expect(config).toBeDefined();
  });

  it('should export all service functions', () => {
    // Auth services
    expect(SDK.authServiceRegister).toBeDefined();
    expect(SDK.authServiceLogin).toBeDefined();
    expect(SDK.authServiceLogout).toBeDefined();
    expect(SDK.authServiceRefreshToken).toBeDefined();
    
    // Policy services
    expect(SDK.policyServiceCreatePolicy).toBeDefined();
    expect(SDK.policyServiceGetPolicy).toBeDefined();
    expect(SDK.policyServiceListUserPolicies).toBeDefined();
    
    // Claim services
    expect(SDK.claimServiceSubmitClaim).toBeDefined();
    expect(SDK.claimServiceGetClaim).toBeDefined();
    
    // Product services
    expect(SDK.productServiceListProducts).toBeDefined();
    expect(SDK.productServiceCalculatePremium).toBeDefined();
  });

  it('should export all types', () => {
    // This test ensures types are exported (TypeScript compilation check)
    const typeCheck: SDK.AuthnRegisterRequest = {
      mobile_number: '+8801712345678',
      password: 'test',
      device_id: 'test',
      device_type: 'DEVICE_TYPE_MOBILE',
    };
    expect(typeCheck).toBeDefined();
  });
});
