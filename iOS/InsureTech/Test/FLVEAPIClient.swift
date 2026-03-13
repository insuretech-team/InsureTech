//
//  FLVEAPIClient.swift
//  InsureTech
//
//  Created by LifeplusBD on 20/1/26.
//

//import Foundation
//
//
//final class FLVEAPIClient {
//
//    private let baseURL = URL(string: "https://farukhannan-flve.hf.space")!
//
//    func submitFrame(_ jpeg: Data,
//                     completion: @escaping (FLVEResponse) -> Void) {
//
//        var request = URLRequest(url: baseURL.appendingPathComponent("/liveness"))
//        request.httpMethod = "POST"
//        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
//
//        let body: [String: Any] = [
//            "image_data": jpeg.base64EncodedString(),
//            "format": "IMAGE_FORMAT_JPEG"
//        ]
//
//        request.httpBody = try? JSONSerialization.data(withJSONObject: body)
//
//        URLSession.shared.dataTask(with: request) { data, _, _ in
//            guard let data,
//                  let decoded = try? JSONDecoder().decode(FLVEResponse.self, from: data)
//            else { return }
//
//            completion(decoded)
//        }.resume()
//    }
//}


import Foundation

final class FLVEAPIClient {

    private let baseURL = URL(string: "https://farukhannan-flve.hf.space")!
    private let token = "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"

    func submitFrame(_ jpeg: Data,
                     completion: @escaping (FLVEResponse) -> Void) {

        var request = URLRequest(url: baseURL.appendingPathComponent("/liveness"))
        request.httpMethod = "POST"

        // Correct headers
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        // JSON body with base64
        let body: [String: Any] = [
            "image_data": jpeg.base64EncodedString(),
            "format": "IMAGE_FORMAT_JPEG"
        ]

        do {
            request.httpBody = try JSONSerialization.data(withJSONObject: body, options: [])
        } catch {
            print("Failed to encode JSON body:", error)
            return
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Network error:", error)
                return
            }

            guard let data else { return }

            do {
                let decoded = try JSONDecoder().decode(FLVEResponse.self, from: data)
                DispatchQueue.main.async {
                    completion(decoded)
                }
            } catch {
                print("Failed to decode response:", error)
            }

        }.resume()
    }
}
