// Utility functions for generating demo data with Bangladesh context

export const bangladeshLocations = [
  'Dhaka',
  'Chittagong',
  'Sylhet',
  'Rajshahi',
  'Khulna',
  'Barisal',
  'Rangpur',
  'Mymensingh',
  'Gazipur',
  'Narayanganj',
];

export const bangladeshBanks = [
  'Dutch Bangla Bank',
  'Sonali Bank',
  'Janata Bank',
  'Agrani Bank',
  'BRAC Bank',
  'City Bank',
  'Eastern Bank',
  'Islami Bank',
  'Prime Bank',
  'Standard Chartered Bank',
];

export const bangladeshNames = {
  male: [
    'Mohammad Rahman',
    'Abdul Karim',
    'Md. Hasan',
    'Rafiqul Islam',
    'Kamal Hossain',
    'Jahangir Alam',
    'Mizanur Rahman',
    'Shamsul Haque',
    'Nurul Islam',
    'Abdur Razzak',
  ],
  female: [
    'Fatema Begum',
    'Nasrin Akter',
    'Rahima Khatun',
    'Salma Begum',
    'Ayesha Siddika',
    'Roksana Parvin',
    'Shahana Akter',
    'Taslima Khatun',
    'Farzana Rahman',
    'Sultana Begum',
  ],
  organizations: [
    'LabAid Hospital',
    'Square Hospital',
    'United Hospital',
    'Apollo Hospital',
    'Popular Diagnostic Centre',
    'Ibn Sina Hospital',
    'Dhaka Medical College Hospital',
    'Holy Family Hospital',
    'Evercare Hospital',
    'Delta Medical College Hospital',
  ],
};

/**
 * Generate a random Bangladesh phone number in +880 format
 */
export function generateBDPhone(): string {
  const operators = ['13', '14', '15', '16', '17', '18', '19'];
  const operator = operators[Math.floor(Math.random() * operators.length)];
  const number = Math.floor(10000000 + Math.random() * 90000000);
  return `+880${operator}${number}`;
}

/**
 * Generate a random Bangladesh NID (10, 13, or 17 digits)
 */
export function generateNID(): string {
  const lengths = [10, 13, 17];
  const length = lengths[Math.floor(Math.random() * lengths.length)];
  let nid = '';
  for (let i = 0; i < length; i++) {
    nid += Math.floor(Math.random() * 10);
  }
  return nid;
}

/**
 * Generate a random TIN (12 digits)
 */
export function generateTIN(): string {
  let tin = '';
  for (let i = 0; i < 12; i++) {
    tin += Math.floor(Math.random() * 10);
  }
  return tin;
}

/**
 * Generate a random trade license number
 */
export function generateTradeLicense(location: string): string {
  const year = new Date().getFullYear();
  const number = Math.floor(100000 + Math.random() * 900000);
  return `TL-${location.substring(0, 3).toUpperCase()}-${year}-${number}`;
}

/**
 * Generate a random bank account number
 */
export function generateBankAccount(): string {
  let account = '';
  for (let i = 0; i < 13; i++) {
    account += Math.floor(Math.random() * 10);
  }
  return account;
}

/**
 * Format currency in BDT with thousand separators
 */
export function formatBDT(amount: number): string {
  return `৳ ${amount.toLocaleString('en-BD')}`;
}

/**
 * Format date in DD/MM/YYYY format
 */
export function formatDate(date: Date): string {
  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();
  return `${day}/${month}/${year}`;
}

/**
 * Generate a random date within the last N days
 */
export function generateRecentDate(daysAgo: number): string {
  const date = new Date();
  date.setDate(date.getDate() - Math.floor(Math.random() * daysAgo));
  return date.toISOString();
}

/**
 * Generate a random date in the future within N days
 */
export function generateFutureDate(daysAhead: number): string {
  const date = new Date();
  date.setDate(date.getDate() + Math.floor(Math.random() * daysAhead));
  return date.toISOString();
}

/**
 * Generate a random ID
 */
export function generateId(prefix: string = ''): string {
  const timestamp = Date.now();
  const random = Math.floor(Math.random() * 10000);
  return prefix ? `${prefix}-${timestamp}-${random}` : `${timestamp}-${random}`;
}

/**
 * Get a random item from an array
 */
export function randomItem<T>(array: T[]): T {
  return array[Math.floor(Math.random() * array.length)];
}

/**
 * Get multiple random items from an array
 */
export function randomItems<T>(array: T[], count: number): T[] {
  const shuffled = [...array].sort(() => 0.5 - Math.random());
  return shuffled.slice(0, Math.min(count, array.length));
}

/**
 * Generate a random number between min and max (inclusive)
 */
export function randomNumber(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

/**
 * Generate a random percentage (0-100)
 */
export function randomPercentage(): number {
  return Math.floor(Math.random() * 101);
}

/**
 * Generate a random amount in BDT within a range
 */
export function randomAmount(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

/**
 * Calculate days between two dates
 */
export function daysBetween(date1: string, date2: string): number {
  const d1 = new Date(date1);
  const d2 = new Date(date2);
  const diffTime = Math.abs(d2.getTime() - d1.getTime());
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
}
