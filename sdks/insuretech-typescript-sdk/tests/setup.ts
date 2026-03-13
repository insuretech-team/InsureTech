// Test setup file - runs before all tests
import { beforeAll, afterEach, afterAll } from 'vitest';
import { server } from './helpers/mock-server';
import { client } from '../src/client.gen';

// Configure default client for tests
client.setConfig({
  baseUrl: 'http://localhost:3000',
});

// Start mock server before all tests
beforeAll(() => {
  server.listen({ onUnhandledRequest: 'warn' });
});

// Reset handlers after each test to ensure test isolation
afterEach(() => {
  server.resetHandlers();
});

// Clean up after all tests
afterAll(() => {
  server.close();
});
