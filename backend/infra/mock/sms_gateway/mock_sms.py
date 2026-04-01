"""
Mock SSL Wireless SMS Gateway for local dev/testing.

Implements:
  POST /api/v3/send-sms  — accepts SMS, extracts OTP from message, stores it
  GET  /mock/last-otp    — returns the last OTP sent (for test automation)
  GET  /mock/messages    — returns all messages sent
  DELETE /mock/reset     — clears all stored messages
  GET  /health           — health check
"""

import re
import uuid
from datetime import datetime
from flask import Flask, request, jsonify

app = Flask(__name__)

# In-memory store of sent messages
messages = []


def extract_otp(text: str) -> str | None:
    """Extract a 6-digit OTP from SMS message text."""
    match = re.search(r'\b(\d{6})\b', text)
    return match.group(1) if match else None


@app.route('/api/v3/send-sms', methods=['POST'])
def send_sms():
    data = request.get_json(force=True, silent=True) or {}
    msisdn = data.get('msisdn', '')
    sms_text = data.get('sms', '')
    sender_id = data.get('sender_id', 'MOCK')
    csms_id = data.get('csms_id', '')
    reference_id = f"mock-{uuid.uuid4().hex[:12]}"

    otp = extract_otp(sms_text)

    msg = {
        'reference_id': reference_id,
        'msisdn': msisdn,
        'sender_id': sender_id,
        'sms': sms_text,
        'otp': otp,
        'csms_id': csms_id,
        'sent_at': datetime.utcnow().isoformat() + 'Z',
    }
    messages.append(msg)

    print(f"[MOCK SMS] To={msisdn} OTP={otp} Ref={reference_id} Text={sms_text!r}", flush=True)

    return jsonify({
        "status": "SUCCESS",
        "status_code": 200,
        "error_message": "",
        "smsinfo": [
            {
                "sms_status": "SUCCESS",
                "status_message": "Message sent successfully",
                "msisdn": msisdn,
                "csms_id": csms_id,
                "reference_id": reference_id,
            }
        ]
    }), 200


@app.route('/mock/last-otp', methods=['GET'])
def last_otp():
    """Return the last OTP sent, optionally filtered by msisdn."""
    msisdn = request.args.get('msisdn')
    filtered = [m for m in reversed(messages) if m.get('otp')]
    if msisdn:
        filtered = [m for m in filtered if msisdn in m.get('msisdn', '')]
    if not filtered:
        return jsonify({"otp": None, "message": "No OTP found"}), 404
    last = filtered[0]
    return jsonify({
        "otp": last['otp'],
        "msisdn": last['msisdn'],
        "reference_id": last['reference_id'],
        "sent_at": last['sent_at'],
    }), 200


@app.route('/mock/messages', methods=['GET'])
def list_messages():
    return jsonify({"count": len(messages), "messages": list(reversed(messages))}), 200


@app.route('/mock/reset', methods=['DELETE'])
def reset():
    messages.clear()
    return jsonify({"cleared": True}), 200


@app.route('/health', methods=['GET'])
def health():
    return jsonify({"status": "ok", "service": "mock-sms-gateway", "messages_sent": len(messages)}), 200


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8600, debug=False)
