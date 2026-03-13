//
//  FeaturedPoliciesCell.swift
//  InsureTech
//
//  Created by LifeplusBD on 15/2/26.
//


import UIKit

class FeaturedPoliciesCell: UICollectionViewCell {
    
    @IBOutlet weak var bannerImageView: UIImageView!
    @IBOutlet weak var titleLabel: UILabel!
    @IBOutlet weak var subtitleLabel: UILabel!
    
    override func awakeFromNib() {
        super.awakeFromNib()
        bannerImageView.layer.cornerRadius = 12
        bannerImageView.clipsToBounds = true
    }
    
    func configure(with model: FeaturedPolicy) {
        bannerImageView.image = UIImage(named: model.bannerName)
        titleLabel.text = model.title
        subtitleLabel.text = model.subtitle
    }
}
