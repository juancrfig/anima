#!/usr/bin/env python3
"""
CSV to Markdown Converter for Anima Entries
Reads a CSV file and creates individual .md files in ~/.anima/entries/
"""

import csv
import os
import sys
from pathlib import Path


def main():
    # Check if CSV file path is provided
    if len(sys.argv) != 2:
        print("Usage: python script.py <path_to_csv_file>")
        sys.exit(1)
    
    csv_file = sys.argv[1]
    
    # Verify CSV file exists
    if not os.path.isfile(csv_file):
        print(f"Error: File '{csv_file}' not found.")
        sys.exit(1)
    
    # Define target directory
    entries_dir = Path.home() / '.anima' / 'entries'
    
    # Create directory if it doesn't exist (won't affect existing files)
    entries_dir.mkdir(parents=True, exist_ok=True)
    
    # Process CSV file
    files_created = 0
    files_updated = 0
    
    try:
        with open(csv_file, 'r', encoding='utf-8-sig') as f:  # utf-8-sig handles BOM
            # Read first to check delimiter
            sample = f.read(1024)
            f.seek(0)
            
            # Try to detect delimiter
            sniffer = csv.Sniffer()
            try:
                dialect = sniffer.sniff(sample)
                delimiter = dialect.delimiter
            except:
                delimiter = ','
            
            reader = csv.DictReader(f, delimiter=delimiter)
            
            # Get actual fieldnames and create a mapping to clean names
            original_fields = reader.fieldnames
            clean_fields = {field: field.strip() for field in original_fields}
            
            # Debug: Show what columns were found
            print(f"Found columns: {list(clean_fields.values())}")
            
            # Find the actual column names (with potential spaces)
            full_date_col = None
            note_col = None
            
            for orig, clean in clean_fields.items():
                if clean == 'full_date':
                    full_date_col = orig
                elif clean == 'note':
                    note_col = orig
            
            # Verify required columns exist
            if full_date_col is None or note_col is None:
                print(f"\nError: CSV must contain 'full_date' and 'note' columns.")
                print(f"Found: {list(clean_fields.values())}")
                sys.exit(1)
            
            for row in reader:
                # Access columns using the original field names
                full_date = row.get(full_date_col, '').strip()
                note = row.get(note_col, '')
                
                # Skip rows with empty dates
                if not full_date:
                    continue
                
                # Create filename
                filename = f"{full_date}.md"
                filepath = entries_dir / filename
                
                # Check if file already exists
                file_existed = filepath.exists()
                
                # Write the note content to the file
                with open(filepath, 'w', encoding='utf-8') as md_file:
                    md_file.write(note)
                
                if file_existed:
                    files_updated += 1
                else:
                    files_created += 1
        
        # Summary
        print(f"✓ Processing complete!")
        print(f"  Files created: {files_created}")
        print(f"  Files updated: {files_updated}")
        print(f"  Location: {entries_dir}")
        
    except Exception as e:
        print(f"Error processing CSV: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
