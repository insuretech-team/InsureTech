import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  client: '@hey-api/client-fetch',
  input: '../../../api/openapi.yaml',
  output: {
    path: '../../insuretech-typescript-sdk/src',
    format: 'prettier',
    lint: 'eslint',
  },
  types: {
    enums: 'javascript',
  },
  services: {
    asClass: true,
    name: '{{name}}Service', // Groups methods by OpenAPI tags (e.g., AuthService, PolicyService)
  },
  schemas: false,
});
