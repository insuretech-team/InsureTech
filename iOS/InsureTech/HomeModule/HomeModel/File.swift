//
//  File.swift
//  InsureTech
//
//  Created by LifeplusBD on 15/2/26.
//


var categories: [InsuranceCategory] = [
    InsuranceCategory(iconName: "car.fill", name: "Car Insurance"),
    InsuranceCategory(iconName: "cross.fill", name: "Health Insurance"),
    InsuranceCategory(iconName: "house.fill", name: "Home Insurance"),
    InsuranceCategory(iconName: "airplane", name: "Travel Insurance")
]

var policies: [Policy] = [
    Policy(bannerName: "policy1",
           policyName: "Premium Care",
           coverageType1: "Accidental Cover",
           coverageType2: "Medical Support",
           iconName: "shield.fill"),
    
    Policy(bannerName: "policy2",
           policyName: "Family Secure",
           coverageType1: "Life Cover",
           coverageType2: "Critical Illness",
           iconName: "heart.fill"),
    
    Policy(bannerName: "policy3",
           policyName: "Vehicle Shield",
           coverageType1: "Collision Cover",
           coverageType2: "Roadside Assist",
           iconName: "car.fill")
]

var featuredPolicies: [FeaturedPolicy] = [
    FeaturedPolicy(bannerName: "featured1",
                   title: "Best Health Plan 2026",
                   subtitle: "Up to 10L coverage"),
    
    FeaturedPolicy(bannerName: "featured2",
                   title: "Smart Travel Protect",
                   subtitle: "Global Coverage"),
    
    FeaturedPolicy(bannerName: "featured3",
                   title: "Secure Home Plus",
                   subtitle: "Fire & Theft Protection")
]
