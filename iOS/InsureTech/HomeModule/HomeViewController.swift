//
//  HomeViewController.swift
//  InsureTech
//
//  Created by LifeplusBD on 6/1/26.
//

import Foundation
import UIKit

class HomeViewController: UIViewController, UICollectionViewDataSource, UICollectionViewDelegate, UICollectionViewDelegateFlowLayout {

    var categories: [InsuranceCategory] = [
        InsuranceCategory(iconName: "car.fill", name: "Car Insurance"),
        InsuranceCategory(iconName: "cross.fill", name: "Health Insurance"),
        InsuranceCategory(iconName: "house.fill", name: "Home Insurance"),
        InsuranceCategory(iconName: "airplane", name: "Travel Insurance")
    ]

    var policies: [Policy] = [
        Policy(bannerName: "policy1",
               policyName: "Premium Care",
               coverageType1: "Accidental Cover",
               coverageType2: "Medical Support",
               iconName: "shield.fill"),
        
        Policy(bannerName: "policy2",
               policyName: "Family Secure",
               coverageType1: "Life Cover",
               coverageType2: "Critical Illness",
               iconName: "heart.fill"),
        
        Policy(bannerName: "policy3",
               policyName: "Vehicle Shield",
               coverageType1: "Collision Cover",
               coverageType2: "Roadside Assist",
               iconName: "car.fill")
    ]

    var featuredPolicies: [FeaturedPolicy] = [
        FeaturedPolicy(bannerName: "featured1",
                       title: "Best Health Plan 2026",
                       subtitle: "Up to 10L coverage"),
        
        FeaturedPolicy(bannerName: "featured2",
                       title: "Smart Travel Protect",
                       subtitle: "Global Coverage"),
        
        FeaturedPolicy(bannerName: "featured3",
                       title: "Secure Home Plus",
                       subtitle: "Fire & Theft Protection")
    ]

    
//    @IBOutlet weak var imgBanner: UIImageView!
    @IBOutlet weak var imgUser: UIImageView!
    @IBOutlet weak var lblStaticWelcome: UILabel!
    @IBOutlet weak var lblUser: UILabel!
    
    @IBOutlet weak var stackView: UIStackView!
    @IBOutlet weak var view1: UIView!
    @IBOutlet weak var collectionView1: UICollectionView!
    @IBOutlet weak var collectionView2: UICollectionView!
    @IBOutlet weak var collectionView3: UICollectionView!
    
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        collectionView1.dataSource = self
        collectionView1.delegate = self

        collectionView2.dataSource = self
        collectionView2.delegate = self

        collectionView3.dataSource = self
        collectionView3.delegate = self
        
        
        collectionView1.reloadData()
        collectionView2.reloadData()
        collectionView3.reloadData()

        
        if let layout2 = collectionView2.collectionViewLayout as? UICollectionViewFlowLayout {
            layout2.scrollDirection = .vertical
        }

        if let layout3 = collectionView3.collectionViewLayout as? UICollectionViewFlowLayout {
            layout3.scrollDirection = .vertical
        }

    }
    
    func collectionView(_ collectionView: UICollectionView,
                        layout collectionViewLayout: UICollectionViewLayout,
                        sizeForItemAt indexPath: IndexPath) -> CGSize {
        
        if collectionView == collectionView1 {
            // 2 items per row (so 2 rows if you have 4 items)
            let padding: CGFloat = 16
            let totalSpacing: CGFloat = padding * 3   // left + right + middle
            let width = (collectionView.frame.width - 18) / 2
            return CGSize(width: width, height: 146)
        }
        
        else if collectionView == collectionView2 {
            // Horizontal scrolling card
            return CGSize(width: collectionView.frame.width - 28,
                          height: 278)
        }
        
        else {
            // Featured policies card
            return CGSize(width: collectionView.frame.width - 28,
                          height: 273)
        }
    }
    
    func collectionView(_ collectionView: UICollectionView,
                        numberOfItemsInSection section: Int) -> Int {
        
        if collectionView == collectionView1 {
            return categories.count
        } else if collectionView == collectionView2 {
            return policies.count
        } else {
            return featuredPolicies.count
        }
    }

    
    func collectionView(_ collectionView: UICollectionView,
                        cellForItemAt indexPath: IndexPath) -> UICollectionViewCell {

        if collectionView == collectionView1 {
            let cell = collectionView.dequeueReusableCell(
                withReuseIdentifier: "InsuranceCategoryCell",
                for: indexPath
            ) as! InsuranceCategoryCell
            
            let category = categories[indexPath.item]
            cell.configure(with: category)
            return cell
        }
        
        else if collectionView == collectionView2 {
            let cell = collectionView.dequeueReusableCell(
                withReuseIdentifier: "PoliciesCell",
                for: indexPath
            ) as! PoliciesCell
            
            let policy = policies[indexPath.item]
            cell.configure(with: policy)
            return cell
        }
        
        else {
            let cell = collectionView.dequeueReusableCell(
                withReuseIdentifier: "FeaturedPoliciesCell",
                for: indexPath
            ) as! FeaturedPoliciesCell
            
            let featured = featuredPolicies[indexPath.item]
            cell.configure(with: featured)
            return cell
        }
    }

    @IBAction func gotoProfile() {
        DispatchQueue.main.async {
            let storyboard = UIStoryboard(name: "Profile", bundle: nil)
            let profileVC = storyboard.instantiateViewController(withIdentifier: "ProfileViewController") as! ProfileViewController
            profileVC.modalPresentationStyle = .fullScreen
            self.present(profileVC, animated: false, completion: nil)
        }
    }
    
    
    
    
}
