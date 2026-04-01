#!/usr/bin/env python3
from __future__ import annotations

import argparse
import copy
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile
import uuid
from collections import defaultdict
from pathlib import Path
from typing import Any

import requests
import yaml
from dotenv import dotenv_values

ROOT_DIR = Path(__file__).resolve().parents[2]
API_DIR = ROOT_DIR / "api"
DEFAULT_SPEC_PATH = API_DIR / "openapi.yaml"
DEFAULT_OUTPUT_DIR = API_DIR / "postman"
DEFAULT_DOTENV_PATH = ROOT_DIR / ".env"
POSTMAN_API_BASE = "https://api.getpostman.com"
POSTMAN_COLLECTION_FILENAME = "InsureTech.postman_collection.json"
AUTH_SMOKE_FILENAME = "auth_smoke.postman_collection.json"
B2C_SUITE_FILENAME = "b2c_authn_authz_suite.postman_collection.json"

COLLECTION_NAME = "InsureTech API"
COLLECTION_SCHEMA = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"

HTTP_METHODS = ("get", "post", "put", "patch", "delete", "options", "head")
PLACEHOLDER_RE = re.compile(r"\{\{([^{}]+)\}\}")
SPACED_ACRONYM_RE = re.compile(r"\b(?:[A-Za-z]\s+){1,7}[A-Za-z]\b")

SECRET_KEYS = {
    "access_token",
    "refresh_token",
    "session_token",
    "csrf_token",
    "mfa_session_token",
    "device_credential",
    "totp_secret",
    "biometric_token",
    "login_password",
    "old_password",
    "new_password",
    "totp_code",
    "mobile_otp_code",
    "email_otp_code",
    "email_login_otp_code",
    "email_verification_otp_code",
    "otp_code",
    "api_key",
}

ENVIRONMENT_SPECS = [
    ("local", "InsureTech — Local"),
    ("staging", "InsureTech — Staging"),
    ("production", "InsureTech — Production"),
    ("mock", "InsureTech — Mock"),
    ("newman_test", "InsureTech — Newman Test"),
]

BASE_POSTMAN_VARIABLES: list[dict[str, str]] = [
    {"key": "base_url", "value": ""},
    {"key": "access_token", "value": ""},
    {"key": "refresh_token", "value": ""},
    {"key": "session_token", "value": ""},
    {"key": "csrf_token", "value": ""},
    {"key": "session_id", "value": ""},
    {"key": "session_type", "value": ""},
    {"key": "user_id", "value": ""},
    {"key": "user_mobile_number", "value": ""},
    {"key": "user_email", "value": ""},
    {"key": "full_name", "value": "Postman Test User"},
    {"key": "tenant_id", "value": ""},
    {"key": "device_id", "value": "postman-device"},
    {"key": "device_name", "value": "Postman Desktop"},
    {"key": "login_password", "value": ""},
    {"key": "old_password", "value": ""},
    {"key": "new_password", "value": ""},
    {"key": "login_device_type", "value": "ANDROID"},
    {"key": "login_expect_session_type", "value": "JWT"},
    {"key": "otp_id", "value": ""},
    {"key": "otp_code", "value": ""},
    {"key": "mfa_session_token", "value": ""},
    {"key": "device_credential", "value": ""},
    {"key": "totp_secret", "value": ""},
    {"key": "totp_code", "value": ""},
    {"key": "provisioning_uri", "value": ""},
    {"key": "biometric_token", "value": ""},
    {"key": "mobile_otp_type", "value": "login"},
    {"key": "mobile_otp_channel", "value": "sms"},
    {"key": "mobile_otp_use_masking", "value": "true"},
    {"key": "mobile_otp_id", "value": ""},
    {"key": "mobile_otp_code", "value": ""},
    {"key": "mobile_otp_verified", "value": "false"},
    {"key": "last_b2c_otp_login", "value": "false"},
    {"key": "last_b2c_device_credential_login", "value": "false"},
    {"key": "profile_full_name", "value": "B2C Test User"},
    {"key": "profile_full_name_updated", "value": "B2C Test User Updated"},
    {"key": "profile_date_of_birth", "value": "1990-01-01T00:00:00Z"},
    {"key": "profile_gender", "value": "MALE"},
    {"key": "profile_address_line1", "value": "House 10, Road 5"},
    {"key": "profile_address_line2", "value": "Apt 4B"},
    {"key": "profile_city", "value": "Dhaka"},
    {"key": "profile_city_updated", "value": "Chittagong"},
    {"key": "profile_district", "value": "Dhaka"},
    {"key": "profile_division", "value": "Dhaka"},
    {"key": "profile_country", "value": "Bangladesh"},
    {"key": "profile_nid_number", "value": "1234567890123"},
    {"key": "email_otp_type", "value": "email_login"},
    {"key": "email_otp_id", "value": ""},
    {"key": "email_otp_code", "value": ""},
    {"key": "email_login_otp_id", "value": ""},
    {"key": "email_login_otp_code", "value": ""},
    {"key": "email_verification_otp_id", "value": ""},
    {"key": "email_verification_otp_code", "value": ""},
    {"key": "email_verified", "value": "false"},
    {"key": "authz_domain", "value": "system:root"},
    {"key": "authz_object", "value": "svc:dashboard/read"},
    {"key": "authz_resource", "value": "dashboard"},
    {"key": "authz_action", "value": "read"},
    {"key": "api_key", "value": ""},
    {"key": "mock_sms_url", "value": ""},
    {"key": "policy_id", "value": ""},
    {"key": "claim_id", "value": ""},
    {"key": "payment_id", "value": ""},
    {"key": "order_id", "value": ""},
    {"key": "product_id", "value": ""},
    {"key": "quote_id", "value": ""},
    {"key": "ticket_id", "value": ""},
    {"key": "partner_id", "value": ""},
    {"key": "kyc_id", "value": ""},
    {"key": "invoice_id", "value": ""},
    {"key": "document_id", "value": ""},
    {"key": "proposal_id", "value": ""},
]

NON_STRING_PLACEHOLDER_KEYS = {
    "mobile_otp_use_masking",
}

DIRECT_FIELD_PLACEHOLDERS = {
    "access_token": "access_token",
    "refresh_token": "refresh_token",
    "session_token": "session_token",
    "csrf_token": "csrf_token",
    "session_id": "session_id",
    "session_type": "session_type",
    "user_id": "user_id",
    "tenant_id": "tenant_id",
    "otp_id": "otp_id",
    "device_id": "device_id",
    "device_name": "device_name",
    "mfa_session_token": "mfa_session_token",
    "device_credential": "device_credential",
    "totp_secret": "totp_secret",
    "totp_code": "totp_code",
    "provisioning_uri": "provisioning_uri",
    "biometric_token": "biometric_token",
    "api_key": "api_key",
    "api_key_id": "api_key_id",
    "domain": "authz_domain",
    "object": "authz_object",
    "policy_id": "policy_id",
    "claim_id": "claim_id",
    "payment_id": "payment_id",
    "order_id": "order_id",
    "product_id": "product_id",
    "quote_id": "quote_id",
    "ticket_id": "ticket_id",
    "partner_id": "partner_id",
    "kyc_id": "kyc_id",
    "invoice_id": "invoice_id",
    "document_id": "document_id",
    "proposal_id": "proposal_id",
    "email": "user_email",
    "mobile_number": "user_mobile_number",
    "phone_number": "user_mobile_number",
    "password": "login_password",
    "old_password": "old_password",
    "new_password": "new_password",
    "full_name": "full_name",
}

PROFILE_FIELD_PLACEHOLDERS = {
    "full_name": "profile_full_name",
    "date_of_birth": "profile_date_of_birth",
    "gender": "profile_gender",
    "address_line1": "profile_address_line1",
    "address_line2": "profile_address_line2",
    "city": "profile_city",
    "district": "profile_district",
    "division": "profile_division",
    "country": "profile_country",
    "nid_number": "profile_nid_number",
}

AUTH_SMOKE_PLAN = [
    {
        "name": "1. B2C Mobile OTP -> JWT",
        "description": "OTP Send -> OTP Verify -> Login -> Session -> Refresh -> Logout",
        "requests": [
            {
                "key": ("post", "/v1/auth/otp:send"),
                "name": "01 OTP Send",
                "body": {
                    "recipient": "{{user_mobile_number}}",
                    "type": "{{mobile_otp_type}}",
                    "channel": "{{mobile_otp_channel}}",
                    "use_masking": "{{mobile_otp_use_masking}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('user_mobile_number'), 'user_mobile_number must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('OTP send returned success', function () { pm.expect(pm.response.code).to.be.oneOf([200, 201, 202]); });",
                    "const j = pm.response.json();",
                    "if (j && j.data && j.data.otp_id) { pm.environment.set('mobile_otp_id', j.data.otp_id); pm.environment.set('otp_id', j.data.otp_id); }",
                ],
            },
            {
                "key": ("post", "/v1/auth/otp:verify"),
                "name": "02 OTP Verify",
                "body": {
                    "otp_id": "{{mobile_otp_id}}",
                    "code": "{{mobile_otp_code}}",
                    "device_id": "{{device_id}}",
                    "device_type": "{{login_device_type}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('mobile_otp_id'), 'mobile_otp_id must be set by OTP send').to.not.be.empty;",
                    "// Auto-fetch OTP from mock SMS gateway if running locally (no-op in production)",
                    "const mockSmsUrl = pm.environment.get('mock_sms_url');",
                    "if (mockSmsUrl) {",
                    "  const mobile = pm.environment.get('user_mobile_number') || '';",
                    "  const msisdn = mobile.replace('+880', '880');",
                    "  pm.sendRequest({ url: mockSmsUrl + '/mock/last-otp?msisdn=' + msisdn, method: 'GET' }, function(err, res) {",
                    "    if (!err && res && res.code === 200) {",
                    "      const j = res.json();",
                    "      if (j && j.otp) { pm.environment.set('mobile_otp_code', j.otp); console.log('Mock OTP fetched for verify:', j.otp); }",
                    "    }",
                    "  });",
                    "} else {",
                    "  pm.expect(pm.environment.get('mobile_otp_code'), 'mobile_otp_code must be set from delivery channel').to.not.be.empty;",
                    "}",
                ],
                "tests": [
                    "pm.test('OTP verify returned success', function () { pm.expect(pm.response.code).to.be.oneOf([200, 201, 204]); });",
                    "const j = pm.response.json();",
                    "pm.environment.set('mobile_otp_verified', 'true');",
                    "if (j && j.data && j.data.otp_id) { pm.environment.set('mobile_otp_id', j.data.otp_id); pm.environment.set('otp_id', j.data.otp_id); }",
                    "if (j && j.data && j.data.device_credential) { pm.environment.set('device_credential', j.data.device_credential); pm.environment.set('login_password', j.data.device_credential); }",
                    "if (j && j.data && j.data.user_id) { pm.environment.set('user_id', j.data.user_id); }",
                ],
            },
            {
                "key": ("post", "/v1/auth/login"),
                "name": "03 Login Passwordless (OTP)",
                "body": {
                    "mobile_number": "{{user_mobile_number}}",
                    "otp_id": "{{mobile_otp_id}}",
                    "device_id": "{{device_id}}",
                    "device_type": "{{login_device_type}}",
                    "device_name": "{{device_name}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('mobile_otp_id'), 'mobile_otp_id must be set by OTP send/verify').to.not.be.empty;",
                    "pm.expect(pm.environment.get('mobile_otp_verified'), 'Run OTP verify first').to.equal('true');",
                    "const deviceType = pm.environment.get('login_device_type');",
                    "pm.expect(['ANDROID', 'IOS', 'API']).to.include(deviceType);",
                    "// Auto-fetch OTP from mock SMS gateway if running locally (no-op in production)",
                    "const mockSmsUrl = pm.environment.get('mock_sms_url');",
                    "if (mockSmsUrl) {",
                    "  const mobile = pm.environment.get('user_mobile_number') || '';",
                    "  const msisdn = mobile.replace('+880', '880');",
                    "  pm.sendRequest({ url: mockSmsUrl + '/mock/last-otp?msisdn=' + msisdn, method: 'GET' }, function(err, res) {",
                    "    if (!err && res && res.code === 200) {",
                    "      const j = res.json();",
                    "      if (j && j.otp) { pm.environment.set('mobile_otp_code', j.otp); console.log('Mock OTP fetched:', j.otp); }",
                    "    }",
                    "  });",
                    "}",
                ],
                "tests": [
                    "pm.test('Login returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                    "const j = pm.response.json();",
                    "pm.test('Session type is JWT', function () { pm.expect(j.data.session_type).to.eql('JWT'); });",
                    "if (j && j.data && j.data.access_token) { pm.environment.set('access_token', j.data.access_token); }",
                    "if (j && j.data && j.data.refresh_token) { pm.environment.set('refresh_token', j.data.refresh_token); }",
                    "if (j && j.data && j.data.session_id) { pm.environment.set('session_id', j.data.session_id); }",
                    "if (j && j.data && j.data.user_id) { pm.environment.set('user_id', j.data.user_id); }",
                    "if (j && j.data && j.data.session_type) { pm.environment.set('session_type', j.data.session_type); }",
                    "pm.environment.set('last_b2c_otp_login', 'true');",
                ],
            },
            {
                "key": ("post", "/v1/auth/login"),
                "name": "04 Login Passwordless (Device-Bound)",
                "body": {
                    "mobile_number": "{{user_mobile_number}}",
                    "password": "{{device_credential}}",
                    "device_id": "{{device_id}}",
                    "device_type": "{{login_device_type}}",
                    "device_name": "{{device_name}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('device_credential'), 'device_credential must be set by OTP verify').to.not.be.empty;",
                    "const deviceType = pm.environment.get('login_device_type');",
                    "pm.expect(['ANDROID', 'IOS', 'API']).to.include(deviceType);",
                ],
                "tests": [
                    "pm.test('Device-bound login returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                    "const j = pm.response.json();",
                    "pm.test('Session type is JWT', function () { pm.expect(j.data.session_type).to.eql('JWT'); });",
                    "if (j && j.data && j.data.access_token) { pm.environment.set('access_token', j.data.access_token); }",
                    "if (j && j.data && j.data.refresh_token) { pm.environment.set('refresh_token', j.data.refresh_token); }",
                    "if (j && j.data && j.data.session_id) { pm.environment.set('session_id', j.data.session_id); }",
                    "if (j && j.data && j.data.user_id) { pm.environment.set('user_id', j.data.user_id); }",
                    "if (j && j.data && j.data.session_type) { pm.environment.set('session_type', j.data.session_type); }",
                    "pm.environment.set('last_b2c_device_credential_login', 'true');",
                ],
            },
            {
                "key": ("get", "/v1/auth/session/current"),
                "name": "05 Current Session",
                "tests": [
                    "pm.test('Current session returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                ],
            },
            {
                "key": ("post", "/v1/auth/token:refresh"),
                "name": "06 Refresh Token",
                "body": {
                    "refresh_token": "{{refresh_token}}",
                    "device_id": "{{device_id}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('refresh_token'), 'refresh_token must be set by login').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Refresh returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                    "const j = pm.response.json();",
                    "if (j && j.data && j.data.access_token) { pm.environment.set('access_token', j.data.access_token); }",
                    "if (j && j.data && j.data.refresh_token) { pm.environment.set('refresh_token', j.data.refresh_token); }",
                    "if (j && j.data && j.data.session_id) { pm.environment.set('session_id', j.data.session_id); }",
                ],
            },
            {
                "key": ("post", "/v1/auth/logout"),
                "name": "07 Logout",
                "body": {
                    "session_id": "{{session_id}}",
                    "access_token": "{{access_token}}",
                    "logout_reason": "manual_test",
                },
                "tests": [
                    "pm.test('Logout completed', function () { pm.expect(pm.response.code).to.be.oneOf([200, 204]); });",
                    "['access_token', 'refresh_token', 'session_token', 'csrf_token'].forEach((key) => pm.environment.unset(key));",
                ],
            },
        ],
    },
    {
        "name": "2. Web Portal Password -> Server Session",
        "description": "Password login with WEB device_type; session_token + csrf_token are auto-captured",
        "requests": [
            {
                "key": ("post", "/v1/auth/login"),
                "name": "01 Web Login",
                "body": {
                    "mobile_number": "{{user_mobile_number}}",
                    "password": "{{login_password}}",
                    "device_id": "{{device_id}}",
                    "device_type": "WEB",
                    "device_name": "{{device_name}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('user_mobile_number'), 'user_mobile_number must be set').to.not.be.empty;",
                    "pm.expect(pm.environment.get('login_password'), 'login_password must be set for WEB flow').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('WEB login returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                    "const j = pm.response.json();",
                    "pm.test('Session type is SERVER_SIDE', function () { pm.expect(j.data.session_type).to.eql('SERVER_SIDE'); });",
                ],
            },
            {
                "key": ("get", "/v1/auth/session/current"),
                "name": "02 Current Session",
                "tests": [
                    "pm.test('Current session returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                ],
            },
            {
                "key": ("post", "/v1/auth/logout"),
                "name": "03 Logout",
                "body": {
                    "session_id": "{{session_id}}",
                    "logout_reason": "manual_test",
                },
                "tests": [
                    "pm.test('WEB logout completed', function () { pm.expect(pm.response.code).to.be.oneOf([200, 204]); });",
                    "['session_token', 'csrf_token', 'access_token', 'refresh_token'].forEach((key) => pm.environment.unset(key));",
                ],
            },
        ],
    },
    {
        "name": "3. Email OTP -> Server Session",
        "description": "Email OTP send -> email login -> logout",
        "requests": [
            {
                "key": ("post", "/v1/auth/email/otp:send"),
                "name": "01 Send Email OTP",
                "body": {
                    "email": "{{user_email}}",
                    "type": "{{email_otp_type}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('user_email'), 'user_email must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Email OTP send returned success', function () { pm.expect(pm.response.code).to.be.oneOf([200, 201, 202]); });",
                    "const j = pm.response.json();",
                    "if (j && j.data && j.data.otp_id) { pm.environment.set('email_otp_id', j.data.otp_id); pm.environment.set('email_login_otp_id', j.data.otp_id); }",
                ],
            },
            {
                "key": ("post", "/v1/auth/email/login"),
                "name": "02 Email Login",
                "body": {
                    "email": "{{user_email}}",
                    "otp_id": "{{email_otp_id}}",
                    "code": "{{email_login_otp_code}}",
                    "device_id": "{{device_id}}",
                    "device_name": "{{device_name}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('email_otp_id'), 'email_otp_id must be set by email OTP send').to.not.be.empty;",
                    "pm.expect(pm.environment.get('email_login_otp_code'), 'email_login_otp_code must be set from inbox').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Email login returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                    "const j = pm.response.json();",
                    "pm.test('Session type is SERVER_SIDE', function () { pm.expect(j.data.session_type).to.eql('SERVER_SIDE'); });",
                ],
            },
            {
                "key": ("post", "/v1/auth/logout"),
                "name": "03 Logout",
                "body": {
                    "session_id": "{{session_id}}",
                    "logout_reason": "manual_test",
                },
                "tests": [
                    "pm.test('Email logout completed', function () { pm.expect(pm.response.code).to.be.oneOf([200, 204]); });",
                    "['session_token', 'csrf_token', 'access_token', 'refresh_token'].forEach((key) => pm.environment.unset(key));",
                ],
            },
        ],
    },
    {
        "name": "4. AuthZ Checks",
        "description": "Validate authz/check with both JWT and server-side session contexts",
        "requests": [
            {
                "key": ("post", "/v1/authz/check"),
                "name": "01 AuthZ Check (JWT/API Key)",
                "body": {
                    "user_id": "{{user_id}}",
                    "domain": "{{authz_domain}}",
                    "object": "{{authz_object}}",
                    "action": "{{authz_action}}",
                },
                "tests": [
                    "pm.test('AuthZ check returned a valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 401, 403]); });",
                ],
            },
            {
                "key": ("post", "/v1/authz/check"),
                "name": "02 AuthZ Check (Server Session)",
                "body": {
                    "user_id": "{{user_id}}",
                    "domain": "{{authz_domain}}",
                    "object": "{{authz_object}}",
                    "action": "{{authz_action}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('session_token'), 'session_token must be set by a WEB or email login flow').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Session authz check returned a valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 401, 403]); });",
                ],
            },
        ],
    },
]

B2C_SUITE_PLAN = [
    {
        "name": "1. B2C AuthN",
        "description": "OTP -> OTP login -> device-bound login -> current session -> refresh",
        "requests": AUTH_SMOKE_PLAN[0]["requests"][:6],
    },
    {
        "name": "2. B2C Profile & Sessions",
        "description": "Session lookups plus user profile create/read/update for the authenticated B2C user",
        "requests": [
            {
                "key": ("get", "/v1/auth/sessions/{session_id}"),
                "name": "Session-01 Get Session By ID",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set by login').to.not.be.empty;",
                    "pm.expect(pm.environment.get('session_id'), 'session_id must be set by login').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Get session by ID returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                ],
            },
            {
                "key": ("post", "/v1/auth/users/{user_id}/profile"),
                "name": "Profile-01 Create Profile",
                "body": {
                    "user_id": "{{user_id}}",
                    "full_name": "{{profile_full_name}}",
                    "date_of_birth": "{{profile_date_of_birth}}",
                    "gender": "{{profile_gender}}",
                    "address_line1": "{{profile_address_line1}}",
                    "city": "{{profile_city}}",
                    "district": "{{profile_district}}",
                    "division": "{{profile_division}}",
                    "country": "{{profile_country}}",
                    "nid_number": "{{profile_nid_number}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set by login').to.not.be.empty;",
                    "pm.expect(pm.environment.get('user_id'), 'user_id must be set by login').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Create profile returned created or already exists', function () { pm.expect(pm.response.code).to.be.oneOf([200, 201, 409]); });",
                ],
            },
            {
                "key": ("get", "/v1/auth/users/{user_id}/profile"),
                "name": "Profile-02 Get Profile",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set by login').to.not.be.empty;",
                    "pm.expect(pm.environment.get('user_id'), 'user_id must be set by login').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Get profile returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                ],
            },
            {
                "key": ("patch", "/v1/auth/users/{user_id}/profile"),
                "name": "Profile-03 Update Profile",
                "body": {
                    "user_id": "{{user_id}}",
                    "full_name": "{{profile_full_name_updated}}",
                    "city": "{{profile_city_updated}}",
                    "address_line2": "{{profile_address_line2}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set by login').to.not.be.empty;",
                    "pm.expect(pm.environment.get('user_id'), 'user_id must be set by login').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Update profile returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                ],
            },
            {
                "key": ("get", "/v1/auth/users/{user_id}/sessions"),
                "name": "Session-02 List User Sessions",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set by login').to.not.be.empty;",
                    "pm.expect(pm.environment.get('user_id'), 'user_id must be set by login').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('List user sessions returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                ],
            },
        ],
    },
    {
        "name": "3. B2C AuthZ",
        "description": "Authorization checks plus logout and same-device relogin after logout",
        "requests": [
            {
                "key": ("post", "/v1/authz/check"),
                "name": "AuthZ-01 Check Access",
                "body": {
                    "user_id": "{{user_id}}",
                    "domain": "{{authz_domain}}",
                    "object": "{{authz_object}}",
                    "action": "{{authz_action}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set by B2C login').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('AuthZ response is valid', function () { pm.expect(pm.response.code).to.be.oneOf([200, 401, 403]); });",
                ],
            },
            {
                "key": ("post", "/v1/auth/logout"),
                "name": "AuthZ-02 Logout",
                "body": {
                    "session_id": "{{session_id}}",
                    "access_token": "{{access_token}}",
                    "logout_reason": "manual_test",
                },
                "tests": [
                    "pm.test('Logout completed', function () { pm.expect(pm.response.code).to.be.oneOf([200, 204]); });",
                    "['access_token', 'refresh_token', 'session_token', 'csrf_token'].forEach((key) => pm.environment.unset(key));",
                ],
            },
            {
                "key": ("post", "/v1/auth/login"),
                "name": "AuthZ-03 Re-Login Passwordless (Device-Bound)",
                "body": {
                    "mobile_number": "{{user_mobile_number}}",
                    "password": "{{device_credential}}",
                    "device_id": "{{device_id}}",
                    "device_type": "{{login_device_type}}",
                    "device_name": "{{device_name}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('device_credential'), 'device_credential must be set by OTP verify').to.not.be.empty;",
                    "const deviceType = pm.environment.get('login_device_type');",
                    "pm.expect(['ANDROID', 'IOS', 'API']).to.include(deviceType);",
                ],
                "tests": [
                    "pm.test('Post-logout device-bound login returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                    "const j = pm.response.json();",
                    "pm.test('Session type is JWT', function () { pm.expect(j.data.session_type).to.eql('JWT'); });",
                    "if (j && j.data && j.data.access_token) { pm.environment.set('access_token', j.data.access_token); }",
                    "if (j && j.data && j.data.refresh_token) { pm.environment.set('refresh_token', j.data.refresh_token); }",
                    "if (j && j.data && j.data.session_id) { pm.environment.set('session_id', j.data.session_id); }",
                    "if (j && j.data && j.data.user_id) { pm.environment.set('user_id', j.data.user_id); }",
                    "if (j && j.data && j.data.session_type) { pm.environment.set('session_type', j.data.session_type); }",
                ],
            },
            {
                "key": ("get", "/v1/auth/session/current"),
                "name": "AuthZ-04 Current Session After Re-Login",
                "tests": [
                    "pm.test('Current session after relogin returned success', function () { pm.expect(pm.response.code).to.equal(200); });",
                ],
            },
            {
                "key": ("post", "/v1/auth/logout"),
                "name": "AuthZ-05 Final Logout",
                "body": {
                    "session_id": "{{session_id}}",
                    "access_token": "{{access_token}}",
                    "logout_reason": "final_test_cleanup",
                },
                "tests": [
                    "pm.test('Final logout completed', function () { pm.expect(pm.response.code).to.be.oneOf([200, 204]); });",
                    "['access_token', 'refresh_token', 'session_token', 'csrf_token'].forEach((key) => pm.environment.unset(key));",
                ],
            },
        ],
    },
    {
        "name": "4. B2C AuthZ Extended",
        "description": "JWKS, roles, policies, batch authz checks, session revocation",
        "requests": [
            {
                "key": ("post", "/v1/auth/login"),
                "name": "AuthZ-Ext-00 Re-Login (Device-Bound)",
                "body": {
                    "mobile_number": "{{user_mobile_number}}",
                    "password": "{{device_credential}}",
                    "device_id": "{{device_id}}",
                    "device_type": "{{login_device_type}}",
                    "device_name": "{{device_name}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('device_credential'), 'device_credential must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Re-login for authz extended section', function () { pm.expect(pm.response.code).to.be.oneOf([200, 429]); });",
                    "const j = pm.response.json();",
                    "if (j && j.data && j.data.access_token) { pm.environment.set('access_token', j.data.access_token); }",
                    "if (j && j.data && j.data.refresh_token) { pm.environment.set('refresh_token', j.data.refresh_token); }",
                    "if (j && j.data && j.data.session_id) { pm.environment.set('session_id', j.data.session_id); }",
                    "if (j && j.data && j.data.user_id) { pm.environment.set('user_id', j.data.user_id); }",
                ],
            },
            {
                "key": ("get", "/.well-known/jwks.json"),
                "name": "AuthZ-Ext-01 JWKS Public Keys",
                "tests": [
                    "pm.test('JWKS returned successfully', function () { pm.expect(pm.response.code).to.equal(200); });",
                    "const j = pm.response.json();",
                    "pm.test('JWKS has keys array', function () { pm.expect(j).to.have.property('keys'); pm.expect(j.keys).to.be.an('array').that.is.not.empty; });",
                ],
            },
            {
                "key": ("post", "/v1/authz/check:batch"),
                "name": "AuthZ-Ext-02 Batch Access Check",
                "body": {
                    "checks": [
                        {
                            "user_id": "{{user_id}}",
                            "domain": "{{authz_domain}}",
                            "object": "{{authz_object}}",
                            "action": "{{authz_action}}",
                        }
                    ]
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Batch authz check returned valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 400, 401, 403]); });",
                ],
            },
            {
                "key": ("get", "/v1/authz/roles"),
                "name": "AuthZ-Ext-03 List Roles",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('List roles returned valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 403, 503]); });",
                    "const j = pm.response.json();",
                    "if (pm.response.code === 200 && j && j.data && j.data.roles && j.data.roles.length > 0) { pm.environment.set('role_id', j.data.roles[0].role_id); }",
                ],
            },
            {
                "key": ("get", "/v1/authz/users/{user_id}/permissions"),
                "name": "AuthZ-Ext-04 Get User Permissions",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                    "pm.expect(pm.environment.get('user_id'), 'user_id must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('User permissions returned valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 403, 503]); });",
                ],
            },
            {
                "key": ("post", "/v1/auth/token:validate"),
                "name": "AuthZ-Ext-05 Validate Token",
                "body": {
                    "token": "{{access_token}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Token validation returned valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 400, 401]); });",
                ],
            },
            {
                "key": ("delete", "/v1/auth/sessions/{session_id}"),
                "name": "AuthZ-Ext-06 Revoke Session",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                    "pm.expect(pm.environment.get('session_id'), 'session_id must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Session revocation returned valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 204, 404]); });",
                ],
            },
        ],
    },
    {
        "name": "5. B2C Notifications",
        "description": "Get notifications, mark as read, update preferences",
        "requests": [
            {
                "key": ("post", "/v1/auth/login"),
                "name": "Notif-00 Re-Login (Device-Bound)",
                "body": {
                    "mobile_number": "{{user_mobile_number}}",
                    "password": "{{device_credential}}",
                    "device_id": "{{device_id}}",
                    "device_type": "{{login_device_type}}",
                    "device_name": "{{device_name}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('device_credential'), 'device_credential must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Re-login for notifications section', function () { pm.expect(pm.response.code).to.be.oneOf([200, 429]); });",
                    "const j = pm.response.json();",
                    "if (j && j.data && j.data.access_token) { pm.environment.set('access_token', j.data.access_token); }",
                    "if (j && j.data && j.data.user_id) { pm.environment.set('user_id', j.data.user_id); }",
                ],
            },
            {
                "key": ("get", "/v1/notifications"),
                "name": "Notif-01 Get User Notifications",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Get notifications returned success', function () { pm.expect(pm.response.code).to.be.oneOf([200, 403]); });",
                    "const j = pm.response.json();",
                    "if (pm.response.code === 200 && j && j.data && j.data.notifications && j.data.notifications.length > 0) { pm.environment.set('notification_id', j.data.notifications[0].notification_id); }",
                ],
            },
            {
                "key": ("patch", "/v1/auth/users/{user_id}/notification-preferences"),
                "name": "Notif-02 Update Notification Preferences",
                "body": {
                    "user_id": "{{user_id}}",
                    "sms_enabled": True,
                    "email_enabled": True,
                    "push_enabled": False,
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                    "pm.expect(pm.environment.get('user_id'), 'user_id must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Update notification preferences returned success', function () { pm.expect(pm.response.code).to.be.oneOf([200, 401, 403, 404, 503]); });",
                ],
            },
        ],
    },
    {
        "name": "6. B2C User Documents",
        "description": "Upload, list, get and update user identity documents",
        "requests": [
            {
                "key": ("post", "/v1/auth/login"),
                "name": "Doc-00 Re-Login (Device-Bound)",
                "body": {
                    "mobile_number": "{{user_mobile_number}}",
                    "password": "{{device_credential}}",
                    "device_id": "{{device_id}}",
                    "device_type": "{{login_device_type}}",
                    "device_name": "{{device_name}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('device_credential'), 'device_credential must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Re-login for documents section', function () { pm.expect(pm.response.code).to.be.oneOf([200, 429]); });",
                    "const j = pm.response.json();",
                    "if (j && j.data && j.data.access_token) { pm.environment.set('access_token', j.data.access_token); }",
                    "if (j && j.data && j.data.user_id) { pm.environment.set('user_id', j.data.user_id); }",
                ],
            },
            {
                "key": ("get", "/v1/auth/document-types"),
                "name": "Doc-01 List Document Types",
                "tests": [
                    "pm.test('Document types returned success', function () { pm.expect(pm.response.code).to.be.oneOf([200, 401]); });",
                    "const j = pm.response.json();",
                    "pm.test('Has document types', function () { if (pm.response.code === 200) { pm.expect(j.data).to.exist; } });",
                ],
            },
            {
                "key": ("post", "/v1/auth/users/{user_id}/documents"),
                "name": "Doc-02 Upload User Document",
                "body": {
                    "user_id": "{{user_id}}",
                    "document_type": "NATIONAL_ID",
                    "document_number": "{{profile_nid_number}}",
                    "file_url": "https://example.com/test-document.jpg",
                    "issued_at": "2020-01-01T00:00:00Z",
                    "expires_at": "2030-01-01T00:00:00Z",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                    "pm.expect(pm.environment.get('user_id'), 'user_id must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Document upload returned success', function () { pm.expect(pm.response.code).to.be.oneOf([200, 201, 401, 409]); });",
                    "const j = pm.response.json();",
                    "if (j && j.data && j.data.document_id) { pm.environment.set('user_document_id', j.data.document_id); }",
                    "if (j && j.data && j.data.user_document && j.data.user_document.user_document_id) { pm.environment.set('user_document_id', j.data.user_document.user_document_id); }",
                ],
            },
            {
                "key": ("get", "/v1/auth/users/{user_id}/documents"),
                "name": "Doc-03 List User Documents",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                    "pm.expect(pm.environment.get('user_id'), 'user_id must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('List documents returned success', function () { pm.expect(pm.response.code).to.be.oneOf([200, 401]); });",
                    "const j = pm.response.json();",
                    "if (j && j.data && j.data.documents && j.data.documents.length > 0) { pm.environment.set('user_document_id', j.data.documents[0].user_document_id); }",
                ],
            },
            {
                "key": ("get", "/v1/auth/documents/{user_document_id}"),
                "name": "Doc-04 Get Document By ID",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                    "if (!pm.environment.get('user_document_id')) { console.warn('user_document_id not set, skipping'); }",
                ],
                "tests": [
                    "pm.test('Get document returned success', function () { pm.expect(pm.response.code).to.be.oneOf([200, 401, 404]); });",
                ],
            },
            {
                "key": ("patch", "/v1/auth/documents/{user_document_id}"),
                "name": "Doc-05 Update Document",
                "body": {
                    "document_number": "{{profile_nid_number}}",
                    "verification_status": "PENDING",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                    "if (!pm.environment.get('user_document_id')) { console.warn('user_document_id not set, skipping'); }",
                ],
                "tests": [
                    "pm.test('Update document returned success', function () { pm.expect(pm.response.code).to.be.oneOf([200, 401, 403, 404]); });",
                ],
            },
        ],
    },
    {
        "name": "7. B2C Orders & Payments",
        "description": "Create order, initiate payment, verify payment status",
        "requests": [
            {
                "key": ("post", "/v1/auth/login"),
                "name": "Order-00 Re-Login (Device-Bound)",
                "body": {
                    "mobile_number": "{{user_mobile_number}}",
                    "password": "{{device_credential}}",
                    "device_id": "{{device_id}}",
                    "device_type": "{{login_device_type}}",
                    "device_name": "{{device_name}}",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('device_credential'), 'device_credential must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Re-login for orders section', function () { pm.expect(pm.response.code).to.be.oneOf([200, 429]); });",
                    "const j = pm.response.json();",
                    "if (j && j.data && j.data.access_token) { pm.environment.set('access_token', j.data.access_token); }",
                    "if (j && j.data && j.data.user_id) { pm.environment.set('user_id', j.data.user_id); }",
                ],
            },
            {
                "key": ("get", "/v1/products"),
                "name": "Order-01 List Products",
                "tests": [
                    "pm.test('List products returned success', function () { pm.expect(pm.response.code).to.be.oneOf([200, 401, 503]); });",
                    "const j = pm.response.json();",
                    "if (pm.response.code === 200 && j && j.data && j.data.products && j.data.products.length > 0) { pm.environment.set('product_id', j.data.products[0].product_id); }",
                ],
            },
            {
                "key": ("post", "/v1/orders"),
                "name": "Order-02 Create Order",
                "body": {
                    "user_id": "{{user_id}}",
                    "product_id": "{{product_id}}",
                    "quantity": 1,
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                    "pm.expect(pm.environment.get('user_id'), 'user_id must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Create order returned valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 201, 400, 401, 422, 503]); });",
                    "const j = pm.response.json();",
                    "if (j && j.data && j.data.order_id) { pm.environment.set('order_id', j.data.order_id); }",
                    "if (j && j.data && j.data.order && j.data.order.order_id) { pm.environment.set('order_id', j.data.order.order_id); }",
                ],
            },
            {
                "key": ("get", "/v1/orders"),
                "name": "Order-03 List My Orders",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('List orders returned valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 401, 503]); });",
                    "const j = pm.response.json();",
                    "if (pm.response.code === 200 && j && j.data && j.data.orders && j.data.orders.length > 0) { pm.environment.set('order_id', j.data.orders[0].order_id); }",
                ],
            },
            {
                "key": ("get", "/v1/orders/{order_id}"),
                "name": "Order-04 Get Order By ID",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                    "if (!pm.environment.get('order_id')) { console.warn('order_id not set, skipping'); }",
                ],
                "tests": [
                    "pm.test('Get order returned valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 404, 503]); });",
                ],
            },
            {
                "key": ("post", "/v1/payments"),
                "name": "Payment-01 Initiate Payment",
                "body": {
                    "order_id": "{{order_id}}",
                    "user_id": "{{user_id}}",
                    "amount": 10000,
                    "currency": "BDT",
                    "provider": "sslcommerz",
                    "return_url": "http://localhost:3000/payment/return",
                    "cancel_url": "http://localhost:3000/payment/cancel",
                },
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('Initiate payment returned valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 201, 400, 401, 422, 503]); });",
                    "const j = pm.response.json();",
                    "if (j && j.data && j.data.payment_id) { pm.environment.set('payment_id', j.data.payment_id); }",
                ],
            },
            {
                "key": ("get", "/v1/payments"),
                "name": "Payment-02 List Payments",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                ],
                "tests": [
                    "pm.test('List payments returned valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 401, 503]); });",
                    "const j = pm.response.json();",
                    "if (pm.response.code === 200 && j && j.data && j.data.payments && j.data.payments.length > 0) { pm.environment.set('payment_id', j.data.payments[0].payment_id); }",
                ],
            },
            {
                "key": ("get", "/v1/payments/{payment_id}"),
                "name": "Payment-03 Get Payment By ID",
                "prerequest": [
                    "pm.expect(pm.environment.get('access_token'), 'access_token must be set').to.not.be.empty;",
                    "if (!pm.environment.get('payment_id')) { console.warn('payment_id not set, skipping'); }",
                ],
                "tests": [
                    "pm.test('Get payment returned valid response', function () { pm.expect(pm.response.code).to.be.oneOf([200, 404, 503]); });",
                ],
            },
        ],
    },
]


def log(message: str) -> None:
    print(message, flush=True)


def stable_uuid(name: str) -> str:
    return str(uuid.uuid5(uuid.NAMESPACE_URL, f"insuretech.postman::{name}"))


def safe_json_dump(value: Any) -> str:
    return json.dumps(value, indent=2, ensure_ascii=False)


def parse_args(argv: list[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Generate coherent Postman collections/environments from api/openapi.yaml and optionally sync them to Postman."
    )
    parser.add_argument("--spec", default=str(DEFAULT_SPEC_PATH), help="Path to OpenAPI YAML file")
    parser.add_argument("--output-dir", "--output", dest="output_dir", default=str(DEFAULT_OUTPUT_DIR), help="Output directory")
    parser.add_argument("--dotenv", "--env", dest="dotenv_path", default=str(DEFAULT_DOTENV_PATH), help="Path to local .env file")
    parser.add_argument("--upload", action="store_true", help="Upload generated artifacts to Postman")
    parser.add_argument("--collection-name", default=COLLECTION_NAME, help="Postman collection name")
    parser.add_argument("--postman-api-key", default="", help="Override POSTMAN_API_KEY")
    parser.add_argument("--workspace-id", default="", help="Override POSTMAN_WORKSPACE_ID")
    parser.add_argument("--collection-id", default="", help="Override POSTMAN_COLLECTION_ID")
    return parser.parse_args(argv)


def load_runtime_config(dotenv_path: Path) -> dict[str, str]:
    values: dict[str, str] = {}
    if dotenv_path.exists():
        for key, value in dotenv_values(dotenv_path).items():
            if value is None:
                continue
            values[key] = str(value)
    for key, value in os.environ.items():
        values[key] = value
    return values


def read_openapi_spec(spec_path: Path) -> dict[str, Any]:
    with spec_path.open("r", encoding="utf-8") as handle:
        return yaml.safe_load(handle)


def resolve_converter_command() -> list[str]:
    explicit = os.environ.get("POSTMAN_OPENAPI_CLI", "").strip()
    if explicit:
        return [explicit]

    for candidate in ("openapi2postmanv2", "openapi-to-postmanv2"):
        resolved = shutil.which(candidate)
        if resolved:
            return [resolved]

    npx = shutil.which("npx") or shutil.which("npx.cmd")
    if npx:
        return [npx, "--yes", "openapi-to-postmanv2"]

    raise RuntimeError(
        "Could not find Postman's OpenAPI converter. Install Node.js/npx or set POSTMAN_OPENAPI_CLI."
    )


def run_openapi_converter(spec_path: Path) -> dict[str, Any]:
    converter_command = resolve_converter_command()
    with tempfile.TemporaryDirectory(prefix="insuretech-postman-") as temp_dir:
        output_path = Path(temp_dir) / "collection.json"
        command = [*converter_command, "-s", str(spec_path), "-o", str(output_path), "-p"]
        log("[1/5] Converting OpenAPI spec with Postman's openapi-to-postmanv2...")
        completed = subprocess.run(
            command,
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace",
            check=False,
        )
        if completed.returncode != 0:
            raise RuntimeError(
                "openapi-to-postmanv2 failed.\n"
                f"Command: {' '.join(command)}\n"
                f"stdout:\n{completed.stdout}\n"
                f"stderr:\n{completed.stderr}"
            )
        with output_path.open("r", encoding="utf-8") as handle:
            return json.load(handle)


def normalize_path_param(segment: str) -> str:
    segment = str(segment)
    if segment.startswith(":"):
        return f"{{{segment[1:]}}}"
    placeholder = PLACEHOLDER_RE.fullmatch(segment)
    if placeholder:
        return f"{{{placeholder.group(1)}}}"
    return segment


def postman_path_to_openapi(path_value: Any) -> str:
    if isinstance(path_value, list):
        segments = [normalize_path_param(segment) for segment in path_value]
        return "/" + "/".join(segment.strip("/") for segment in segments if str(segment).strip("/"))
    if isinstance(path_value, str):
        return "/" + "/".join(normalize_path_param(part) for part in path_value.split("/") if part)
    return "/"


def extract_path_from_item(item: dict[str, Any]) -> str:
    request = item.get("request", {})
    url = request.get("url", {})
    if isinstance(url, dict) and url.get("path"):
        return postman_path_to_openapi(url["path"])
    raw = url.get("raw") if isinstance(url, dict) else url
    if isinstance(raw, str):
        raw = raw.replace("{{baseUrl}}", "").replace("{{base_url}}", "")
        raw = re.sub(r"^https?://[^/]+", "", raw)
        return postman_path_to_openapi(raw.split("?")[0])
    return "/"


def normalize_url_path_variables(url: dict[str, Any]) -> None:
    if not isinstance(url, dict):
        return

    def env_segment(segment: Any) -> str:
        text = str(segment)
        if text.startswith(":") and len(text) > 1:
            return f"{{{{{text[1:]}}}}}"
        brace_match = re.fullmatch(r"\{([^{}]+)\}", text)
        if brace_match:
            return f"{{{{{brace_match.group(1)}}}}}"
        return text

    path = url.get("path")
    if isinstance(path, list):
        normalized_path = [env_segment(segment) for segment in path]
        url["path"] = normalized_path
        host = url.get("host")
        if isinstance(host, list):
            host_parts = [str(part).strip("/") for part in host if str(part).strip("/")]
            if host_parts:
                url["raw"] = "/".join(host_parts) + "/" + "/".join(segment.strip("/") for segment in normalized_path)

    raw = url.get("raw")
    if isinstance(raw, str):
        raw = raw.replace("{{baseUrl}}", "{{base_url}}")
        raw = re.sub(r":([A-Za-z_][A-Za-z0-9_]*)", r"{{\1}}", raw)
        raw = re.sub(r"(?<!\{)\{([A-Za-z_][A-Za-z0-9_]*)\}(?!\})", r"{{\1}}", raw)
        url["raw"] = raw

    variables = url.get("variable")
    if isinstance(variables, list):
        for variable in variables:
            key = str(variable.get("key", "")).strip()
            if key:
                variable["value"] = f"{{{{{key}}}}}"


def replace_placeholders(value: Any) -> Any:
    if isinstance(value, str):
        return value.replace("{{baseUrl}}", "{{base_url}}").replace("{{bearerToken}}", "{{access_token}}")
    if isinstance(value, list):
        return [replace_placeholders(entry) for entry in value]
    if isinstance(value, dict):
        return {key: replace_placeholders(entry) for key, entry in value.items()}
    return value


def flatten_collection_items(items: list[dict[str, Any]]) -> list[dict[str, Any]]:
    leaves: list[dict[str, Any]] = []
    for item in items:
        if "request" in item:
            leaves.append(item)
        for child in item.get("item", []) or []:
            leaves.extend(flatten_collection_items([child]))
    return leaves


def singularize(name: str) -> str:
    if name.endswith("ies") and len(name) > 3:
        return f"{name[:-3]}y"
    if name.endswith("ses") and len(name) > 3:
        return name[:-2]
    if name.endswith("s") and not name.endswith("ss") and len(name) > 1:
        return name[:-1]
    return name


def sanitize_display_text(text: str) -> str:
    if not text:
        return text

    def collapse_spaced_acronym(match: re.Match[str]) -> str:
        token = match.group(0)
        letters = token.split()
        if len(letters) < 2:
            return token
        return "".join(letters).upper()

    sanitized = SPACED_ACRONYM_RE.sub(collapse_spaced_acronym, text)
    acronym_replacements = {
        "Api": "API",
        "Otp": "OTP",
        "Totp": "TOTP",
        "Kyc": "KYC",
        "Jwt": "JWT",
        "Csrf": "CSRF",
        "Jwks": "JWKS",
        "Mfa": "MFA",
        "Id": "ID",
    }
    for source, target in acronym_replacements.items():
        sanitized = re.sub(rf"\b{source}\b", target, sanitized)
    return sanitized


def resolve_ref(ref: str, spec: dict[str, Any]) -> dict[str, Any]:
    if not ref.startswith("#/"):
        return {}
    current: Any = spec
    for part in ref[2:].split("/"):
        decoded = part.replace("~1", "/").replace("~0", "~")
        if not isinstance(current, dict) or decoded not in current:
            return {}
        current = current[decoded]
    return copy.deepcopy(current) if isinstance(current, dict) else {}


def merge_schema_parts(parts: list[dict[str, Any]]) -> dict[str, Any]:
    merged: dict[str, Any] = {}
    required: list[str] = []
    properties: dict[str, Any] = {}

    for part in parts:
        if not isinstance(part, dict):
            continue
        for key, value in part.items():
            if key == "required":
                for entry in value or []:
                    if entry not in required:
                        required.append(entry)
            elif key == "properties":
                properties.update(value or {})
            elif key == "description":
                if not merged.get("description") and value:
                    merged["description"] = value
            elif key in {"type", "format", "example", "default", "enum", "items", "additionalProperties"}:
                merged[key] = copy.deepcopy(value)
            elif key not in {"allOf", "anyOf", "oneOf"}:
                merged[key] = copy.deepcopy(value)

    if properties:
        merged["properties"] = properties
    if required:
        merged["required"] = required
    return merged


def resolve_schema(schema: dict[str, Any] | None, spec: dict[str, Any], seen_refs: set[str] | None = None) -> dict[str, Any]:
    if not isinstance(schema, dict):
        return {}

    seen_refs = seen_refs or set()
    if "$ref" in schema:
        ref = str(schema["$ref"])
        if ref in seen_refs:
            return {}
        return resolve_schema(resolve_ref(ref, spec), spec, seen_refs | {ref})

    resolved = copy.deepcopy(schema)

    if "allOf" in resolved:
        merged = merge_schema_parts([resolve_schema(part, spec, seen_refs) for part in resolved.get("allOf", []) or []])
        inline = {key: value for key, value in resolved.items() if key != "allOf"}
        resolved = merge_schema_parts([merged, resolve_schema(inline, spec, seen_refs)])

    for union_key in ("oneOf", "anyOf"):
        if union_key in resolved:
            branches = resolved.get(union_key) or []
            if branches:
                chosen = resolve_schema(branches[0], spec, seen_refs)
                inline = {key: value for key, value in resolved.items() if key != union_key}
                resolved = merge_schema_parts([chosen, resolve_schema(inline, spec, seen_refs)])
            else:
                resolved.pop(union_key, None)

    properties = {}
    for key, value in (resolved.get("properties") or {}).items():
        properties[key] = resolve_schema(value, spec, seen_refs)
    if properties:
        resolved["properties"] = properties

    if isinstance(resolved.get("items"), dict):
        resolved["items"] = resolve_schema(resolved["items"], spec, seen_refs)
    if isinstance(resolved.get("additionalProperties"), dict):
        resolved["additionalProperties"] = resolve_schema(resolved["additionalProperties"], spec, seen_refs)
    return resolved


def pick_request_media_type(content: dict[str, Any]) -> tuple[str, dict[str, Any]] | tuple[None, None]:
    if not isinstance(content, dict):
        return None, None
    preferred_order = [
        "application/json",
        "application/problem+json",
        "application/merge-patch+json",
        "multipart/form-data",
        "application/x-www-form-urlencoded",
    ]
    for preferred in preferred_order:
        if preferred in content:
            return preferred, content[preferred]
    for media_type, media in content.items():
        if "json" in media_type:
            return str(media_type), media
    for media_type, media in content.items():
        return str(media_type), media
    return None, None


def extract_request_body_schema(operation: dict[str, Any], spec: dict[str, Any]) -> tuple[str, dict[str, Any] | None]:
    request_body = operation.get("requestBody")
    if not isinstance(request_body, dict):
        return "", None
    if "$ref" in request_body:
        request_body = resolve_ref(str(request_body["$ref"]), spec)
    content_type, media = pick_request_media_type(request_body.get("content", {}) or {})
    if not content_type or not isinstance(media, dict):
        return "", None
    schema = resolve_schema(media.get("schema"), spec)
    return content_type, schema or None


def infer_placeholder_variable(
    field_name: str,
    field_path: tuple[str, ...],
    operation: dict[str, Any],
) -> str:
    name = str(field_name).strip()
    lower_name = name.lower()
    context = " ".join([operation.get("path", ""), operation.get("operation_id", ""), *field_path]).lower()

    if "/auth/users/" in context and "/profile" in context and lower_name in PROFILE_FIELD_PLACEHOLDERS:
        return PROFILE_FIELD_PLACEHOLDERS[lower_name]

    if lower_name in DIRECT_FIELD_PLACEHOLDERS:
        return DIRECT_FIELD_PLACEHOLDERS[lower_name]
    if lower_name == "recipient":
        return "user_email" if "email" in context else "user_mobile_number"
    if lower_name in {"resource", "resource_name"}:
        return "authz_resource"
    if lower_name == "action":
        return "authz_action"
    if lower_name == "type":
        if "/auth/email/otp:send" in operation.get("path", ""):
            return "email_otp_type"
        if "/auth/otp:send" in operation.get("path", ""):
            return "mobile_otp_type"
    if lower_name == "channel" and "/auth/otp:send" in operation.get("path", ""):
        return "mobile_otp_channel"
    if lower_name == "use_masking" and "/auth/otp:send" in operation.get("path", ""):
        return "mobile_otp_use_masking"
    if lower_name == "device_type" and any(token in context for token in ("login", "otp:verify", "biometric")):
        return "login_device_type"
    if lower_name in {"code", "otp_code"}:
        if "totp" in context:
            return "totp_code"
        if "email" in context and "verify" in context:
            return "email_verification_otp_code"
        if "email" in context and "login" in context:
            return "email_login_otp_code"
        if "otp" in context:
            return "mobile_otp_code"
        return "otp_code"
    if lower_name == "id" and operation.get("resource_hint"):
        return f"{operation['resource_hint']}_id"
    if lower_name.endswith("_id") or lower_name.endswith("_token") or lower_name.endswith("_secret") or lower_name.endswith("_credential"):
        return lower_name
    return ""


def build_schema_template_value(schema: dict[str, Any], operation: dict[str, Any], field_path: tuple[str, ...] = ()) -> Any:
    if not isinstance(schema, dict):
        return ""

    field_name = field_path[-1] if field_path else ""
    placeholder_variable = infer_placeholder_variable(field_name, field_path, operation) if field_name else ""
    if placeholder_variable:
        return f"{{{{{placeholder_variable}}}}}"

    if "example" in schema:
        return copy.deepcopy(schema["example"])
    if "default" in schema:
        return copy.deepcopy(schema["default"])
    if schema.get("enum"):
        return copy.deepcopy(schema["enum"][0])

    schema_type = schema.get("type")
    schema_format = str(schema.get("format", "")).lower()

    if schema_type == "object" or "properties" in schema:
        payload: dict[str, Any] = {}
        for key, value in (schema.get("properties") or {}).items():
            if isinstance(value, dict) and value.get("readOnly") is True:
                continue
            payload[key] = build_schema_template_value(value, operation, (*field_path, key))
        return payload
    if schema_type == "array":
        item_schema = schema.get("items") if isinstance(schema.get("items"), dict) else {}
        item_value = build_schema_template_value(item_schema, operation, (*field_path, "item"))
        return [item_value] if item_value not in ("", None, {}) else []
    if schema_type == "boolean":
        return False
    if schema_type in {"integer", "number"}:
        return 1
    if schema_format == "date-time":
        return "2024-01-15T10:30:00Z"
    if schema_format == "date":
        return "2024-01-15"
    if schema_format == "email":
        return "{{user_email}}"
    if schema_format == "uri":
        return "https://example.com"
    return "example_value"


def render_postman_json(value: Any, indent: int = 0) -> str:
    spacing = "  " * indent
    child_spacing = "  " * (indent + 1)

    if isinstance(value, dict):
        if not value:
            return "{}"
        lines = ["{"]
        entries = list(value.items())
        for index, (key, entry) in enumerate(entries):
            suffix = "," if index < len(entries) - 1 else ""
            lines.append(f"{child_spacing}{json.dumps(key, ensure_ascii=False)}: {render_postman_json(entry, indent + 1)}{suffix}")
        lines.append(f"{spacing}}}")
        return "\n".join(lines)

    if isinstance(value, list):
        if not value:
            return "[]"
        lines = ["["]
        for index, entry in enumerate(value):
            suffix = "," if index < len(value) - 1 else ""
            lines.append(f"{child_spacing}{render_postman_json(entry, indent + 1)}{suffix}")
        lines.append(f"{spacing}]")
        return "\n".join(lines)

    if isinstance(value, str):
        match = PLACEHOLDER_RE.fullmatch(value)
        if match and match.group(1) in NON_STRING_PLACEHOLDER_KEYS:
            return value
        return json.dumps(value, ensure_ascii=False)

    if isinstance(value, bool):
        return "true" if value else "false"
    if value is None:
        return "null"
    return json.dumps(value, ensure_ascii=False)


def collect_placeholders(value: Any, found: set[str] | None = None) -> set[str]:
    found = found or set()
    if isinstance(value, str):
        found.update(match.group(1) for match in PLACEHOLDER_RE.finditer(value))
        return found
    if isinstance(value, list):
        for entry in value:
            collect_placeholders(entry, found)
        return found
    if isinstance(value, dict):
        for entry in value.values():
            collect_placeholders(entry, found)
    return found


def build_operation_index(spec: dict[str, Any]) -> dict[tuple[str, str], dict[str, Any]]:
    index: dict[tuple[str, str], dict[str, Any]] = {}
    global_security = spec.get("security", [])
    paths = spec.get("paths", {})
    for openapi_path, path_item in paths.items():
        if not isinstance(path_item, dict):
            continue
        common_parameters = path_item.get("parameters", []) or []
        for method in HTTP_METHODS:
            operation = path_item.get(method)
            if not isinstance(operation, dict):
                continue
            parameters = [*common_parameters, *(operation.get("parameters", []) or [])]
            query_params = {
                parameter.get("name"): parameter
                for parameter in parameters
                if isinstance(parameter, dict) and parameter.get("in") == "query"
            }
            responses = operation.get("responses", {}) or {}
            success_codes = [int(code) for code in responses if str(code).isdigit() and int(code) < 400]
            error_codes = [int(code) for code in responses if str(code).isdigit() and int(code) >= 400]
            summary = operation.get("summary") or operation.get("operationId") or f"{method.upper()} {openapi_path}"
            operation_id = operation.get("operationId") or ""
            resource_hint = derive_resource_hint(openapi_path, operation_id)
            security = operation.get("security", global_security)
            request_content_type, request_schema = extract_request_body_schema(operation, spec)
            index[(method.upper(), openapi_path)] = {
                "method": method.upper(),
                "path": openapi_path,
                "summary": summary,
                "description": operation.get("description", ""),
                "operation_id": operation_id,
                "responses": responses,
                "success_codes": success_codes or [200],
                "error_codes": error_codes or [400, 401, 403, 404, 422, 429, 500, 503],
                "query_params": query_params,
                "security": security,
                "tags": operation.get("tags") or [],
                "resource_hint": resource_hint,
                "request_content_type": request_content_type,
                "request_schema": request_schema,
            }
    return index


def derive_resource_hint(openapi_path: str, operation_id: str) -> str:
    segments = [segment for segment in openapi_path.split("/") if segment and not segment.startswith("v")]
    resource_candidates = [segment for segment in segments if not segment.startswith("{")]
    if resource_candidates:
        hint = singularize(resource_candidates[-1].split(":")[0].replace("-", "_"))
        if hint in {"auth", "authz"} and len(resource_candidates) > 1:
            hint = singularize(resource_candidates[-2].split(":")[0].replace("-", "_"))
        return hint
    if operation_id:
        pieces = operation_id.split("_", 1)
        if len(pieces) == 2:
            return singularize(re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", pieces[1]).lower())
    return ""


def derive_service_name(operation: dict[str, Any]) -> str:
    operation_id = operation.get("operation_id") or ""
    if "_" in operation_id:
        return operation_id.split("_", 1)[0]
    if operation.get("tags"):
        tag = str(operation["tags"][0]).strip()
        return re.sub(r"\W+", "", tag.title()) or "General"
    segments = [segment for segment in operation["path"].split("/") if segment]
    if len(segments) >= 2:
        return f"{segments[1].replace('-', '_').title().replace('_', '')}Service"
    return "General"


def build_request_description(operation: dict[str, Any]) -> str:
    auth_label = "Public" if not operation["security"] else "Protected"
    lines = [
        sanitize_display_text(operation["summary"]),
        "",
        f"Operation ID: `{operation['operation_id'] or 'n/a'}`",
        f"Path: `{operation['method']} {operation['path']}`",
        f"Auth: {auth_label}",
    ]
    if operation["description"]:
        lines.extend(["", sanitize_display_text(operation["description"].strip())])
    if operation["responses"]:
        lines.extend(["", "Responses:"])
        for code, response in operation["responses"].items():
            description = ""
            if isinstance(response, dict):
                description = response.get("description", "")
            lines.append(sanitize_display_text(f"- `{code}` {description}".rstrip()))
    return "\n".join(lines).strip()


def remove_header(headers: list[dict[str, Any]], header_name: str) -> list[dict[str, Any]]:
    return [header for header in headers if str(header.get("key", "")).lower() != header_name.lower()]


def normalize_query_params(item: dict[str, Any], operation: dict[str, Any]) -> None:
    url = item.get("request", {}).get("url", {})
    if not isinstance(url, dict):
        return
    query = url.get("query")
    if not isinstance(query, list):
        return
    for parameter in query:
        name = parameter.get("key")
        metadata = operation["query_params"].get(name, {})
        if not metadata.get("required", False):
            parameter["disabled"] = True


def make_collection_prerequest_script() -> list[str]:
    script = """
const method = (pm.request.method || 'GET').toUpperCase();
const ensureHeader = (key, value) => {
  if (value === undefined || value === null || value === '') {
    return;
  }
  pm.request.headers.upsert({ key, value: String(value) });
};
const getEnv = (...keys) => {
  for (const key of keys) {
    const value = pm.environment.get(key);
    if (value !== undefined && value !== null && String(value).trim() !== '') {
      return String(value).trim();
    }
  }
  return '';
};
const validDeviceTypes = ['ANDROID', 'IOS', 'API', 'WEB'];
const configuredDeviceType = getEnv('login_device_type');
if (configuredDeviceType && !validDeviceTypes.includes(configuredDeviceType)) {
  console.warn('Unexpected login_device_type:', configuredDeviceType);
}

ensureHeader('Accept', 'application/json');
if (['POST', 'PUT', 'PATCH', 'DELETE'].includes(method) && pm.request.body && pm.request.body.mode === 'raw') {
  ensureHeader('Content-Type', 'application/json');
}

const accessToken = getEnv('access_token');
const sessionToken = getEnv('session_token');
const csrfToken = getEnv('csrf_token');
const apiKey = getEnv('api_key');

if (accessToken) {
  ensureHeader('Authorization', `Bearer ${accessToken}`);
}
if (apiKey) {
  ensureHeader('X-API-Key', apiKey);
}
if (sessionToken && !accessToken) {
  ensureHeader('X-Session-Token', sessionToken);
}
if (csrfToken && ['POST', 'PUT', 'PATCH', 'DELETE'].includes(method)) {
  ensureHeader('X-CSRF-Token', csrfToken);
}

ensureHeader(
  'X-Request-ID',
  'req_' + pm.variables.replaceIn('{{$randomUUID}}').replace(/-/g, '').slice(0, 16)
);

pm.environment.set('baseUrl', getEnv('base_url'));
pm.environment.set('bearerToken', accessToken);
if (apiKey) {
  pm.environment.set('api_key_id', apiKey);
}
"""
    return [line.rstrip() for line in script.strip().splitlines()]


def make_test_script(operation: dict[str, Any]) -> list[str]:
    operation_name = operation["operation_id"] or operation["summary"]
    success_codes = json.dumps(sorted(set(operation["success_codes"])))
    error_codes = json.dumps(sorted(set(operation["error_codes"])))
    resource_hint = json.dumps(operation["resource_hint"] or "")
    script = f"""
const operationName = {json.dumps(operation_name)};
const successCodes = {success_codes};
const errorCodes = {error_codes};
const resourceHint = {resource_hint};
let body = null;
try {{
  body = pm.response.json();
}} catch (error) {{
  body = null;
}}

const looksLikeEnvelope = body && typeof body === 'object' && 'success' in body && 'data' in body && 'error' in body;

pm.test(`[${{operationName}}] Status code is documented`, function () {{
  if (pm.response.code < 400) {{
    pm.expect(successCodes).to.include(pm.response.code);
  }} else {{
    pm.expect(errorCodes).to.include(pm.response.code);
  }}
}});

pm.test(`[${{operationName}}] Response time < 10000ms`, function () {{
  pm.expect(pm.response.responseTime).to.be.below(10000);
}});

if (pm.response.code !== 204 && body !== null) {{
  pm.test(`[${{operationName}}] Content-Type is JSON`, function () {{
    pm.expect(pm.response.headers.get('Content-Type') || '').to.include('application/json');
  }});
}}

if (looksLikeEnvelope) {{
  pm.test(`[${{operationName}}] ApiResponse envelope is present`, function () {{
    pm.expect(body).to.have.property('success');
    pm.expect(body).to.have.property('data');
    pm.expect(body).to.have.property('error');
  }});

  pm.test(`[${{operationName}}] success/error fields are coherent`, function () {{
    if (body.success === true) {{
      pm.expect(body.error).to.satisfy((value) => value === null || value === undefined);
    }}
    if (body.success === false) {{
      pm.expect(body.error).to.not.be.null;
      pm.expect(body.error.code).to.be.a('string').and.not.empty;
      pm.expect(body.error.message).to.be.a('string').and.not.empty;
    }}
  }});
}}

if (pm.response.code === 422 && looksLikeEnvelope) {{
  pm.test(`[${{operationName}}] 422 includes field_violations`, function () {{
    pm.expect(body.error).to.have.property('field_violations');
    pm.expect(body.error.field_violations).to.be.an('array');
  }});
}}

const toStringValue = (value) => {{
  if (value === undefined || value === null || value === '') {{
    return '';
  }}
  if (typeof value === 'string') {{
    return value;
  }}
  if (typeof value === 'number' || typeof value === 'boolean') {{
    return String(value);
  }}
  return JSON.stringify(value);
}};

const setEnvIfPresent = (key, value, aliases = []) => {{
  const normalized = toStringValue(value);
  if (!normalized) {{
    return;
  }}
  pm.environment.set(key, normalized);
  aliases.forEach((alias) => pm.environment.set(alias, normalized));
}};

const getNestedValue = (source, path) => {{
  if (!source || typeof source !== 'object') {{
    return undefined;
  }}
  let current = source;
  for (const part of path.split('.')) {{
    if (current === undefined || current === null || !(part in current)) {{
      return undefined;
    }}
    current = current[part];
  }}
  return current;
}};

const aliasMap = [
  {{ paths: ['data.access_token', 'access_token'], key: 'access_token' }},
  {{ paths: ['data.refresh_token', 'refresh_token'], key: 'refresh_token' }},
  {{ paths: ['data.session_token', 'session_token'], key: 'session_token' }},
  {{ paths: ['data.csrf_token', 'csrf_token'], key: 'csrf_token' }},
  {{ paths: ['data.mfa_session_token', 'mfa_session_token'], key: 'mfa_session_token' }},
  {{ paths: ['data.session_id', 'session_id', 'data.session.session_id'], key: 'session_id' }},
  {{ paths: ['data.session_type', 'session_type'], key: 'session_type' }},
  {{ paths: ['data.user_id', 'user_id', 'data.user.user_id', 'data.user.id', 'user.id'], key: 'user_id' }},
  {{ paths: ['data.tenant_id', 'tenant_id'], key: 'tenant_id' }},
  {{ paths: ['data.otp_id', 'otp_id'], key: 'otp_id' }},
  {{ paths: ['data.mobile_number', 'data.user.mobile_number', 'data.user.phone_number', 'mobile_number'], key: 'user_mobile_number' }},
  {{ paths: ['data.email', 'data.user.email', 'email'], key: 'user_email' }},
  {{ paths: ['data.device_credential', 'device_credential'], key: 'device_credential' }},
  {{ paths: ['data.totp_secret', 'totp_secret'], key: 'totp_secret' }},
  {{ paths: ['data.provisioning_uri', 'provisioning_uri'], key: 'provisioning_uri' }},
  {{ paths: ['data.biometric_token', 'biometric_token'], key: 'biometric_token' }},
];

aliasMap.forEach((entry) => {{
  for (const path of entry.paths) {{
    const value = getNestedValue(body, path);
    if (value !== undefined && value !== null && value !== '') {{
      const aliases = [];
      if (entry.key.endsWith('_id')) {{
        aliases.push(`last_${{entry.key}}`);
      }}
      setEnvIfPresent(entry.key, value, aliases);
      break;
    }}
  }}
}});
"""
    script += """
if (pm.request.url.getPath().includes('/auth/otp:send')) {
  const otpId = getNestedValue(body, 'data.otp_id');
  if (otpId) {
    setEnvIfPresent('mobile_otp_id', otpId, ['otp_id', 'last_mobile_otp_id']);
  }
}
if (pm.request.url.getPath().includes('/auth/email/otp:send')) {
  const otpId = getNestedValue(body, 'data.otp_id');
  if (otpId) {
    setEnvIfPresent('email_otp_id', otpId, ['email_login_otp_id', 'otp_id']);
  }
}
if ((pm.request.url.getPath().includes('/auth/email:verify') || pm.request.url.getPath().includes('/auth/email/verify')) && pm.response.code < 400) {
  pm.environment.set('email_verified', 'true');
}
if (pm.request.url.getPath().includes('/auth/otp:verify') && pm.response.code < 400) {
  pm.environment.set('mobile_otp_verified', 'true');
}
if (pm.request.url.getPath().includes('/auth/login') && pm.response.code < 400) {
  const sessionType = getNestedValue(body, 'data.session_type');
  if (sessionType === 'JWT') {
    pm.environment.set('last_b2c_otp_login', 'true');
  }
}
if (pm.request.url.getPath().includes('/auth/logout') && pm.response.code < 400) {
  ['access_token', 'refresh_token', 'session_token', 'csrf_token'].forEach((key) => pm.environment.unset(key));
}

const walk = (node, depth = 0) => {
  if (!node || typeof node !== 'object' || depth > 5) {
    return;
  }
  if (Array.isArray(node)) {
    node.forEach((entry) => walk(entry, depth + 1));
    return;
  }
  for (const [key, value] of Object.entries(node)) {
    if (value === undefined || value === null || value === '') {
      continue;
    }
    if (key === 'id' && resourceHint && (typeof value === 'string' || typeof value === 'number')) {
      setEnvIfPresent(`${resourceHint}_id`, value, [`last_${resourceHint}_id`]);
    }
    if (/_id$/.test(key) && (typeof value === 'string' || typeof value === 'number')) {
      setEnvIfPresent(key, value, [`last_${key}`]);
    }
    if ((/_token$/.test(key) || /_secret$/.test(key) || /_credential$/.test(key)) && (typeof value === 'string' || typeof value === 'number')) {
      setEnvIfPresent(key, value);
    }
    walk(value, depth + 1);
  }
};

if (body && typeof body === 'object') {
  walk(body);
}
"""
    return [line.rstrip() for line in script.strip().splitlines()]


def normalize_leaf_item(base_item: dict[str, Any], operation: dict[str, Any]) -> dict[str, Any]:
    item = replace_placeholders(copy.deepcopy(base_item))
    request = item.setdefault("request", {})
    request["name"] = sanitize_display_text(operation["summary"])
    request["description"] = build_request_description(operation)
    normalize_url_path_variables(request.get("url", {}))
    request.pop("auth", None)
    request["header"] = remove_header(request.get("header", []) or [], "Authorization")
    request["header"] = remove_header(request["header"], "X-Request-ID")
    request["header"] = remove_header(request["header"], "X-API-Key")
    request["header"] = remove_header(request["header"], "X-Session-Token")
    request["header"] = remove_header(request["header"], "X-CSRF-Token")
    normalize_query_params(item, operation)
    if operation.get("request_content_type", "").startswith("application/json") and operation.get("request_schema"):
        set_request_body(item, build_schema_template_value(operation["request_schema"], operation))
    item["name"] = sanitize_display_text(f"{operation['method']} {operation['summary']}")
    item["event"] = [
        {"listen": "prerequest", "script": {"type": "text/javascript", "exec": make_collection_prerequest_script()}},
        {"listen": "test", "script": {"type": "text/javascript", "exec": make_test_script(operation)}},
    ]
    return item


def build_operation_variants(item: dict[str, Any], operation: dict[str, Any]) -> list[dict[str, Any]]:
    if operation["method"] == "POST" and operation["path"] == "/v1/auth/login":
        generic_item = copy.deepcopy(item)

        b2c_item = copy.deepcopy(item)
        b2c_item["name"] = "POST Login with credentials (B2C OTP/JWT)"
        b2c_item["request"]["name"] = "Login with credentials (B2C OTP/JWT)"
        b2c_item["request"]["description"] = (
            f"{item['request']['description']}\n\n"
            "Flow variant: B2C OTP login. Use this after `/v1/auth/otp:verify`.\n"
            "Leave `password` empty and keep `device_type` on a JWT/mobile-style client such as `ANDROID`, `IOS`, or `API`."
        ).strip()
        set_request_body(
            b2c_item,
            {
                "mobile_number": "{{user_mobile_number}}",
                "password": "",
                "device_id": "{{device_id}}",
                "device_type": "{{login_device_type}}",
                "device_name": "{{device_name}}",
            },
        )

        password_item = copy.deepcopy(item)
        password_item["name"] = "POST Login with credentials (Password/B2B)"
        password_item["request"]["name"] = "Login with credentials (Password/B2B)"
        password_item["request"]["description"] = (
            f"{item['request']['description']}\n\n"
            "Flow variant: password login. Use this for password-based sign-in and set `login_password` in the environment."
        ).strip()
        set_request_body(
            password_item,
            {
                "mobile_number": "{{user_mobile_number}}",
                "password": "{{login_password}}",
                "device_id": "{{device_id}}",
                "device_type": "{{login_device_type}}",
                "device_name": "{{device_name}}",
            },
        )
        return [generic_item, b2c_item, password_item]

    return [item]


def build_main_collection(
    converted_collection: dict[str, Any],
    spec: dict[str, Any],
    collection_name: str,
) -> dict[str, Any]:
    operation_index = build_operation_index(spec)
    grouped_items: dict[str, list[tuple[str, dict[str, Any]]]] = defaultdict(list)
    for leaf in flatten_collection_items(converted_collection.get("item", []) or []):
        method = str(leaf.get("request", {}).get("method", "GET")).upper()
        path = extract_path_from_item(leaf)
        operation = operation_index.get((method, path))
        if not operation:
            continue
        enriched = normalize_leaf_item(leaf, operation)
        for variant_index, variant in enumerate(build_operation_variants(enriched, operation)):
            grouped_items[derive_service_name(operation)].append((f"{operation['path']}::{variant_index}", variant))

    top_level_items: list[dict[str, Any]] = []
    for service_name in sorted(grouped_items):
        folder_items = [item for _, item in sorted(grouped_items[service_name], key=lambda entry: entry[0])]
        top_level_items.append(
            {
                "name": service_name,
                "description": f"Operations for {service_name} ({len(folder_items)} endpoints)",
                "item": folder_items,
            }
        )

    collection_description = "\n".join(
        [
            spec.get("info", {}).get("description", "OpenAPI-generated collection."),
            "",
            "Generated from `api/openapi.yaml` via Postman's `openapi-to-postmanv2`, then normalized for InsureTech auth flows.",
            "",
            "Setup",
            "1. Import one of the generated environments from `api/postman/`.",
            "2. Set `user_mobile_number`, `user_email`, and `login_password` as needed.",
            "3. Run `auth_smoke.postman_collection.json` first to prime auth/session variables.",
            "4. Use the main collection for broader endpoint exploration with the same environment.",
        ]
    ).strip()

    return {
        "info": {
            "_postman_id": stable_uuid(collection_name),
            "name": collection_name,
            "description": collection_description,
            "schema": COLLECTION_SCHEMA,
        },
        "event": [
            {"listen": "prerequest", "script": {"type": "text/javascript", "exec": make_collection_prerequest_script()}},
        ],
        "variable": [
            {"key": "base_url", "value": "{{base_url}}", "type": "string"},
        ],
        "item": top_level_items,
    }


def find_request_item(collection: dict[str, Any], key: tuple[str, str]) -> dict[str, Any] | None:
    method, path = key
    for item in flatten_collection_items(collection.get("item", []) or []):
        item_method = str(item.get("request", {}).get("method", "GET")).upper()
        item_path = extract_path_from_item(item)
        if item_method == method.upper() and item_path == path:
            return item
    return None


def set_request_body(item: dict[str, Any], payload: Any) -> None:
    request = item.setdefault("request", {})
    request["body"] = {
        "mode": "raw",
        "raw": render_postman_json(payload),
        "options": {"raw": {"language": "json"}},
    }


def overlay_request_scripts(item: dict[str, Any], prerequest_lines: list[str], test_lines: list[str]) -> None:
    item["event"] = [
        {
            "listen": "prerequest",
            "script": {"type": "text/javascript", "exec": [*make_collection_prerequest_script(), *prerequest_lines]},
        },
        {"listen": "test", "script": {"type": "text/javascript", "exec": test_lines}},
    ]


def build_focus_collection(
    main_collection: dict[str, Any],
    plan: list[dict[str, Any]],
    name: str,
    description: str,
) -> dict[str, Any]:
    folders: list[dict[str, Any]] = []
    for section in plan:
        requests: list[dict[str, Any]] = []
        for entry in section["requests"]:
            base_item = find_request_item(main_collection, entry["key"])
            if not base_item:
                continue
            item = copy.deepcopy(base_item)
            item["name"] = entry["name"]
            if "body" in entry:
                set_request_body(item, entry["body"])
            overlay_request_scripts(item, entry.get("prerequest", []), entry.get("tests", []))
            requests.append(item)
        folders.append({"name": section["name"], "description": section["description"], "item": requests})
    return {
        "info": {
            "_postman_id": stable_uuid(name),
            "name": name,
            "description": description,
            "schema": COLLECTION_SCHEMA,
        },
        "event": [
            {"listen": "prerequest", "script": {"type": "text/javascript", "exec": make_collection_prerequest_script()}},
        ],
        "item": folders,
    }


def existing_environment_values(path: Path) -> dict[str, str]:
    if not path.exists():
        return {}
    try:
        with path.open("r", encoding="utf-8") as handle:
            document = json.load(handle)
    except Exception:
        return {}
    values: dict[str, str] = {}
    for entry in document.get("values", []) or []:
        key = entry.get("key")
        if key:
            values[str(key)] = str(entry.get("value", ""))
    return values


def env_default_url(config: dict[str, str], spec: dict[str, Any], env_key: str) -> str:
    suffix = env_key.upper()
    explicit = config.get(f"POSTMAN_{suffix}_BASE_URL", "").strip()
    if explicit:
        return explicit

    servers = spec.get("servers", []) or []
    if env_key == "local":
        newman_base = config.get("NEWMAN_BASE_URL", "").strip()
        if newman_base:
            return newman_base
        server_port = config.get("SERVER_PORT", "8080").strip() or "8080"
        return f"http://localhost:{server_port}"
    if env_key == "staging":
        for server in servers:
            url = str(server.get("url", ""))
            if "staging" in url:
                return url
    if env_key == "production":
        for server in servers:
            url = str(server.get("url", ""))
            if url and "staging" not in url:
                return url
    if env_key == "mock":
        return config.get("POSTMAN_MOCK_SERVER_URL", "").strip() or "http://localhost:4010"
    if env_key == "newman_test":
        return config.get("NEWMAN_BASE_URL", "").strip() or env_default_url(config, spec, "local")
    return ""


def derive_postman_value(key: str, env_key: str, config: dict[str, str], spec: dict[str, Any]) -> str:
    if key == "base_url":
        return env_default_url(config, spec, env_key)

    env_overrides = {
        f"POSTMAN_ENV_{key.upper()}",
        f"POSTMAN_{env_key.upper()}_{key.upper()}",
    }
    for override_key in env_overrides:
        value = config.get(override_key, "").strip()
        if value:
            return value

    direct_aliases = {
        "user_email": [f"POSTMAN_{env_key.upper()}_USER_EMAIL", "POSTMAN_USER_EMAIL", "USER_EMAIL", "EMAIL_FROM"],
        "user_mobile_number": [f"POSTMAN_{env_key.upper()}_USER_MOBILE_NUMBER", "POSTMAN_USER_MOBILE_NUMBER", "USER_MOBILE_NUMBER"],
        "full_name": [f"POSTMAN_{env_key.upper()}_FULL_NAME", "POSTMAN_FULL_NAME", "FULL_NAME"],
        "login_password": [f"POSTMAN_{env_key.upper()}_LOGIN_PASSWORD", "POSTMAN_LOGIN_PASSWORD"],
        "old_password": [f"POSTMAN_{env_key.upper()}_OLD_PASSWORD", "POSTMAN_OLD_PASSWORD"],
        "new_password": [f"POSTMAN_{env_key.upper()}_NEW_PASSWORD", "POSTMAN_NEW_PASSWORD"],
        "device_id": [f"POSTMAN_{env_key.upper()}_DEVICE_ID", "POSTMAN_DEVICE_ID"],
        "device_name": [f"POSTMAN_{env_key.upper()}_DEVICE_NAME", "POSTMAN_DEVICE_NAME"],
        "login_device_type": [f"POSTMAN_{env_key.upper()}_LOGIN_DEVICE_TYPE", "POSTMAN_LOGIN_DEVICE_TYPE"],
        "totp_code": [f"POSTMAN_{env_key.upper()}_TOTP_CODE", "POSTMAN_TOTP_CODE", "TOTP_CODE"],
        "api_key": [f"POSTMAN_{env_key.upper()}_API_KEY", "POSTMAN_API_KEY_VALUE", "API_KEY"],
        "authz_domain": [f"POSTMAN_{env_key.upper()}_AUTHZ_DOMAIN", "POSTMAN_AUTHZ_DOMAIN"],
        "authz_object": [f"POSTMAN_{env_key.upper()}_AUTHZ_OBJECT", "POSTMAN_AUTHZ_OBJECT"],
        "authz_resource": [f"POSTMAN_{env_key.upper()}_AUTHZ_RESOURCE", "POSTMAN_AUTHZ_RESOURCE"],
        "authz_action": [f"POSTMAN_{env_key.upper()}_AUTHZ_ACTION", "POSTMAN_AUTHZ_ACTION"],
    }
    for alias in direct_aliases.get(key, []):
        value = config.get(alias, "").strip()
        if value:
            return value

    if env_key == "local":
        direct_value = config.get(key, "").strip()
        if direct_value:
            return direct_value
    return ""


def build_environment_document(
    env_key: str,
    env_name: str,
    config: dict[str, str],
    spec: dict[str, Any],
    output_dir: Path,
    placeholder_keys: set[str],
) -> dict[str, Any]:
    path = output_dir / f"InsureTech_{env_key}.postman_environment.json"
    existing_values = existing_environment_values(path)
    merged_values: list[dict[str, Any]] = []
    seen_keys: set[str] = set()

    for variable in BASE_POSTMAN_VARIABLES:
        key = variable["key"]
        generated = derive_postman_value(key, env_key, config, spec)
        preserved = existing_values.get(key, "")
        value = generated if generated != "" else preserved if preserved != "" else variable["value"]
        merged_values.append(
            {
                "key": key,
                "value": value,
                "enabled": True,
                "type": "secret" if key in SECRET_KEYS else "default",
            }
        )
        seen_keys.add(key)

    for key in sorted(placeholder_keys):
        if key in seen_keys:
            continue
        generated = derive_postman_value(key, env_key, config, spec)
        merged_values.append(
            {
                "key": key,
                "value": generated,
                "enabled": True,
                "type": "secret" if key in SECRET_KEYS else "default",
            }
        )
        seen_keys.add(key)

    if env_key == "local":
        for key in sorted(config):
            if key in seen_keys:
                continue
            value = config[key]
            if value in ("", None):
                continue
            merged_values.append(
                {
                    "key": key,
                    "value": value,
                    "enabled": True,
                    "type": "secret" if "KEY" in key or "TOKEN" in key or "PASSWORD" in key else "default",
                }
            )

    return {
        "id": stable_uuid(env_name),
        "name": env_name,
        "values": merged_values,
        "_postman_variable_scope": "environment",
        "_postman_exported_at": "",
        "_postman_exported_using": "insuretech-postman-sync",
    }


def write_json(path: Path, payload: dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8", newline="\n") as handle:
        json.dump(payload, handle, indent=2, ensure_ascii=False)
        handle.write("\n")


def generate_artifacts(spec_path: Path, output_dir: Path, config: dict[str, str], collection_name: str) -> dict[str, Path]:
    spec = read_openapi_spec(spec_path)
    converted_collection = run_openapi_converter(spec_path)

    log("[2/5] Normalizing main collection structure and scripts...")
    main_collection = build_main_collection(converted_collection, spec, collection_name)
    auth_smoke = build_focus_collection(
        main_collection,
        AUTH_SMOKE_PLAN,
        "InsureTech Auth Smoke",
        "Ordered auth smoke flows generated from the OpenAPI-backed collection.",
    )
    b2c_suite = build_focus_collection(
        main_collection,
        B2C_SUITE_PLAN,
        "B2C AuthN + AuthZ Test Suite",
        "Focused B2C OTP/JWT authentication and authorization regression suite.",
    )
    placeholder_keys = set()
    collect_placeholders(main_collection, placeholder_keys)
    collect_placeholders(auth_smoke, placeholder_keys)
    collect_placeholders(b2c_suite, placeholder_keys)

    log("[3/5] Building Postman environments from defaults + local .env...")
    output_dir.mkdir(parents=True, exist_ok=True)
    paths = {
        "collection": output_dir / POSTMAN_COLLECTION_FILENAME,
        "auth_smoke": output_dir / AUTH_SMOKE_FILENAME,
        "b2c_suite": output_dir / B2C_SUITE_FILENAME,
    }
    write_json(paths["collection"], main_collection)
    write_json(paths["auth_smoke"], auth_smoke)
    write_json(paths["b2c_suite"], b2c_suite)

    for env_key, env_name in ENVIRONMENT_SPECS:
        env_document = build_environment_document(env_key, env_name, config, spec, output_dir, placeholder_keys)
        env_path = output_dir / f"InsureTech_{env_key}.postman_environment.json"
        write_json(env_path, env_document)
        paths[f"env_{env_key}"] = env_path

    test_env_alias = output_dir / "_test_env.json"
    shutil.copyfile(paths["env_newman_test"], test_env_alias)
    paths["test_env_alias"] = test_env_alias
    return paths


def postman_headers(api_key: str) -> dict[str, str]:
    return {
        "X-Api-Key": api_key,
        "Content-Type": "application/json",
    }


def postman_request(
    method: str,
    path: str,
    api_key: str,
    *,
    params: dict[str, str] | None = None,
    payload: dict[str, Any] | None = None,
) -> dict[str, Any]:
    response = requests.request(
        method=method,
        url=f"{POSTMAN_API_BASE}{path}",
        headers=postman_headers(api_key),
        params=params,
        data=json.dumps(payload) if payload is not None else None,
        timeout=60,
    )
    if response.status_code >= 400:
        raise RuntimeError(f"Postman API {method} {path} failed ({response.status_code}): {response.text}")
    if response.content:
        return response.json()
    return {}


def first_present(document: dict[str, Any], *keys: str) -> str:
    for key in keys:
        value = document.get(key)
        if value:
            return str(value)
    return ""


def find_remote_asset(api_key: str, asset_type: str, name: str) -> str:
    listing_key = "collections" if asset_type == "collection" else "environments"
    response = postman_request("GET", f"/{listing_key}", api_key)
    for candidate in response.get(listing_key, []) or []:
        if candidate.get("name") == name:
            return first_present(candidate, "uid", "id")
    return ""


def delete_remote_asset(api_key: str, asset_type: str, asset_id: str) -> None:
    if not asset_id:
        return
    listing_key = "collections" if asset_type == "collection" else "environments"
    postman_request("DELETE", f"/{listing_key}/{asset_id}", api_key)


def get_default_workspace_id(api_key: str) -> str:
    response = postman_request("GET", "/workspaces", api_key)
    workspaces = response.get("workspaces", []) or []
    for workspace in workspaces:
        if workspace.get("name") == "Default workspace":
            return first_present(workspace, "id", "uid")
    if workspaces:
        return first_present(workspaces[0], "id", "uid")
    return ""


def fetch_remote_collection(api_key: str, collection_id: str) -> dict[str, Any]:
    response = postman_request("GET", f"/collections/{collection_id}", api_key)
    return response.get("collection", {}) if isinstance(response, dict) else {}


def postman_request_signature(collection: dict[str, Any]) -> set[tuple[str, str, str]]:
    signature: set[tuple[str, str, str]] = set()
    for item in flatten_collection_items(collection.get("item", []) or []):
        request = item.get("request", {}) if isinstance(item, dict) else {}
        method = str(request.get("method", "GET")).upper()
        path = extract_path_from_item(item)
        name = str(item.get("name", ""))
        signature.add((method, path, name))
    return signature


def create_collection(api_key: str, collection: dict[str, Any], workspace_id: str) -> str:
    payload = {"collection": collection}
    params = {"workspace": workspace_id} if workspace_id else None
    response = postman_request("POST", "/collections", api_key, params=params, payload=payload)
    return first_present(response.get("collection", {}), "uid", "id")


def upsert_collection(api_key: str, collection: dict[str, Any], workspace_id: str, explicit_id: str) -> str:
    target_workspace_id = workspace_id or get_default_workspace_id(api_key)
    collection_id = explicit_id or find_remote_asset(api_key, "collection", collection["info"]["name"])
    payload = {"collection": collection}
    if collection_id:
        response = postman_request("PUT", f"/collections/{collection_id}", api_key, payload=payload)
        remote_collection_id = first_present(response.get("collection", {}), "uid", "id") or collection_id
        remote_collection = fetch_remote_collection(api_key, remote_collection_id)
        if postman_request_signature(remote_collection) != postman_request_signature(collection):
            log("  Remote collection remained stale after PUT; recreating collection in workspace...")
            delete_remote_asset(api_key, "collection", collection_id)
            recreated_id = create_collection(api_key, collection, target_workspace_id)
            return recreated_id
        return remote_collection_id
    return create_collection(api_key, collection, target_workspace_id)


def upsert_environment(api_key: str, environment: dict[str, Any], workspace_id: str, explicit_id: str) -> str:
    environment_id = explicit_id or find_remote_asset(api_key, "environment", environment["name"])
    payload = {"environment": environment}
    params = {"workspace": workspace_id} if workspace_id else None
    if environment_id:
        response = postman_request("PUT", f"/environments/{environment_id}", api_key, payload=payload)
    else:
        response = postman_request("POST", "/environments", api_key, params=params, payload=payload)
    return first_present(response.get("environment", {}), "uid", "id")


def upload_artifacts(paths: dict[str, Path], config: dict[str, str], workspace_id: str, collection_id: str) -> None:
    api_key = (config.get("POSTMAN_API_KEY", "") or "").strip()
    if not api_key:
        raise RuntimeError("POSTMAN_API_KEY is required for --upload.")

    log("[5/5] Uploading collection + environments to Postman API...")
    with paths["collection"].open("r", encoding="utf-8") as handle:
        collection = json.load(handle)
    remote_collection_id = upsert_collection(api_key, collection, workspace_id, collection_id)
    log(f"  Collection synced: {collection['info']['name']} ({remote_collection_id or 'created'})")

    for env_key, env_name in ENVIRONMENT_SPECS:
        path = paths[f"env_{env_key}"]
        with path.open("r", encoding="utf-8") as handle:
            environment = json.load(handle)
        explicit_env_id = config.get(f"POSTMAN_ENVIRONMENT_ID_{env_key.upper()}", "").strip()
        remote_env_id = upsert_environment(api_key, environment, workspace_id, explicit_env_id)
        log(f"  Environment synced: {env_name} ({remote_env_id or 'created'})")


def main(argv: list[str] | None = None) -> int:
    args = parse_args(argv)
    spec_path = Path(args.spec).resolve()
    output_dir = Path(args.output_dir).resolve()
    dotenv_path = Path(args.dotenv_path).resolve()

    if not spec_path.exists():
        raise SystemExit(f"OpenAPI spec not found: {spec_path}")

    config = load_runtime_config(dotenv_path)
    if args.postman_api_key:
        config["POSTMAN_API_KEY"] = args.postman_api_key
    if args.workspace_id:
        config["POSTMAN_WORKSPACE_ID"] = args.workspace_id
    if args.collection_id:
        config["POSTMAN_COLLECTION_ID"] = args.collection_id

    paths = generate_artifacts(spec_path, output_dir, config, args.collection_name)
    log("[4/5] Wrote Postman artifacts:")
    for label, path in paths.items():
        if label == "test_env_alias":
            continue
        log(f"  - {path}")

    if args.upload:
        upload_artifacts(
            paths,
            config,
            workspace_id=config.get("POSTMAN_WORKSPACE_ID", "").strip(),
            collection_id=config.get("POSTMAN_COLLECTION_ID", "").strip(),
        )
    else:
        log("[5/5] Local Postman generation complete. Re-run with --upload to sync to Postman cloud.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
