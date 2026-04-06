#!/usr/bin/env python3
"""
Test Insurance Engine API with Excel data
"""
import requests
import json

BASE_URL = "http://localhost:5001"

def test_health():
    """Test health endpoint"""
    resp = requests.get(f"{BASE_URL}/")
    print(f"Health: {resp.json()}")
    return resp.ok

def test_create_policy_overseas_mediclaim():
    """Test CreatePolicy using Sheet1 - Overseas Mediclaim data"""
    data = {
        "productId": "OMC-PRODUCT-001",
        "customerId": "CUST-PRAGATI-001",
        "premiumAmount": 1239,
        "sumInsured": 50000,
        "tenureMonths": 1,
        "startDate": "2026-04-15T00:00:00Z",
        "proposerDetails": "Name: Md. Zubayed Ur Rahman, Address: N.B Tower Level-5, 40/7 North Avenue, Gulshan-2, Dhaka-1212, Mobile: 01985700011, Email: Zubayer@ymail.com, Occupation: Service at Medland Bank Plc, Passport: GA-18-6525, Plan: Business & Holiday"
    }
    resp = requests.post(f"{BASE_URL}/v1/policies", json=data)
    print(f"\n=== Create Policy (Overseas Mediclaim - Sheet1) ===")
    print(f"Status: {resp.status_code}")
    print(f"Response: {resp.text[:500]}")
    return resp

def test_create_policy_vehicle():
    """Test CreatePolicy using Sheet8 - Vehicle Insurance data"""
    data = {
        "productId": "VEHICLE-PRODUCT-001",
        "customerId": "CUST-PRAGATI-002",
        "premiumAmount": 5000,
        "sumInsured": 500000,
        "tenureMonths": 12,
        "startDate": "2026-04-15T00:00:00Z",
        "proposerDetails": "Vehicle: Dhaka Metro GA-18-6525, Chassis: NKE165-7216292, Engine: G4NAEM48921, Make: Hyundai, Model: Tucson, Year: 2024"
    }
    resp = requests.post(f"{BASE_URL}/v1/policies", json=data)
    print(f"\n=== Create Policy (Vehicle Insurance - Sheet8) ===")
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

def test_get_policy(policy_id):
    """Test GetPolicy"""
    resp = requests.get(f"{BASE_URL}/v1/policies/{policy_id}")
    print(f"\n=== Get Policy {policy_id} ===")
    print(f"Status: {resp.status_code}")
    print(f"Response: {resp.text[:500]}")
    return resp

if __name__ == "__main__":
    print("=" * 60)
    print("Insurance Engine API Test with Excel Data")
    print("=" * 60)
    
    if not test_health():
        print("Health check failed!")
        exit(1)
    
    # Test create policy - Overseas Mediclaim
    resp1 = test_create_policy_overseas_mediclaim()
    
    # Test create policy - Vehicle
    resp2 = test_create_policy_vehicle()
    
    # Test list policies
    test_list_policies()
    
    print("\n" + "=" * 60)
    print("Tests completed!")
    print("=" * 60)
