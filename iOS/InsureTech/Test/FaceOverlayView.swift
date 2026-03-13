//
//  FaceOverlayView.swift
//  InsureTech
//
//  Created by LifeplusBD on 20/1/26.
//


import UIKit

final class FaceOverlayView: UIView {

    private var box: [CGFloat]?
    private var isLive = false

    func update(_ response: FLVEResponse) {
        self.box = response.box
        self.isLive = response.isLive ?? false
        setNeedsDisplay()
    }

    override func draw(_ rect: CGRect) {
        guard let ctx = UIGraphicsGetCurrentContext() else { return }

        ctx.clear(rect)

        // Guide oval
        ctx.setStrokeColor(UIColor.white.withAlphaComponent(0.3).cgColor)
        ctx.setLineWidth(3)
        ctx.addEllipse(in: rect.insetBy(dx: 60, dy: 120))
        ctx.strokePath()

        guard let box else { return }

        ctx.setStrokeColor(isLive ? UIColor.systemGreen.cgColor : UIColor.systemRed.cgColor)
        ctx.setLineWidth(2)

        let faceRect = CGRect(
            x: box[0],
            y: box[1],
            width: box[2],
            height: box[3]
        )

        ctx.stroke(faceRect)
    }
}
