//
//  OTPViewController.swift
//  oneTapService
//
//  Created by Joy Biswas - joybiswas101626@gmail.com 01689639998 on 7/1/25.
//

import UIKit
import Foundation

class OTPViewController: UIViewController, UITextFieldDelegate {
    
    var phoneNumber: String = ""
    private var countdownTimer: Timer?
    private var remainingSeconds: Int = 60
    
    @IBOutlet weak var otpView1:UIView!
    @IBOutlet weak var otpView2:UIView!
    @IBOutlet weak var otpView3:UIView!
    @IBOutlet weak var otpView4:UIView!
    @IBOutlet weak var otpText1:UITextField!
    @IBOutlet weak var otpText2:UITextField!
    @IBOutlet weak var otpText3:UITextField!
    @IBOutlet weak var otpText4:UITextField!
    
    @IBOutlet weak var lblNumber: UILabel!
    @IBOutlet weak var lblTimer: UILabel!
    @IBOutlet weak var lblStatus: UILabel!
    @IBOutlet weak var btnVerify: UIButton!
    @IBOutlet weak var resendView: UIView!
        
    override func viewDidLoad() {
        super.viewDidLoad()

        lblTimer.isHidden = false
        
        DispatchQueue.main.asyncAfter(deadline: .now() + 5) { [self] in
            [otpText1, otpText2, otpText3, otpText4].forEach {
                $0.delegate = self
                $0.keyboardType = .numberPad
                $0.textContentType = .oneTimeCode
                $0.isSecureTextEntry = false
            }
        }
        
        otpText1.returnKeyType = UIReturnKeyType.next
        otpText2.returnKeyType = UIReturnKeyType.next
        otpText3.returnKeyType = UIReturnKeyType.next
        otpText4.returnKeyType = UIReturnKeyType.done
        otpText1.delegate = self
        otpText2.delegate = self
        otpText3.delegate = self
        otpText4.delegate = self
        
        [otpView1, otpView2, otpView3, otpView4].forEach {
            setInactiveBorder(for: $0)
        }

        let userNumber = UserDefaults.standard.string(forKey: "phoneNumber")
        lblNumber.text = "A 4 digit code has been sent to your no. +88\(userNumber ?? "")"

        let tap = UITapGestureRecognizer(target: self, action: #selector(UIInputViewController.dismissKeyboard))
        view.addGestureRecognizer(tap)
        
        btnVerify.layer.cornerRadius = 8
        btnVerify.clipsToBounds = true
        
        btnVerify.backgroundColor = UIColor(named: "EFEFEF")
        btnVerify.setTitleColor(UIColor(named: "AEAEAE"), for: .normal)
        
        lblStatus.isHidden = true
        resendView.isHidden = true
        startTimer(seconds: 60)
    }
    
    private func startTimer(seconds: Int) {
        countdownTimer?.invalidate()
        remainingSeconds = seconds

        lblTimer.isHidden = false
        resendView.isHidden = true

        updateTimerLabel()

        countdownTimer = Timer.scheduledTimer(withTimeInterval: 1, repeats: true) { [weak self] _ in
            guard let self = self else { return }

            self.remainingSeconds -= 1
            self.updateTimerLabel()

            if self.remainingSeconds <= 0 {
                self.countdownTimer?.invalidate()
                self.countdownTimer = nil

                self.lblTimer.isHidden = true
                self.resendView.isHidden = false
                lblStatus.text = "Didn’t receive the code?"
                lblStatus.isHidden = false
            }
        }
    }

    private func updateTimerLabel() {
        let minutes = remainingSeconds / 60
        let seconds = remainingSeconds % 60
        lblTimer.text = String(format: "%02dm : %02ds", minutes, seconds)
    }
    
//    @objc func textFieldDidChange(_ textField: UITextField) {
//        if let text = textField.text, text.count >= 1 {
//            switch textField {
//            case otpText1: otpText2.becomeFirstResponder()
//            case otpText2: otpText3.becomeFirstResponder()
//            case otpText3: otpText4.becomeFirstResponder()
//            case otpText4:
//                otpText4.resignFirstResponder()
//                verifyTapped()
//            default: break
//            }
//        }
//    }

    @objc func dismissKeyboard() {
        view.endEditing(true)
    }

    override func viewWillLayoutSubviews() {
        [otpView1, otpView2, otpView3, otpView4].forEach {
            $0.layer.cornerRadius = 8
            $0.layer.borderWidth = 1
        }

        otpText1.attributedPlaceholder = NSAttributedString(
            string: "_",
            attributes: [NSAttributedString.Key.foregroundColor: UIColor(named: "565656") as Any]
        )
        
        otpText2.attributedPlaceholder = NSAttributedString(
            string: "_",
            attributes: [NSAttributedString.Key.foregroundColor: UIColor(named: "565656") as Any]
        )
        
        otpText3.attributedPlaceholder = NSAttributedString(
            string: "_",
            attributes: [NSAttributedString.Key.foregroundColor: UIColor(named: "565656") as Any]
        )

        otpText4.attributedPlaceholder = NSAttributedString(
            string: "_",
            attributes: [NSAttributedString.Key.foregroundColor: UIColor(named: "565656") as Any]
        )
    }
    
    @IBAction func verifyTapped(){

        let captcha = otpText1.text! + otpText2.text! + otpText3.text! + otpText4.text!
        
        if captcha.count != 4 {
            view.showWarningToast(message: "Enter 4 digit OTP.")
            
            return
        } else if captcha != "1122" {
            otpView1.layer.borderColor = UIColor(named: "FF383C")?.cgColor
            otpView2.layer.borderColor = UIColor(named: "FF383C")?.cgColor
            otpView3.layer.borderColor = UIColor(named: "FF383C")?.cgColor
            otpView4.layer.borderColor = UIColor(named: "FF383C")?.cgColor
//            lblStatus.isHidden = false
//            lblStatus.text = "Invalid OTP by Default"
//            lblStatus.textColor = UIColor(named: "FF383C")
            
            view.showToast(message: "Invalid OTP by Default")
            
            return
        }

//        let userNumber = UserDefaults.standard.string(forKey: "phoneNumber")
//        SVProgressHUD.show()
//        NetworkManager.shared.callVerifyOtp(phone: userNumber!, otp: Int(captcha)!) { result in
//            SVProgressHUD.dismiss()
//            switch result {
//            case .success(let otpModel):
//                print("OTP Verified, Token: \(otpModel.data?.token ?? "")")
//                
//                let success = otpModel.success
//                if success == true {
//                    let tokenType = otpModel.data?.tokenType
//                    let userToken = otpModel.data?.token
//                    
//                    UserDefaults.standard.set(tokenType, forKey: "tokenType")
//                    UserDefaults.standard.set(userToken, forKey: "userToken")
                    view.showToast(message: "Valid OTP by Default")
        DispatchQueue.main.asyncAfter(deadline: .now() + 4.0) {
            self.gotoHome()
        }
//                } else {
//                    self.view.showWarningToast(message: otpModel.message ?? "Wrong Input", duration: 3.0)
//                }
//            case .failure(let error):
//                print("OTP Verification failed: \(error.localizedDescription)")
//                self.view.showWarningToast(message: "Wrong Input", duration: 3.0)
//            }
//        }
    }
    
    func gotoHome(){
        DispatchQueue.main.async {
            let storyboard = UIStoryboard(name: "Home", bundle: nil)
            let homeVC = storyboard.instantiateViewController(withIdentifier: "HomeViewController") as! HomeViewController
            homeVC.modalPresentationStyle = .fullScreen
            self.present(homeVC, animated: false, completion: nil)
        }
        
        
//        let faceVC = FaceDetectionViewController()
//        faceVC.modalPresentationStyle = .fullScreen // optional but recommended
//        present(faceVC, animated: true)
        

    }

    @IBAction func resendTapped() {
        resendView.isHidden = true
        lblStatus.text = "Code sent again."
        startTimer(seconds: 300)
//        let userNumber = UserDefaults.standard.string(forKey: "phoneNumber")
//        SVProgressHUD.show()
//        NetworkManager.shared.callLogin(phone: userNumber!) { result in
//            SVProgressHUD.dismiss()
//            switch result {
//            case .success(let loginModel):
//                print("Login successful: \(loginModel.message ?? "")")
//                self.view.showToast(message: loginModel.message!, duration: 3.0)
//            case .failure(let error):
//                print("Login failed: \(error.localizedDescription)")
//            }
//        }
    }
    
    @IBAction func backTapped() {
        self.dismiss(animated: true)
    }
    
    func textField(_ textField: UITextField,
                   shouldChangeCharactersIn range: NSRange,
                   replacementString string: String) -> Bool {

        // BACKSPACE
        if string.isEmpty {
            textField.text = ""

            switch textField {
            case otpText4:
                setInactiveBorder(for: otpView4)
                otpText3.becomeFirstResponder()
            case otpText3:
                setInactiveBorder(for: otpView3)
                otpText2.becomeFirstResponder()
            case otpText2:
                setInactiveBorder(for: otpView2)
                otpText1.becomeFirstResponder()
            case otpText1:
                setInactiveBorder(for: otpView1)
            default:
                break
            }
            
            updateVerifyButtonState()

            return false
        }

        guard CharacterSet.decimalDigits.isSuperset(of: CharacterSet(charactersIn: string)) else {
            return false
        }

        textField.text = string

        switch textField {
        case otpText1:
            setActiveBorder(for: otpView1)
            otpText2.becomeFirstResponder()
        case otpText2:
            setActiveBorder(for: otpView2)
            otpText3.becomeFirstResponder()
        case otpText3:
            setActiveBorder(for: otpView3)
            otpText4.becomeFirstResponder()
        case otpText4:
            setActiveBorder(for: otpView4)
            otpText4.resignFirstResponder()
//            verifyTapped()
            
        default:
            break
        }

        updateVerifyButtonState()

        return false
    }

    private func updateVerifyButtonState() {
        let otp = (otpText1.text ?? "") + (otpText2.text ?? "") + (otpText3.text ?? "") + (otpText4.text ?? "")
        if otp.count == 4 {
            btnVerify.backgroundColor = UIColor(named: "8C34C7")
            btnVerify.setTitleColor(UIColor(named: "FFFFFF"), for: .normal)
        } else {
            btnVerify.backgroundColor = UIColor(named: "EFEFEF")
            btnVerify.setTitleColor(UIColor(named: "AEAEAE"), for: .normal)
        }
    }

    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        if(textField == self.otpText1){
            self.otpText1.resignFirstResponder()
            self.otpText2.becomeFirstResponder()
        }else if(textField == self.otpText2){
            self.otpText2.resignFirstResponder()
            self.otpText3.becomeFirstResponder()
        }else if(textField == self.otpText3){
            self.otpText3.resignFirstResponder()
            self.otpText4.becomeFirstResponder()
        }else if(textField == self.otpText4){
            self.otpText4.resignFirstResponder()
//            self.verifyTapped()
        }
        
        updateVerifyButtonState()

        return true
    }
    
    func setActiveBorder(for view: UIView) {
        view.layer.borderColor = UIColor(named: "322F70")?.cgColor
    }

    func setInactiveBorder(for view: UIView) {
        view.layer.borderColor = UIColor(named: "AEAEAE")?.cgColor
    }
    
    deinit {
        countdownTimer?.invalidate()
    }
}


