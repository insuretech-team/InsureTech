import fs from 'fs';
import path from 'path';

function walk(dir) {
    let results = [];
    const list = fs.readdirSync(dir);
    list.forEach(file => {
        file = path.resolve(dir, file);
        const stat = fs.statSync(file);
        if (stat && stat.isDirectory()) {
            results = results.concat(walk(file));
        } else if (file.match(/\.(svelte|ts)$/)) {
            results.push(file);
        }
    });
    return results;
}

const files = walk('./src');
files.forEach(file => {
    let content = fs.readFileSync(file, 'utf8');
    let changed = false;

    if (content.includes('$lib/generated/insuretech')) {
        // Replace imports
        // e.g. import { ClaimStatus, ClaimType } from '$lib/generated/insuretech/claims/entity...';
        content = content.replace(/import\s*\{([^}]+)\}\s*from\s*['"]\$lib\/generated\/insuretech.*?['"];?/g, "import { $1 } from '$lib/types';");
        changed = true;
    }

    if (content.match(/new\s+(Claim|ClaimDocument|ClaimApproval|FraudCheckResult|Product|ProductRider|Policy|PolicyDocument)\s*\(/g)) {
        // Replace `new Claim({` with `({`
        content = content.replace(/new\s+(Claim|ClaimDocument|ClaimApproval|FraudCheckResult|Product|ProductRider|Policy|PolicyDocument)\s*\(/g, "<$1> (");
        changed = true;
    }

    if (changed) {
        fs.writeFileSync(file, content);
        console.log('Updated:', file);
    }
});
