// Test utility functions
import { expect } from 'vitest';
import { createInsureTechClient } from '../../src/client-wrapper';

/**
 * Create a test client with mock configuration
 */
export function createTestClient() {
  return createInsureTechClient({
    apiKey: 'test_api_key',
    baseUrl: 'http://localhost:3000',
  });
}

/**
 * Validate JWT token format
 */
export function expectValidJWT(token: string) {
  expect(token).toBeDefined();
  expect(typeof token).toBe('string');
  expect(token.split('.')).toHaveLength(3);
}

/**
 * Validate UUID format
 */
export function expectValidUUID(id: string) {
  expect(id).toBeDefined();
  expect(id).toMatch(/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i);
}

/**
 * Validate phone number format (Bangladesh)
 */
export function expectValidPhoneNumber(phone: string) {
  expect(phone).toBeDefined();
  expect(phone).toMatch(/^\+8801[3-9]\d{8}$/);
}

/**
 * Validate email format
 */
export function expectValidEmail(email: string) {
  expect(email).toBeDefined();
  expect(email).toMatch(/^[^\s@]+@[^\s@]+\.[^\s@]+$/);
}

/**
 * Validate ISO 8601 timestamp
 */
export function expectValidTimestamp(timestamp: string) {
  expect(timestamp).toBeDefined();
  expect(new Date(timestamp).toISOString()).toBe(timestamp);
}

/**
 * Validate error response structure
 */
export function expectValidError(error: any) {
  expect(error).toBeDefined();
  expect(error).toHaveProperty('code');
  expect(error).toHaveProperty('message');
  expect(typeof error.code).toBe('string');
  expect(typeof error.message).toBe('string');
}

/**
 * Wait for a specified time (for testing async operations)
 */
export function wait(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms));
}

/**
 * Generate random test data
 */
export function generateTestData() {
  const timestamp = Date.now();
  return {
    userId: `user_test_${timestamp}`,
    sessionId: `session_test_${timestamp}`,
    deviceId: `device_test_${timestamp}`,
    otpId: `otp_test_${timestamp}`,
    policyId: `pol_test_${timestamp}`,
    claimId: `claim_test_${timestamp}`,
  };
}

/**
 * Mock successful response
 */
export function mockSuccess<T>(data: T) {
  return {
    data,
    status: 200,
    statusText: 'OK',
  };
}

/**
 * Mock error response
 */
export function mockError(code: string, message: string, status: number = 400) {
  return {
    error: {
      code,
      message,
    },
    status,
    statusText: getStatusText(status),
  };
}

function getStatusText(status: number): string {
  const statusTexts: Record<number, string> = {
    400: 'Bad Request',
    401: 'Unauthorized',
    403: 'Forbidden',
    404: 'Not Found',
    409: 'Conflict',
    500: 'Internal Server Error',
  };
  return statusTexts[status] || 'Unknown';
}
