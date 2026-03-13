"""
Validate description quality and coverage
"""
import yaml
from collections import defaultdict

class DescriptionValidator:
    def __init__(self, openapi_file):
        self.openapi_file = openapi_file
        self.spec = None
        self.schemas = None
    
    def load_spec(self):
        """Load OpenAPI spec"""
        with open(self.openapi_file, 'r', encoding='utf-8') as f:
            self.spec = yaml.safe_load(f)
        self.schemas = self.spec['components']['schemas']
    
    def validate_all(self):
        """Validate all schemas"""
        results = {
            'total': 0,
            'with_description': 0,
            'without_description': [],
            'short_description': [],  # < 20 chars
            'good_description': [],   # >= 50 chars
            'with_required': 0,
            'without_required': [],
            'by_category': defaultdict(lambda: {'total': 0, 'with_desc': 0})
        }
        
        for schema_name, schema_def in self.schemas.items():
            if not isinstance(schema_def, dict):
                continue
            
            results['total'] += 1
            
            # Categorize
            if schema_name.endswith('Request'):
                category = 'Request DTOs'
            elif schema_name.endswith('Response'):
                category = 'Response DTOs'
            elif schema_name.endswith('Event'):
                category = 'Events'
            elif schema_def.get('type') == 'string' and 'enum' in schema_def:
                category = 'Enums'
            else:
                category = 'Entities'
            
            results['by_category'][category]['total'] += 1
            
            # Check description
            desc = schema_def.get('description', '').strip()
            if desc:
                results['with_description'] += 1
                results['by_category'][category]['with_desc'] += 1
                
                if len(desc) < 20:
                    results['short_description'].append(schema_name)
                elif len(desc) >= 50:
                    results['good_description'].append(schema_name)
            else:
                results['without_description'].append(schema_name)
            
            # Check required fields
            if schema_def.get('type') == 'object' and 'properties' in schema_def:
                if 'required' in schema_def and schema_def['required']:
                    results['with_required'] += 1
                else:
                    results['without_required'].append(schema_name)
        
        return results
    
    def print_report(self, results):
        """Print validation report"""
        print("\n" + "="*60)
        print("DESCRIPTION VALIDATION REPORT")
        print("="*60)
        
        print(f"\n📊 Overall Coverage:")
        total = results['total']
        with_desc = results['with_description']
        coverage = (with_desc / total * 100) if total > 0 else 0
        print(f"  Total Schemas: {total}")
        print(f"  With Descriptions: {with_desc} ({coverage:.1f}%)")
        print(f"  Without Descriptions: {len(results['without_description'])}")
        
        print(f"\n📝 Description Quality:")
        print(f"  Good (≥50 chars): {len(results['good_description'])}")
        print(f"  Short (<20 chars): {len(results['short_description'])}")
        
        print(f"\n✅ Required Fields:")
        print(f"  With Required: {results['with_required']}")
        print(f"  Without Required: {len(results['without_required'])}")
        
        print(f"\n📂 By Category:")
        for category, stats in sorted(results['by_category'].items()):
            cat_total = stats['total']
            cat_with_desc = stats['with_desc']
            cat_coverage = (cat_with_desc / cat_total * 100) if cat_total > 0 else 0
            print(f"  {category:20s} {cat_with_desc:3d}/{cat_total:3d} ({cat_coverage:5.1f}%)")
        
        if results['without_description']:
            print(f"\n⚠ Schemas Without Descriptions ({len(results['without_description'])}):")
            for name in results['without_description'][:10]:
                print(f"  - {name}")
            if len(results['without_description']) > 10:
                print(f"  ... and {len(results['without_description']) - 10} more")
        
        if results['short_description']:
            print(f"\n⚠ Schemas With Short Descriptions ({len(results['short_description'])}):")
            for name in results['short_description'][:5]:
                desc = self.schemas[name].get('description', '')
                print(f"  - {name}: '{desc}'")
            if len(results['short_description']) > 5:
                print(f"  ... and {len(results['short_description']) - 5} more")
        
        print("\n" + "="*60)
        
        # Overall grade
        if coverage >= 95:
            print("✅ EXCELLENT: 95%+ coverage achieved!")
        elif coverage >= 80:
            print("✓ GOOD: 80%+ coverage achieved")
        elif coverage >= 50:
            print("⚠ FAIR: 50%+ coverage, needs improvement")
        else:
            print("❌ POOR: <50% coverage, significant work needed")
        
        print("="*60 + "\n")

if __name__ == '__main__':
    validator = DescriptionValidator('../openapi.yaml')
    
    print("Loading OpenAPI spec...")
    validator.load_spec()
    
    print(f"Validating {len(validator.schemas)} schemas...")
    results = validator.validate_all()
    
    validator.print_report(results)
