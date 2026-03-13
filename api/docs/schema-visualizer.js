/**
 * Schema Visualizer for InsureTech API
 * Uses Mermaid.js to create interactive schema diagrams
 */

class SchemaVisualizer {
    constructor() {
        this.schemas = {};
        this.entities = {};
        this.entitiesBySchema = {}; // Group entities by schema (authn, policy, etc.)
        this.dtos = {};
        this.events = {};
        this.enums = {};
        this.currentSchema = null;
        this.currentCategory = 'entities';
    }

    /**
     * Initialize the visualizer by loading the OpenAPI spec
     * This is called by the main initialization function
     */
    async init() {
        // This method is kept for compatibility but initialization
        // is now handled by initializeVisualizer() in index.html
        console.log('Schema visualizer initialized');
    }

    /**
     * Simple YAML parser for OpenAPI spec
     */
    parseYAML(yamlText) {
        // For a production system, you'd use a proper YAML parser
        // For now, we'll load it as JSON if available, or parse basic YAML
        try {
            // Try to convert YAML to JSON (simplified approach)
            // In practice, the pipeline can generate a JSON version
            return this.basicYAMLParse(yamlText);
        } catch (error) {
            console.error('YAML parsing error:', error);
            return { components: { schemas: {} } };
        }
    }

    /**
     * Basic YAML to JS object parser (simplified for OpenAPI structure)
     */
    basicYAMLParse(yaml) {
        // This is a simplified parser - in production, use js-yaml library
        // For now, we'll work with a JSON version of the spec
        return { components: { schemas: {} } };
    }

    /**
     * Load OpenAPI spec from JSON (fallback method)
     */
    async loadFromJSON() {
        try {
            const response = await fetch('../openapi.json');
            if (response.ok) {
                const spec = await response.json();
                if (spec.components && spec.components.schemas) {
                    this.schemas = spec.components.schemas;
                    this.renderSchemaList();
                    return true;
                }
            }
        } catch (error) {
            console.error('JSON loading failed:', error);
        }
        return false;
    }

    /**
     * Render the list of available schemas with category tabs
     */
    renderSchemaList() {
        // Create category tabs
        const tabsHtml = `
            <div class="visualizer-tabs">
                <button class="visualizer-tab active" onclick="visualizer.switchCategory('all')">
                    All (${Object.keys(this.schemas).length})
                </button>
                <button class="visualizer-tab" onclick="visualizer.switchCategory('entities')">
                    🗄️ Entities (${Object.keys(this.entities).length})
                </button>
                <button class="visualizer-tab" onclick="visualizer.switchCategory('dtos')">
                    📋 DTOs (${Object.keys(this.dtos).length})
                </button>
                <button class="visualizer-tab" onclick="visualizer.switchCategory('events')">
                    ⚡ Events (${Object.keys(this.events).length})
                </button>
                <button class="visualizer-tab" onclick="visualizer.switchCategory('enums')">
                    🔢 Enums (${Object.keys(this.enums).length})
                </button>
            </div>
        `;
        
        const searchHtml = `
            <div class="visualizer-search">
                <input type="text" id="schema-search" placeholder="Search schemas..." 
                       onkeyup="visualizer.filterSchemas(this.value)">
                <span class="search-icon">🔍</span>
            </div>
        `;

        document.getElementById('visualizer-sidebar').innerHTML = tabsHtml + searchHtml + 
            '<div id="schema-list" class="schema-list"></div>';
        
        // Render the current category
        this.renderCategoryList(this.currentCategory);
    }
    
    /**
     * Switch between schema categories
     */
    switchCategory(category) {
        this.currentCategory = category;
        
        // Update active tab
        document.querySelectorAll('.visualizer-tab').forEach(tab => {
            tab.classList.remove('active');
        });
        event.target.classList.add('active');
        
        // Clear search
        document.getElementById('schema-search').value = '';
        
        // Render the category list
        this.renderCategoryList(category);
    }
    
    /**
     * Render schemas for a specific category
     */
    renderCategoryList(category) {
        let schemasToShow = {};
        
        switch(category) {
            case 'entities':
                schemasToShow = this.entities;
                break;
            case 'dtos':
                schemasToShow = this.dtos;
                break;
            case 'events':
                schemasToShow = this.events;
                break;
            case 'enums':
                schemasToShow = this.enums;
                break;
            default:
                schemasToShow = this.schemas;
        }
        
        const schemaNames = Object.keys(schemasToShow).sort();
        const listHtml = schemaNames.map(name => {
            const schema = schemasToShow[name];
            const badge = this.getSchemaBadge(name, schema);
            
            return `
                <div class="schema-item" onclick="visualizer.showSchema('${name}')">
                    <div class="schema-name">${name}</div>
                    <div class="schema-type">${badge}</div>
                </div>
            `;
        }).join('');

        document.getElementById('schema-list').innerHTML = listHtml;
    }
    
    /**
     * Get badge label for schema
     */
    getSchemaBadge(name, schema) {
        if (schema.type === 'string' && schema.enum) {
            return '🔢 enum';
        }
        if (name.endsWith('Event')) {
            return '⚡ event';
        }
        if (name.endsWith('Request')) {
            return '📤 request';
        }
        if (name.endsWith('Response')) {
            return '📥 response';
        }
        return '🗄️ entity';
    }

    /**
     * Filter schemas based on search input
     */
    filterSchemas(searchTerm) {
        const schemaItems = document.querySelectorAll('.schema-item');
        const term = searchTerm.toLowerCase();
        
        schemaItems.forEach(item => {
            const name = item.querySelector('.schema-name').textContent.toLowerCase();
            item.style.display = name.includes(term) ? 'block' : 'none';
        });
    }

    /**
     * Show a specific schema diagram
     */
    showSchema(schemaName) {
        this.currentSchema = schemaName;
        const schema = this.schemas[schemaName];
        
        if (!schema) {
            document.getElementById('visualizer-content').innerHTML = 
                '<div class="error">Schema not found</div>';
            return;
        }

        // Generate Mermaid diagram
        const mermaidCode = this.generateMermaidDiagram(schemaName, schema);
        
        // Create unique ID for this diagram
        const diagramId = `mermaid-${Date.now()}`;
        
        // Create the visualization HTML
        const html = `
            <div class="schema-header">
                <h2>${schemaName}</h2>
                <div class="schema-description">${schema.description || 'No description available'}</div>
            </div>
            <div class="diagram-container" id="diagram-container">
                <pre class="mermaid" id="${diagramId}">${mermaidCode}</pre>
            </div>
            <div class="schema-details">
                <h3>Properties</h3>
                <div class="properties-table">
                    ${this.generatePropertiesTable(schema)}
                </div>
            </div>
        `;

        document.getElementById('visualizer-content').innerHTML = html;
        
        // Render Mermaid diagram
        if (typeof mermaid !== 'undefined') {
            try {
                mermaid.run({
                    nodes: [document.getElementById(diagramId)]
                });
            } catch (error) {
                console.error('Mermaid rendering error:', error);
                document.getElementById(diagramId).innerHTML = 
                    `<div style="color: red; padding: 20px;">Error rendering diagram: ${error.message}</div>`;
            }
        }
    }

    /**
     * Generate Mermaid class diagram from schema
     */
    generateMermaidDiagram(schemaName, schema) {
        // Special handling for enums
        if (schema.type === 'string' && schema.enum) {
            let diagram = 'classDiagram\n';
            diagram += `    class ${schemaName} {\n`;
            diagram += `        <<enumeration>>\n`;
            
            // Add enum values
            for (const value of schema.enum) {
                diagram += `        ${value}\n`;
            }
            
            diagram += '    }\n';
            return diagram;
        }
        
        // Regular class diagram for objects
        let diagram = 'classDiagram\n';
        diagram += `    class ${schemaName} {\n`;

        if (schema.properties) {
            for (const [propName, propDef] of Object.entries(schema.properties)) {
                const type = this.getPropertyType(propDef);
                const required = schema.required && schema.required.includes(propName) ? '*' : '';
                diagram += `        ${required}${type} ${propName}\n`;
            }
        }

        diagram += '    }\n\n';

        // Add relationships to other schemas
        if (schema.properties) {
            for (const [propName, propDef] of Object.entries(schema.properties)) {
                const refSchema = this.getRefSchema(propDef);
                if (refSchema) {
                    diagram += `    ${schemaName} --> ${refSchema} : ${propName}\n`;
                } else if (propDef.type === 'array' && propDef.items) {
                    const arrayRefSchema = this.getRefSchema(propDef.items);
                    if (arrayRefSchema) {
                        diagram += `    ${schemaName} --> ${arrayRefSchema} : ${propName}[]\n`;
                    }
                }
            }
        }

        return diagram;
    }

    /**
     * Get property type for display
     */
    getPropertyType(propDef) {
        if (propDef.$ref) {
            return this.extractRefName(propDef.$ref);
        }
        if (propDef.type === 'array' && propDef.items) {
            if (propDef.items.$ref) {
                return this.extractRefName(propDef.items.$ref) + '[]';
            }
            return (propDef.items.type || 'any') + '[]';
        }
        return propDef.type || 'object';
    }

    /**
     * Get referenced schema name
     */
    getRefSchema(propDef) {
        if (propDef.$ref) {
            return this.extractRefName(propDef.$ref);
        }
        return null;
    }

    /**
     * Extract schema name from $ref
     */
    extractRefName(ref) {
        return ref.split('/').pop();
    }

    /**
     * Generate properties table
     */
    generatePropertiesTable(schema) {
        // Special handling for enums
        if (schema.type === 'string' && schema.enum) {
            let html = '<div style="margin-bottom: 15px;"><strong>Enum Values:</strong></div>';
            html += '<table><thead><tr><th>Value</th></tr></thead><tbody>';
            
            for (const value of schema.enum) {
                html += `<tr><td><code>${value}</code></td></tr>`;
            }
            
            html += '</tbody></table>';
            return html;
        }
        
        // Regular properties table for objects
        if (!schema.properties) {
            return '<p>No properties defined</p>';
        }

        let html = '<table><thead><tr><th>Property</th><th>Type</th><th>Required</th><th>Description</th></tr></thead><tbody>';
        
        for (const [propName, propDef] of Object.entries(schema.properties)) {
            const type = this.getPropertyType(propDef);
            const required = schema.required && schema.required.includes(propName) ? '✓' : '';
            const description = propDef.description || '-';
            
            html += `<tr>
                <td><code>${propName}</code></td>
                <td><code>${type}</code></td>
                <td>${required}</td>
                <td>${description}</td>
            </tr>`;
        }
        
        html += '</tbody></table>';
        return html;
    }

    /**
     * Show JSON representation of schema
     */
    showJSON(schemaName) {
        const schema = this.schemas[schemaName];
        const jsonStr = JSON.stringify(schema, null, 2);
        
        const html = `
            <div class="json-viewer">
                <div class="json-header">
                    <h3>JSON Schema: ${schemaName}</h3>
                    <button onclick="visualizer.copyJSON('${schemaName}')" class="action-btn">📋 Copy</button>
                </div>
                <pre><code class="language-json">${this.escapeHtml(jsonStr)}</code></pre>
            </div>
        `;
        
        const modal = document.createElement('div');
        modal.className = 'modal';
        modal.innerHTML = `
            <div class="modal-content">
                <span class="modal-close" onclick="this.parentElement.parentElement.remove()">&times;</span>
                ${html}
            </div>
        `;
        document.body.appendChild(modal);
    }

    /**
     * Copy JSON to clipboard
     */
    copyJSON(schemaName) {
        const schema = this.schemas[schemaName];
        const jsonStr = JSON.stringify(schema, null, 2);
        
        navigator.clipboard.writeText(jsonStr).then(() => {
            // Show success message
            this.showToast('✅ JSON copied to clipboard!', 'success');
        }).catch(err => {
            console.error('Failed to copy:', err);
            // Fallback: show in textarea for manual copy
            this.showCopyFallback(jsonStr);
        });
    }

    /**
     * Export diagram as SVG
     */
    async exportSVG() {
        // Find the SVG element in the diagram container
        const svgElement = document.querySelector('#diagram-container svg');
        
        if (!svgElement) {
            this.showToast('❌ No diagram to export. Please visualize a schema first.', 'error');
            return;
        }

        try {
            const svgData = new XMLSerializer().serializeToString(svgElement);
            const svgBlob = new Blob([svgData], { type: 'image/svg+xml;charset=utf-8' });
            const svgUrl = URL.createObjectURL(svgBlob);
            
            const downloadLink = document.createElement('a');
            downloadLink.href = svgUrl;
            downloadLink.download = `${this.currentSchema}_diagram.svg`;
            document.body.appendChild(downloadLink);
            downloadLink.click();
            document.body.removeChild(downloadLink);
            URL.revokeObjectURL(svgUrl);
            
            this.showToast('✅ SVG exported successfully!', 'success');
        } catch (error) {
            console.error('SVG export failed:', error);
            this.showToast('❌ SVG export failed: ' + error.message, 'error');
        }
    }

    /**
     * Export as PNG
     */
    async exportPNG() {
        // Find the SVG element in the diagram container
        const svgElement = document.querySelector('#diagram-container svg');
        
        if (!svgElement) {
            this.showToast('❌ No diagram to export. Please visualize a schema first.', 'error');
            return;
        }

        try {
            const svgData = new XMLSerializer().serializeToString(svgElement);
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');
            const img = new Image();
            
            // Get SVG dimensions
            const svgRect = svgElement.getBoundingClientRect();
            canvas.width = svgRect.width * 2; // 2x for better quality
            canvas.height = svgRect.height * 2;
            
            return new Promise((resolve, reject) => {
                img.onload = () => {
                    ctx.fillStyle = 'white';
                    ctx.fillRect(0, 0, canvas.width, canvas.height);
                    ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
                    
                    canvas.toBlob((blob) => {
                        const url = URL.createObjectURL(blob);
                        const downloadLink = document.createElement('a');
                        downloadLink.href = url;
                        downloadLink.download = `${this.currentSchema}_diagram.png`;
                        document.body.appendChild(downloadLink);
                        downloadLink.click();
                        document.body.removeChild(downloadLink);
                        URL.revokeObjectURL(url);
                        
                        this.showToast('✅ PNG exported successfully!', 'success');
                        resolve();
                    }, 'image/png');
                };
                
                img.onerror = (error) => {
                    reject(new Error('Failed to load SVG image'));
                };
                
                const svgBlob = new Blob([svgData], { type: 'image/svg+xml;charset=utf-8' });
                const url = URL.createObjectURL(svgBlob);
                img.src = url;
            });
        } catch (error) {
            console.error('PNG export failed:', error);
            this.showToast('❌ PNG export failed: ' + error.message, 'error');
        }
    }

    /**
     * Show toast notification
     */
    showToast(message, type = 'info') {
        const toast = document.createElement('div');
        toast.className = `toast toast-${type}`;
        toast.textContent = message;
        document.body.appendChild(toast);
        
        setTimeout(() => {
            toast.classList.add('toast-show');
        }, 10);
        
        setTimeout(() => {
            toast.classList.remove('toast-show');
            setTimeout(() => {
                document.body.removeChild(toast);
            }, 300);
        }, 3000);
    }

    /**
     * Fallback copy method using textarea
     */
    showCopyFallback(text) {
        const textarea = document.createElement('textarea');
        textarea.value = text;
        textarea.style.position = 'fixed';
        textarea.style.opacity = '0';
        document.body.appendChild(textarea);
        textarea.select();
        
        try {
            document.execCommand('copy');
            this.showToast('✅ JSON copied to clipboard!', 'success');
        } catch (err) {
            this.showToast('❌ Failed to copy. Please copy manually.', 'error');
        }
        
        document.body.removeChild(textarea);
    }

    /**
     * Escape HTML for safe display
     */
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    /**
     * Load schemas from parsed OpenAPI data and categorize them
     */
    loadFromData(openApiData) {
        console.log('Loading schemas from OpenAPI data...', openApiData);
        
        if (!openApiData) {
            console.error('No OpenAPI data provided');
            // Only try to update UI elements if they exist (sidebar layout)
            const sidebar = document.getElementById('visualizer-sidebar');
            const content = document.getElementById('visualizer-content');
            if (sidebar) {
                sidebar.innerHTML = '<div class="error">No OpenAPI data available</div>';
            }
            if (content) {
                content.innerHTML = '<div class="error">Failed to load schemas. Check console for details.</div>';
            }
            return;
        }
        
        if (openApiData.components && openApiData.components.schemas) {
            this.schemas = openApiData.components.schemas;
            
            // Categorize schemas
            this.categorizeSchemas();
            
            const totalCount = Object.keys(this.schemas).length;
            console.log(`Loaded ${totalCount} schemas:`);
            console.log(`  • Entities (DB Tables): ${Object.keys(this.entities).length}`);
            console.log(`  • DTOs: ${Object.keys(this.dtos).length}`);
            console.log(`  • Events: ${Object.keys(this.events).length}`);
            console.log(`  • Enums: ${Object.keys(this.enums).length}`);
            
            // Only render schema list if sidebar exists (sidebar layout)
            if (document.getElementById('visualizer-sidebar')) {
                this.renderSchemaList();
            }
        } else {
            console.error('No schemas found in OpenAPI data');
            // Only try to update UI elements if they exist (sidebar layout)
            const sidebar = document.getElementById('visualizer-sidebar');
            const content = document.getElementById('visualizer-content');
            if (sidebar) {
                sidebar.innerHTML = '<div class="error">No schemas found in specification</div>';
            }
            if (content) {
                content.innerHTML = '<div class="error">The OpenAPI specification does not contain any schemas.</div>';
            }
        }
    }
    
    /**
     * Categorize schemas into entities, DTOs, events, and enums
     */
    categorizeSchemas() {
        for (const [name, schema] of Object.entries(this.schemas)) {
            // Check if it's an enum
            if (schema.type === 'string' && schema.enum) {
                this.enums[name] = schema;
            }
            // Check if it's an event
            else if (name.endsWith('Event')) {
                this.events[name] = schema;
            }
            // Check if it's a DTO (Request or Response)
            else if (name.endsWith('Request') || name.endsWith('Response')) {
                this.dtos[name] = schema;
            }
            // Otherwise it's an entity (DB table)
            else if (schema.type === 'object') {
                this.entities[name] = schema;
                
                // Group entities by schema (extract schema name from table name)
                // Example: "AuthnUser" -> "authn", "PolicyClaim" -> "policy"
                const schemaGroup = this.extractSchemaGroup(name);
                if (!this.entitiesBySchema[schemaGroup]) {
                    this.entitiesBySchema[schemaGroup] = [];
                }
                this.entitiesBySchema[schemaGroup].push(name);
            }
        }
        
        // Sort entities within each schema group
        for (const group in this.entitiesBySchema) {
            this.entitiesBySchema[group].sort();
        }
    }
    
    /**
     * Extract schema group from entity name
     * Example: "AuthnUser" -> "authn", "PolicyClaim" -> "policy"
     */
    extractSchemaGroup(entityName) {
        // Common schema prefixes in InsureTech
        const prefixes = [
            'Authn', 'Authz', 'Policy', 'Claim', 'Payment', 'Product',
            'Partner', 'Commission', 'Document', 'Notification', 'Workflow',
            'Task', 'Audit', 'Fraud', 'Kyc', 'Mfs', 'Iot', 'Ai', 'Analytics',
            'Beneficiary', 'Endorsement', 'Renewal', 'Refund', 'Underwriting'
        ];
        
        for (const prefix of prefixes) {
            if (entityName.startsWith(prefix)) {
                return prefix.toLowerCase();
            }
        }
        
        // Default: use first word before capital letter
        const match = entityName.match(/^([A-Z][a-z]+)/);
        return match ? match[1].toLowerCase() : 'other';
    }
    
    /**
     * Switch between category tabs
     */
    switchCategory(category) {
        this.currentCategory = category;
        
        // Update tab styles
        document.querySelectorAll('.visualizer-tab').forEach(tab => {
            tab.style.background = '#f0f0f0';
            tab.style.color = '#333';
        });
        event.target.style.background = '#667eea';
        event.target.style.color = 'white';
        
        // Render the schema list for this category
        this.renderSchemaList();
    }
    
    /**
     * Render schema list in sidebar based on current category
     */
    renderSchemaList() {
        const listContainer = document.getElementById('schema-list');
        if (!listContainer) return;
        
        let html = '';
        
        if (this.currentCategory === 'entities') {
            // Group by schema, then show tables
            const schemaGroups = Object.keys(this.entitiesBySchema).sort();
            
            for (const group of schemaGroups) {
                const tables = this.entitiesBySchema[group];
                html += `
                    <div style="margin-bottom: 15px;">
                        <div style="font-weight: 600; color: #667eea; margin-bottom: 5px; padding: 5px; background: #f0f0ff; border-radius: 4px;">
                            📁 ${group.toUpperCase()} (${tables.length})
                        </div>
                `;
                
                for (const tableName of tables) {
                    html += `
                        <div onclick="visualizer.showSchema('${tableName}')" 
                             style="padding: 8px 12px; margin: 2px 0; cursor: pointer; border-radius: 4px; transition: all 0.2s;"
                             onmouseover="this.style.background='#f0f0ff'"
                             onmouseout="this.style.background='transparent'">
                            ${tableName}
                        </div>
                    `;
                }
                
                html += `</div>`;
            }
        } else if (this.currentCategory === 'dtos') {
            const dtoNames = Object.keys(this.dtos).sort();
            for (const name of dtoNames) {
                html += `
                    <div onclick="visualizer.showSchema('${name}')" 
                         style="padding: 8px 12px; margin: 2px 0; cursor: pointer; border-radius: 4px; transition: all 0.2s;"
                         onmouseover="this.style.background='#f0f0ff'"
                         onmouseout="this.style.background='transparent'">
                        ${name}
                    </div>
                `;
            }
        } else if (this.currentCategory === 'events') {
            const eventNames = Object.keys(this.events).sort();
            for (const name of eventNames) {
                html += `
                    <div onclick="visualizer.showSchema('${name}')" 
                         style="padding: 8px 12px; margin: 2px 0; cursor: pointer; border-radius: 4px; transition: all 0.2s;"
                         onmouseover="this.style.background='#f0f0ff'"
                         onmouseout="this.style.background='transparent'">
                        ${name}
                    </div>
                `;
            }
        } else if (this.currentCategory === 'enums') {
            const enumNames = Object.keys(this.enums).sort();
            for (const name of enumNames) {
                html += `
                    <div onclick="visualizer.showSchema('${name}')" 
                         style="padding: 8px 12px; margin: 2px 0; cursor: pointer; border-radius: 4px; transition: all 0.2s;"
                         onmouseover="this.style.background='#f0f0ff'"
                         onmouseout="this.style.background='transparent'">
                        ${name}
                    </div>
                `;
            }
        }
        
        if (html === '') {
            html = '<p style="text-align: center; color: #999; padding: 20px;">No schemas in this category</p>';
        }
        
        listContainer.innerHTML = html;
    }

    /**
     * Show schema in dropdown layout (for index.html visualizer tab)
     * Also updates the schema details section
     */
    showSchemaInDropdown(schemaName) {
        // Use the main showSchema method which handles everything
        this.showSchema(schemaName);
        
        // Also update the dropdown selector to match
        const selector = document.getElementById('schema-selector');
        if (selector) {
            selector.value = schemaName;
        }
    }
}

// Global instances
let visualizer;
let openApiLoader;

/**
 * Populate the schema selector dropdown with loaded schemas
 */
function populateSchemaSelector(schemas) {
    const selector = document.getElementById('schema-selector');
    if (!selector) {
        console.error('Schema selector element not found');
        return;
    }
    
    // Clear existing options
    selector.innerHTML = '<option value="">Select a schema...</option>';
    
    // Add schema options sorted alphabetically
    const schemaNames = Object.keys(schemas).sort();
    schemaNames.forEach(name => {
        const option = document.createElement('option');
        option.value = name;
        option.textContent = name;
        selector.appendChild(option);
    });
    
    // Add change event listener
    selector.addEventListener('change', (e) => {
        const schemaName = e.target.value;
        if (schemaName && visualizer) {
            visualizer.showSchemaInDropdown(schemaName);
        }
    });
    
    console.log(`✅ Populated selector with ${schemaNames.length} schemas`);
}

/**
 * Initialize the schema visualizer with OpenAPI data
 */
async function initializeSchemaVisualizer() {
    try {
        console.log('🚀 Initializing Schema Visualizer...');
        
        // Create instances
        visualizer = new SchemaVisualizer();
        openApiLoader = new OpenAPILoader();
        
        console.log('✅ Instances created');
        
        // Load OpenAPI spec
        const spec = await openApiLoader.loadSpec();
        
        if (!spec || !spec.components || !spec.components.schemas) {
            throw new Error('Failed to load OpenAPI spec or no schemas found');
        }
        
        console.log('✅ OpenAPI spec loaded');
        
        // Load data into visualizer
        visualizer.loadFromData(spec);
        
        console.log('✅ Data loaded into visualizer');
        
        // Populate the dropdown selector
        populateSchemaSelector(visualizer.schemas);
        
        console.log('✅ Schema Visualizer initialization complete!');
        
    } catch (error) {
        console.error('❌ Failed to initialize Schema Visualizer:', error);
        
        // Show error in UI
        const selector = document.getElementById('schema-selector');
        if (selector) {
            selector.innerHTML = '<option>Error loading schemas</option>';
        }
        
        const details = document.getElementById('schema-details');
        if (details) {
            details.style.display = 'block';
            details.innerHTML = `
                <div style="color: #d32f2f; padding: 20px; background: #ffebee; border-radius: 8px;">
                    <h3>❌ Failed to Load Schemas</h3>
                    <p>${error.message}</p>
                    <p style="font-size: 0.9em; color: #666;">Check the browser console for more details.</p>
                </div>
            `;
        }
    }
}

/**
 * Export diagram in specified format (called from HTML buttons)
 */
function exportDiagram(format) {
    if (!visualizer || !visualizer.currentSchema) {
        alert('Please select a schema first');
        return;
    }
    
    if (format === 'png') {
        visualizer.exportPNG();
    } else if (format === 'svg') {
        visualizer.exportSVG();
    } else if (format === 'json') {
        visualizer.showJSON(visualizer.currentSchema);
    }
}

// Initialize when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeSchemaVisualizer);
} else {
    initializeSchemaVisualizer();
}
