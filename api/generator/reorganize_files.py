"""
Reorganize OpenAPI Schema Files
Moves 458 DTOs from api/schemas/ to api/dtos/
Moves common components to api/components/
Implements proper folder structure as per apiplan.md
"""

import os
import shutil
import yaml
from pathlib import Path
from file_organizer import FileOrganizer


class FileReorganizer:
    """Reorganizes existing generated files into proper structure"""
    
    def __init__(self, api_dir: str, dry_run: bool = True):
        """
        Initialize reorganizer
        
        Args:
            api_dir: Path to api/ directory
            dry_run: If True, only print what would be done
        """
        self.api_dir = Path(api_dir)
        self.schemas_dir = self.api_dir / 'schemas'
        self.dtos_dir = self.api_dir / 'dtos'
        self.components_dir = self.api_dir / 'components'
        self.dry_run = dry_run
        
        self.organizer = FileOrganizer(str(self.api_dir))
        
        self.stats = {
            'dtos_moved': 0,
            'components_moved': 0,
            'entities_kept': 0,
            'events_kept': 0,
            'errors': 0
        }
    
    def scan_files(self):
        """Scan all YAML files in schemas/ directory"""
        if not self.schemas_dir.exists():
            print(f"❌ Error: {self.schemas_dir} does not exist")
            return []
        
        files = []
        for yaml_file in self.schemas_dir.rglob('*.yaml'):
            # Skip google api files
            if 'google' in str(yaml_file):
                continue
            
            files.append(yaml_file)
        
        return files
    
    def get_message_info(self, yaml_file: Path) -> dict:
        """
        Extract message information from YAML file
        
        Args:
            yaml_file: Path to YAML file
            
        Returns:
            Dict with message_name and full_name
        """
        try:
            with open(yaml_file, 'r', encoding='utf-8') as f:
                data = yaml.safe_load(f)
            
            if not data:
                return None
            
            # YAML files have format: {MessageName: {...}}
            message_name = list(data.keys())[0]
            
            # Reconstruct full name from file path
            # schemas/insuretech/policy/services/v1/CreatePolicyRequest.yaml
            # -> insuretech.policy.services.v1.CreatePolicyRequest
            rel_path = yaml_file.relative_to(self.schemas_dir)
            parts = list(rel_path.parts[:-1])  # Remove filename
            parts.append(message_name)
            full_name = '.'.join(parts)
            
            return {
                'message_name': message_name,
                'full_name': full_name,
                'current_path': yaml_file
            }
        except Exception as e:
            print(f"⚠️  Warning: Could not parse {yaml_file}: {e}")
            return None
    
    def reorganize_file(self, file_info: dict):
        """
        Move a file to its correct location
        
        Args:
            file_info: Dict with message info
        """
        message_name = file_info['message_name']
        full_name = file_info['full_name']
        current_path = file_info['current_path']
        
        # Determine target location
        target_path, category = self.organizer.get_output_path(full_name, message_name)
        
        # Check if file needs to move
        if current_path == target_path:
            # File is already in correct location
            if category == 'entity':
                self.stats['entities_kept'] += 1
            elif category == 'event':
                self.stats['events_kept'] += 1
            return
        
        # Move file
        if self.dry_run:
            print(f"\n[DRY RUN] Would move {category}:")
            print(f"  From: {current_path.relative_to(self.api_dir)}")
            print(f"  To:   {target_path.relative_to(self.api_dir)}")
        else:
            # Create target directory
            target_path.parent.mkdir(parents=True, exist_ok=True)
            
            # Move file
            try:
                shutil.move(str(current_path), str(target_path))
                print(f"✓ Moved {category}: {message_name}")
            except Exception as e:
                print(f"❌ Error moving {message_name}: {e}")
                self.stats['errors'] += 1
                return
        
        # Update stats
        if category == 'dto':
            self.stats['dtos_moved'] += 1
        elif category == 'component':
            self.stats['components_moved'] += 1
    
    def cleanup_empty_dirs(self):
        """Remove empty directories after reorganization"""
        if self.dry_run:
            return
        
        # Walk bottom-up to remove empty dirs
        for dirpath, dirnames, filenames in os.walk(self.schemas_dir, topdown=False):
            # Skip if has files
            if filenames:
                continue
            
            # Skip if has subdirectories with files
            if dirnames:
                continue
            
            # Remove empty directory
            dir_path = Path(dirpath)
            if dir_path != self.schemas_dir:  # Don't remove schemas/ itself
                try:
                    dir_path.rmdir()
                    print(f"🗑️  Removed empty directory: {dir_path.relative_to(self.api_dir)}")
                except:
                    pass
    
    def run(self):
        """Execute reorganization"""
        print("=" * 80)
        if self.dry_run:
            print("DRY RUN - No files will be moved")
        else:
            print("REORGANIZING FILES - Files will be moved")
        print("=" * 80)
        
        # Ensure target directories exist
        if not self.dry_run:
            self.organizer.ensure_directories()
        
        # Scan all files
        print(f"\nScanning {self.schemas_dir}...")
        files = self.scan_files()
        print(f"Found {len(files)} YAML files")
        
        # Process each file
        print("\nProcessing files...")
        for yaml_file in files:
            file_info = self.get_message_info(yaml_file)
            if file_info:
                self.reorganize_file(file_info)
        
        # Cleanup empty directories
        if not self.dry_run:
            print("\nCleaning up empty directories...")
            self.cleanup_empty_dirs()
        
        # Print summary
        self.print_summary()
    
    def print_summary(self):
        """Print reorganization summary"""
        print("\n" + "=" * 80)
        print("REORGANIZATION SUMMARY")
        print("=" * 80)
        
        total_moved = self.stats['dtos_moved'] + self.stats['components_moved']
        total_kept = self.stats['entities_kept'] + self.stats['events_kept']
        
        print(f"\n📁 Files Moved: {total_moved}")
        print(f"  • DTOs moved to api/dtos/: {self.stats['dtos_moved']}")
        print(f"  • Components moved to api/components/: {self.stats['components_moved']}")
        
        print(f"\n📁 Files Kept in api/schemas/: {total_kept}")
        print(f"  • Entities: {self.stats['entities_kept']}")
        print(f"  • Events: {self.stats['events_kept']}")
        
        if self.stats['errors'] > 0:
            print(f"\n❌ Errors: {self.stats['errors']}")
        
        print(f"\n📊 Final Structure:")
        print(f"  • api/dtos/: {self.stats['dtos_moved']} files")
        print(f"  • api/components/schemas/: {self.stats['components_moved']} files")
        print(f"  • api/schemas/: {total_kept} files (entities + events)")
        
        if self.dry_run:
            print("\n⚠️  This was a DRY RUN - no files were actually moved")
            print("   Run with --execute to perform the reorganization")
        else:
            print("\n✅ Reorganization complete!")
        
        print("=" * 80)


def main():
    """CLI entry point"""
    import argparse
    
    parser = argparse.ArgumentParser(
        description='Reorganize OpenAPI schema files into proper structure'
    )
    parser.add_argument(
        '--api-dir',
        default='C:\\_DEV\\GO\\InsureTech\\api',
        help='Path to api/ directory'
    )
    parser.add_argument(
        '--execute',
        action='store_true',
        help='Actually move files (default is dry-run)'
    )
    
    args = parser.parse_args()
    
    # Determine if dry run
    dry_run = not args.execute
    
    # Create reorganizer
    reorganizer = FileReorganizer(args.api_dir, dry_run=dry_run)
    
    # Run reorganization
    reorganizer.run()


if __name__ == '__main__':
    main()
