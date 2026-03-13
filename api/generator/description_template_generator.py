#!/usr/bin/env python3
"""
Description Template Generator

Generates markdown description templates for OpenAPI schemas.
Populates with proto comments and creates structured sections for manual completion.

Usage:
    python description_template_generator.py --descriptor input/descriptors.pb --output descriptions/
"""

import argparse
import os
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).parent))

from proto_parser import ProtoParser


class DescriptionTemplateGenerator:
    """Generates markdown description templates for schemas"""
    
    def __init__(self, output_dir: str):
        self.output_dir = Path(output_dir)
        self.output_dir.mkdir(parents=True, exist_ok=True)
    
    def generate_schema_description(self, message_data: dict, category: str = "schema") -> str:
        """
        Generate markdown description for a schema
        
        Args:
            message_data: Parsed message data from proto_parser
            category: Type of schema (dto, entity, event)
        
        Returns:
            Markdown formatted description
        """
        name = message_data['descriptor'].name
        full_name = message_data.get('full_name', name)
        proto_comment = message_data.get('comment', '').strip()
        fields = message_data.get('fields', [])
        
        # Determine purpose based on category and name
        if category == 'dto':
            if name.endswith('Request'):
                purpose = f"Request payload for {self._extract_operation(name)} operation"
            elif name.endswith('Response'):
                purpose = f"Response payload for {self._extract_operation(name)} operation"
            else:
                purpose = "Data transfer object"
        elif category == 'event':
            purpose = f"Event emitted when {self._event_to_description(name)}"
        else:
            purpose = f"Domain entity representing {self._humanize_name(name)}"
        
        template = f"""# {name}

## Overview

**Type**: {category.capitalize()}  
**Proto**: `{full_name}`  
**Purpose**: {purpose}

{f'**Proto Comment**: {proto_comment}' if proto_comment else ''}

## Description

<!-- Add detailed description here -->
{self._generate_description_placeholder(name, category)}

## Fields

"""
        
        # Add field descriptions
        for field in fields:
            field_name = field.get('name', 'unknown')
            field_type = field.get('type', 'unknown')
            field_comment = field.get('comment', '').strip()
            is_repeated = field.get('repeated', False)
            is_required = self._is_required(field)
            
            requirement = "**Required**" if is_required else "*Optional*"
            cardinality = " (repeated)" if is_repeated else ""
            
            template += f"""### `{field_name}`

- **Type**: `{field_type}`{cardinality}
- **Requirement**: {requirement}
{f'- **Proto Comment**: {field_comment}' if field_comment else ''}

<!-- Add detailed field description here -->

"""
        
        # Add usage examples section
        template += f"""
## Usage Examples

<!-- Add usage examples here -->

"""
        
        if category == 'dto' and name.endswith('Request'):
            template += f"""### Example Request

```json
{{
  // Add example request payload
}}
```

"""
        
        if category == 'dto' and name.endswith('Response'):
            template += f"""### Example Response

```json
{{
  // Add example response payload
}}
```

"""
        
        # Add notes section
        template += f"""
## Notes

<!-- Add additional notes, constraints, or business rules here -->

## Related

<!-- Link to related schemas, endpoints, or documentation -->

---

**Generated**: Auto-generated template  
**Last Updated**: <!-- Update date when modified -->
"""
        
        return template
    
    def _extract_operation(self, name: str) -> str:
        """Extract operation name from DTO name"""
        # Remove Request/Response suffix
        operation = name.replace('Request', '').replace('Response', '')
        
        # Convert PascalCase to words
        import re
        words = re.findall(r'[A-Z][a-z]*', operation)
        return ' '.join(words).lower()
    
    def _event_to_description(self, name: str) -> str:
        """Convert event name to readable description"""
        # Remove Event suffix
        base = name.replace('Event', '')
        
        # Convert PascalCase to words
        import re
        words = re.findall(r'[A-Z][a-z]*', base)
        return ' '.join(words).lower()
    
    def _humanize_name(self, name: str) -> str:
        """Convert schema name to human-readable form"""
        import re
        words = re.findall(r'[A-Z][a-z]*', name)
        return ' '.join(words).lower()
    
    def _generate_description_placeholder(self, name: str, category: str) -> str:
        """Generate context-specific description placeholder"""
        if category == 'dto' and 'Creation' in name:
            return """This schema defines the data required to create a new resource.
Ensure all required fields are provided and validated before submission.

**Validation Rules**:
- <!-- Add validation rules -->

**Business Logic**:
- <!-- Add business logic notes -->
"""
        elif category == 'event':
            return """This event is emitted when specific conditions are met in the system.
Consumers can subscribe to this event for asynchronous processing.

**Event Data**:
- <!-- Describe what data is included -->

**When Emitted**:
- <!-- Describe triggering conditions -->

**Consumers**:
- <!-- List typical consumers -->
"""
        else:
            return """<!-- Provide a comprehensive description of this schema -->

**Key Characteristics**:
- <!-- Add key points -->

**Use Cases**:
- <!-- Describe common use cases -->
"""
    
    def _is_required(self, field: dict) -> bool:
        """Check if field is required based on proto annotations"""
        # Check for field_behavior = REQUIRED
        options = field.get('options', {})
        if isinstance(options, dict):
            field_behavior = options.get('field_behavior', [])
            return 'REQUIRED' in field_behavior or 2 in field_behavior
        
        # Check comment
        comment = field.get('comment', '').lower()
        return 'required' in comment
    
    def generate_all_templates(self, proto_parser: ProtoParser, limit: int = None):
        """
        Generate description templates for all or top N schemas
        
        Args:
            proto_parser: Initialized proto parser
            limit: Maximum number of templates to generate (None = all)
        """
        messages = proto_parser.get_messages()
        
        # Categorize messages
        dtos = []
        entities = []
        events = []
        
        for msg in messages:
            name = msg['descriptor'].name
            full_name = msg.get('full_name', '')
            
            # Skip google.protobuf types
            if full_name.startswith('google.protobuf'):
                continue
            
            if name.endswith('Request') or name.endswith('Response'):
                dtos.append(msg)
            elif name.endswith('Event'):
                events.append(msg)
            else:
                entities.append(msg)
        
        print(f"Found {len(dtos)} DTOs, {len(entities)} entities, {len(events)} events")
        print()
        
        # Sort by importance (DTOs first, then events, then entities)
        priority_list = dtos[:limit//2 if limit else None] + events[:limit//4 if limit else None] + entities[:limit//4 if limit else None]
        
        if limit:
            priority_list = priority_list[:limit]
        
        print(f"Generating templates for {len(priority_list)} schemas...")
        print()
        
        generated_count = 0
        for msg in priority_list:
            # Determine category
            name = msg['descriptor'].name
            if name.endswith('Request') or name.endswith('Response'):
                category = 'dto'
            elif name.endswith('Event'):
                category = 'event'
            else:
                category = 'entity'
            
            # Generate template
            template = self.generate_schema_description(msg, category)
            
            # Determine output path
            full_name = msg.get('full_name', name)
            package_parts = full_name.split('.')[:-1]  # Remove message name
            
            output_path = self.output_dir / category / '/'.join(package_parts) / f"{name}.md"
            output_path.parent.mkdir(parents=True, exist_ok=True)
            
            # Write template
            with open(output_path, 'w', encoding='utf-8') as f:
                f.write(template)
            
            generated_count += 1
            if generated_count % 10 == 0:
                print(f"  Generated {generated_count} templates...")
        
        print()
        print(f"✅ Generated {generated_count} description templates in {self.output_dir}")


def main():
    parser = argparse.ArgumentParser(description='Generate description templates for OpenAPI schemas')
    parser.add_argument('--descriptor', required=True, help='Path to proto descriptor file')
    parser.add_argument('--output', default='../descriptions', help='Output directory for descriptions')
    parser.add_argument('--limit', type=int, default=50, help='Maximum number of templates to generate')
    parser.add_argument('--all', action='store_true', help='Generate for all schemas (ignore limit)')
    
    args = parser.parse_args()
    
    # Parse proto descriptors
    print("Parsing proto descriptors...")
    proto_parser = ProtoParser()
    proto_parser.load_descriptor_set(args.descriptor)
    print(f"✅ Loaded {len(proto_parser.get_messages())} messages")
    print()
    
    # Generate templates
    generator = DescriptionTemplateGenerator(args.output)
    limit = None if args.all else args.limit
    generator.generate_all_templates(proto_parser, limit)
    
    print()
    print("=" * 80)
    print("TEMPLATE GENERATION COMPLETE")
    print("=" * 80)
    print()
    print("Next steps:")
    print(f"  1. Review generated templates in {args.output}")
    print("  2. Fill in detailed descriptions and examples")
    print("  3. Run regeneration to load descriptions")
    print()


if __name__ == '__main__':
    main()
