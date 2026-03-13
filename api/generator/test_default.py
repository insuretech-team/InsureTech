from proto_source_parser import ProtoSourceParser

p = ProtoSourceParser('../../proto')
p.scan_all_protos()

ann = p.get_field_annotations('insuretech.apikey.entity.v1.ApiKey', 'status')
print(f'Annotations: {ann}')
if 'default_value' in ann:
    print(f'Default value: {ann["default_value"]}')
    print(f'Default value repr: {repr(ann["default_value"])}')
