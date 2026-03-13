//
//  InsuranceCategoryCell.swift
//  InsureTech
//
//  Created by LifeplusBD on 15/2/26.
//


import UIKit

class InsuranceCategoryCell: UICollectionViewCell {
    
    @IBOutlet weak var backView: UIView!
    @IBOutlet weak var iconImageView: UIImageView!
    @IBOutlet weak var nameLabel: UILabel!
    
    override func awakeFromNib() {
        super.awakeFromNib()
        backView.layer.cornerRadius = 12
        backView.clipsToBounds = true
    }
    
    func configure(with model: InsuranceCategory) {
        iconImageView.image = UIImage(systemName: model.iconName)
        nameLabel.text = model.name
    }
}
