SSL Wireless SMS API

OTP Integration & BTRC Compliance Guide

Comprehensive Implementation Documentation


Version 2.0 - Updated January 2026
 
1. Executive Summary
This comprehensive guide provides complete technical documentation for integrating SSL Wireless SMS Gateway services to send One-Time Passwords (OTP) and transactional messages in Bangladesh. It incorporates BTRC (Bangladesh Telecommunication Regulatory Commission) regulatory requirements, SSL Wireless API specifications, and industry best practices for secure OTP delivery.
1.1 Key Information
•	Region: Bangladesh (Country Code: +880)
•	Protocol: HTTPS POST with JSON
•	Provider: SSL Wireless (BTRC-registered SMS aggregator)
•	Primary Operators: Grameenphone, Robi, Banglalink, Teletalk
•	Platforms: ISMSPlus (modern API, recommended) and legacy ISMS
•	Authentication: API Token + Service ID (SID)
1.2 Document Scope
This document covers:
•	BTRC regulatory compliance requirements
•	Sender ID registration process (Masking vs Non-Masking)
•	Complete SSL Wireless API specification
•	Request/response formats and error handling
•	Delivery report (DLR) webhook implementation
•	Security best practices for OTP systems
•	Production-ready code examples in Node.js and Python
 
2. BTRC Regulatory Compliance
Bangladesh Telecommunication Regulatory Commission (BTRC) enforces strict regulations on Application-to-Person (A2P) SMS traffic. Compliance is mandatory to avoid penalties, service suspension, and ensure message delivery.
2.1 Aggregator Mandate
BTRC requires all commercial SMS traffic to route through registered aggregators. Direct connection to mobile operators (Grameenphone, Robi, Banglalink, Teletalk) is not permitted for most businesses.

•	SSL Wireless is a BTRC-enlisted aggregator authorized to provide A2P SMS services
•	All SMS must comply with BTRC Memorandum 14.32.0000.600.43.005.21.434 (dated May 20, 2021)
•	Aggregators maintain direct connections to all major mobile network operators
•	Compliance ensures legal operation and reliable message delivery
2.2 Sender ID Types: Masking vs Non-Masking
Choosing the correct Sender ID type is critical for compliance, user trust, and message delivery rates. Bangladesh supports two sender ID types:

Feature	Masking (Alphanumeric)	Non-Masking (Numeric)
Display Name	Brand Name (e.g., "MYAPP", "BKASH")	Virtual Number (e.g., "88096...")
Best For	OTP, Transaction Alerts, High Trust	Two-Way Messaging, Customer Support
User Reply	Not Possible (One-Way Only)	Possible (Two-Way)
Approval Process	Strict: Operator vetting required	Route reservation (faster)
Timeline	3-7 days (new) / 15-30 days (transfer)	1-3 days
Format	Max 11 chars; A-Z, 0-9 only	Standard numeric (880...)

2.2.1 Masking (Alphanumeric Sender ID)
Recommended for OTP and transactional messages due to high user trust and brand recognition.
•	Maximum 11 characters (A-Z, 0-9 only)
•	No special characters allowed (!@#$%^&*()-_)
•	Must be in CAPITAL LETTERS if representing an individual's name
•	One-way communication only (users cannot reply)
•	Requires approval from all major operators (GP, Robi, Banglalink, Teletalk)
•	Higher trust perception increases OTP conversion rates

2.2.2 Non-Masking (Numeric Sender ID)
Used when two-way communication is needed or for faster setup.
•	Displays as virtual mobile number (880XXXXXXXXX)
•	Allows recipients to reply to messages
•	Faster registration process (1-3 days)
•	Lower trust perception compared to branded masking
•	Suitable for customer service and support use cases
 
2.3 Sender ID Registration Process
To use a masking sender ID, you must complete a formal registration process with SSL Wireless, who will submit your application to all major operators.
2.3.1 Required Documents
1.	Valid Trade License or Business Incorporation Certificate
2.	Company TIN (Tax Identification Number) or BIN
3.	Authorized Signatory Letter on company letterhead
4.	NID/Passport copy of the authorized person
5.	Company contact details (address, phone, email)
6.	For government/educational institutions: Ministry authorization or UGC authorization

2.3.2 Registration Steps
7.	Submit documentation to SSL Wireless (online portal or email)
8.	SSL Wireless verifies KYC (Know Your Customer) compliance
9.	SSL submits masking request to operators (GP, Robi, Banglalink, Teletalk)
10.	Each operator independently reviews and approves
11.	SSL Wireless notifies you upon approval (3-7 working days for new, 15-30 days for transfers)
12.	Masking is activated in your SSL account configuration

2.3.3 Important Notes
•	The same masking name can be registered with multiple aggregators
•	Transferring existing masking from another provider takes 15-30 days
•	Request release from previous provider to expedite transfer
•	Only local (domestic) SMS traffic is allowed with masking
•	Masking approval can be rejected if it's too generic or misleading
 
2.4 Content Regulations & NDNC Rules
BTRC enforces strict content regulations to protect consumers and maintain social harmony.
2.4.1 Prohibited Content
•	Political campaign messages (except with special permission during elections)
•	Gambling, betting, or lottery promotions
•	Religious provocative content
•	Anti-state content or content that disturbs social peace
•	Hate speech, discrimination, or harassment
•	Adult or explicit content

2.4.2 Language Requirements
•	Promotional SMS must be sent in Bangla (Bengali) language
•	Transactional and OTP messages can be in English
•	Bangla messages use UCS-2 encoding (70 characters per SMS segment)
•	English messages use GSM-7 encoding (160 characters per SMS segment)
•	Unicode characters count as multiple characters in GSM-7

2.4.3 Opt-Out and NDNC Compliance
•	Include STOP/HELP keywords for promotional messages
•	Honor National Do Not Call (NDNC) registry
•	Transactional OTPs generally bypass DND timing restrictions
•	Marketing SMS limited to 2 messages per day per recipient
•	Send promotional SMS only between 9 AM - 9 PM Bangladesh Standard Time (BST, GMT+6)
•	Respect user opt-out requests immediately

2.4.4 Transactional vs Promotional Classification
SSL Wireless (as an aggregator) must tag and flag messages appropriately:

Transactional: OTP, PIN codes, order confirmations, account alerts, password resets, booking confirmations
Promotional: Marketing, offers, discounts, brand awareness campaigns, product launches
 
3. SSL Wireless API Specification
SSL Wireless provides two API platforms: ISMSPlus (modern, recommended) and legacy ISMS. This guide focuses on ISMSPlus, which uses JSON-based REST API with token authentication.
3.1 Authentication & Credentials
SSL Wireless provides the following credentials upon account creation:
•	SID (Service ID): Unique identifier for your account
•	API Token: Authentication token for API requests
•	Domain: API base URL (typically https://smsplus.sslwireless.com)

Access credentials from SSL Wireless portal:
Login: https://ismsplus.sslwireless.com/login
Navigate to: Profile > API Hash Token

3.2 API Endpoints
Platform	Endpoint URL	Method
ISMSPlus (Recommended)	https://smsplus.sslwireless.com/api/v3/send-sms	POST
Legacy ISMS	https://sms.sslwireless.com/pushapi/dynamic/server.php	POST

Note: This guide focuses on ISMSPlus. Legacy ISMS uses different authentication (username/password instead of API token).
 
3.3 Request Format
3.3.1 HTTP Headers
Content-Type: application/json
Accept: application/json

3.3.2 Request Body Parameters
Parameter	Required	Type/Format	Description
api_token	Yes	String	API authentication token from SSL Wireless
sid	Yes	String	Service ID for your account (English or Bangla SID)
msisdn	Yes	String	Recipient mobile number in international format (880XXXXXXXXXX)
sms	Yes	String (max 1000 chars)	SMS message body (English: 160 chars, Bangla: 70 chars per segment)
csms_id	Optional	String (unique)	Client-generated unique reference ID for tracking (highly recommended)
sms_type	Optional	String: "EN" or "BN"	Language indicator (EN=English/GSM-7, BN=Bangla/UCS-2)
sender	Optional	String (max 11 chars)	Override default sender ID (masking name). Must be pre-registered.
dlr_url	Optional	String (valid URL)	Webhook URL for delivery reports (can also be configured in portal)
batch	Optional	Array	For bulk SMS: array of objects with msisdn, text, csms_id
 
3.3.3 Example Request (Single SMS)
POST https://smsplus.sslwireless.com/api/v3/send-sms
Content-Type: application/json

{
  "api_token": "your-api-token-here",
  "sid": "YOURSIDHERE",
  "msisdn": "8801712345678",
  "sms": "Your OTP is 123456. Valid for 5 minutes. Do not share.",
  "csms_id": "OTP_20260130_1234567890",
  "sms_type": "EN"
}

3.3.4 Example Request (Bulk SMS - Same Message)
{
  "api_token": "your-api-token-here",
  "sid": "YOURSIDHERE",
  "msisdn": "8801712345678,8801812345678,8801912345678",
  "sms": "Notification message for all recipients",
  "csms_id": "BULK_20260130_1234567890"
}

3.3.5 Example Request (Dynamic SMS - Different Messages)
{
  "api_token": "your-api-token-here",
  "sid": "YOURSIDHERE",
  "sms": [
    {
      "msisdn": "8801712345678",
      "text": "Hello John, your OTP is 123456",
      "csms_id": "OTP_JOHN_123456"
    },
    {
      "msisdn": "8801812345678",
      "text": "Hello Jane, your OTP is 654321",
      "csms_id": "OTP_JANE_654321"
    }
  ]
}
 
3.4 Response Format
3.4.1 Success Response (HTTP 200)
{
  "status": "SUCCESS",
  "status_code": 200,
  "error_message": "",
  "smsinfo": [
    {
      "sms_status": "SUCCESS",
      "status_message": "Success",
      "msisdn": "8801712345678",
      "sms_type": "EN",
      "sms_body": "Your OTP is 123456. Valid for 5 minutes.",
      "csms_id": "150120261212010966449",
      "reference_id": "OTP_20260130_1234567890"
    }
  ]
}

3.4.2 Response Fields Explanation
Field	Description
status	Overall request status: "SUCCESS" or "FAILED"
status_code	HTTP status code (200 = success, 4XXX = client error, 5XXX = server error)
error_message	Error description (empty if successful)
csms_id	SSL Wireless gateway message ID. CRITICAL: Save this for matching delivery reports.
reference_id	Your original csms_id from the request (client reference)
 
3.5 Error Responses
Understanding error codes is critical for proper error handling and debugging.
3.5.1 Common Error Codes
Code	Error	Description	Action
4001	Invalid Request	Missing required parameters or invalid JSON format	Validate request payload
4003	Authentication Failed	Invalid API token or SID, or IP whitelist issue	Verify credentials and IP
4005	Invalid Mobile Number	Malformed phone number or unsupported operator	Do not retry. Log error.
4006	Insufficient Balance	Account balance too low to send SMS	Alert admin. Recharge.
4007	Invalid Sender ID	Sender ID not registered or inactive	Contact SSL support
5000	Internal Server Error	Gateway temporary issue	Retry with backoff
5003	Operator Down	Mobile operator network issue	Retry later (15-30 min)

3.5.2 Example Error Response
{
  "status": "FAILED",
  "status_code": 4003,
  "error_message": "Authentication Failed - Invalid API Token"
}
 
4. Delivery Reports (DLR)
Delivery reports confirm that the SMS was successfully delivered to the recipient's phone. Implementing DLR webhooks is essential for tracking message status and debugging delivery issues.
4.1 DLR Configuration
•	Method 1: Pass dlr_url parameter in each API request
•	Method 2: Configure a global DLR URL in the SSL Wireless portal (Settings > API Configuration)

4.2 DLR Webhook Format
SSL Wireless sends delivery reports via HTTP POST to your configured webhook URL:

POST https://your-server.com/dlr-webhook
Content-Type: application/x-www-form-urlencoded

message_id=150120261212010966449&
status=DELIVRD&
phone=8801712345678&
submit_time=2026-01-30 14:30:00&
done_time=2026-01-30 14:30:05&
operator=GRAMEENPHONE

4.3 DLR Parameters
Parameter	Description
message_id	Gateway message ID (matches csms_id from success response)
status	Delivery status: DELIVRD, FAILED, EXPIRED, REJECTED, UNDELIV
phone	Recipient phone number
submit_time	Time when message was submitted to operator
done_time	Time when final status was received
operator	Mobile operator name (GRAMEENPHONE, ROBI, BANGLALINK, TELETALK)

4.4 DLR Status Codes
Status	Meaning
DELIVRD	Successfully delivered to recipient
FAILED	Delivery failed (invalid number, network issue)
EXPIRED	Message expired before delivery (phone off for too long)
REJECTED	Rejected by operator (spam filter, invalid sender ID)
UNDELIV	Undeliverable (permanent failure)
 
4.5 DLR Webhook Implementation Example
// Node.js Express webhook endpoint
const express = require('express');
const app = express();

app.use(express.urlencoded({ extended: true }));

app.post('/dlr-webhook', async (req, res) => {
  try {
    const {
      message_id,
      status,
      phone,
      submit_time,
      done_time,
      operator
    } = req.body;
    
    console.log('DLR Received:', {
      message_id,
      status,
      phone,
      operator
    });
    
    // Update database
    await db.otps.update(
      { gateway_message_id: message_id },
      { 
        delivery_status: status,
        delivered_at: status === 'DELIVRD' ? done_time : null,
        operator: operator
      }
    );
    
    // Respond with 200 OK to acknowledge receipt
    res.status(200).send('OK');
    
  } catch (error) {
    console.error('DLR webhook error:', error);
    res.status(500).send('ERROR');
  }
});

app.listen(3000, () => {
  console.log('DLR webhook listening on port 3000');
});
 
5. Implementation Guide
This section provides production-ready code examples and security best practices for integrating SSL Wireless API into your authentication service.
5.1 Environment Setup
5.1.1 Node.js Dependencies
npm install axios dotenv crypto

5.1.2 Environment Variables (.env)
SSL_API_TOKEN=your_api_token_here
SSL_SID=YOUR_SID_HERE
SSL_API_BASE_URL=https://smsplus.sslwireless.com
SSL_SENDER_ID=MYAPP
DLR_WEBHOOK_URL=https://your-server.com/dlr-webhook
OTP_VALIDITY_MINUTES=5
OTP_LENGTH=6
 
5.2 Phone Number Validation and Normalization
/**
 * Normalize Bangladesh phone numbers to international format
 * Handles various input formats: 01712345678, 8801712345678, +8801712345678
 */
function normalizePhoneNumber(phone) {
  // Remove all non-numeric characters
  let cleaned = phone.replace(/\D/g, '');
  
  // Handle different input formats
  if (cleaned.startsWith('880')) {
    // Already in international format
    return cleaned;
  } else if (cleaned.startsWith('0')) {
    // Remove leading 0 and add country code
    return '880' + cleaned.substring(1);
  } else if (cleaned.length === 10) {
    // Add country code to 10-digit number
    return '880' + cleaned;
  }
  
  throw new Error('Invalid phone number format');
}

/**
 * Validate Bangladesh phone number
 * Valid operators: GP (017/013), Robi (018), Banglalink (019/014), Teletalk (015)
 */
function isValidBangladeshNumber(phone) {
  const normalized = normalizePhoneNumber(phone);
  
  // Check format: 880 followed by valid operator prefix and 8 more digits
  const validPattern = /^880(13|14|15|16|17|18|19)\d{8}$/;
  
  return validPattern.test(normalized);
}

// Example usage
try {
  const phone = normalizePhoneNumber('01712345678');
  if (isValidBangladeshNumber(phone)) {
    console.log('Valid phone:', phone);
  }
} catch (error) {
  console.error('Invalid phone number:', error.message);
}
 
5.3 Secure OTP Generation
const crypto = require('crypto');

/**
 * Generate cryptographically secure OTP
 * Uses crypto.randomInt for true randomness
 */
function generateOTP(length = 6) {
  if (length < 4 || length > 8) {
    throw new Error('OTP length must be between 4 and 8');
  }
  
  let otp = '';
  const digits = '0123456789';
  
  for (let i = 0; i < length; i++) {
    const randomIndex = crypto.randomInt(0, digits.length);
    otp += digits[randomIndex];
  }
  
  return otp;
}

/**
 * Hash OTP for secure storage
 * Never store plain text OTPs in database
 */
function hashOTP(otp) {
  return crypto
    .createHash('sha256')
    .update(otp)
    .digest('hex');
}

/**
 * Verify OTP against hash
 */
function verifyOTP(inputOTP, storedHash) {
  const inputHash = hashOTP(inputOTP);
  return inputHash === storedHash;
}

// Example usage
const otp = generateOTP(6);
const hashed = hashOTP(otp);
console.log('Generated OTP:', otp);
console.log('Hashed:', hashed);

// Verification
const isValid = verifyOTP('123456', hashed);
console.log('Verification:', isValid);
 
5.4 Complete OTP Send Implementation
const axios = require('axios');
const crypto = require('crypto');
require('dotenv').config();

/**
 * Send OTP via SSL Wireless
 * Implements best practices: validation, error handling, retry logic
 */
async function sendOTP(phoneNumber) {
  try {
    // Step 1: Validate and normalize phone number
    const msisdn = normalizePhoneNumber(phoneNumber);
    
    if (!isValidBangladeshNumber(msisdn)) {
      throw new Error('Invalid Bangladesh phone number');
    }
    
    // Step 2: Generate OTP and unique reference ID
    const otpCode = generateOTP(6);
    const timestamp = Date.now();
    const csms_id = `OTP_${timestamp}_${crypto.randomBytes(4).toString('hex')}`;
    
    // Step 3: Prepare SMS message
    const message = `Your verification code is ${otpCode}. Valid for 5 minutes. Do not share this code with anyone.`;
    
    // Step 4: Prepare SSL Wireless request
    const payload = {
      api_token: process.env.SSL_API_TOKEN,
      sid: process.env.SSL_SID,
      msisdn: msisdn,
      sms: message,
      csms_id: csms_id,
      sms_type: 'EN',
      sender: process.env.SSL_SENDER_ID,
      dlr_url: process.env.DLR_WEBHOOK_URL
    };
    
    // Step 5: Call SSL Wireless API with timeout and retry logic
    const response = await axios.post(
      `${process.env.SSL_API_BASE_URL}/api/v3/send-sms`,
      payload,
      {
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
        timeout: 10000 // 10 seconds
      }
    );
    
    // Step 6: Handle response
    if (response.data.status === 'SUCCESS') {
      const gatewayMessageId = response.data.smsinfo[0].csms_id;
      
      // Step 7: Store OTP in database (HASHED - never plain text)
      await storeOTP({
        phone: msisdn,
        otp_hash: hashOTP(otpCode),
        client_reference: csms_id,
        gateway_message_id: gatewayMessageId,
        status: 'SENT',
        attempts: 0,
        created_at: new Date(),
        expires_at: new Date(Date.now() + 5 * 60 * 1000) // 5 minutes
      });
      
      // SECURITY: Never return OTP in response
      return {
        success: true,
        message: 'OTP sent successfully',
        data: {
          phone: msisdn,
          validity: '5 minutes',
          reference: csms_id
        }
      };
      
    } else {
      // Handle SSL Wireless error response
      throw new Error(response.data.error_message || 'Failed to send SMS');
    }
    
  } catch (error) {
    // Error handling
    if (error.response) {
      // SSL Wireless API error
      const statusCode = error.response.data.status_code;
      const errorMessage = error.response.data.error_message;
      
      console.error('SSL Wireless API Error:', {
        code: statusCode,
        message: errorMessage
      });
      
      // Don't retry on client errors (4xxx)
      if (statusCode >= 4000 && statusCode < 5000) {
        return {
          success: false,
          error: 'VALIDATION_ERROR',
          message: errorMessage
        };
      }
      
      // Retry on server errors (5xxx) with exponential backoff
      if (statusCode >= 5000) {
        return {
          success: false,
          error: 'SERVICE_UNAVAILABLE',
          message: 'SMS service temporarily unavailable. Please try again.'
        };
      }
      
    } else if (error.request) {
      // Network error
      console.error('Network error:', error.message);
      return {
        success: false,
        error: 'NETWORK_ERROR',
        message: 'Unable to connect to SMS service'
      };
      
    } else {
      // Other errors (validation, etc.)
      console.error('Error:', error.message);
      return {
        success: false,
        error: 'INTERNAL_ERROR',
        message: error.message
      };
    }
  }
}

/**
 * Database storage function (example with Mongoose)
 */
async function storeOTP(otpData) {
  // Example: MongoDB with Mongoose
  // const OTP = require('./models/OTP');
  // await OTP.create(otpData);
  
  // Example: PostgreSQL with Knex
  // await knex('otps').insert(otpData);
  
  console.log('OTP stored in database:', {
    phone: otpData.phone,
    reference: otpData.client_reference
  });
}

// Export for use in routes
module.exports = { sendOTP };
 
5.5 OTP Verification Implementation
/**
 * Verify OTP entered by user
 * Implements rate limiting and attempt tracking
 */
async function verifyOTP(phoneNumber, inputOTP) {
  try {
    const msisdn = normalizePhoneNumber(phoneNumber);
    
    // Fetch OTP record from database
    const otpRecord = await db.otps.findOne({
      phone: msisdn,
      status: 'SENT'
    }).sort({ created_at: -1 }); // Get latest OTP
    
    if (!otpRecord) {
      return {
        success: false,
        error: 'OTP_NOT_FOUND',
        message: 'No OTP found for this number'
      };
    }
    
    // Check expiration
    if (new Date() > otpRecord.expires_at) {
      await db.otps.update(
        { _id: otpRecord._id },
        { status: 'EXPIRED' }
      );
      
      return {
        success: false,
        error: 'OTP_EXPIRED',
        message: 'OTP has expired. Please request a new one.'
      };
    }
    
    // Check max attempts (prevent brute force)
    if (otpRecord.attempts >= 5) {
      await db.otps.update(
        { _id: otpRecord._id },
        { status: 'BLOCKED' }
      );
      
      return {
        success: false,
        error: 'MAX_ATTEMPTS_EXCEEDED',
        message: 'Too many failed attempts. Please request a new OTP.'
      };
    }
    
    // Verify OTP
    const isValid = verifyOTP(inputOTP, otpRecord.otp_hash);
    
    // Increment attempts
    await db.otps.update(
      { _id: otpRecord._id },
      { 
        $inc: { attempts: 1 },
        last_attempt_at: new Date()
      }
    );
    
    if (isValid) {
      // Mark as verified
      await db.otps.update(
        { _id: otpRecord._id },
        { 
          status: 'VERIFIED',
          verified_at: new Date()
        }
      );
      
      return {
        success: true,
        message: 'OTP verified successfully',
        data: {
          phone: msisdn,
          verified_at: new Date()
        }
      };
      
    } else {
      return {
        success: false,
        error: 'INVALID_OTP',
        message: `Invalid OTP. ${5 - otpRecord.attempts - 1} attempts remaining.`
      };
    }
    
  } catch (error) {
    console.error('OTP verification error:', error);
    return {
      success: false,
      error: 'VERIFICATION_ERROR',
      message: 'Failed to verify OTP'
    };
  }
}

module.exports = { verifyOTP };
 
5.6 Express.js Route Implementation
const express = require('express');
const router = express.Router();
const { sendOTP, verifyOTP } = require('./otpService');

// Rate limiting middleware (example with express-rate-limit)
const rateLimit = require('express-rate-limit');

const otpLimiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 3, // 3 requests per window
  message: 'Too many OTP requests. Please try again later.',
  standardHeaders: true,
  legacyHeaders: false
});

/**
 * POST /api/otp/send
 * Send OTP to phone number
 */
router.post('/api/otp/send', otpLimiter, async (req, res) => {
  try {
    const { phone } = req.body;
    
    // Validation
    if (!phone) {
      return res.status(400).json({
        success: false,
        error: 'MISSING_PARAMETER',
        message: 'Phone number is required'
      });
    }
    
    // Send OTP
    const result = await sendOTP(phone);
    
    if (result.success) {
      return res.status(200).json(result);
    } else {
      return res.status(400).json(result);
    }
    
  } catch (error) {
    console.error('Send OTP route error:', error);
    return res.status(500).json({
      success: false,
      error: 'INTERNAL_ERROR',
      message: 'Failed to send OTP'
    });
  }
});

/**
 * POST /api/otp/verify
 * Verify OTP
 */
router.post('/api/otp/verify', async (req, res) => {
  try {
    const { phone, otp } = req.body;
    
    // Validation
    if (!phone || !otp) {
      return res.status(400).json({
        success: false,
        error: 'MISSING_PARAMETERS',
        message: 'Phone number and OTP are required'
      });
    }
    
    // Verify OTP
    const result = await verifyOTP(phone, otp);
    
    if (result.success) {
      // Create session or JWT token here
      // const token = generateJWT({ phone: result.data.phone });
      
      return res.status(200).json({
        ...result,
        // token: token
      });
    } else {
      return res.status(400).json(result);
    }
    
  } catch (error) {
    console.error('Verify OTP route error:', error);
    return res.status(500).json({
      success: false,
      error: 'INTERNAL_ERROR',
      message: 'Failed to verify OTP'
    });
  }
});

module.exports = router;
 
6. Security Best Practices
6.1 Critical Security Rules
•	NEVER return OTP code in API responses  (prevents interception)
•	ALWAYS hash OTPs before storing in database  (use SHA-256 or bcrypt)
•	Implement rate limiting  (max 3 OTP requests per 15 minutes per phone)
•	Set short expiration times  (5 minutes maximum)
•	Limit verification attempts  (max 5 attempts before blocking)
•	Use HTTPS for all API communications
•	Validate phone numbers server-side  (never trust client validation)
•	Store credentials in environment variables  (never in code)
•	Implement CSRF protection for web applications
•	Log all OTP activities for audit trails

6.2 Database Schema Best Practices
-- PostgreSQL schema example
CREATE TABLE otps (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(15) NOT NULL,
    otp_hash VARCHAR(64) NOT NULL,  -- SHA-256 hash
    client_reference VARCHAR(100) UNIQUE,
    gateway_message_id VARCHAR(100),
    status VARCHAR(20) DEFAULT 'SENT',  -- SENT, VERIFIED, EXPIRED, BLOCKED
    delivery_status VARCHAR(20),  -- DELIVRD, FAILED, etc.
    attempts INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    verified_at TIMESTAMP,
    last_attempt_at TIMESTAMP,
    operator VARCHAR(50),
    INDEX idx_phone_status (phone, status),
    INDEX idx_gateway_msg (gateway_message_id),
    INDEX idx_expires (expires_at)
);

-- Cleanup expired OTPs periodically
DELETE FROM otps WHERE expires_at < NOW() - INTERVAL '1 day';
 
6.3 Rate Limiting Implementation
const rateLimit = require('express-rate-limit');
const RedisStore = require('rate-limit-redis');
const redis = require('redis');

// Create Redis client
const redisClient = redis.createClient({
  host: process.env.REDIS_HOST,
  port: process.env.REDIS_PORT
});

// OTP send rate limiter (3 requests per 15 minutes per IP)
const otpSendLimiter = rateLimit({
  store: new RedisStore({
    client: redisClient,
    prefix: 'rl:otp:send:'
  }),
  windowMs: 15 * 60 * 1000,  // 15 minutes
  max: 3,  // 3 requests per window
  message: {
    success: false,
    error: 'RATE_LIMIT_EXCEEDED',
    message: 'Too many OTP requests. Please try again after 15 minutes.'
  },
  standardHeaders: true,
  legacyHeaders: false,
  // Custom key generator (use phone number instead of IP for better control)
  keyGenerator: (req) => {
    return req.body.phone || req.ip;
  }
});

// OTP verify rate limiter (5 attempts per 5 minutes per phone)
const otpVerifyLimiter = rateLimit({
  store: new RedisStore({
    client: redisClient,
    prefix: 'rl:otp:verify:'
  }),
  windowMs: 5 * 60 * 1000,  // 5 minutes
  max: 5,  // 5 attempts
  message: {
    success: false,
    error: 'RATE_LIMIT_EXCEEDED',
    message: 'Too many verification attempts. Please request a new OTP.'
  },
  keyGenerator: (req) => {
    return req.body.phone;
  }
});

// Apply to routes
router.post('/api/otp/send', otpSendLimiter, sendOTPHandler);
router.post('/api/otp/verify', otpVerifyLimiter, verifyOTPHandler);
 
6.4 Monitoring and Alerting
Implement monitoring for the following metrics:
•	OTP send success rate (should be > 95%)
•	SMS delivery rate (track via DLR, should be > 90%)
•	Average delivery time (should be < 10 seconds)
•	Failed authentication attempts
•	SSL Wireless API errors (4xxx, 5xxx codes)
•	Account balance alerts (alert at 20% remaining)
•	Rate limit violations
•	Suspicious patterns (multiple failed verifications)

// Example monitoring with Prometheus
const promClient = require('prom-client');

const otpSentCounter = new promClient.Counter({
  name: 'otp_sent_total',
  help: 'Total number of OTPs sent',
  labelNames: ['status', 'operator']
});

const otpVerifyCounter = new promClient.Counter({
  name: 'otp_verify_total',
  help: 'Total number of OTP verifications',
  labelNames: ['result']
});

const smsDeliveryGauge = new promClient.Gauge({
  name: 'sms_delivery_rate',
  help: 'SMS delivery success rate'
});

// Track metrics in your code
async function sendOTP(phone) {
  const result = await sslWirelessSend(phone);
  
  otpSentCounter.labels({
    status: result.success ? 'success' : 'failed',
    operator: result.operator || 'unknown'
  }).inc();
  
  return result;
}
 
7. Testing and Debugging
7.1 Test Accounts
SSL Wireless provides test credentials for development. Contact support to obtain:
•	Test API token and SID
•	Test sender ID (usually TEST or your company name)
•	Test phone numbers that don't consume credits
•	Sandbox environment URL (if available)

7.2 Common Issues and Solutions
Issue	Solution
Authentication Failed (4003)	Verify API token and SID are correct. Check if your IP is whitelisted in SSL portal.
Invalid Sender ID (4007)	Ensure sender ID is registered and approved. Check if using correct case (UPPERCASE).
SMS not delivered	Check DLR status. Verify phone number format. Ensure sufficient balance.
Timeout errors	Increase timeout to 15-30 seconds. Check network connectivity. Retry with exponential backoff.
DLR webhook not receiving	Verify webhook URL is publicly accessible. Check firewall rules. Test with ngrok for local dev.
Message truncated	English: 160 chars/segment. Bangla: 70 chars/segment. Split long messages or use concatenation.

7.3 Debugging Checklist
•	Verify SSL Wireless credentials are loaded from environment
•	Check phone number format (must be 880XXXXXXXXXX)
•	Ensure sender ID is registered and active
•	Confirm sufficient account balance
•	Test with a known working phone number first
•	Check API response for error messages
•	Monitor DLR webhook for delivery status
•	Review application logs for errors
•	Test network connectivity to SSL Wireless servers
•	Verify request payload matches API specification exactly
 
7.4 Testing Script
#!/usr/bin/env node
/**
 * SSL Wireless Integration Test Script
 * Run: node test-ssl-wireless.js
 */

require('dotenv').config();
const axios = require('axios');

async function testSSLWireless() {
  console.log('=== SSL Wireless Integration Test ===\n');
  
  // 1. Environment Check
  console.log('1. Checking environment variables...');
  const required = ['SSL_API_TOKEN', 'SSL_SID', 'SSL_SENDER_ID'];
  const missing = required.filter(key => !process.env[key]);
  
  if (missing.length > 0) {
    console.error('❌ Missing environment variables:', missing.join(', '));
    process.exit(1);
  }
  console.log('✅ All environment variables present\n');
  
  // 2. API Connectivity
  console.log('2. Testing API connectivity...');
  const testPayload = {
    api_token: process.env.SSL_API_TOKEN,
    sid: process.env.SSL_SID,
    msisdn: '8801712345678',  // Replace with your test number
    sms: 'Test message from integration test',
    csms_id: 'TEST_' + Date.now()
  };
  
  try {
    const response = await axios.post(
      'https://smsplus.sslwireless.com/api/v3/send-sms',
      testPayload,
      { timeout: 10000 }
    );
    
    if (response.data.status === 'SUCCESS') {
      console.log('✅ API connectivity successful');
      console.log('   Gateway Message ID:', response.data.smsinfo[0].csms_id);
    } else {
      console.error('❌ API returned error:', response.data.error_message);
    }
  } catch (error) {
    console.error('❌ API request failed:', error.message);
    if (error.response) {
      console.error('   Status:', error.response.data.status_code);
      console.error('   Error:', error.response.data.error_message);
    }
  }
  
  console.log('\n=== Test Complete ===');
}

testSSLWireless();
 
8. Production Deployment Checklist
8.1 Pre-Deployment
•	SSL Wireless account fully verified and active
•	Sender ID(s) registered and approved by all operators
•	Account funded with sufficient balance (recommended: 1 month buffer)
•	Production API credentials secured in environment variables
•	IP whitelist configured in SSL Wireless portal (if required)
•	DLR webhook URL configured and tested
•	Rate limiting implemented and tested
•	OTP hashing implemented (never store plain text)
•	Database schema deployed with proper indexes
•	Monitoring and alerting configured
•	Error handling and retry logic implemented
•	Load testing completed (simulate 1000 OTP/hour)
•	Security audit completed

8.2 Post-Deployment
•	Monitor OTP send success rate (target > 95%)
•	Monitor SMS delivery rate via DLR (target > 90%)
•	Set up balance alerts (alert at 20% remaining)
•	Configure log aggregation and search
•	Implement automated balance recharge (if available)
•	Schedule regular security audits
•	Document runbook for common issues
•	Train support team on OTP troubleshooting
•	Set up on-call rotation for critical issues

8.3 Maintenance Schedule
Frequency	Task	Owner
Daily	Check balance and delivery rates	Operations
Weekly	Review error logs and failed OTPs	Engineering
Monthly	Security audit and dependency updates	Security Team
Quarterly	Review and optimize OTP flow	Product Team
 
9. Appendix
9.1 Useful Resources
•	SSL Wireless Portal: https://ismsplus.sslwireless.com
•	SSL Wireless Support: support@sslwireless.com
•	BTRC Official Website: http://www.btrc.gov.bd
•	Bangladesh Mobile Operators: GP, Robi, Banglalink, Teletalk

9.2 SMS Character Encoding Reference
Encoding	Characters per Segment	Use Case
GSM-7 (English)	160 chars (single) / 153 chars (concatenated)	English text, numbers, basic symbols
UCS-2 (Bangla)	70 chars (single) / 67 chars (concatenated)	Bangla text, Unicode characters

9.3 HTTP Status Codes Quick Reference
•	200 - Success - Request processed successfully
•	400 - Bad Request - Invalid parameters
•	401 - Unauthorized - Invalid credentials
•	403 - Forbidden - IP not whitelisted or account suspended
•	429 - Too Many Requests - Rate limit exceeded
•	500 - Internal Server Error - SSL Wireless temporary issue
•	503 - Service Unavailable - Gateway maintenance

9.4 Support and Contact Information
For technical support and assistance:

•	SSL Wireless Support Email: support@sslwireless.com
•	Technical Hotline: +880-2-XXXXXXXX (contact SSL for number)
•	Business Hours: Saturday - Thursday, 9:00 AM - 6:00 PM (BST)
•	Emergency Support: Available for critical production issues


End of Document
Version 2.0 - January 2026
