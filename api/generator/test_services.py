import sys
sys.path.insert(0, '.')
from proto_parser import ProtoParser

parser = ProtoParser()
parser.load_descriptor_set('../input/descriptors.pb')
services = parser.get_services()

print(f'Found {len(services)} services')
for i, s in enumerate(services[:10]):
    svc_name = s['descriptor'].name
    method_count = len(s['methods'])
    print(f'  {i+1}. {svc_name} ({method_count} methods)')
    
    # Check first method for http annotation
    if s['methods']:
        first_method = s['methods'][0]
        print(f'     First method: {first_method.get("name", "unknown")}')
