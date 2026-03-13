// Test data fixtures for reusable test data

export const testUsers = {
  mobile: {
    mobile_number: '+8801712345678',
    password: 'SecurePass123!',
    device_id: 'test_device_mobile_123',
    device_type: 'DEVICE_TYPE_MOBILE',
    device_name: 'iPhone 14 Pro',
  },
  web: {
    mobile_number: '+8801798765432',
    password: 'WebPass456!',
    device_id: 'test_device_web_456',
    device_type: 'DEVICE_TYPE_WEB',
    device_name: 'Chrome on Windows',
  },
  email: {
    email: 'test@business.com',
    password: 'EmailPass789!',
    device_id: 'test_device_email_789',
    full_name: 'Test Business User',
    user_type: 'BUSINESS_BENEFICIARY',
  },
};

export const testResponses = {
  registration: {
    success: {
      user_id: 'user_reg_123',
      message: 'Registration successful',
      otp_sent: true,
      otp_id: 'otp_reg_456',
      otp_expires_in_seconds: 300,
    },
    duplicate: {
      error: {
        code: 'ALREADY_EXISTS',
        message: 'User already exists',
      },
    },
  },
  login: {
    jwt: {
      user_id: 'user_123',
      session_id: 'session_jwt_456',
      access_token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyXzEyMyJ9.test',
      refresh_token: 'refresh_jwt_token_xyz',
      access_token_expires_in: 3600,
      refresh_token_expires_in: 86400,
      session_type: 'JWT',
      user: {
        user_id: 'user_123',
        mobile_number: '+8801712345678',
        email: 'user@example.com',
      },
    },
    serverSide: {
      user_id: 'user_456',
      session_id: 'session_web_789',
      session_token: 'server_session_token_abc',
      csrf_token: 'csrf_token_def',
      session_type: 'SERVER_SIDE',
      user: {
        user_id: 'user_456',
        mobile_number: '+8801798765432',
      },
    },
    invalid: {
      error: {
        code: 'UNAUTHENTICATED',
        message: 'Invalid credentials',
      },
    },
  },
  otp: {
    sent: {
      otp_id: 'otp_123',
      message: 'OTP sent successfully',
      expires_in_seconds: 300,
      sender_id: 'LABAIDINS',
      cooldown_seconds: 60,
    },
    verified: {
      verified: true,
      user_id: 'user_123',
      message: 'OTP verified successfully',
    },
    invalid: {
      error: {
        code: 'INVALID_ARGUMENT',
        message: 'Invalid OTP code',
      },
    },
    expired: {
      error: {
        code: 'DEADLINE_EXCEEDED',
        message: 'OTP has expired',
      },
    },
  },
  session: {
    details: {
      session_id: 'session_123',
      user_id: 'user_123',
      device_id: 'device_123',
      device_type: 'DEVICE_TYPE_MOBILE',
      session_type: 'JWT',
      created_at: '2024-01-01T00:00:00Z',
      expires_at: '2024-01-02T00:00:00Z',
      last_activity_at: '2024-01-01T12:00:00Z',
    },
    list: {
      sessions: [
        {
          session_id: 'session_1',
          device_name: 'iPhone 14',
          device_type: 'DEVICE_TYPE_MOBILE',
          last_activity_at: '2024-01-01T12:00:00Z',
        },
        {
          session_id: 'session_2',
          device_name: 'Chrome on Mac',
          device_type: 'DEVICE_TYPE_WEB',
          last_activity_at: '2024-01-01T10:00:00Z',
        },
      ],
    },
    revoked: {
      message: 'Session revoked successfully',
      session_revoked: true,
    },
  },
  token: {
    refreshed: {
      access_token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyXzEyMyJ9.new_signature',
      refresh_token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyXzEyMyJ9.refresh_signature',
      access_token_expires_in: 3600,
      refresh_token_expires_in: 86400,
      session_id: 'session_123',
      session_expires_at: '2024-01-02T00:00:00Z',
    },
    validated: {
      valid: true,
      user_id: 'user_123',
      session_id: 'session_123',
      expires_at: '2024-01-02T00:00:00Z',
      session_type: 'JWT',
    },
  },
  password: {
    changed: {
      message: 'Password changed successfully',
    },
    reset: {
      message: 'Password reset successfully',
    },
  },
  policy: {
    created: {
      policy_id: 'pol_123',
      policy_number: 'POL-2024-001',
      status: 'DRAFT',
      user_id: 'user_123',
      product_id: 'prod_456',
    },
    details: {
      policy_id: 'pol_123',
      policy_number: 'POL-2024-001',
      status: 'ACTIVE',
      user_id: 'user_123',
      product_id: 'prod_456',
      premium_amount: 5000,
      coverage_amount: 100000,
    },
  },
  product: {
    list: {
      products: [
        {
          product_id: 'prod_1',
          name: 'Life Insurance Basic',
          category: 'LIFE',
          base_premium: 5000,
        },
        {
          product_id: 'prod_2',
          name: 'Health Insurance Premium',
          category: 'HEALTH',
          base_premium: 8000,
        },
      ],
    },
    premium: {
      premium_amount: 5500,
      breakdown: {
        base_premium: 5000,
        tax: 500,
      },
    },
  },
};

export const testErrors = {
  badRequest: {
    error: {
      code: 'INVALID_ARGUMENT',
      message: 'Invalid request parameters',
    },
  },
  unauthorized: {
    error: {
      code: 'UNAUTHENTICATED',
      message: 'Authentication required',
    },
  },
  forbidden: {
    error: {
      code: 'PERMISSION_DENIED',
      message: 'Insufficient permissions',
    },
  },
  notFound: {
    error: {
      code: 'NOT_FOUND',
      message: 'Resource not found',
    },
  },
  conflict: {
    error: {
      code: 'ALREADY_EXISTS',
      message: 'Resource already exists',
    },
  },
  serverError: {
    error: {
      code: 'INTERNAL',
      message: 'Internal server error',
    },
  },
};
