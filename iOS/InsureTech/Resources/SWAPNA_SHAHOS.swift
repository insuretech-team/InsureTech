
//
//  SWAPNA_SHAHOS.swift
//
//  Created by Joy Biswas Shahos - joybiswas101626@gmail.com +8801689639998 on 4/12/22.
//

import UIKit

public extension UIView {
    
    func showText(message: String, duration: TimeInterval = 3.0) {
        let toastContainer = UIView()
        toastContainer.backgroundColor = UIColor.black.withAlphaComponent(0.8)
        toastContainer.layer.cornerRadius = 8.0
        toastContainer.clipsToBounds = true

        let messageLabel = UILabel()
        messageLabel.text = message
        messageLabel.textColor = .white
        messageLabel.numberOfLines = 0
        messageLabel.textAlignment = .center
        messageLabel.font = UIFont.systemFont(ofSize: 14)
        messageLabel.translatesAutoresizingMaskIntoConstraints = false

        toastContainer.addSubview(messageLabel)
        self.addSubview(toastContainer)

        toastContainer.translatesAutoresizingMaskIntoConstraints = false

        NSLayoutConstraint.activate([
            toastContainer.centerXAnchor.constraint(equalTo: self.centerXAnchor),
            toastContainer.bottomAnchor.constraint(equalTo: self.bottomAnchor, constant: -120),
            toastContainer.widthAnchor.constraint(lessThanOrEqualToConstant: 300),
            messageLabel.leadingAnchor.constraint(equalTo: toastContainer.leadingAnchor, constant: 20),
            messageLabel.trailingAnchor.constraint(equalTo: toastContainer.trailingAnchor, constant: -20),
            messageLabel.topAnchor.constraint(equalTo: toastContainer.topAnchor, constant: 10),
            messageLabel.bottomAnchor.constraint(equalTo: toastContainer.bottomAnchor, constant: -10)
        ])

        toastContainer.alpha = 0.0
        UIView.animate(withDuration: 0.5, animations: {
            toastContainer.alpha = 1.0
        }) { _ in
            UIView.animate(withDuration: 0.5, delay: duration, options: [], animations: {
                toastContainer.alpha = 0.0
            }) { _ in
                toastContainer.removeFromSuperview()
            }
        }
    }

    func showToast(message: String, duration: TimeInterval = 3.0) {
        let toastContainer = UIView()
        toastContainer.backgroundColor = UIColor.black.withAlphaComponent(0.8)
        toastContainer.layer.cornerRadius = 8.0
        toastContainer.clipsToBounds = true
        
        guard let image = UIImage(named: "success") else {
            print("Image 'congratulation' not found")
            return
        }
        
        let imageView = UIImageView(image: image)
        imageView.contentMode = .scaleAspectFit
        imageView.translatesAutoresizingMaskIntoConstraints = false
        
        let messageLabel = UILabel()
        messageLabel.text = message
        messageLabel.textColor = .white
        messageLabel.numberOfLines = 0
        messageLabel.font = UIFont.systemFont(ofSize: 14)
        messageLabel.translatesAutoresizingMaskIntoConstraints = false
        
        toastContainer.addSubview(imageView)
        toastContainer.addSubview(messageLabel)
        
        self.addSubview(toastContainer)
        
        toastContainer.translatesAutoresizingMaskIntoConstraints = false
        
        NSLayoutConstraint.activate([
            toastContainer.centerXAnchor.constraint(equalTo: self.centerXAnchor),
            toastContainer.bottomAnchor.constraint(equalTo: self.bottomAnchor, constant: -120),
            toastContainer.widthAnchor.constraint(lessThanOrEqualToConstant: 300),
            toastContainer.heightAnchor.constraint(greaterThanOrEqualToConstant: 50)
        ])
        
        NSLayoutConstraint.activate([
            imageView.leadingAnchor.constraint(equalTo: toastContainer.leadingAnchor, constant: 20),
            imageView.centerYAnchor.constraint(equalTo: toastContainer.centerYAnchor),
            imageView.widthAnchor.constraint(equalToConstant: 24),
            imageView.heightAnchor.constraint(equalToConstant: 24)
        ])
        
        NSLayoutConstraint.activate([
            messageLabel.leadingAnchor.constraint(equalTo: imageView.trailingAnchor, constant: 20),
            messageLabel.centerYAnchor.constraint(equalTo: toastContainer.centerYAnchor),
            messageLabel.trailingAnchor.constraint(equalTo: toastContainer.trailingAnchor, constant: -20),
            messageLabel.topAnchor.constraint(greaterThanOrEqualTo: toastContainer.topAnchor, constant: 10),
            messageLabel.bottomAnchor.constraint(lessThanOrEqualTo: toastContainer.bottomAnchor, constant: -10)
        ])
        
        toastContainer.alpha = 0.0
        UIView.animate(withDuration: 0.5, animations: {
            toastContainer.alpha = 1.0
        }) { _ in
            UIView.animate(withDuration: 0.5, delay: duration, options: [], animations: {
                toastContainer.alpha = 0.0
            }) { _ in
                toastContainer.removeFromSuperview()
            }
        }
    }
    
    func showWarningToast(message: String, duration: TimeInterval = 3.0) {
        let toastContainer = UIView()
        toastContainer.backgroundColor = UIColor.black.withAlphaComponent(0.8)
        toastContainer.layer.cornerRadius = 8.0
        toastContainer.clipsToBounds = true
        
        guard let image = UIImage(named: "warning") else {
            print("Image 'warning' not found")
            return
        }
        
        let imageView = UIImageView(image: image)
        imageView.contentMode = .scaleAspectFit
        imageView.translatesAutoresizingMaskIntoConstraints = false
        
        let messageLabel = UILabel()
        messageLabel.text = message
        messageLabel.textColor = .white
        messageLabel.numberOfLines = 0
        messageLabel.font = UIFont.systemFont(ofSize: 14)
        messageLabel.translatesAutoresizingMaskIntoConstraints = false
        
        toastContainer.addSubview(imageView)
        toastContainer.addSubview(messageLabel)
        
        self.addSubview(toastContainer)
        
        toastContainer.translatesAutoresizingMaskIntoConstraints = false
        
        NSLayoutConstraint.activate([
            toastContainer.centerXAnchor.constraint(equalTo: self.centerXAnchor),
            toastContainer.bottomAnchor.constraint(equalTo: self.bottomAnchor, constant: -120),
            toastContainer.widthAnchor.constraint(lessThanOrEqualToConstant: 300),
            toastContainer.heightAnchor.constraint(greaterThanOrEqualToConstant: 50)
        ])
        
        NSLayoutConstraint.activate([
            imageView.leadingAnchor.constraint(equalTo: toastContainer.leadingAnchor, constant: 20),
            imageView.centerYAnchor.constraint(equalTo: toastContainer.centerYAnchor),
            imageView.widthAnchor.constraint(equalToConstant: 24),
            imageView.heightAnchor.constraint(equalToConstant: 24)
        ])
        
        NSLayoutConstraint.activate([
            messageLabel.leadingAnchor.constraint(equalTo: imageView.trailingAnchor, constant: 20),
            messageLabel.centerYAnchor.constraint(equalTo: toastContainer.centerYAnchor),
            messageLabel.trailingAnchor.constraint(equalTo: toastContainer.trailingAnchor, constant: -20),
            messageLabel.topAnchor.constraint(greaterThanOrEqualTo: toastContainer.topAnchor, constant: 10),
            messageLabel.bottomAnchor.constraint(lessThanOrEqualTo: toastContainer.bottomAnchor, constant: -10)
        ])
        
        toastContainer.alpha = 0.0
        UIView.animate(withDuration: 0.5, animations: {
            toastContainer.alpha = 1.0
        }) { _ in
            UIView.animate(withDuration: 0.5, delay: duration, options: [], animations: {
                toastContainer.alpha = 0.0
            }) { _ in
                toastContainer.removeFromSuperview()
            }
        }
    }
}
