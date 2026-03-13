//
//  RemoteConfigService.swift
//  InsureTech
//
//  Created by LifeplusBD on 20/1/26.
//


final class RemoteConfigService {

    static let shared = RemoteConfigService()
    private init() {}

    func fetch() async {
        // Later: Firebase / custom backend
        try? await Task.sleep(nanoseconds: 300_000_000)
    }
}
