//
//  FLVEAPI.swift
//  InsureTech
//
//  Created by LifeplusBD on 18/1/26.
//


import UIKit

enum FLVEAPI {
    static func startEKYC(
        userId: String,
        challenges: [EKYCChallenge]
    ) async throws -> EKYCStartResponse {

        return EKYCStartResponse(
            sessionId: UUID().uuidString,
            currentChallenge: challenges.first ?? .blink
        )
    }

    static func submitFrame(
        sessionId: String,
        image: UIImage
    ) async throws -> EKYCFrameResponse {

        return EKYCFrameResponse(
            stepCompleted: true,
            nextChallenge: nil
        )
    }

    static func completeEKYC(
        sessionId: String
    ) async throws -> EKYCCompleteResponse {

        return EKYCCompleteResponse(success: true)
    }
}
