import { describe, it, expect } from 'vitest';
import {
  authServiceRegister,
  authServiceVerifyOtp,
  authServiceLogin,
  productServiceListProducts,
  productServiceCalculatePremium,
  policyServiceCreatePolicy,
  policyServiceIssuePolicy,
  claimServiceSubmitClaim,
  authServiceLogout,
} from '../../src/sdk.gen';
import { createTestClient, expectValidJWT } from '../helpers/test-utils';
import { testUsers } from '../helpers/test-data';

describe('E2E - Complete User Journey', () => {
  it('should complete full insurance purchase flow', async () => {
    const client = createTestClient();

    // Step 1: Register new user
    const registerResponse = await authServiceRegister({
      client,
      body: testUsers.mobile,
    });
    expect(registerResponse.data?.user_id).toBeDefined();
    expect(registerResponse.data?.otp_sent).toBe(true);
    const userId = registerResponse.data?.user_id;
    const otpId = registerResponse.data?.otp_id;

    // Step 2: Verify OTP
    const verifyResponse = await authServiceVerifyOtp({
      client,
      body: {
        otp_id: otpId!,
        otp_code: '123456',
        mobile_number: testUsers.mobile.mobile_number,
      },
    });
    expect(verifyResponse.data?.verified).toBe(true);

    // Step 3: Login
    const loginResponse = await authServiceLogin({
      client,
      body: {
        mobile_number: testUsers.mobile.mobile_number,
        password: testUsers.mobile.password,
        device_id: testUsers.mobile.device_id,
      },
    });
    expect(loginResponse.data?.access_token).toBeDefined();
    expectValidJWT(loginResponse.data?.access_token!);
    const sessionId = loginResponse.data?.session_id;

    // Step 4: Browse products
    const productsResponse = await productServiceListProducts({ client });
    expect(productsResponse.data?.products).toBeDefined();
    expect(productsResponse.data?.products?.length).toBeGreaterThan(0);
    const productId = productsResponse.data?.products?.[0]?.product_id;

    // Step 5: Calculate premium
    const premiumResponse = await productServiceCalculatePremium({
      client,
      path: { product_id: productId! },
      body: {
        coverage_amount: 100000,
        age: 30,
        gender: 'MALE',
      },
    });
    expect(premiumResponse.data?.premium_amount).toBeGreaterThan(0);

    // Step 6: Create policy
    const policyResponse = await policyServiceCreatePolicy({
      client,
      body: {
        user_id: userId!,
        product_id: productId!,
      },
    });
    expect(policyResponse.data?.policy_id).toBeDefined();
    const policyId = policyResponse.data?.policy_id;

    // Step 7: Issue policy
    const issueResponse = await policyServiceIssuePolicy({
      client,
      path: { policy_id: policyId! },
      body: {},
    });
    expect(issueResponse.data?.issued).toBe(true);

    // Step 8: Submit claim
    const claimResponse = await claimServiceSubmitClaim({
      client,
      body: {
        policy_id: policyId!,
        claim_type: 'HEALTH',
        claim_amount: 50000,
        incident_date: '2024-01-01',
        description: 'Medical treatment',
      },
    });
    expect(claimResponse.data?.claim_id).toBeDefined();

    // Step 9: Logout
    const logoutResponse = await authServiceLogout({
      client,
      body: { session_id: sessionId! },
    });
    expect(logoutResponse.data?.message).toBe('Logged out successfully');
  });

  it('should handle authentication flow with session management', async () => {
    const client = createTestClient();

    // Login
    const loginResponse = await authServiceLogin({
      client,
      body: {
        mobile_number: testUsers.mobile.mobile_number,
        password: testUsers.mobile.password,
        device_id: testUsers.mobile.device_id,
      },
    });
    expect(loginResponse.data?.session_id).toBeDefined();

    // Browse products while authenticated
    const productsResponse = await productServiceListProducts({ client });
    expect(productsResponse.data?.products).toBeDefined();

    // Logout
    const logoutResponse = await authServiceLogout({
      client,
      body: { session_id: loginResponse.data?.session_id! },
    });
    expect(logoutResponse.data?.message).toBeDefined();
  });
});
