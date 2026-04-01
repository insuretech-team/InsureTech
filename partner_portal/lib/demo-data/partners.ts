import type { Partner, PartnerType, PartnerStatus, PolicyType } from '../types';
import {
  bangladeshLocations,
  bangladeshBanks,
  bangladeshNames,
  generateBDPhone,
  generateTIN,
  generateTradeLicense,
  generateBankAccount,
  generateRecentDate,
  generateId,
  randomItem,
  randomItems,
  randomNumber,
  randomPercentage,
} from './utils';

const partnerTypeConfig: Record<
  PartnerType,
  {
    defaultPolicyTypes: PolicyType[];
    cashlessDefault: boolean;
    discountDefault: boolean;
    namePrefix: string;
  }
> = {
  HOSPITAL: {
    defaultPolicyTypes: ['HEALTH', 'LIFE'],
    cashlessDefault: true,
    discountDefault: false,
    namePrefix: 'Hospital',
  },
  CLINIC: {
    defaultPolicyTypes: ['HEALTH'],
    cashlessDefault: true,
    discountDefault: true,
    namePrefix: 'Clinic',
  },
  PHARMACY: {
    defaultPolicyTypes: ['HEALTH'],
    cashlessDefault: false,
    discountDefault: true,
    namePrefix: 'Pharmacy',
  },
  AUTO_REPAIR: {
    defaultPolicyTypes: ['MOTOR'],
    cashlessDefault: true,
    discountDefault: false,
    namePrefix: 'Auto Service Center',
  },
  PET_CLINIC: {
    defaultPolicyTypes: ['PET'],
    cashlessDefault: true,
    discountDefault: false,
    namePrefix: 'Pet Clinic',
  },
  FIRE_INSPECTOR: {
    defaultPolicyTypes: ['FIRE'],
    cashlessDefault: false,
    discountDefault: false,
    namePrefix: 'Fire Inspector',
  },
  LAPTOP_REPAIR: {
    defaultPolicyTypes: ['DEVICE'],
    cashlessDefault: true,
    discountDefault: false,
    namePrefix: 'Laptop Repair',
  },
  MOBILE_REPAIR: {
    defaultPolicyTypes: ['DEVICE'],
    cashlessDefault: true,
    discountDefault: false,
    namePrefix: 'Mobile Repair',
  },
  AMBULANCE: {
    defaultPolicyTypes: ['HEALTH'],
    cashlessDefault: true,
    discountDefault: false,
    namePrefix: 'Ambulance Service',
  },
  MFS: {
    defaultPolicyTypes: ['HEALTH', 'MOTOR', 'LIFE', 'FIRE', 'PET', 'DEVICE', 'TRAVEL'],
    cashlessDefault: false,
    discountDefault: false,
    namePrefix: 'MFS Provider',
  },
  ECOMMERCE: {
    defaultPolicyTypes: ['HEALTH', 'MOTOR', 'LIFE', 'FIRE', 'PET', 'DEVICE', 'TRAVEL'],
    cashlessDefault: false,
    discountDefault: false,
    namePrefix: 'E-commerce',
  },
  AGENT_NETWORK: {
    defaultPolicyTypes: ['HEALTH', 'MOTOR', 'LIFE', 'FIRE', 'PET', 'DEVICE', 'TRAVEL'],
    cashlessDefault: false,
    discountDefault: false,
    namePrefix: 'Agent Network',
  },
  CORPORATE: {
    defaultPolicyTypes: ['HEALTH', 'LIFE'],
    cashlessDefault: false,
    discountDefault: false,
    namePrefix: 'Corporate',
  },
};

function generatePartner(
  partnerType: PartnerType,
  status: PartnerStatus = 'ACTIVE'
): Partner {
  const config = partnerTypeConfig[partnerType];
  const location = randomItem(bangladeshLocations);
  const orgName = randomItem(bangladeshNames.organizations);
  
  const cashlessEnabled = config.cashlessDefault && Math.random() > 0.2;
  const discountEnabled = config.discountDefault && Math.random() > 0.3;

  return {
    id: generateId('partner'),
    organizationName: `${orgName} ${location}`,
    partnerType,
    status,
    tradeLicense: generateTradeLicense(location),
    tin: generateTIN(),
    contactEmail: `contact@${orgName.toLowerCase().replace(/\s+/g, '')}.com`,
    contactPhone: generateBDPhone(),
    bankAccount: generateBankAccount(),
    bankName: randomItem(bangladeshBanks),
    bankBranch: `${location} Branch`,
    policyTypes: config.defaultPolicyTypes,
    cashlessEnabled,
    cashlessLimit: cashlessEnabled ? randomNumber(100000, 500000) : undefined,
    discountEnabled,
    discountPercentage: discountEnabled ? randomNumber(5, 25) : undefined,
    autoApprovalThreshold: cashlessEnabled ? randomNumber(5000, 15000) : undefined,
    serviceLocations: randomItems(bangladeshLocations, randomNumber(1, 3)),
    nationwideCoverage: Math.random() > 0.7,
    acquisitionRate: randomNumber(10, 25),
    renewalRate: randomNumber(5, 15),
    claimsAssistanceRate: Math.random() > 0.5 ? randomNumber(3, 8) : undefined,
    agentCount: randomNumber(5, 50),
    activePolicies: randomNumber(50, 500),
    totalClaims: randomNumber(10, 200),
    conversionRate: randomNumber(15, 85),
    createdAt: generateRecentDate(365),
  };
}

export function generatePartners(count: number = 50): Partner[] {
  const partners: Partner[] = [];
  const partnerTypes: PartnerType[] = [
    'HOSPITAL',
    'CLINIC',
    'PHARMACY',
    'AUTO_REPAIR',
    'PET_CLINIC',
    'FIRE_INSPECTOR',
    'LAPTOP_REPAIR',
    'MOBILE_REPAIR',
    'AMBULANCE',
    'MFS',
    'ECOMMERCE',
    'AGENT_NETWORK',
    'CORPORATE',
  ];

  const statuses: PartnerStatus[] = [
    'ACTIVE',
    'ACTIVE',
    'ACTIVE',
    'ACTIVE',
    'ACTIVE',
    'PENDING_VERIFICATION',
    'SUSPENDED',
  ];

  for (let i = 0; i < count; i++) {
    const partnerType = randomItem(partnerTypes);
    const status = randomItem(statuses);
    partners.push(generatePartner(partnerType, status));
  }

  return partners;
}

// Generate a default set of partners
export const demoPartners = generatePartners(50);

// Get partner by ID
export function getPartnerById(id: string): Partner | undefined {
  return demoPartners.find((p) => p.id === id);
}

// Filter partners by type
export function getPartnersByType(type: PartnerType): Partner[] {
  return demoPartners.filter((p) => p.partnerType === type);
}

// Filter partners by status
export function getPartnersByStatus(status: PartnerStatus): Partner[] {
  return demoPartners.filter((p) => p.status === status);
}

// Filter partners by policy type
export function getPartnersByPolicyType(policyType: PolicyType): Partner[] {
  return demoPartners.filter((p) => p.policyTypes.includes(policyType));
}

// Search partners by name
export function searchPartners(query: string): Partner[] {
  const lowerQuery = query.toLowerCase();
  return demoPartners.filter(
    (p) =>
      p.organizationName.toLowerCase().includes(lowerQuery) ||
      p.tradeLicense.toLowerCase().includes(lowerQuery) ||
      p.tin.includes(lowerQuery)
  );
}
