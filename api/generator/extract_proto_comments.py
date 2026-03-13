"""
Extract comments from .proto files and populate markdown description files
"""
import os
import re
from collections import defaultdict

class ProtoCommentExtractor:
    def __init__(self, proto_root, descriptions_root):
        self.proto_root = proto_root
        self.descriptions_root = descriptions_root
        self.extracted = defaultdict(dict)
    
    def extract_from_file(self, proto_file):
        """Extract comments from a single proto file"""
        with open(proto_file, 'r', encoding='utf-8') as f:
            content = f.read()
        
        results = {}
        
        # Extract message comments
        # Pattern: // Comment\nmessage MessageName {
        pattern = r'//\s*(.+?)\s*\n\s*message\s+(\w+)\s*\{'
        for match in re.finditer(pattern, content):
            comment = match.group(1).strip()
            message_name = match.group(2)
            if comment and len(comment) > 5:
                results[message_name] = comment
        
        # Also try multi-line comments
        pattern = r'/\*\s*(.+?)\s*\*/\s*\n\s*message\s+(\w+)\s*\{'
        for match in re.finditer(pattern, content, re.DOTALL):
            comment = match.group(1).strip()
            message_name = match.group(2)
            # Clean up comment
            comment = ' '.join(comment.split())
            if comment and len(comment) > 5:
                results[message_name] = comment
        
        # Extract enum comments
        pattern = r'//\s*(.+?)\s*\n\s*enum\s+(\w+)\s*\{'
        for match in re.finditer(pattern, content):
            comment = match.group(1).strip()
            enum_name = match.group(2)
            if comment and len(comment) > 5:
                results[enum_name] = comment
        
        return results
    
    def extract_all(self):
        """Extract comments from all proto files"""
        print("Extracting proto comments...")
        count = 0
        
        for root, dirs, files in os.walk(self.proto_root):
            for file in files:
                if file.endswith('.proto'):
                    proto_file = os.path.join(root, file)
                    comments = self.extract_from_file(proto_file)
                    
                    if comments:
                        rel_path = os.path.relpath(proto_file, self.proto_root)
                        self.extracted[rel_path] = comments
                        count += len(comments)
                        
                        for name, comment in comments.items():
                            print(f"  {name}: {comment[:60]}...")
        
        print(f"\nExtracted {count} comments from {len(self.extracted)} proto files")
        return count
    
    def update_markdown_file(self, md_file, schema_name, comment):
        """Update a markdown file with proto comment"""
        if not os.path.exists(md_file):
            return False
        
        with open(md_file, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Check if already has proto comment
        if '**Proto Comment**:' in content and not content.count('**Proto Comment**:') > 1:
            # Update existing
            pattern = r'\*\*Proto Comment\*\*:\s*.*'
            replacement = f'**Proto Comment**: {comment}'
            content = re.sub(pattern, replacement, content)
        else:
            # Add after Purpose line
            pattern = r'(\*\*Purpose\*\*:.*?\n)'
            replacement = f'\\1**Proto Comment**: {comment}\n'
            content = re.sub(pattern, replacement, content)
        
        with open(md_file, 'w', encoding='utf-8') as f:
            f.write(content)
        
        return True
    
    def populate_markdown_files(self):
        """Populate markdown files with extracted comments"""
        print("\nPopulating markdown files...")
        updated = 0
        
        for proto_path, comments in self.extracted.items():
            # Determine category from proto path
            # insuretech/ai/entity/v1/ai.proto -> entity
            # insuretech/ai/services/v1/ai_service.proto -> dto
            # insuretech/ai/events/v1/ai_events.proto -> event
            
            parts = proto_path.replace('\\', '/').split('/')
            
            if 'entity' in parts:
                category = 'entity'
            elif 'services' in parts:
                category = 'dto'
            elif 'events' in parts:
                category = 'event'
            else:
                continue
            
            # Build markdown path
            # Remove 'proto/' prefix if exists
            rel_parts = [p for p in parts if p != 'proto']
            
            # Build path: descriptions/{category}/insuretech/...
            md_dir = os.path.join(self.descriptions_root, category)
            for part in rel_parts[:-1]:  # Exclude filename
                md_dir = os.path.join(md_dir, part)
            
            # Update each schema's markdown
            for schema_name, comment in comments.items():
                md_file = os.path.join(md_dir, f"{schema_name}.md")
                if self.update_markdown_file(md_file, schema_name, comment):
                    updated += 1
                    print(f"  Updated: {md_file}")
        
        print(f"\nUpdated {updated} markdown files")
        return updated

if __name__ == '__main__':
    proto_root = '../../proto/insuretech'
    descriptions_root = '../descriptions'
    
    extractor = ProtoCommentExtractor(proto_root, descriptions_root)
    extractor.extract_all()
    extractor.populate_markdown_files()
    
    print("\n✓ Proto comment extraction complete!")
