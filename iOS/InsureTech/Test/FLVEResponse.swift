//
//  FLVEResponse.swift
//  InsureTech
//
//  Created by LifeplusBD on 20/1/26.
//

import CoreFoundation


struct FLVEResponse: Decodable {
    let isLive: Bool?
    let confidence: Double?
    let box: [CGFloat]?
    let error: String?
}
