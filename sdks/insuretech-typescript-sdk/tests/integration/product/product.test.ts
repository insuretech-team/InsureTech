import { describe, it, expect } from 'vitest';
import {
  productServiceListProducts,
  productServiceGetProduct,
  productServiceCalculatePremium,
  productServiceSearchProducts,
} from '../../../src/sdk.gen';
import { createTestClient } from '../../helpers/test-utils';
import { testResponses } from '../../helpers/test-data';

describe('Product Service', () => {
  describe('List Products', () => {
    it('should list all products', async () => {
      const client = createTestClient();
      const response = await productServiceListProducts({
        client,
      });

      expect(response.data).toBeDefined();
      expect(response.data?.products).toBeDefined();
      expect(Array.isArray(response.data?.products)).toBe(true);
      expect(response.data?.products?.length).toBeGreaterThan(0);
    });

    it('should include product details', async () => {
      const client = createTestClient();
      const response = await productServiceListProducts({
        client,
      });

      const firstProduct = response.data?.products?.[0];
      expect(firstProduct?.product_id).toBeDefined();
      expect(firstProduct?.name).toBeDefined();
      expect(firstProduct?.category).toBeDefined();
      expect(firstProduct?.base_premium).toBeDefined();
    });
  });

  describe('Get Product', () => {
    it('should get product details', async () => {
      const client = createTestClient();
      const response = await productServiceGetProduct({
        client,
        path: {
          product_id: 'prod_1',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.product_id).toBeDefined();
      expect(response.data?.name).toBeDefined();
    });
  });

  describe('Calculate Premium', () => {
    it('should calculate premium', async () => {
      const client = createTestClient();
      const response = await productServiceCalculatePremium({
        client,
        path: {
          product_id: 'prod_1',
        },
        body: {
          coverage_amount: 100000,
          age: 30,
          gender: 'MALE',
        },
      });

      expect(response.data).toBeDefined();
      expect(response.data?.premium_amount).toBe(testResponses.product.premium.premium_amount);
      expect(response.data?.breakdown).toBeDefined();
      expect(response.data?.breakdown?.base_premium).toBeDefined();
      expect(response.data?.breakdown?.tax).toBeDefined();
    });

    it('should handle different coverage amounts', async () => {
      const client = createTestClient();
      const response = await productServiceCalculatePremium({
        client,
        path: {
          product_id: 'prod_1',
        },
        body: {
          coverage_amount: 200000,
          age: 35,
          gender: 'FEMALE',
        },
      });

      expect(response.data?.premium_amount).toBeGreaterThan(0);
    });
  });

  describe('Search Products', () => {
    it('should search products by category', async () => {
      const client = createTestClient();
      const response = await productServiceSearchProducts({
        client,
        query: {
          category: 'LIFE',
        },
      });

      expect(response.data).toBeDefined();
      expect(Array.isArray(response.data?.products)).toBe(true);
    });

    it('should search products by name', async () => {
      const client = createTestClient();
      const response = await productServiceSearchProducts({
        client,
        query: {
          query: 'Life Insurance',
        },
      });

      expect(response.data).toBeDefined();
    });
  });
});
