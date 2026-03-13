//
//  ProfileViewController.swift
//  InsureTech
//
//  Created by LifeplusBD on 15/2/26.
//

import UIKit

class ProfileViewController: UIViewController {
    
    @IBOutlet weak var imgUser: UIImageView!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
    }
    
    @IBAction func backTapped() {
        self.dismiss(animated: true)
    }
    
}
