/**
 * OpenAPI Loader - Fetches and parses OpenAPI spec for the visualizer
 * This script loads the OpenAPI YAML and makes it available to the schema visualizer
 */

class OpenAPILoader {
    constructor() {
        this.spec = null;
    }

    /**
     * Load and parse OpenAPI specification from YAML
     */
    async loadSpec() {
        try {
            // First, try to load from JSON if available
            const jsonLoaded = await this.loadFromJSON();
            if (jsonLoaded) {
                return this.spec;
            }

            // Fallback: Load YAML and parse it
            const yamlPaths = [
                './openapi.yaml',
                '../openapi.yaml',
                '/openapi.yaml',
                'openapi.yaml'
            ];
            
            for (const path of yamlPaths) {
                try {
                    console.log(`Trying to load YAML from: ${path}`);
                    const response = await fetch(path);
                    if (response.ok) {
                        const yamlText = await response.text();
                        this.spec = await this.parseYAML(yamlText);
                        console.log('Successfully loaded OpenAPI spec from YAML:', path);
                        return this.spec;
                    }
                } catch (error) {
                    console.log(`Failed to load YAML from ${path}:`, error.message);
                }
            }
            
            throw new Error('Failed to load OpenAPI spec from any location');
        } catch (error) {
            console.error('Error loading OpenAPI spec:', error);
            throw error;
        }
    }

    /**
     * Try to load from JSON version
     */
    async loadFromJSON() {
        // Try multiple paths for different hosting scenarios
        const paths = [
            './openapi.json',           // Same directory (GitHub Pages)
            '../openapi.json',          // Parent directory (local server)
            '/openapi.json',            // Root (some servers)
            'openapi.json'              // Relative (fallback)
        ];
        
        for (const path of paths) {
            try {
                console.log(`Trying to load from: ${path}`);
                const response = await fetch(path);
                if (response.ok) {
                    this.spec = await response.json();
                    console.log('Successfully loaded OpenAPI spec from:', path);
                    return true;
                }
            } catch (error) {
                console.log(`Failed to load from ${path}:`, error.message);
            }
        }
        
        console.log('No JSON version available in any location');
        return false;
    }

    /**
     * Parse YAML using js-yaml library (if available) or manual parsing
     */
    async parseYAML(yamlText) {
        // Check if js-yaml is loaded
        if (typeof jsyaml !== 'undefined') {
            return jsyaml.load(yamlText);
        }

        // Fallback: basic manual parsing for OpenAPI structure
        return this.manualYAMLParse(yamlText);
    }

    /**
     * Manual YAML parser (simplified for OpenAPI structure)
     * This is a basic implementation - for production, use js-yaml library
     */
    manualYAMLParse(yamlText) {
        const lines = yamlText.split('\n');
        const result = {
            openapi: '',
            info: {},
            servers: [],
            paths: {},
            components: {
                schemas: {}
            }
        };

        let currentPath = [];
        let currentIndent = 0;
        let inSchemas = false;
        let currentSchema = null;
        let schemaLines = [];

        for (let i = 0; i < lines.length; i++) {
            const line = lines[i];
            const indent = line.search(/\S/);
            
            if (indent === -1) continue; // Empty line

            const content = line.trim();
            
            // Check if we're in the schemas section
            if (content.startsWith('schemas:')) {
                inSchemas = true;
                continue;
            }

            // If we're in schemas, collect schema definitions
            if (inSchemas) {
                // Detect when we leave schemas section (lower or equal indent to 'components')
                if (indent <= 2 && !content.startsWith('-') && content.includes(':') && !content.startsWith(' ')) {
                    inSchemas = false;
                    continue;
                }

                // New schema definition (indent 4 from 'schemas')
                if (indent === 4 && content.includes(':') && !content.startsWith('-')) {
                    if (currentSchema) {
                        // Parse previous schema
                        result.components.schemas[currentSchema] = this.parseSchemaBlock(schemaLines);
                    }
                    currentSchema = content.replace(':', '').trim();
                    schemaLines = [];
                } else if (currentSchema && indent > 4) {
                    schemaLines.push(line);
                }
            }
        }

        // Parse last schema
        if (currentSchema && schemaLines.length > 0) {
            result.components.schemas[currentSchema] = this.parseSchemaBlock(schemaLines);
        }

        return result;
    }

    /**
     * Parse a schema block from YAML lines
     */
    parseSchemaBlock(lines) {
        const schema = {
            type: 'object',
            properties: {},
            required: []
        };

        let currentProperty = null;
        let inProperties = false;
        let inRequired = false;

        for (const line of lines) {
            const indent = line.search(/\S/);
            const content = line.trim();

            if (content.startsWith('type:')) {
                schema.type = content.split(':')[1].trim();
            } else if (content.startsWith('description:')) {
                schema.description = content.substring(12).trim().replace(/^["']|["']$/g, '');
            } else if (content === 'properties:') {
                inProperties = true;
                inRequired = false;
            } else if (content === 'required:') {
                inRequired = true;
                inProperties = false;
            } else if (inRequired && content.startsWith('- ')) {
                schema.required.push(content.substring(2).trim());
            } else if (inProperties && indent === 8 && content.includes(':')) {
                currentProperty = content.split(':')[0].trim();
                schema.properties[currentProperty] = {};
            } else if (currentProperty && indent > 8) {
                if (content.startsWith('type:')) {
                    schema.properties[currentProperty].type = content.split(':')[1].trim();
                } else if (content.startsWith('description:')) {
                    schema.properties[currentProperty].description = 
                        content.substring(12).trim().replace(/^["']|["']$/g, '');
                } else if (content.startsWith('$ref:')) {
                    schema.properties[currentProperty].$ref = content.split(':')[1].trim().replace(/["']/g, '');
                } else if (content.startsWith('format:')) {
                    schema.properties[currentProperty].format = content.split(':')[1].trim();
                }
            }
        }

        return schema;
    }

    /**
     * Get all schemas from the loaded spec
     */
    getSchemas() {
        return this.spec?.components?.schemas || {};
    }

    /**
     * Get a specific schema by name
     */
    getSchema(name) {
        return this.spec?.components?.schemas?.[name] || null;
    }
}

// Export for use in other scripts
window.OpenAPILoader = OpenAPILoader;
