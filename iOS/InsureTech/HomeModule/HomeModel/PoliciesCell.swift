//
//  PoliciesCell.swift
//  InsureTech
//
//  Created by LifeplusBD on 15/2/26.
//


import UIKit

class PoliciesCell: UICollectionViewCell {
    
    @IBOutlet weak var bannerImageView: UIImageView!
    @IBOutlet weak var policyNameLabel: UILabel!
    @IBOutlet weak var coverage1Label: UILabel!
    @IBOutlet weak var coverage2Label: UILabel!
    @IBOutlet weak var iconImageView: UIImageView!
    
    override func awakeFromNib() {
        super.awakeFromNib()
        bannerImageView.layer.cornerRadius = 12
        bannerImageView.clipsToBounds = true
    }
    
    func configure(with model: Policy) {
        bannerImageView.image = UIImage(named: model.bannerName)
        policyNameLabel.text = model.policyName
        coverage1Label.text = model.coverageType1
        coverage2Label.text = model.coverageType2
        iconImageView.image = UIImage(systemName: model.iconName)
    }
}
