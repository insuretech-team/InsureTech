//
//  HealthResponse.swift
//  InsureTech
//
//  Created by LifeplusBD on 19/1/26.
//

import UIKit


// MARK: - Simple Health
struct HealthResponse: Codable {
    let status: String?
    let device: String?
    let models: HealthModels?
}

struct HealthModels: Codable {
    let detector: Bool?
    let embedder: Bool?
    let liveness: Bool?
}

// MARK: - Debug Health
struct DebugHealthResponse: Codable {
    let pythonVersion: String?
    let engineInitialized: Bool?
    let startupError: String?
    let modelsDirExists: Bool?
    let modelsFound: [String]?
    let modelSizesMB: [String: Double]?
    let onnxruntimeProviders: [String]?
    let mediapipeAvailable: Bool?
    let mediapipeError: String?
    let device: String?
    let detectorLoaded: Bool?
    let embedderLoaded: Bool?
    let livenessLoaded: Bool?

    enum CodingKeys: String, CodingKey {
        case pythonVersion = "python_version"
        case engineInitialized = "engine_initialized"
        case startupError = "startup_error"
        case modelsDirExists = "models_dir_exists"
        case modelsFound = "models_found"
        case modelSizesMB = "model_sizes_mb"
        case onnxruntimeProviders = "onnxruntime_providers"
        case mediapipeAvailable = "mediapipe_available"
        case mediapipeError = "mediapipe_error"
        case device
        case detectorLoaded = "detector_loaded"
        case embedderLoaded = "embedder_loaded"
        case livenessLoaded = "liveness_loaded"
    }
}

// MARK: - Liveness
struct LivenessResponse: Codable {
    let box: [Int]?
    let confidence: Double?
    let error: String?
    let isLive: Bool?
    let livenessScore: Double?

    enum CodingKeys: String, CodingKey {
        case box, confidence, error
        case isLive = "is_live"
        case livenessScore = "liveness_score"
    }
}



