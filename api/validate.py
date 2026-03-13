from openapi_spec_validator import validate_spec
from openapi_spec_validator.readers import read_from_filename

try:
    spec_dict, base_uri = read_from_filename('openapi.yaml')
    validate_spec(spec_dict)
    print('✅ OpenAPI spec is VALID - 0 ERRORS')
except Exception as e:
    print(f'❌ Validation error:')
    print(str(e))
