#!/usr/bin/env python3
"""
Insurance Engine gRPC API Test Script
Uses test data from Test.xlsx
"""
import json
import requests

BASE_URL = "http://localhost:5001"

def test_health():
    """Test health endpoint"""
    resp = requests.get(f"{BASE_URL}/")
    print(f"Health Check: {resp.json()}")
    return resp.ok

def test_create_policy():
    """Test CreatePolicy using Excel data from Sheet1 (Overseas Mediclaim)"""
    # Data from Sheet1 - Overseas Mediclaim Proposal
    # Proposer: Md. Zubayed Ur Rahman
    data = {
        "productId": "OMC-PRODUCT-001",
        "customerId": "CUST-PRAGATI-001",
        "premiumAmount": 1239,
        "sumInsured": 50000,
        "tenureMonths": 1,
        "startDate": "2026-04-15T00:00:00Z",
        "proposerDetails": "Name: Md. Zubayed Ur Rahman, Age: 35, Trip: Business, Passport: GA-18-6525"
    }
    resp = requests.post(f"{BASE_URL}/v1/policies", json=data)
    print(f"\n=== Create Policy (Overseas Mediclaim) ===")
    print(f"Status: {resp.status_code}")
    print(f"Response: {resp.text[:500]}")
    return resp

def test_create_policy_vehicle():
    """Test CreatePolicy using Excel data from Sheet8 (Private Vehicle Insurance)"""
    # Data from Sheet8 - Private Vehicle Insurance
    data = {
        "productId": "VEHICLE-PRODUCT-001",
        "customerId": "CUST-PRAGATI-002",
        "premiumAmount": 5000,
        "sumInsured": 500000,
        "tenureMonths": 12,
        "startDate": "2026-04-15T00:00:00Z",
        "proposerDetails": "Name: Md. Zubayed Ur Rahman, Vehicle: Dhaka Metro GA-18-6525, Chassis: NKE165-7216292"
    }
    resp = requests.post(f"{BASE_URL}/v1/policies", json=data)
    print(f"\n=== Create Policy (Vehicle Insurance) ===")
    print(f"Status: {resp.status_code}")
    print(f"Response: {resp.text[:500]}")
    return resp

def test_create_policy_health():
    """Test CreatePolicy using Excel data from Sheet14 (Health Insurance Claim)"""
    data = {
        "productId": "HEALTH-PRODUCT-001",
        "customerId": "CUST-PRAGATI-003",
        "premiumAmount": 10000,
        "sumInsured": 100000,
        "tenureMonths": 12,
        "startDate": "2026-04-15T00:00:00Z",
        "proposerDetails": "Name: Md. Zubayed Ur Rahman, Org: Medland Bank Plc, Health Insurance Claim Form"
    }
    resp = requests.post(f"{BASE_URL}/v1/policies", json=data)
    print(f"\n=== Create Policy (Health Insurance) ===")
    print(f"Status: {resp.status_code}")
    print(f"Response: {resp.text[:500]}")
    return resp

def test_list_policies():
    """Test ListPolicies"""
    resp = requests.get(f"{BASE_URL}/v1/policies?page=1&pageSize=10")
    print(f"\n=== List Policies ===")
    print(f"Status: {resp.status_code}")
    print(f"Response: {resp.text[:500]}")
    return resp

def test_create_claim():
    """Test CreateClaim from Sheet14 Health Insurance Claim Form"""
    data = {
        "policyId": "test-policy-001",
        "customerId": "CUST-PRAGATI-001",
        "type": "HEALTH",
        "claimedAmount": 25000,
        "incidentDate": "2026-04-01T00:00:00Z",
        "incidentDescription": "Hospitalization - Medland Bank Plc Employee"
    }
    resp = requests.post(f"{BASE_URL}/v1/claims", json=data)
    print(f"\n=== Submit Claim ===")
    print(f"Status: {resp.status_code}")
    print(f"Response: {resp.text[:500]}")
    return resp

def test_premium_quote():
    """Test Premium Quote - using data from Sheet4/Sheet6 (Premium Rates)"""
    # Sheet4: Non-Schengen Business & Holiday
    # Age 41-50, Period 14 days = 1860 BDT
    data = {
        "productId": "OMC-PRODUCT-001",
        "age": 35,
        "tenureDays": 14,
        "plan": "A"
    }
    resp = requests.post(f"{BASE_URL}/v1/quotes", json=data)
    print(f"\n=== Premium Quote ===")
    print(f"Status: {resp.status_code}")
    print(f"Response: {resp.text[:500]}")
    return resp

if __name__ == "__main__":
    print("Insurance Engine API Test")
    print("=" * 50)
    print(f"Testing against: {BASE_URL}")
    print("=" * 50)

    try:
        if not test_health():
            print("Health check failed!")
            exit(1)
    except Exception as e:
        print(f"Health check error: {e}")

    try:
        test_create_policy()
    except Exception as e:
        print(f"Create policy error: {e}")

    try:
        test_list_policies()
    except Exception as e:
        print(f"List policies error: {e}")

    try:
        test_create_claim()
    except Exception as e:
        print(f"Create claim error: {e}")
