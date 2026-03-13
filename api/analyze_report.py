import json
from collections import defaultdict

with open('validation_report.json') as f:
    report = json.load(f)

print('=== Validation Report Analysis ===\n')

# Summary
summary = report.get('summary', {})
print('Summary:')
print(f'  Total Issues: {summary.get("total_issues", 0)}')
print(f'  Errors: {summary.get("errors", 0)}')
print(f'  Warnings: {summary.get("warnings", 0)}')
print(f'  Info: {summary.get("info", 0)}\n')

# Metrics
metrics = report.get('metrics', {})
print('Metrics:')
print(f'  Total Schemas: {metrics.get("total_schemas", 0)}')
print(f'  Schemas with Descriptions: {metrics.get("schemas_with_descriptions", 0)}')
print(f'  Description Coverage: {metrics.get("description_coverage", 0):.1f}%')
print(f'  Example Coverage: {metrics.get("example_coverage", 0):.1f}%\n')

# Group issues by category
issues_data = report.get('issues', {})
warnings = issues_data.get('warnings', [])
infos = issues_data.get('info', [])

print(f'Warnings ({len(warnings)}):')
by_category = defaultdict(list)
for warning in warnings:
    by_category[warning.get('category', 'Unknown')].append(warning)

for category, cat_warnings in sorted(by_category.items(), key=lambda x: -len(x[1])):
    print(f'  {category}: {len(cat_warnings)}')
    if cat_warnings:
        print(f'    Example: {cat_warnings[0].get("message", "")}')

print(f'\nInfo ({len(infos)}):')
by_category_info = defaultdict(list)
for info in infos:
    by_category_info[info.get('category', 'Unknown')].append(info)

for category, cat_infos in sorted(by_category_info.items(), key=lambda x: -len(x[1])):
    print(f'  {category}: {len(cat_infos)}')
