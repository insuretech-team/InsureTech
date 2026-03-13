import { describe, it, expect } from 'vitest';
import { createInsureTechClient } from '../../src/client-wrapper';

describe('Client', () => {
  describe('createInsureTechClient', () => {
    it('should create client with API key', () => {
      const client = createInsureTechClient({
        apiKey: 'test_api_key',
      });
      expect(client).toBeDefined();
    });

    it('should create client with custom baseUrl', () => {
      const client = createInsureTechClient({
        apiKey: 'test_api_key',
        baseUrl: 'http://localhost:3000',
      });
      expect(client).toBeDefined();
    });

    it('should create client with custom headers', () => {
      const client = createInsureTechClient({
        apiKey: 'test_api_key',
        headers: {
          'X-Custom-Header': 'value',
        },
      });
      expect(client).toBeDefined();
    });
  });
});
