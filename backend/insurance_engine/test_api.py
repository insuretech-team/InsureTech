import grpc
import sys
import os
from datetime import datetime

sys.path.insert(0, 'D:/InsureTech/gen/csharp')
os.chdir('D:/InsureTech/backend/insurance_engine')

import Insuretech.Policy.Services.V1.PolicyService_pb2 as policy_pb2
import Insuretech.Policy.Services.V1.PolicyServiceGrpc as policy_pb2_grpc
import Insuretech.Products.Services.V1.ProductService_pb2 as product_pb2
import Insuretech.Products.Services.V1.ProductServiceGrpc as product_pb2_grpc

def test_product_service():
    print("=== Testing ProductService gRPC ===")
    
    channel = grpc.insecure_channel('localhost:5001')
    stub = product_pb2_grpc.ProductServiceStub(channel)
    
    try:
        request = product_pb2.SearchProductsRequest(query="OMC", page=1, pageSize=10)
        response = stub.SearchProducts(request, timeout=5)
        print(f"Search Products: Found {len(response.products)} products")
        for p in response.products[:3]:
            print(f"  - {p.product_id}: {p.name} ({p.code})")
    except grpc.RpcError as e:
        print(f"Error: {e.code()}: {e.details()}")
    
    channel.close()

def test_policy_service():
    print("\n=== Testing PolicyService gRPC ===")
    
    channel = grpc.insecure_channel('localhost:5001')
    stub = policy_pb2_grpc.PolicyServiceStub(channel)
    
    try:
        request = policy_pb2.CreatePolicyRequest(
            product_id="OMC-PRODUCT-001",
            customer_id="CUST-001",
            premium_amount=1239,
            sum_insured=50000,
            tenure_months=1,
            start_date="2026-04-15T00:00:00Z"
        )
        response = stub.CreatePolicy(request, timeout=10)
        print(f"Create Policy Response:")
        print(f"  PolicyId: {response.policy_id}")
        print(f"  PolicyNumber: {response.policy_number}")
        print(f"  Message: {response.message}")
        if response.error:
            print(f"  Error: {response.error.code} - {response.error.message}")
    except grpc.RpcError as e:
        print(f"Error: {e.code()}: {e.details()}")
    
    channel.close()

def test_policy_with_nominees():
    print("\n=== Testing Policy with Nominees (from Excel Data) ===")
    
    channel = grpc.insecure_channel('localhost:5001')
    stub = policy_pb2_grpc.PolicyServiceStub(channel)
    
    try:
        from Insuretech.Policy.Entity.V1 import Nominee
        
        nominees = [
            Nominee(
                full_name="Fatema Begum",
                relationship="Wife",
                share_percentage=50,
                nid_number="198515678900001",
                phone_number="+88017111234567"
            ),
            Nominee(
                full_name="Ahmed Reza",
                relationship="Son",
                share_percentage=50,
                nid_number="",
                phone_number=""
            )
        ]
        
        request = policy_pb2.CreatePolicyRequest(
            product_id="OMC-PRODUCT-001",
            customer_id="CUST-002",
            premium_amount=2499,
            sum_insured=50000,
            tenure_months=1,
            start_date="2026-04-20T00:00:00Z",
            proposer_details="Name: Md. Zubayed Ur Rahman, Age: 35, Trip: Business",
            nominees=nominees
        )
        response = stub.CreatePolicy(request, timeout=10)
        print(f"Policy with Nominees Created:")
        print(f"  PolicyId: {response.policy_id}")
        print(f"  PolicyNumber: {response.policy_number}")
    except grpc.RpcError as e:
        print(f"Error: {e.code()}: {e.details()}")
    except Exception as e:
        print(f"Error: {e}")
    
    channel.close()

if __name__ == '__main__':
    print("Insurance Engine API Test")
    print("=" * 50)
    print(f"Testing against: localhost:5001")
    print(f"Time: {datetime.now()}")
    print("=" * 50)
    
    try:
        test_product_service()
    except Exception as e:
        print(f"Product service test failed: {e}")
    
    try:
        test_policy_service()
    except Exception as e:
        print(f"Policy service test failed: {e}")
    
    try:
        test_policy_with_nominees()
    except Exception as e:
        print(f"Policy with nominees test failed: {e}")
