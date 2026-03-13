import os
import re
import collections

orders = collections.defaultdict(list)

# Regex to roughly find migration_order: <num>
pattern = re.compile(r'migration_order\s*:\s*(\d+)')

def scan_dir(d):
    for root, _, files in os.walk(d):
        for f in files:
            if f.endswith('.proto'):
                filepath = os.path.join(root, f)
                with open(filepath, 'r', encoding='utf-8') as file:
                    content = file.read()
                    matches = pattern.findall(content)
                    for m in matches:
                        orders[int(m)].append(filepath)

scan_dir('e:\\Projects\\InsureTech\\proto')

print("=== MIGRATION ORDER REPORT ===")
duplicates = {k: v for k, v in orders.items() if len(v) > 1}

if duplicates:
    print("WARNING: FOUND DUPLICATE MIGRATION ORDERS!")
    for order in sorted(duplicates.keys()):
        print(f"Order {order} is used in:")
        for path in duplicates[order]:
            # Print just the path relative to proto
            rel = os.path.relpath(path, 'e:\\Projects\\InsureTech\\proto')
            print(f"  - {rel}")
else:
    print("No duplicate migration orders found.")

print("\nAll parsed migration orders:")
for order in sorted(orders.keys()):
    files = ", ".join([os.path.relpath(p, 'e:\\Projects\\InsureTech\\proto') for p in orders[order]])
    print(f"Order {order}: {files}")
