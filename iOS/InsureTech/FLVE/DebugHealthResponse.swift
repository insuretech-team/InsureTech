//
//  DebugHealthResponse.swift
//  InsureTech
//
//  Created by LifeplusBD on 19/1/26.
//

import UIKit

struct LivenessRequest: Codable {
    let imageData: String
    let format: String

    enum CodingKeys: String, CodingKey {
        case imageData = "image_data"
        case format
    }
}

extension UIImage {
    func toBase64JPEG(compression: CGFloat = 0.8) -> String? {
        guard let data = self.jpegData(compressionQuality: compression) else {
            return nil
        }
        return data.base64EncodedString()
    }
}
