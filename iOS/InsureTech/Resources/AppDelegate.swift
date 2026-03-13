//
//  Untitled.swift
//  InsureTech
//
//  Created by LifeplusBD on 4/1/26.
//

import UIKit
import IQKeyboardManagerSwift

@main
class AppDelegate: UIResponder, UIApplicationDelegate {
    
    var window: UIWindow?
    
    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
        
        IQKeyboardManager.shared.isEnabled = true
        IQKeyboardManager.shared.resignOnTouchOutside = true
        
        window = UIWindow(frame: UIScreen.main.bounds)
        window?.rootViewController = SplashViewController()
        window?.makeKeyAndVisible()
        
        return true
    }
}
