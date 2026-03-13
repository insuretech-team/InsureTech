const fs = require('fs');

const services = [
    { name: 'tenant', grpc: 50050, http: 50051 },
    { name: 'authn', grpc: 50060, http: 50061, deps: ['redis', 'kafka'] },
    { name: 'authz', grpc: 50070, http: 50071 },
    { name: 'audit', grpc: 50080, http: 50081 },
    { name: 'kyc', grpc: 50090, http: 50091 },
    { name: 'partner', grpc: 50100, http: 50101 },
    { name: 'beneficiary', grpc: 50110, http: 50111 },
    { name: 'b2b', grpc: 50112, http: 50113 },
    { name: 'workflow', grpc: 50180, http: 50181 },
    { name: 'payment', grpc: 50190, http: 50191 },
    { name: 'fraud', grpc: 50220, http: 50221 },
    { name: 'notification', grpc: 50230, http: 50231 },
    { name: 'support', grpc: 50240, http: 50241 },
    { name: 'webrtc', grpc: 50250, http: 50251 },
    { name: 'media', grpc: 50260, http: 50261 },
    { name: 'ocr', grpc: 50270, http: 50271 },
    { name: 'docgen', grpc: 50280, http: 50281 },
    { name: 'storage', grpc: 50290, http: 50291 },
    { name: 'iot', grpc: 50300, http: 50301 },
    { name: 'analytics', grpc: 50310, http: 50311 },
    { name: 'ai', grpc: 50320, http: 50321 },
];

let yaml = '';

for (const s of services) {
    yaml += `  ${s.name}:\n`;
    yaml += `    build:\n`;
    yaml += `      context: .\n`;
    yaml += `      dockerfile: backend/infra/docker/${s.name}/Dockerfile\n`;
    yaml += `    container_name: insuretech-${s.name}\n`;
    yaml += `    profiles: ["${s.name}", "full"]\n`;
    yaml += `    env_file:\n`;
    yaml += `      - .env\n`;
    if (s.deps) {
        yaml += `    depends_on:\n`;
        for (const d of s.deps) {
            yaml += `      ${d}:\n`;
            yaml += `        condition: service_healthy\n`;
        }
    }
    yaml += `    ports:\n`;
    yaml += `      - "${s.grpc}:${s.grpc}"\n`;
    yaml += `      - "${s.http}:${s.http}"\n`;
    yaml += `    restart: unless-stopped\n\n`;
}

// Add B2B portal
yaml += `  b2b_portal:\n`;
yaml += `    build:\n`;
yaml += `      context: .\n`;
yaml += `      dockerfile: backend/infra/docker/b2b_portal/Dockerfile\n`;
yaml += `    container_name: insuretech-b2b-portal\n`;
yaml += `    profiles: ["frontend", "full"]\n`;
yaml += `    env_file:\n`;
yaml += `      - .env\n`;
yaml += `    ports:\n`;
yaml += `      - "3000:3000"\n`;
yaml += `    restart: unless-stopped\n\n`;

const og = fs.readFileSync('docker-compose.yml', 'utf-8');
// The authn block starts at `  authn:` and goes to `restart: unless-stopped` before polisync block
const authnIndex = og.indexOf('  authn:');
const polisyncIndex = og.indexOf('  # ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n  # PoliSync — C# .NET 8 Insurance Commerce & Policy Engine');

const before = og.slice(0, authnIndex);
const after = og.slice(polisyncIndex);

const finalYaml = before + yaml + after;
fs.writeFileSync('docker-compose.yml', finalYaml);

console.log("Successfully wired all microservices into docker-compose.yml!");
