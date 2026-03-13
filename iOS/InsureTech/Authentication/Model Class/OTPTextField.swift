//
//  OTPTextField.swift
//  InsureTech
//
//  Created by LifeplusBD on 4/1/26.
//


import UIKit

final class OTPTextField: UITextField {

    override func caretRect(for position: UITextPosition) -> CGRect {
        return .zero
    }

    override func selectionRects(for range: UITextRange) -> [UITextSelectionRect] {
        return []
    }

    override func canPerformAction(_ action: Selector, withSender sender: Any?) -> Bool {
        return false
    }
}
