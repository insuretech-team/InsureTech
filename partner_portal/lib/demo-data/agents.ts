import type { Agent, AgentStatus } from '../types';
import { demoPartners } from './partners';
import {
  bangladeshNames,
  generateBDPhone,
  generateNID,
  generateRecentDate,
  generateId,
  randomItem,
  randomNumber,
} from './utils';

function generateAgent(partnerId: string, status: AgentStatus = 'ACTIVE'): Agent {
  const isMale = Math.random() > 0.3;
  const fullName = isMale
    ? randomItem(bangladeshNames.male)
    : randomItem(bangladeshNames.female);
  
  const phone = generateBDPhone();
  const email = Math.random() > 0.3 
    ? `${fullName.toLowerCase().replace(/\s+/g, '.')}@agent.com`
    : undefined;

  const territories = [
    'Dhaka North',
    'Dhaka South',
    'Chittagong',
    'Sylhet',
    'Rajshahi',
    'Khulna',
    'Barisal',
    'Rangpur',
  ];

  return {
    id: generateId('agent'),
    partnerId,
    fullName,
    phone,
    email,
    nid: generateNID(),
    status,
    commissionRate: randomNumber(5, 20),
    territory: randomItem(territories),
    policiesSold: randomNumber(10, 150),
    commissionEarned: randomNumber(50000, 500000),
    referralCount: randomNumber(20, 200),
    conversionRate: randomNumber(20, 80),
    createdAt: generateRecentDate(180),
  };
}

export function generateAgents(count: number = 100): Agent[] {
  const agents: Agent[] = [];
  const statuses: AgentStatus[] = ['ACTIVE', 'ACTIVE', 'ACTIVE', 'ACTIVE', 'INACTIVE', 'SUSPENDED'];

  // Distribute agents across partners
  const activePartners = demoPartners.filter(p => p.status === 'ACTIVE');
  
  for (let i = 0; i < count; i++) {
    const partner = randomItem(activePartners);
    const status = randomItem(statuses);
    agents.push(generateAgent(partner.id, status));
  }

  return agents;
}

// Generate a default set of agents
export const demoAgents = generateAgents(100);

// Get agent by ID
export function getAgentById(id: string): Agent | undefined {
  return demoAgents.find((a) => a.id === id);
}

// Get agents by partner ID
export function getAgentsByPartnerId(partnerId: string): Agent[] {
  return demoAgents.filter((a) => a.partnerId === partnerId);
}

// Filter agents by status
export function getAgentsByStatus(status: AgentStatus): Agent[] {
  return demoAgents.filter((a) => a.status === status);
}

// Filter agents by territory
export function getAgentsByTerritory(territory: string): Agent[] {
  return demoAgents.filter((a) => a.territory === territory);
}

// Search agents by name, phone, or NID
export function searchAgents(query: string): Agent[] {
  const lowerQuery = query.toLowerCase();
  return demoAgents.filter(
    (a) =>
      a.fullName.toLowerCase().includes(lowerQuery) ||
      a.phone.includes(lowerQuery) ||
      a.nid.includes(lowerQuery)
  );
}

// Get top performing agents
export function getTopAgents(limit: number = 5): Agent[] {
  return [...demoAgents]
    .filter(a => a.status === 'ACTIVE')
    .sort((a, b) => b.policiesSold - a.policiesSold)
    .slice(0, limit);
}
