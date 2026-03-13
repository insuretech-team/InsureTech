//
//  ViewController.swift
//  oneTapService
//
//  Created by Joy Biswas - joybiswas101626@gmail.com 01689639998 on 2/12/24.
//

import UIKit
import SVProgressHUD
import SWAPNA_SHAHOS

class LoginViewController: UIViewController,  UITextFieldDelegate {
    
    
    @IBOutlet weak var textFieldView: UIView!
    @IBOutlet weak var view1: UIView!
    @IBOutlet weak var view2: UIView!
    @IBOutlet weak var lblPhone1: UILabel!
    @IBOutlet weak var tfPhone: UITextField!
    @IBOutlet weak var btnTextField: UIButton!
    @IBOutlet weak var btnGetCode: UIButton!

    override func viewDidLoad() {
        super.viewDidLoad()
        updateUI()
        
        let placeholderText = "01689639998"
        let attributes: [NSAttributedString.Key: Any] = [
            .foregroundColor: UIColor(named: "AEAEAE") ?? UIColor.gray,
            .font: UIFont.systemFont(ofSize: 14)
        ]
        tfPhone.attributedPlaceholder = NSAttributedString(string: placeholderText, attributes: attributes)
        tfPhone.delegate = self
    }

    func updateUI() {
        btnGetCode.layer.cornerRadius = 8
        
        updateButton()
        
        self.textFieldView.layer.cornerRadius = 8
        self.textFieldView.layer.borderWidth = 1
        self.textFieldView.layer.borderColor = UIColor(named: "AEAEAE")?.cgColor
        
        self.view2.isHidden = true
    }
    
    @objc func doneButtonTapped() {
        tfPhone.resignFirstResponder()
    }
    
    func textField(_ textField: UITextField, shouldChangeCharactersIn range: NSRange, replacementString string: String) -> Bool {
        let allowedCharacters = CharacterSet.decimalDigits
        let characterSet = CharacterSet(charactersIn: string)
        if !allowedCharacters.isSuperset(of: characterSet) {
            return false
        }
        
        let currentText = textField.text ?? ""
        let newText = (currentText as NSString).replacingCharacters(in: range, with: string)
        let digitsOnly = newText.replacingOccurrences(of: "\\D", with: "", options: .regularExpression)
        
        if digitsOnly.count > 11 {
            return false
        }
        
        let formattedText = formatPhoneNumber(digitsOnly)
        textField.text = formattedText
        
        updateButton()
        
        return false
    }
    
    func updateButton() {
        if tfPhone.text!.count == 13 {
            btnGetCode.backgroundColor = UIColor(named: "8C34C7")
            btnGetCode.setTitleColor(UIColor(named: "FFFFFF"), for: .normal)
            textFieldView.layer.borderColor = UIColor(named: "AEAEAE")?.cgColor
            doneButtonTapped()
        } else {
            btnGetCode.backgroundColor = UIColor(named: "EFEFEF")
            btnGetCode.setTitleColor(UIColor(named: "AEAEAE"), for: .normal)
        }
    }
    
    func formatPhoneNumber(_ number: String) -> String {
        var formattedString = ""
        let maxDigits = 11
        let digits = number.prefix(maxDigits)
        let firstGroup = digits.prefix(3)
        formattedString += firstGroup
        
        let secondGroup = digits.dropFirst(3).prefix(4)
        if !secondGroup.isEmpty {
            formattedString += " " + secondGroup
        }
        
        let thirdGroup = digits.dropFirst(7).prefix(4)
        if !thirdGroup.isEmpty {
            formattedString += " " + thirdGroup
        }
        return formattedString
    }
    
    @IBAction func textBtnTapped(_ sender: UIButton) {
        btnTextField.isHidden = true
        view1.isHidden = true
        view2.isHidden = false
        
        tfPhone.becomeFirstResponder()
    }
    
    @IBAction func getCodeTapped() {
        
        let numberWithSpaces = tfPhone.text ?? ""
        let number = numberWithSpaces.replacingOccurrences(of: " ", with: "")
        
        if number.count == 0 {
            print("empty")
            view.showWarningToast(message: "Please, enter your phone number.")
            textFieldView.layer.borderColor = UIColor(named: "FF383C")?.cgColor
        } else if number.count != 11 {
            print("not 11")
            view.showWarningToast(message: "Please, enter your 11 digits phone number.")
            textFieldView.layer.borderColor = UIColor(named: "FF383C")?.cgColor
        } else {
            if ["013", "014", "015", "016", "017", "018", "019"].contains(String(number.prefix(3))) {
                print("The first three digits are valid.")
//                SVProgressHUD.show()
//                NetworkManager.shared.callLogin(phone: number) { result in
//                    SVProgressHUD.dismiss()
//                    switch result {
//                    case .success(let loginModel):
//                        if let data = loginModel.data {
                            UserDefaults.standard.set(number, forKey: "phoneNumber")
////                            UserDefaults.standard.set(data.otp, forKey: "OTPNumber")
//                        }
                        self.gotoOTPView()
//                        print("Login successful: \(loginModel.message ?? "")")
//                    case .failure(let error):
//                        print("Login failed: \(error.localizedDescription)")
//                    }
//                }
            } else {
                print("The first three digits are not valid.")
                textFieldView.layer.borderColor = UIColor(named: "FF383C")?.cgColor
                view.showWarningToast(message: "Please, enter a valid number.")
            }
        }
    }
    
    func gotoOTPView() {
        DispatchQueue.main.async {
            let storyboard = UIStoryboard(name: "Main", bundle: nil)
            let OTPVC = storyboard.instantiateViewController(withIdentifier: "OTPViewController") as! OTPViewController
            OTPVC.modalPresentationStyle = .fullScreen
            self.present(OTPVC, animated: false, completion: nil)
        }
    }

    @IBAction func skipTapped() {
//        UserDefaults.standard.set(true, forKey: "isGuest")
//        DispatchQueue.main.async {
//            let storyboard = UIStoryboard(name: "Main", bundle: nil)
//            let tabBarVC = storyboard.instantiateViewController(withIdentifier: "CustomTabBarViewController") as! CustomTabBarViewController
//            tabBarVC.modalPresentationStyle = .fullScreen
//            self.present(tabBarVC, animated: false, completion: nil)
//        }
    }
}

