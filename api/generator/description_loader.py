"""
Description Loader
Loads rich descriptions from markdown files for OpenAPI schemas and operations
Implements 3-tier priority: MD file > proto comments > auto-generated
"""

import os
import re
from typing import Optional, Dict
from pathlib import Path


class DescriptionLoader:
    """Loads descriptions from markdown files and proto comments"""
    
    def __init__(self, descriptions_dir: str, fallback_to_generated: bool = True):
        """
        Initialize description loader
        
        Args:
            descriptions_dir: Path to api/descriptions directory
            fallback_to_generated: Generate description if not found
        """
        self.descriptions_dir = Path(descriptions_dir)
        self.fallback_to_generated = fallback_to_generated
        self._cache = {}
    
    def load_schema_description(self, full_name: str, proto_comment: str = "") -> str:
        """
        Load description for a schema
        Priority: MD file > proto comment > generated
        
        Args:
            full_name: Full proto name (e.g., insuretech.policy.entity.v1.Policy)
            proto_comment: Comment from proto file
            
        Returns:
            Description text
        """
        # Check cache
        cache_key = f"schema:{full_name}"
        if cache_key in self._cache:
            return self._cache[cache_key]
        
        # Try loading from MD file
        md_path = self._get_schema_md_path(full_name)
        if md_path.exists():
            desc = self._load_md_file(md_path)
            self._cache[cache_key] = desc
            return desc
        
        # Fallback to proto comment
        if proto_comment and proto_comment.strip():
            desc = proto_comment.strip()
            self._cache[cache_key] = desc
            return desc
        
        # Generate description
        if self.fallback_to_generated:
            desc = self._generate_schema_description(full_name)
            self._cache[cache_key] = desc
            return desc
        
        return ""
    
    def load_operation_description(
        self, 
        service_name: str, 
        operation_name: str,
        proto_comment: str = ""
    ) -> Dict[str, str]:
        """
        Load operation description (summary + description)
        
        Args:
            service_name: Service name (e.g., PolicyService)
            operation_name: Operation name (e.g., CreatePolicy)
            proto_comment: Proto comment
            
        Returns:
            Dict with 'summary' and 'description' keys
        """
        cache_key = f"operation:{service_name}.{operation_name}"
        if cache_key in self._cache:
            return self._cache[cache_key]
        
        # Try loading from MD file
        md_path = self._get_operation_md_path(service_name, operation_name)
        if md_path.exists():
            content = self._load_md_file(md_path)
            result = self._parse_operation_md(content)
            self._cache[cache_key] = result
            return result
        
        # Fallback to proto comment
        if proto_comment and proto_comment.strip():
            result = {
                'summary': self._extract_summary(proto_comment),
                'description': proto_comment.strip()
            }
            self._cache[cache_key] = result
            return result
        
        # Generate
        if self.fallback_to_generated:
            result = self._generate_operation_description(operation_name)
            self._cache[cache_key] = result
            return result
        
        return {'summary': '', 'description': ''}
    
    def load_field_description(
        self,
        schema_name: str,
        field_name: str,
        proto_comment: str = ""
    ) -> str:
        """Load description for a specific field"""
        # For now, prioritize proto comment, can extend to MD files later
        if proto_comment and proto_comment.strip():
            return proto_comment.strip()
        
        if self.fallback_to_generated:
            return self._generate_field_description(field_name)
        
        return ""
    
    def _get_schema_md_path(self, full_name: str) -> Path:
        """Get path to schema markdown file"""
        # insuretech.policy.entity.v1.Policy → schemas/policy/entity/v1/Policy.md
        parts = full_name.split('.')
        
        # Skip 'insuretech' prefix if present
        if parts[0] == 'insuretech':
            parts = parts[1:]
        
        filename = parts[-1] + '.md'
        subpath = Path(*parts[:-1]) / filename
        
        return self.descriptions_dir / 'schemas' / subpath
    
    def _get_operation_md_path(self, service_name: str, operation_name: str) -> Path:
        """Get path to operation markdown file"""
        # PolicyService, CreatePolicy → endpoints/PolicyService/CreatePolicy.md
        return self.descriptions_dir / 'endpoints' / service_name / f"{operation_name}.md"
    
    def _load_md_file(self, path: Path) -> str:
        """Load content from markdown file"""
        try:
            with open(path, 'r', encoding='utf-8') as f:
                return f.read()
        except Exception as e:
            print(f"Warning: Could not load {path}: {e}")
            return ""
    
    def _parse_operation_md(self, content: str) -> Dict[str, str]:
        """Parse operation markdown into summary and description"""
        lines = content.strip().split('\n')
        
        summary = ""
        description = content
        
        # Look for first heading as summary
        for line in lines:
            if line.startswith('# '):
                summary = line[2:].strip()
                break
        
        # If no heading, use first sentence as summary
        if not summary:
            first_line = lines[0] if lines else ""
            summary = first_line.split('.')[0].strip()
        
        return {
            'summary': summary,
            'description': description.strip()
        }
    
    def _extract_summary(self, text: str) -> str:
        """Extract one-line summary from text"""
        lines = [l.strip() for l in text.split('\n') if l.strip()]
        if lines:
            # Take first sentence
            first = lines[0]
            sentences = re.split(r'[.!?]', first)
            return sentences[0].strip()
        return ""
    
    def _generate_schema_description(self, full_name: str) -> str:
        """Generate a default description for schema"""
        parts = full_name.split('.')
        name = parts[-1]
        
        # Convert PascalCase to readable text
        readable = re.sub(r'([A-Z])', r' \1', name).strip().lower()
        
        return f"Represents a {readable} entity"
    
    def _generate_operation_description(self, operation_name: str) -> Dict[str, str]:
        """Generate default operation description"""
        # CreatePolicy → Create a policy
        readable = re.sub(r'([A-Z])', r' \1', operation_name).strip().lower()
        
        return {
            'summary': f"{readable.capitalize()}",
            'description': f"Performs {readable} operation"
        }
    
    def _generate_field_description(self, field_name: str) -> str:
        """Generate default field description"""
        # policy_id → Policy identifier
        readable = field_name.replace('_', ' ')
        return f"{readable.capitalize()}"
    
    def create_template(self, full_name: str, message_data: dict) -> str:
        """
        Create a markdown template for a schema
        
        Args:
            full_name: Full proto name
            message_data: Message data from proto parser
            
        Returns:
            Markdown template content
        """
        parts = full_name.split('.')
        name = parts[-1]
        
        template = f"""# {name}

## Overview
{message_data.get('comment', 'TODO: Add overview description')}

## Purpose
TODO: Describe the purpose and role of this entity in the system

## Lifecycle
TODO: Document the lifecycle states if applicable

## Key Fields
TODO: Document important fields and their business meaning

## Business Rules
TODO: Document business rules and constraints

## Related Entities
TODO: List related entities and their relationships

## Examples
```json
{{
  // TODO: Add example JSON
}}
```

## Compliance & Regulatory Notes
TODO: Add any compliance requirements (IDRA, Bangladesh Insurance Act, etc.)

## Changelog
- Initial version
"""
        return template


def main():
    """CLI for testing description loader"""
    import argparse
    
    parser = argparse.ArgumentParser(
        description='Test description loader or generate templates'
    )
    parser.add_argument(
        '--descriptions-dir',
        default='../descriptions',
        help='Path to descriptions directory'
    )
    parser.add_argument(
        '--test',
        help='Test loading for a full name (e.g., insuretech.policy.entity.v1.Policy)'
    )
    parser.add_argument(
        '--generate-template',
        help='Generate template for a schema'
    )
    
    args = parser.parse_args()
    
    loader = DescriptionLoader(args.descriptions_dir)
    
    if args.test:
        desc = loader.load_schema_description(args.test)
        print(f"\nDescription for {args.test}:")
        print("=" * 60)
        print(desc)
        print("=" * 60)
    
    if args.generate_template:
        template = loader.create_template(args.generate_template, {'comment': ''})
        print(f"\nTemplate for {args.generate_template}:")
        print("=" * 60)
        print(template)


if __name__ == '__main__':
    main()
