import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  client: '@hey-api/client-fetch',
  input: '../../../api/openapi.yaml',
  output: {
    path: '../../insuretech-typescript-sdk/src',
    // format and lint are disabled during pipeline generation for speed.
    // The generated files are immediately post-processed by generator.go anyway.
    // Run `npx prettier --write src/` and `npx eslint src/ --fix` manually if needed.
    format: false,
    lint: false,
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
