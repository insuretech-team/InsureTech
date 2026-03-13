//
//  SplashViewController.swift
//  oneTapService
//
//  Created by Joy Biswas - joybiswas101626@gmail.com 01689639998 on 9/1/25.
//

import UIKit

class SplashViewController: UIViewController {

    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(animated)
        
//        Task {
//            await RemoteConfigService.shared.fetch()
//
//            DispatchQueue.main.asyncAfter(deadline: .now() + 0.3) {
//                let ekycVC = FLVETestVC()
//                ekycVC.modalPresentationStyle = .fullScreen
//                self.present(ekycVC, animated: true)
//            }
//        }
        
        
//        DispatchQueue.main.asyncAfter(deadline: .now() + 2.0) {
//            let storyboard = UIStoryboard(name: "Main", bundle: nil)
//            let auth = storyboard.instantiateViewController(withIdentifier: "LoginViewController") as! LoginViewController
//            auth.modalPresentationStyle = .fullScreen
//            self.present(auth, animated: true, completion: nil)
//        }
        
        DispatchQueue.main.async {
            let storyboard = UIStoryboard(name: "Home", bundle: nil)
            let homeVC = storyboard.instantiateViewController(withIdentifier: "HomeViewController") as! HomeViewController
            homeVC.modalPresentationStyle = .fullScreen
            self.present(homeVC, animated: false, completion: nil)
        }
    }
}
