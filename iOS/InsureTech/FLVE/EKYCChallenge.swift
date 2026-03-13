//
//  EKYCChallenge.swift
//  InsureTech
//
//  Created by LifeplusBD on 18/1/26.
//


import Foundation

enum EKYCChallenge: String, Codable {
    case blink
    case lookLeft = "look_left"
    case lookRight = "look_right"
    case capture
}

// MARK: - eKYC

struct EKYCStartResponse: Codable {
    let sessionId: String
    let currentChallenge: EKYCChallenge
}

struct EKYCFrameResponse: Codable {
    let stepCompleted: Bool
    let nextChallenge: EKYCChallenge?
}

struct EKYCCompleteResponse: Codable {
    let success: Bool
}
