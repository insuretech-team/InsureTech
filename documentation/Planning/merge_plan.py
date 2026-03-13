#!/usr/bin/env python3
import os
import re
from pathlib import Path
from datetime import datetime

try:
    import matplotlib.pyplot as plt
    import matplotlib.patches as mpatches
    MATPLOTLIB_AVAILABLE = True
except ImportError:
    MATPLOTLIB_AVAILABLE = False
    print("Warning: matplotlib not available. Will use ASCII chart fallback.")

def generate_capacity_chart_png(m1_avail, m1_req, m2_avail, m2_req, m3_avail, m3_req, output_dir):
    """Generate PNG chart using matplotlib"""
    if not MATPLOTLIB_AVAILABLE:
        return None
    
    # Calculate percentages
    m1_util = int((m1_req / m1_avail) * 100)
    m2_util = int((m2_req / m2_avail) * 100)
    m3_util = int((m3_req / m3_avail) * 100)
    
    # Calculate buffers
    m1_buffer = m1_avail - m1_req
    m2_buffer = m2_avail - m2_req
    m3_buffer = m3_avail - m3_req
    
    # Create figure with white background
    fig, ax = plt.subplots(figsize=(16, 9), facecolor='white')
    ax.set_facecolor('white')
    
    # Data
    phases = ['M1 (10 weeks)\nDec 20 - Mar 1', 'M2 (5 weeks)\nMar 2 - Apr 14', 'M3 (15 weeks)\nApr 15 - Aug 1']
    available = [m1_avail, m2_avail, m3_avail]
    required = [m1_req, m2_req, m3_req]
    utilization = [m1_util, m2_util, m3_util]
    
    # Bar positions - side by side
    x = [0, 2.5, 5]  # Phase positions
    width = 0.5      # Bar width
    
    # Colors - high contrast for visibility
    color_available = '#B0E57C'  # Light green for available
    color_required = '#5B9BD5'   # Blue for required
    
    # Create side-by-side bars
    x_available = [pos - width/2 - 0.05 for pos in x]  # Left position
    x_required = [pos + width/2 + 0.05 for pos in x]   # Right position
    
    bars_available = ax.bar(x_available, available, width, label='Available Capacity', 
                            color=color_available, edgecolor='black', linewidth=2)
    bars_required = ax.bar(x_required, required, width, label='Required Hours', 
                          color=color_required, edgecolor='black', linewidth=2)
    
    # Add value labels on bars
    for i, (avail, req, util) in enumerate(zip(available, required, utilization)):
        # Available hours label
        ax.text(x_available[i], avail/2, f'{avail:,}\nhrs\nAvailable', 
               ha='center', va='center', fontsize=12, fontweight='bold', color='#2F4F2F')
        # Required hours label  
        ax.text(x_required[i], req/2, f'{req:,}\nhrs\nRequired', 
               ha='center', va='center', fontsize=12, fontweight='bold', color='white')
        # Utilization percentage above bars
        max_height = max(avail, req)
        ax.text(x[i], max_height + 300, f'{util}% Utilized', 
               ha='center', va='bottom', fontsize=13, fontweight='bold', 
               bbox=dict(boxstyle='round,pad=0.5', facecolor='yellow', alpha=0.7))
    
    # Styling
    ax.set_ylabel('Capacity (Hours)', fontsize=16, fontweight='bold')
    ax.set_title('Project Capacity Utilization by Phase - LabAid InsureTech Platform', 
                fontsize=18, fontweight='bold', pad=20, color='#2F4F4F')
    ax.set_xticks(x)
    ax.set_xticklabels(phases, fontsize=13, fontweight='bold')
    ax.set_ylim(0, max(max(available), max(required)) * 1.2)
    ax.legend(loc='upper right', fontsize=14, framealpha=1, edgecolor='black', facecolor='white')
    ax.grid(axis='y', alpha=0.4, linestyle='--', linewidth=1)
    ax.spines['top'].set_visible(False)
    ax.spines['right'].set_visible(False)
    ax.spines['left'].set_linewidth(2)
    ax.spines['bottom'].set_linewidth(2)
    
    # Add summary text
    summary_text = f"""All Phases Comfortable: 68-73% utilization with 27-32% buffer
Total Available: {sum(available):,} hrs  |  Total Required: {sum(required):,} hrs  |  Overall: {int(sum(required)/sum(available)*100)}% utilized"""
    fig.text(0.5, 0.02, summary_text, ha='center', fontsize=11, style='italic', bbox=dict(boxstyle='round', facecolor='wheat', alpha=0.5))
    
    plt.tight_layout(rect=[0, 0.05, 1, 1])
    
    # Save PNG
    png_path = os.path.join(output_dir, 'capacity_utilization_chart.png')
    plt.savefig(png_path, dpi=150, bbox_inches='tight', facecolor='white')
    plt.close()
    
    print(f"  ✓ Generated capacity chart PNG: {png_path}")
    return 'capacity_utilization_chart.png'

def generate_sprint_timeline_png(output_dir):
    """Generate sprint timeline Gantt-style chart"""
    if not MATPLOTLIB_AVAILABLE:
        return None
    
    # Create figure
    fig, ax = plt.subplots(figsize=(16, 10), facecolor='white')
    ax.set_facecolor('white')
    
    # Sprint data
    sprints = [
        ('Sprint 1\nFoundation', 'Dec 20 - Jan 2', 0, 2, '#4472C4'),
        ('Sprint 2\nPolicy Service', 'Jan 3 - Jan 16', 2, 4, '#4472C4'),
        ('Sprint 3\nPayment Complete', 'Jan 17 - Jan 30', 4, 6, '#4472C4'),
        ('Sprint 4\nIntegration', 'Jan 31 - Feb 13', 6, 8, '#4472C4'),
        ('Sprint 5\nTesting', 'Feb 14 - Feb 27', 8, 10, '#4472C4'),
        ('Sprint 5.5\nLaunch', 'Feb 28 - Mar 1', 10, 10.4, '#ED7D31'),  # M1 Launch
        ('Sprint 6\nM2 Start', 'Mar 2 - Mar 15', 11, 13, '#70AD47'),
        ('Sprint 7\nClaims Service', 'Mar 16 - Mar 29', 13, 15, '#70AD47'),
        ('Sprint 8\nM2 Launch', 'Mar 30 - Apr 12', 15, 17, '#70AD47'),
        ('Buffer\nGrand Launch', 'Apr 13 - Apr 14', 17, 17.4, '#FFC000'),  # M2 Launch
        ('Sprints 9-15\nM3 Development', 'Apr 15 - Aug 1', 18, 33, '#5B9BD5'),
    ]
    
    # Milestones
    milestones = [
        ('M1: Beta Launch\nNational Insurance Day', 10.2, '#ED7D31'),
        ('M2: Grand Launch\nPohela Boishakh', 17.2, '#FFC000'),
        ('M3: Complete Platform\nIoT & AI Ready', 33, '#5B9BD5'),
    ]
    
    # Draw sprint bars
    for i, (name, dates, start, end, color) in enumerate(sprints):
        duration = end - start
        ax.barh(i, duration, left=start, height=0.6, color=color, 
                edgecolor='black', linewidth=1.5, alpha=0.8)
        # Sprint name
        ax.text(start + duration/2, i, name, ha='center', va='center', 
                fontsize=9, fontweight='bold', color='white')
        # Dates on right
        ax.text(end + 0.3, i, dates, ha='left', va='center', 
                fontsize=8, style='italic')
    
    # Draw milestone markers
    for label, pos, color in milestones:
        ax.axvline(x=pos, color=color, linewidth=3, linestyle='--', alpha=0.7)
        ax.plot(pos, len(sprints) + 0.5, 'v', markersize=15, color=color, 
                markeredgecolor='black', markeredgewidth=2)
        ax.text(pos, len(sprints) + 1.2, label, ha='center', va='bottom', 
                fontsize=11, fontweight='bold', color=color,
                bbox=dict(boxstyle='round,pad=0.5', facecolor='white', 
                         edgecolor=color, linewidth=2))
    
    # Styling
    ax.set_ylim(-1, len(sprints) + 2.5)
    ax.set_xlim(-1, 35)
    ax.set_xlabel('Project Timeline (Weeks)', fontsize=14, fontweight='bold')
    ax.set_title('Sprint Timeline & Milestone Schedule - LabAid InsureTech Platform', 
                fontsize=18, fontweight='bold', pad=20, color='#2F4F4F')
    ax.set_yticks(range(len(sprints)))
    ax.set_yticklabels([])
    ax.invert_yaxis()
    ax.grid(axis='x', alpha=0.3, linestyle=':')
    ax.spines['top'].set_visible(False)
    ax.spines['right'].set_visible(False)
    ax.spines['left'].set_visible(False)
    
    # Legend
    from matplotlib.patches import Patch
    legend_elements = [
        Patch(facecolor='#4472C4', edgecolor='black', label='M1 Sprints (10 weeks)'),
        Patch(facecolor='#70AD47', edgecolor='black', label='M2 Sprints (5 weeks)'),
        Patch(facecolor='#5B9BD5', edgecolor='black', label='M3 Sprints (15 weeks)'),
    ]
    ax.legend(handles=legend_elements, loc='lower right', fontsize=12, 
             framealpha=1, edgecolor='black')
    
    # Add phase labels at top
    ax.text(5, -0.8, 'M1 DEVELOPMENT\n(10 weeks)', ha='center', va='center', 
           fontsize=13, fontweight='bold', 
           bbox=dict(boxstyle='round,pad=0.7', facecolor='#4472C4', 
                    edgecolor='black', linewidth=2, alpha=0.3))
    ax.text(14, -0.8, 'M2 DEV\n(5 weeks)', ha='center', va='center', 
           fontsize=13, fontweight='bold',
           bbox=dict(boxstyle='round,pad=0.7', facecolor='#70AD47', 
                    edgecolor='black', linewidth=2, alpha=0.3))
    ax.text(25.5, -0.8, 'M3 DEVELOPMENT\n(16 weeks)', ha='center', va='center', 
           fontsize=13, fontweight='bold',
           bbox=dict(boxstyle='round,pad=0.7', facecolor='#5B9BD5', 
                    edgecolor='black', linewidth=2, alpha=0.3))
    
    plt.tight_layout()
    
    # Save PNG
    png_path = os.path.join(output_dir, 'sprint_timeline_chart.png')
    plt.savefig(png_path, dpi=150, bbox_inches='tight', facecolor='white')
    plt.close()
    
    print(f"  ✓ Generated sprint timeline PNG: {png_path}")
    return 'sprint_timeline_chart.png'

def generate_capacity_chart(m1_avail, m1_req, m2_avail, m2_req, m3_avail, m3_req):
    """Generate side-by-side column chart for capacity visualization"""
    
    # Calculate percentages
    m1_util = int((m1_req / m1_avail) * 100)
    m2_util = int((m2_req / m2_avail) * 100)
    m3_util = int((m3_req / m3_avail) * 100)
    
    # Calculate buffer
    m1_buffer = m1_avail - m1_req
    m2_buffer = m2_avail - m2_req
    m3_buffer = m3_avail - m3_req
    
    # Chart height in rows
    chart_height = 16
    bar_width = 12
    
    # Calculate bar heights (proportional to max capacity)
    max_capacity = max(m1_avail, m2_avail, m3_avail)
    
    def get_bar_heights(avail, req):
        avail_height = int((avail / max_capacity) * chart_height)
        req_height = int((req / max_capacity) * chart_height)
        return avail_height, req_height
    
    m1_avail_h, m1_req_h = get_bar_heights(m1_avail, m1_req)
    m2_avail_h, m2_req_h = get_bar_heights(m2_avail, m2_req)
    m3_avail_h, m3_req_h = get_bar_heights(m3_avail, m3_req)
    
    # Generate row by row from top to bottom
    lines = []
    lines.append("```")
    lines.append("Capacity (hours)")
    lines.append("")
    
    for row in range(chart_height, -1, -1):
        # Y-axis label every 2 rows
        if row % 2 == 0:
            label = f"{int((row / chart_height) * max_capacity):5.0f} │"
        else:
            label = "      │"
        
        # M1 bars (Available | Required)
        m1_avail_char = '░' if row <= m1_avail_h else ' '
        m1_req_char = '█' if row <= m1_req_h else ' '
        m1_bar = f"{m1_avail_char * bar_width} {m1_req_char * bar_width}"
        
        # M2 bars
        m2_avail_char = '░' if row <= m2_avail_h else ' '
        m2_req_char = '█' if row <= m2_req_h else ' '
        m2_bar = f"{m2_avail_char * bar_width} {m2_req_char * bar_width}"
        
        # M3 bars
        m3_avail_char = '░' if row <= m3_avail_h else ' '
        m3_req_char = '█' if row <= m3_req_h else ' '
        m3_bar = f"{m3_avail_char * bar_width} {m3_req_char * bar_width}"
        
        lines.append(f"{label}  {m1_bar}   {m2_bar}   {m3_bar}")
    
    # Bottom line
    lines.append("    0 └" + "─" * 100)
    lines.append(f"         Avail  Req       Avail  Req       Avail  Req")
    lines.append(f"           M1 (10 wks)        M2 (5 wks)        M3 (15 wks)")
    lines.append(f"         Dec 20 - Mar 1    Mar 2 - Apr 14   Apr 15 - Aug 1")
    lines.append("")
    lines.append(f"      M1: {m1_avail:,} hrs available ░  |  {m1_req:,} hrs required █  →  {m1_util}% utilized  ({m1_buffer:,} hrs buffer)")
    lines.append(f"      M2: {m2_avail:,} hrs available ░  |  {m2_req:,} hrs required █  →  {m2_util}% utilized  ({m2_buffer:,} hrs buffer)")
    lines.append(f"      M3: {m3_avail:,} hrs available ░  |  {m3_req:,} hrs required █  →  {m3_util}% utilized  ({m3_buffer:,} hrs buffer)")
    lines.append("")
    lines.append(f"      ░ = Available Capacity    █ = Required Hours    ALL PHASES COMFORTABLE!")
    lines.append("```")
    
    return '\n'.join(lines)

def merge_markdown_files(input_dir, output_file, file_order):
    print(f"Starting merge process...")
    merged_content = []
    
    # Don't add header, TOC file has it
    
    for filename in file_order:
        filepath = os.path.join(input_dir, filename)
        if not os.path.exists(filepath):
            print(f"Warning: File not found: {filepath}")
            continue
        
        print(f"Processing: {filename}")
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
            
            # If this is the TOC file, enhance the charts
            if filename == "00_TOC.md":
                # Generate capacity chart
                png_capacity = generate_capacity_chart_png(4040, 2740, 2256, 1641, 6768, 4863, input_dir)
                
                if png_capacity:
                    # Use PNG embedded in markdown with proper formatting for DOCX conversion
                    chart_content = f"""### Capacity Utilization Over Time

<img src="{png_capacity}" alt="Project Capacity Utilization" style="max-width:100%;height:auto;" />

**Chart Description:** Side-by-side comparison showing Available Capacity (green bars) vs Required Hours (blue bars) for each milestone phase. All phases show comfortable 68-73% utilization with adequate buffer."""
                else:
                    # Fallback to ASCII chart
                    chart_content = f'### Capacity Utilization Over Time\n\n{generate_capacity_chart(4040, 2740, 2256, 1641, 6768, 4863)}'
                
                # Replace the capacity chart section
                chart_pattern = r'### Capacity Utilization Over Time\n\n(```.*?```|!\[.*?\]\(.*?\).*|<img.*?/>.*?(?=\n\n###|\n\n\*\*|$))'
                if re.search(chart_pattern, content, re.DOTALL):
                    content = re.sub(
                        chart_pattern, 
                        chart_content,
                        content,
                        flags=re.DOTALL
                    )
                    print(f"  ✓ Embedded capacity chart in TOC")
                
                # Generate sprint timeline chart
                png_sprint = generate_sprint_timeline_png(input_dir)
                
                if png_sprint:
                    sprint_chart_content = f"""### Project Timeline

<img src="{png_sprint}" alt="Sprint Timeline" style="max-width:100%;height:auto;" />

**Timeline Overview:** Complete sprint schedule from December 2025 to August 2026, showing 15 sprints across three major milestones."""
                    
                    # Replace the ASCII timeline section
                    timeline_pattern = r'### Project Timeline\n\n```\n.*?```'
                    if re.search(timeline_pattern, content, re.DOTALL):
                        content = re.sub(
                            timeline_pattern,
                            sprint_chart_content,
                            content,
                            flags=re.DOTALL
                        )
                        print(f"  ✓ Embedded sprint timeline chart in TOC")
            
            merged_content.append(content)
            merged_content.append("\n\n---\n\n")
    
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(''.join(merged_content))
    print(f"\nSuccess! Merged document created: {output_file}")
    print(f"File size: {os.path.getsize(output_file)} bytes")
    return True

script_dir = Path(__file__).parent
file_order = [
    "00_TOC.md",
    "01_TeamAndTechnology.md",
    "02_WorkingDaysChart.md",
    "03_TeamCapacity.md",
    "04_EffortEstimation.md",
    "05_SprintPlanning.md",
    "06_RACIMatrix.md",
    "07_ResourceReassignment.md",
    "08_RisksAndMitigation.md",
    "09_RequirementsByMilestone.md",
    "10_PersonWiseResponsibility.md"
]

output_filename = script_dir / "DetailedProjectPlan.md"
merge_markdown_files(script_dir, output_filename, file_order)
