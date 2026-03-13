// MSW mock server setup
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';
import { testResponses, testErrors } from './test-data';

// Define mock handlers for API endpoints
const handlers = [
  // Auth - Registration
  http.post('http://localhost:3000/v1/auth/register', async ({ request }) => {
    const body = await request.json() as any;
    if (body.mobile_number === '+8801712345678') {
      return HttpResponse.json(testResponses.registration.success, { status: 201 });
    }
    return HttpResponse.json(testResponses.registration.duplicate, { status: 409 });
  }),

  // Auth - Login
  http.post('http://localhost:3000/v1/auth/login', async ({ request }) => {
    const body = await request.json() as any;
    if (body.password === 'SecurePass123!') {
      return HttpResponse.json(testResponses.login.jwt, { status: 200 });
    }
    return HttpResponse.json(testResponses.login.invalid, { status: 401 });
  }),

  // Auth - Email Login
  http.post('http://localhost:3000/v1/auth/email/login', async ({ request }) => {
    const body = await request.json() as any;
    if (body.password === 'EmailPass789!') {
      return HttpResponse.json(testResponses.login.jwt, { status: 200 });
    }
    return HttpResponse.json(testResponses.login.invalid, { status: 401 });
  }),

  // Auth - Send OTP (Google custom method pattern)
  http.post('http://localhost:3000/v1/auth/otp\\:send', () => {
    return HttpResponse.json(testResponses.otp.sent, { status: 200 });
  }),

  // Auth - Verify OTP (Google custom method pattern)
  http.post('http://localhost:3000/v1/auth/otp\\:verify', async ({ request }) => {
    const body = await request.json() as any;
    if (body.otp_code === '123456') {
      return HttpResponse.json(testResponses.otp.verified, { status: 200 });
    }
    if (body.otp_code === '000000') {
      return HttpResponse.json(testResponses.otp.expired, { status: 408 });
    }
    return HttpResponse.json(testResponses.otp.invalid, { status: 400 });
  }),

  // Auth - Refresh Token (Google custom method pattern)
  http.post('http://localhost:3000/v1/auth/token\\:refresh', () => {
    return HttpResponse.json(testResponses.token.refreshed, { status: 200 });
  }),

  // Auth - Logout
  http.post('http://localhost:3000/v1/auth/logout', () => {
    return HttpResponse.json({ message: 'Logged out successfully' }, { status: 200 });
  }),

  // Auth - Change Password (Google custom method pattern)
  http.post('http://localhost:3000/v1/auth/password\\:change', () => {
    return HttpResponse.json(testResponses.password.changed, { status: 200 });
  }),

  // Auth - Reset Password (Google custom method pattern)
  http.post('http://localhost:3000/v1/auth/password\\:reset', () => {
    return HttpResponse.json(testResponses.password.reset, { status: 200 });
  }),

  // Auth - Validate Token (Google custom method pattern)
  http.post('http://localhost:3000/v1/auth/token\\:validate', () => {
    return HttpResponse.json(testResponses.token.validated, { status: 200 });
  }),

  // Auth - Get Session
  http.get('http://localhost:3000/v1/auth/sessions/:sessionId', () => {
    return HttpResponse.json(testResponses.session.details, { status: 200 });
  }),

  // Auth - List Sessions
  http.get('http://localhost:3000/v1/auth/users/:userId/sessions', () => {
    return HttpResponse.json(testResponses.session.list, { status: 200 });
  }),

  // Auth - Get Current Session
  http.get('http://localhost:3000/v1/auth/session/current', () => {
    return HttpResponse.json(testResponses.session.details, { status: 200 });
  }),

  // Auth - Revoke Session (DELETE method)
  http.delete('http://localhost:3000/v1/auth/sessions/:sessionId', () => {
    return HttpResponse.json(testResponses.session.revoked, { status: 200 });
  }),

  // Auth - Revoke Session (POST with custom method - alternative endpoint)
  http.post('http://localhost:3000/v1/auth/sessions/:sessionId\\:revoke', () => {
    return HttpResponse.json(testResponses.session.revoked, { status: 200 });
  }),

  // Auth - Revoke All Sessions
  http.post('http://localhost:3000/v1/auth/users/:userId/sessions:revoke-all', () => {
    return HttpResponse.json({ message: 'All sessions revoked' }, { status: 200 });
  }),

  // Auth - Email Registration
  http.post('http://localhost:3000/v1/auth/email/register', async ({ request }) => {
    console.log('[MSW] Email registration intercepted');
    return HttpResponse.json(testResponses.registration.success, { status: 201 });
  }),

  // Auth - Send Email OTP
  http.post('http://localhost:3000/v1/auth/email/otp:send', () => {
    return HttpResponse.json(testResponses.otp.sent, { status: 200 });
  }),

  // Auth - Verify Email
  http.post('http://localhost:3000/v1/auth/email:verify', () => {
    return HttpResponse.json({ verified: true }, { status: 200 });
  }),

  // Auth - Email Login
  http.post('http://localhost:3000/v1/auth/email:login', () => {
    return HttpResponse.json(testResponses.login.jwt, { status: 200 });
  }),

  // Auth - Request Password Reset
  http.post('http://localhost:3000/v1/auth/email/password:reset-request', () => {
    return HttpResponse.json({ message: 'Reset email sent' }, { status: 200 });
  }),

  // Auth - Reset Password by Email
  http.post('http://localhost:3000/v1/auth/email/password:reset', () => {
    return HttpResponse.json(testResponses.password.reset, { status: 200 });
  }),

  // Auth - Validate CSRF
  http.post('http://localhost:3000/v1/auth/csrf:validate', () => {
    return HttpResponse.json({ valid: true }, { status: 200 });
  }),

  // Policy - Create
  http.post('http://localhost:3000/v1/policies', () => {
    return HttpResponse.json(testResponses.policy.created, { status: 201 });
  }),

  // Policy - Get
  http.get('http://localhost:3000/v1/policies/:policyId', () => {
    return HttpResponse.json(testResponses.policy.details, { status: 200 });
  }),

  // Policy - Update
  http.put('http://localhost:3000/v1/policies/:policyId', () => {
    return HttpResponse.json(testResponses.policy.details, { status: 200 });
  }),

  // Policy - List User Policies
  http.get('http://localhost:3000/v1/users/:customerId/policies', () => {
    return HttpResponse.json({ policies: [testResponses.policy.details] }, { status: 200 });
  }),

  // Policy - Update
  http.patch('http://localhost:3000/v1/policies/:policyId', () => {
    return HttpResponse.json(testResponses.policy.details, { status: 200 });
  }),

  // Policy - Cancel (Google custom method pattern)
  http.post('http://localhost:3000/v1/policies/:policyId\\:cancel', () => {
    return HttpResponse.json({ message: 'Policy cancelled' }, { status: 200 });
  }),

  // Policy - Renew (Google custom method pattern)
  http.post('http://localhost:3000/v1/policies/:policyId\\:renew', () => {
    return HttpResponse.json(testResponses.policy.details, { status: 200 });
  }),

  // Policy - Issue (Google custom method pattern)
  http.post('http://localhost:3000/v1/policies/:policyId\\:issue', () => {
    return HttpResponse.json({ issued: true }, { status: 200 });
  }),

  // Product - List
  http.get('http://localhost:3000/v1/products', () => {
    return HttpResponse.json(testResponses.product.list, { status: 200 });
  }),

  // Product - Get
  http.get('http://localhost:3000/v1/products/:productId', () => {
    return HttpResponse.json(testResponses.product.list.products[0], { status: 200 });
  }),

  // Product - Calculate Premium (Google custom method pattern)
  http.post('http://localhost:3000/v1/products/:productId\\:calculate-premium', () => {
    return HttpResponse.json(testResponses.product.premium, { status: 200 });
  }),

  // Product - Search
  http.get('http://localhost:3000/v1/products:search', () => {
    return HttpResponse.json(testResponses.product.list, { status: 200 });
  }),

  // Claim - Submit
  http.post('http://localhost:3000/v1/claims', () => {
    return HttpResponse.json({ claim_id: 'claim_123', status: 'SUBMITTED' }, { status: 201 });
  }),

  // Claim - Get
  http.get('http://localhost:3000/v1/claims/:claimId', () => {
    return HttpResponse.json({ claim_id: 'claim_123', status: 'APPROVED' }, { status: 200 });
  }),

  // Claim - List User Claims
  http.get('http://localhost:3000/v1/users/:customerId/claims', () => {
    return HttpResponse.json({ claims: [{ claim_id: 'claim_123', status: 'APPROVED' }] }, { status: 200 });
  }),

  // Claim - Approve (Google custom method pattern)
  http.post('http://localhost:3000/v1/claims/:claimId\\:approve', () => {
    return HttpResponse.json({ approved: true }, { status: 200 });
  }),

  // Claim - Reject (Google custom method pattern)
  http.post('http://localhost:3000/v1/claims/:claimId\\:reject', () => {
    return HttpResponse.json({ rejected: true }, { status: 200 });
  }),

  // Claim - Settle (Google custom method pattern)
  http.post('http://localhost:3000/v1/claims/:claimId\\:settle', () => {
    return HttpResponse.json({ settled: true }, { status: 200 });
  }),
];

// Create and export the mock server
export const server = setupServer(...handlers);

// Add unhandled request logging
server.events.on('request:unhandled', ({ request }) => {
  console.log('[MSW] Unhandled request:', request.method, request.url);
});
