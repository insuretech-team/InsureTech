////
////  RemoteAPIConfig.swift
////  InsureTech
////
////  Created by LifeplusBD on 20/1/26.
////
//
//import Foundation
//
//
//struct RemoteAPIConfig: Decodable {
//    let baseURL: String?
//    let authToken: String?
//}
//
//final class RuntimeAPIConfig {
//
//    static let shared = RuntimeAPIConfig()
//    private init() {}
//
//    private let fallbackBaseURL = "https://farukhannan-flve.hf.space"
//    private let fallbackToken = "Bearer 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"
//
//    private(set) var baseURL: String?
//    private(set) var authToken: String?
//
//    func resolvedBaseURL() -> String {
//        return baseURL ?? UserDefaults.standard.string(forKey: "api_base_url") ?? fallbackBaseURL
//    }
//
//    func resolvedAuthToken() -> String? {
//        return authToken
//            ?? UserDefaults.standard.string(forKey: "api_auth_token") ?? fallbackToken
//    }
//
//    fileprivate func update(baseURL: String?, token: String?) {
//        if let baseURL {
//            self.baseURL = baseURL
//            UserDefaults.standard.set(baseURL, forKey: "api_base_url")
//        }
//
//        if let token {
//            self.authToken = token
//            UserDefaults.standard.set(token, forKey: "api_auth_token")
//        }
//    }
//}
//
//final class RemoteConfigService {
//
//    static let shared = RemoteConfigService()
//    private init() {}
//
//    func fetch() async {
//        guard let url = URL(string: "test_test_test") else {
//            return
//        }
//
//        var request = URLRequest(url: url)
//        request.httpMethod = "GET"
//        request.timeoutInterval = 10
//
//        do {
//            let (data, response) = try await URLSession.shared.data(for: request)
//
//            guard let http = response as? HTTPURLResponse,
//                  200..<300 ~= http.statusCode else {
//                return
//            }
//
//            let config = try JSONDecoder().decode(RemoteAPIConfig.self, from: data)
//
//            RuntimeAPIConfig.shared.update(
//                baseURL: config.baseURL,
//                token: config.authToken
//            )
//
//            print("✅ Remote config applied")
//
//        } catch {
//            print("⚠️ Remote config failed, using fallback")
//        }
//    }
//}
//
