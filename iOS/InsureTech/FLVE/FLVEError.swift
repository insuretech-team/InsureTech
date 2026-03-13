//
//  FLVEError.swift
//  InsureTech
//
//  Created by LifeplusBD on 18/1/26.
//


import Foundation

enum FLVEError: Error {
    case systemNotReady
    case modelUnavailable
    case livenessFailed
    case ekycFailed
    case cameraUnavailable
    case unknown
}
