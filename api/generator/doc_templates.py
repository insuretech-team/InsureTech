#!/usr/bin/env python3
"""
HTML Templates for Enhanced Documentation
Modern, responsive UI with tabs and cards
"""

def get_main_template() -> str:
    """Get the main documentation HTML template"""
    # Use raw string and double braces for CSS, single braces for Python format
    return """<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>InsureTech API Documentation</title>
    <style>
        * {{ margin: 0; padding: 0; box-sizing: border-box; }}
        
        body {{
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }}
        
        .container {{
            max-width: 1400px;
            margin: 0 auto;
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }}
        
        .header {{
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px;
            text-align: center;
        }}
        
        .header h1 {{
            font-size: 2.5em;
            margin-bottom: 10px;
        }}
        
        .header p {{
            font-size: 1.2em;
            opacity: 0.9;
        }}
        
        .tabs {{
            display: flex;
            background: #f8f9fa;
            border-bottom: 2px solid #e0e0e0;
            overflow-x: auto;
        }}
        
        .tab {{
            padding: 20px 30px;
            cursor: pointer;
            border: none;
            background: none;
            font-size: 1em;
            font-weight: 500;
            color: #666;
            transition: all 0.3s;
            white-space: nowrap;
            border-bottom: 3px solid transparent;
        }}
        
        .tab:hover {{
            background: rgba(102, 126, 234, 0.1);
            color: #667eea;
        }}
        
        .tab.active {{
            color: #667eea;
            border-bottom-color: #667eea;
            background: white;
        }}
        
        .tab-content {{
            display: none;
            padding: 40px;
            animation: fadeIn 0.3s;
        }}
        
        .tab-content.active {{
            display: block;
        }}
        
        @keyframes fadeIn {{
            from {{ opacity: 0; transform: translateY(10px); }}
            to {{ opacity: 1; transform: translateY(0); }}
        }}
        
        .search-box {{
            margin-bottom: 30px;
            position: relative;
        }}
        
        .search-box input {{
            width: 100%;
            padding: 15px 50px 15px 20px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            font-size: 1em;
            transition: border-color 0.3s;
        }}
        
        .search-box input:focus {{
            outline: none;
            border-color: #667eea;
        }}
        
        .search-icon {{
            position: absolute;
            right: 20px;
            top: 50%;
            transform: translateY(-50%);
            font-size: 1.2em;
            color: #999;
        }}
        
        .domains-grid {{
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
            gap: 20px;
            margin-top: 20px;
        }}
        
        .domain-card {{
            background: white;
            border: 2px solid #e0e0e0;
            border-radius: 12px;
            padding: 25px;
            cursor: pointer;
            transition: all 0.3s;
            position: relative;
        }}
        
        .domain-card:hover {{
            border-color: #667eea;
            box-shadow: 0 5px 20px rgba(102, 126, 234, 0.2);
            transform: translateY(-3px);
        }}
        
        .domain-card.hidden {{
            display: none;
        }}
        
        .domain-icon {{
            font-size: 2.5em;
            margin-bottom: 10px;
        }}
        
        .domain-name {{
            font-size: 1.3em;
            font-weight: 600;
            color: #333;
            margin-bottom: 8px;
        }}
        
        .domain-description {{
            color: #666;
            line-height: 1.5;
            margin-bottom: 12px;
        }}
        
        .domain-count {{
            display: inline-block;
            background: #667eea;
            color: white;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 0.85em;
            font-weight: 500;
        }}
        
        .detail-view {{
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0,0,0,0.5);
            z-index: 1000;
            overflow-y: auto;
        }}
        
        .detail-view.active {{
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }}
        
        .detail-content {{
            background: white;
            border-radius: 20px;
            max-width: 900px;
            width: 100%;
            max-height: 90vh;
            overflow-y: auto;
            position: relative;
        }}
        
        .detail-header {{
            position: sticky;
            top: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            border-radius: 20px 20px 0 0;
            z-index: 10;
        }}
        
        .detail-header h2 {{
            font-size: 2em;
            margin-bottom: 10px;
        }}
        
        .close-btn {{
            position: absolute;
            top: 20px;
            right: 20px;
            background: rgba(255,255,255,0.2);
            border: none;
            color: white;
            width: 40px;
            height: 40px;
            border-radius: 50%;
            cursor: pointer;
            font-size: 1.5em;
            transition: background 0.3s;
        }}
        
        .close-btn:hover {{
            background: rgba(255,255,255,0.3);
        }}
        
        .detail-body {{
            padding: 30px;
        }}
        
        .item-list {{
            display: flex;
            flex-direction: column;
            gap: 15px;
        }}
        
        .item {{
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 20px;
            border-radius: 8px;
            transition: all 0.3s;
        }}
        
        .item:hover {{
            background: #e9ecef;
            transform: translateX(5px);
        }}
        
        .item-name {{
            font-size: 1.2em;
            font-weight: 600;
            color: #333;
            margin-bottom: 5px;
            font-family: 'Courier New', monospace;
        }}
        
        .item-method {{
            display: inline-block;
            padding: 4px 10px;
            border-radius: 4px;
            font-size: 0.85em;
            font-weight: 600;
            margin-right: 10px;
        }}
        
        .method-GET {{ background: #4caf50; color: white; }}
        .method-POST {{ background: #2196f3; color: white; }}
        .method-PUT {{ background: #ff9800; color: white; }}
        .method-DELETE {{ background: #f44336; color: white; }}
        .method-PATCH {{ background: #9c27b0; color: white; }}
        
        .item-description {{
            color: #666;
            line-height: 1.6;
            margin-top: 8px;
        }}
        
        .item-meta {{
            margin-top: 10px;
            font-size: 0.9em;
            color: #999;
        }}
        
        .empty-state {{
            text-align: center;
            padding: 60px 20px;
            color: #999;
        }}
        
        .empty-state-icon {{
            font-size: 4em;
            margin-bottom: 20px;
        }}
        
        .stats-bar {{
            display: flex;
            gap: 30px;
            padding: 20px 40px;
            background: #f8f9fa;
            border-bottom: 1px solid #e0e0e0;
        }}
        
        .stat {{
            text-align: center;
        }}
        
        .stat-value {{
            font-size: 1.8em;
            font-weight: bold;
            color: #667eea;
        }}
        
        .stat-label {{
            font-size: 0.9em;
            color: #666;
            margin-top: 5px;
        }}
        
        .badge {{
            display: inline-block;
            background: #4caf50;
            color: white;
            padding: 3px 8px;
            border-radius: 12px;
            font-size: 0.75em;
            font-weight: 600;
            margin-left: 8px;
        }}
        
        .type-badge {{
            background: #2196f3;
        }}
        
        .enum-badge {{
            background: #ff9800;
        }}
        
        @media (max-width: 768px) {{
            .domains-grid {{
                grid-template-columns: 1fr;
            }}
            
            .tabs {{
                flex-wrap: wrap;
            }}
            
            .stats-bar {{
                flex-direction: column;
                gap: 15px;
            }}
        }}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🏥 InsureTech API</h1>
            <p>Comprehensive Insurance Platform API Documentation</p>
        </div>
        
        <div class="stats-bar">
            <div class="stat">
                <div class="stat-value">{total_schema_groups}</div>
                <div class="stat-label">Schema Groups</div>
            </div>
            <div class="stat">
                <div class="stat-value">{total_tables}</div>
                <div class="stat-label">Tables</div>
            </div>
            <div class="stat">
                <div class="stat-value">{total_apis}</div>
                <div class="stat-label">API Endpoints</div>
            </div>
            <div class="stat">
                <div class="stat-value">{total_enums}</div>
                <div class="stat-label">Enums</div>
            </div>
            <div class="stat">
                <div class="stat-value">{total_dtos}</div>
                <div class="stat-label">DTOs</div>
            </div>
            <div class="stat">
                <div class="stat-value">{total_events}</div>
                <div class="stat-label">Events</div>
            </div>
        </div>
        
        <div class="tabs">
            <button class="tab active" onclick="switchTab('schema-groups')">🗄️ Database Schemas</button>
            <button class="tab" onclick="switchTab('apis')">📡 API Endpoints</button>
            <button class="tab" onclick="switchTab('schemas')">📦 Schemas</button>
            <button class="tab" onclick="switchTab('enums')">🔢 Enums</button>
            <button class="tab" onclick="switchTab('dtos')">📋 DTOs</button>
            <button class="tab" onclick="switchTab('visualizer')">🔍 Visualizer</button>
            <button class="tab" onclick="switchTab('tools')">🛠️ Tools</button>
        </div>
        
        <div id="schema-groups" class="tab-content active">
            <div class="search-box">
                <input type="text" placeholder="Search database schema groups..." onkeyup="searchItems('schema-groups', this.value)">
                <span class="search-icon">🔍</span>
            </div>
            <div class="domains-grid" id="schema-groups-grid">
                {schema_groups_content}
            </div>
        </div>
        
        <div id="apis" class="tab-content">
            <div class="search-box">
                <input type="text" placeholder="Search API endpoints..." onkeyup="searchItems('apis', this.value)">
                <span class="search-icon">🔍</span>
            </div>
            <div class="domains-grid" id="apis-grid">
                {apis_content}
            </div>
        </div>
        
        <div id="schemas" class="tab-content">
            <div class="search-box">
                <input type="text" placeholder="Search schemas..." onkeyup="searchItems('schemas', this.value)">
                <span class="search-icon">🔍</span>
            </div>
            <div class="domains-grid" id="schemas-grid">
                {schemas_content}
            </div>
        </div>
        
        <div id="enums" class="tab-content">
            <div class="search-box">
                <input type="text" placeholder="Search enums..." onkeyup="searchItems('enums', this.value)">
                <span class="search-icon">🔍</span>
            </div>
            <div class="item-list" id="enums-list" style="display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 20px;">
                {enums_content}
            </div>
        </div>
        
        <div id="dtos" class="tab-content">
            <div class="search-box">
                <input type="text" placeholder="Search DTOs..." onkeyup="searchItems('dtos', this.value)">
                <span class="search-icon">🔍</span>
            </div>
            <div class="domains-grid" id="dtos-grid">
                {dtos_content}
            </div>
        </div>
        
        <div id="tools" class="tab-content">
            <h2 style="margin-bottom: 20px;">Interactive Tools</h2>
            <div class="domains-grid">
                <div class="domain-card" onclick="switchTab('visualizer')">
                    <div class="domain-icon">�</div>
                    <div class="domain-name">Schema Visualizer</div>
                    <div class="domain-description">Interactive schema browser with Mermaid diagrams and export</div>
                    <span class="domain-count">Interactive</span>
                </div>
                <div class="domain-card" onclick="window.open('swagger.html', '_blank')">
                    <div class="domain-icon">�</div>
                    <div class="domain-name">Swagger UI</div>
                    <div class="domain-description">Interactive API explorer with try-it-out functionality</div>
                    <span class="domain-count">Interactive</span>
                </div>
                <div class="domain-card" onclick="window.open('redoc.html', '_blank')">
                    <div class="domain-icon">�</div>
                    <div class="domain-name">ReDoc</div>
                    <div class="domain-description">Clean, responsive API reference documentation</div>
                    <span class="domain-count">Reference</span>
                </div>
                <div class="domain-card" onclick="window.open('/openapi.yaml', '_blank')">
                    <div class="domain-icon">📄</div>
                    <div class="domain-name">OpenAPI Spec</div>
                    <div class="domain-description">Download the raw OpenAPI 3.1 specification</div>
                    <span class="domain-count">Download</span>
                </div>
                <div class="domain-card" onclick="window.open('/validation_report.html', '_blank')">
                    <div class="domain-icon">✅</div>
                    <div class="domain-name">Validation Report</div>
                    <div class="domain-description">Detailed validation results and quality metrics</div>
                    <span class="domain-count">Quality</span>
                </div>
            </div>
        </div>
        
        <div id="visualizer" class="tab-content">
            <div style="padding: 20px; background: #f8f9fa; border-radius: 10px; margin-bottom: 20px;">
                <h2 style="margin-bottom: 10px;">🔍 Schema Visualizer</h2>
                <p style="color: #666;">Browse schemas, view relationships, and export diagrams</p>
            </div>
            
            <div id="schema-visualizer-container" style="display: flex; gap: 20px;">
                <!-- Sidebar with category tabs and schema list -->
                <div id="visualizer-sidebar" style="flex: 0 0 300px; background: white; border: 2px solid #e0e0e0; border-radius: 12px; padding: 15px; max-height: 800px; overflow-y: auto;">
                    <div id="visualizer-tabs" style="display: flex; flex-direction: column; gap: 5px; margin-bottom: 15px;">
                        <button class="visualizer-tab active" onclick="visualizer.switchCategory('entities')" style="padding: 10px; border: none; background: #667eea; color: white; border-radius: 8px; cursor: pointer; text-align: left; font-weight: 600;">
                            🗄️ Entities (DB Tables)
                        </button>
                        <button class="visualizer-tab" onclick="visualizer.switchCategory('dtos')" style="padding: 10px; border: none; background: #f0f0f0; color: #333; border-radius: 8px; cursor: pointer; text-align: left; font-weight: 600;">
                            📋 DTOs
                        </button>
                        <button class="visualizer-tab" onclick="visualizer.switchCategory('events')" style="padding: 10px; border: none; background: #f0f0f0; color: #333; border-radius: 8px; cursor: pointer; text-align: left; font-weight: 600;">
                            ⚡ Events
                        </button>
                        <button class="visualizer-tab" onclick="visualizer.switchCategory('enums')" style="padding: 10px; border: none; background: #f0f0f0; color: #333; border-radius: 8px; cursor: pointer; text-align: left; font-weight: 600;">
                            🔢 Enums
                        </button>
                    </div>
                    <div id="schema-list" style="font-size: 0.9em;">
                        <p style="text-align: center; color: #999;">Loading schemas...</p>
                    </div>
                </div>
                
                <!-- Main content area -->
                <div style="flex: 1; display: flex; flex-direction: column; gap: 20px;">
                    <!-- Quick selector dropdown -->
                    <div style="display: flex; gap: 20px;">
                        <div style="flex: 1;">
                            <label style="display: block; margin-bottom: 8px; font-weight: 600;">Quick Select:</label>
                            <select id="schema-selector" style="width: 100%; padding: 10px; border: 2px solid #e0e0e0; border-radius: 8px; font-size: 1em;">
                                <option value="">Select a schema...</option>
                            </select>
                        </div>
                        <div style="flex: 0 0 auto; display: flex; gap: 10px; align-items: flex-end;">
                            <button onclick="exportDiagram('png')" style="padding: 10px 20px; background: #667eea; color: white; border: none; border-radius: 8px; cursor: pointer; font-weight: 600;">Export PNG</button>
                            <button onclick="exportDiagram('svg')" style="padding: 10px 20px; background: #764ba2; color: white; border: none; border-radius: 8px; cursor: pointer; font-weight: 600;">Export SVG</button>
                            <button onclick="exportDiagram('json')" style="padding: 10px 20px; background: #4caf50; color: white; border: none; border-radius: 8px; cursor: pointer; font-weight: 600;">Export JSON</button>
                        </div>
                    </div>
                    
                    <!-- Schema details -->
                    <div id="schema-details" style="background: white; border: 2px solid #e0e0e0; border-radius: 12px; padding: 20px; display: none;">
                        <h3 id="schema-name" style="margin-bottom: 10px;"></h3>
                        <p id="schema-description" style="color: #666; margin-bottom: 15px;"></p>
                        <div id="schema-properties"></div>
                    </div>
                    
                    <!-- Mermaid diagram -->
                    <div id="visualizer-content" style="background: white; border: 2px solid #e0e0e0; border-radius: 12px; padding: 20px; overflow: auto; min-height: 400px;">
                        <p style="text-align: center; color: #999; padding: 40px;">Select a schema from the sidebar to view its diagram</p>
                    </div>
                </div>
            </div>
            
            <script src="https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js"></script>
            <script src="openapi-loader.js"></script>
            <script src="schema-visualizer.js"></script>
        </div>
    </div>
    
    <div id="detail-view" class="detail-view">
        <div class="detail-content">
            <div class="detail-header">
                <button class="close-btn" onclick="closeDetail()">×</button>
                <h2 id="detail-title"></h2>
                <p id="detail-subtitle"></p>
            </div>
            <div class="detail-body" id="detail-body"></div>
        </div>
    </div>
    
    <script>
        const detailData = {detail_data_json};
        
        function switchTab(tabName) {{
            document.querySelectorAll('.tab-content').forEach(tab => {{
                tab.classList.remove('active');
            }});
            document.querySelectorAll('.tab').forEach(tab => {{
                tab.classList.remove('active');
            }});
            
            document.getElementById(tabName).classList.add('active');
            event.target.classList.add('active');
        }}
        
        function showSchemaDetail(schemaName) {{
            const schemaData = detailData.schema_groups && detailData.schema_groups[schemaName];
            if (!schemaData) return;
            
            const icon = schemaData.icon || '📦';
            const description = schemaData.description || 'Database schema group';
            
            document.getElementById('detail-title').innerHTML = icon + ' ' + schemaName.replace('_', ' ').replace(/\\b\\w/g, l => l.toUpperCase());
            document.getElementById('detail-subtitle').textContent = description;
            
            let html = '<div class="item-list">';
            
            if (schemaData.tables && schemaData.tables.length > 0) {{
                schemaData.tables.forEach(table => {{
                    const tableName = table.table_name || table;
                    const messageName = table.message_name || '';
                    const migrationOrder = table.migration_order || 'N/A';
                    
                    html += `
                        <div class="item">
                            <div class="item-name">${{tableName}}</div>
                            <div class="item-description">
                                Message: ${{messageName}}<br>
                                Migration Order: ${{migrationOrder}}
                            </div>
                        </div>
                    `;
                }});
            }} else {{
                html += `
                    <div class="empty-state">
                        <div class="empty-state-icon">📭</div>
                        <p>No tables found</p>
                    </div>
                `;
            }}
            
            html += '</div>';
            
            document.getElementById('detail-body').innerHTML = html;
            document.getElementById('detail-view').classList.add('active');
        }}
        
        function showDetail(domain, type) {{
            const info = detailData[type][domain];
            if (!info || info.items.length === 0) return;
            
            const domainInfo = {domain_info_json}[domain] || {{ name: domain, icon: '📦', description: '' }};
            
            document.getElementById('detail-title').innerHTML = domainInfo.icon + ' ' + domainInfo.name;
            document.getElementById('detail-subtitle').textContent = domainInfo.description;
            
            let html = '<div class="item-list">';
            info.items.forEach(item => {{
                if (type === 'apis') {{
                    // Generate endpoint ID for the page link (ensure .html extension)
                    // Also replace colons for custom actions like /v1/auth/otp:send
                    const cleanPath = item.path.replace('/v1/', '').replace(/\\//g, '_').replace(/[{{}}]/g, '').replace(/-/g, '_').replace(/:/g, '_');
                    const endpointId = `${{cleanPath}}_${{item.method.toLowerCase()}}`;
                    const endpointPage = `endpoint_${{endpointId}}.html`;
                    
                    html += `
                        <div class="item" onclick="window.location.href='${{endpointPage}}'" style="cursor: pointer;">
                            <div>
                                <span class="item-method method-${{item.method}}">${{item.method}}</span>
                                <span class="item-name">${{item.path}}</span>
                            </div>
                            <div class="item-description">${{item.summary || item.description || 'No description'}}</div>
                            <div class="item-meta">Click to view full details →</div>
                        </div>
                    `;
                }} else if (type === 'schemas') {{
                    // Generate schema page link
                    const cleanName = item.name.replace(/\\./g, '_').replace(/:/g, '_').toLowerCase();
                    const schemaPage = `schema_${{cleanName}}.html`;
                    
                    const badge = item.name.endsWith('Request') ? '<span class="badge">Request</span>' : 
                                  item.name.endsWith('Response') ? '<span class="badge">Response</span>' : 
                                  `<span class="badge type-badge">${{item.type || 'object'}}</span>`;
                    
                    html += `
                        <div class="item" onclick="window.location.href='${{schemaPage}}'" style="cursor: pointer;">
                            <div class="item-name">${{item.name}} ${{badge}}</div>
                            <div class="item-description">${{item.description || 'No description'}}</div>
                            <div class="item-meta">Click to view full schema details →</div>
                        </div>
                    `;
                }} else {{
                    // DTOs - same as schemas
                    const cleanName = item.name.replace(/\\./g, '_').replace(/:/g, '_').toLowerCase();
                    const schemaPage = `schema_${{cleanName}}.html`;
                    
                    const badge = item.name.endsWith('Request') ? '<span class="badge">Request</span>' : 
                                  item.name.endsWith('Response') ? '<span class="badge">Response</span>' : 
                                  `<span class="badge type-badge">${{item.type || 'object'}}</span>`;
                    
                    html += `
                        <div class="item" onclick="window.location.href='${{schemaPage}}'" style="cursor: pointer;">
                            <div class="item-name">${{item.name}} ${{badge}}</div>
                            <div class="item-description">${{item.description || 'No description'}}</div>
                            <div class="item-meta">Click to view full DTO details →</div>
                        </div>
                    `;
                }}
            }});
            html += '</div>';
            
            if (info.items.length === 0) {{
                html = `
                    <div class="empty-state">
                        <div class="empty-state-icon">📭</div>
                        <p>No items found</p>
                    </div>
                `;
            }}
            
            document.getElementById('detail-body').innerHTML = html;
            document.getElementById('detail-view').classList.add('active');
        }}
        
        function closeDetail() {{
            document.getElementById('detail-view').classList.remove('active');
        }}
        
        function searchItems(type, query) {{
            query = query.toLowerCase();
            const gridId = type === 'enums' ? 'enums-list' : type + '-grid';
            const items = document.getElementById(gridId).children;
            
            for (let item of items) {{
                const text = item.textContent.toLowerCase();
                if (text.includes(query)) {{
                    item.classList.remove('hidden');
                }} else {{
                    item.classList.add('hidden');
                }}
            }}
        }}
        
        document.getElementById('detail-view').addEventListener('click', function(e) {{
            if (e.target === this) {{
                closeDetail();
            }}
        }});
        
        document.addEventListener('keydown', function(e) {{
            if (e.key === 'Escape') {{
                closeDetail();
            }}
        }});
    </script>
</body>
</html>"""
