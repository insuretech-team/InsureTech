"""
File Organizer for OpenAPI Generation
Determines correct output location for schemas, DTOs, and components
Implements proper separation as per apiplan.md
"""

import os
from typing import Tuple
from pathlib import Path


class FileOrganizer:
    """
    Organizes generated OpenAPI schema files into proper folder structure:
    - api/dtos/ - Request/Response DTOs
    - api/schemas/ - Entity schemas (Policy, Claim, etc.)
    - api/events/ - Event schemas (PolicyCreatedEvent, etc.) - SEPARATE from entities
    - api/enums/ - Enum schemas (FLAT structure, no nesting)
    - api/components/ - Common reusable components (Error, Money, etc.)
    """
    
    def __init__(self, output_dir: str):
        """
        Initialize file organizer
        
        Args:
            output_dir: Base output directory (api/)
        """
        self.output_dir = Path(output_dir)
        self.dtos_dir = self.output_dir / 'dtos'
        self.schemas_dir = self.output_dir / 'schemas'
        self.events_dir = self.output_dir / 'events'  # Separate folder for events
        self.enums_dir = self.output_dir / 'enums'    # Flat folder for enums
        self.components_dir = self.output_dir / 'components'
        
        # Common component types that should go in components/
        self.common_types = [
            'Error',
            'Money',
            'Address',
            'ContactInfo',
            'NIDInfo',
            'TINInfo',
            'PaginationRequest',
            'PaginationResponse',
            'Response',
            'AuditInfo',
            'ApprovalInfo',
            'Document',
        ]
    
    def is_dto(self, message_name: str) -> bool:
        """
        Check if message is a DTO (Request/Response)
        
        Args:
            message_name: Simple message name (e.g., CreatePolicyRequest)
            
        Returns:
            True if DTO, False otherwise
        """
        return message_name.endswith('Request') or message_name.endswith('Response')
    
    def is_event(self, message_name: str) -> bool:
        """
        Check if message is an Event
        
        Args:
            message_name: Simple message name
            
        Returns:
            True if Event, False otherwise
        """
        return message_name.endswith('Event')
    
    def is_common_component(self, full_name: str, message_name: str) -> bool:
        """
        Check if message is a common reusable component
        
        Args:
            full_name: Full proto name (insuretech.common.v1.Error)
            message_name: Simple message name (Error)
            
        Returns:
            True if common component, False otherwise
        """
        # Check if in common.v1 package
        if 'common.v1' in full_name:
            # Check if it's one of the common types
            return message_name in self.common_types
        
        return False
    
    def get_output_path(self, full_name: str, message_name: str) -> Tuple[Path, str]:
        """
        Determine correct output path for a message
        
        Args:
            full_name: Full proto name (e.g., insuretech.policy.services.v1.CreatePolicyRequest)
            message_name: Simple message name (e.g., CreatePolicyRequest)
            
        Returns:
            Tuple of (full_path, category) where category is 'dto', 'entity', or 'component'
        """
        # Check category
        if self.is_common_component(full_name, message_name):
            return self._get_component_path(full_name, message_name), 'component'
        elif self.is_dto(message_name):
            return self._get_dto_path(full_name, message_name), 'dto'
        elif self.is_event(message_name):
            return self._get_event_path(full_name, message_name), 'event'
        else:
            return self._get_entity_path(full_name, message_name), 'entity'
    
    def _get_dto_path(self, full_name: str, message_name: str) -> Path:
        """
        Get path for DTO (Request/Response)
        
        Format: api/dtos/{package}/{message}.yaml
        Example: api/dtos/insuretech/policy/services/v1/CreatePolicyRequest.yaml
        
        Args:
            full_name: Full proto name
            message_name: Simple message name
            
        Returns:
            Full file path
        """
        package_parts = full_name.split('.')[:-1]  # Remove message name
        rel_path = Path(*package_parts) / f"{message_name}.yaml"
        return self.dtos_dir / rel_path
    
    def _get_entity_path(self, full_name: str, message_name: str) -> Path:
        """
        Get path for entity schema
        
        Format: api/schemas/{package}/{message}.yaml
        Example: api/schemas/insuretech/policy/entity/v1/Policy.yaml
        
        Args:
            full_name: Full proto name
            message_name: Simple message name
            
        Returns:
            Full file path
        """
        package_parts = full_name.split('.')[:-1]
        rel_path = Path(*package_parts) / f"{message_name}.yaml"
        return self.schemas_dir / rel_path
    
    def _get_component_path(self, full_name: str, message_name: str) -> Path:
        """
        Get path for common component
        
        Format: api/components/schemas/{message}.yaml
        Example: api/components/schemas/Error.yaml
        
        Args:
            full_name: Full proto name
            message_name: Simple message name
            
        Returns:
            Full file path
        """
        return self.components_dir / 'schemas' / f"{message_name}.yaml"
    
    def _get_event_path(self, full_name: str, message_name: str) -> Path:
        """
        Get path for event schema (SEPARATE from entities)
        
        Format: api/events/{package}/{message}.yaml
        Example: api/events/insuretech/policy/events/v1/PolicyCreatedEvent.yaml
        
        Args:
            full_name: Full proto name
            message_name: Simple message name
            
        Returns:
            Full file path
        """
        package_parts = full_name.split('.')[:-1]
        rel_path = Path(*package_parts) / f"{message_name}.yaml"
        return self.events_dir / rel_path
    
    def ensure_directories(self):
        """Create all necessary directories"""
        # Create base directories
        self.dtos_dir.mkdir(parents=True, exist_ok=True)
        self.schemas_dir.mkdir(parents=True, exist_ok=True)
        self.events_dir.mkdir(parents=True, exist_ok=True)  # Separate events folder
        self.enums_dir.mkdir(parents=True, exist_ok=True)   # Flat enums folder
        self.components_dir.mkdir(parents=True, exist_ok=True)
        
        # Create component subdirectories
        (self.components_dir / 'schemas').mkdir(exist_ok=True)
        (self.components_dir / 'parameters').mkdir(exist_ok=True)
        (self.components_dir / 'responses').mkdir(exist_ok=True)
    
    def get_statistics(self, messages: list) -> dict:
        """
        Get statistics about file organization
        
        Args:
            messages: List of message data from proto parser
            
        Returns:
            Dictionary with counts
        """
        stats = {
            'dtos': 0,
            'entities': 0,
            'events': 0,
            'components': 0,
            'total': len(messages)
        }
        
        for msg in messages:
            full_name = msg.get('full_name', '')
            message_name = msg.get('descriptor').name if msg.get('descriptor') else ''
            
            if self.is_common_component(full_name, message_name):
                stats['components'] += 1
            elif self.is_dto(message_name):
                stats['dtos'] += 1
            elif self.is_event(message_name):
                stats['events'] += 1
            else:
                stats['entities'] += 1
        
        return stats


def main():
    """CLI for testing file organizer"""
    import argparse
    
    parser = argparse.ArgumentParser(
        description='Test file organizer routing'
    )
    parser.add_argument(
        '--output-dir',
        default='../test-output',
        help='Output directory'
    )
    parser.add_argument(
        'names',
        nargs='+',
        help='Full proto names to test (e.g., insuretech.policy.services.v1.CreatePolicyRequest)'
    )
    
    args = parser.parse_args()
    
    organizer = FileOrganizer(args.output_dir)
    
    print("\n" + "=" * 80)
    print("File Organization Test")
    print("=" * 80)
    
    for full_name in args.names:
        message_name = full_name.split('.')[-1]
        path, category = organizer.get_output_path(full_name, message_name)
        
        print(f"\n{full_name}")
        print(f"  Message: {message_name}")
        print(f"  Category: {category}")
        print(f"  Path: {path}")
        
        # Check characteristics
        if organizer.is_dto(message_name):
            print(f"  ✓ Is DTO")
        if organizer.is_event(message_name):
            print(f"  ✓ Is Event")
        if organizer.is_common_component(full_name, message_name):
            print(f"  ✓ Is Common Component")
    
    print("\n" + "=" * 80)


if __name__ == '__main__':
    main()
