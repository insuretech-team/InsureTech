//
//  FLVETestVC.swift
//  InsureTech
//
//  Created by LifeplusBD on 20/1/26.
//


import UIKit
import AVFoundation

final class FLVETestVC: UIViewController {

    private let session = AVCaptureSession()
    private let output = AVCaptureVideoDataOutput()
    private let previewLayer = AVCaptureVideoPreviewLayer()
    private let captureQueue = DispatchQueue(label: "camera.capture.queue")

    private let overlayView = FaceOverlayView()
    private let api = FLVEAPIClient()

    private var lastSent = Date.distantPast
    private let frameInterval: TimeInterval = 0.4

    override func viewDidLoad() {
        super.viewDidLoad()
        view.backgroundColor = .black
        setupCamera()
        setupOverlay()
    }

    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(animated)
        
        // Start session on background queue
        captureQueue.async { [weak self] in
            guard let self = self else { return }
            if !self.session.isRunning {
                self.session.startRunning()
            }
        }
    }


    override func viewWillDisappear(_ animated: Bool) {
        super.viewWillDisappear(animated)
        
        captureQueue.async { [weak self] in
            self?.session.stopRunning()
        }
    }
}

private extension FLVETestVC {

    func setupCamera() {
        session.beginConfiguration()
        session.sessionPreset = .high

        guard
            let device = AVCaptureDevice.default(.builtInWideAngleCamera,
                                                 for: .video,
                                                 position: .front),
            let input = try? AVCaptureDeviceInput(device: device)
        else { return }

        session.addInput(input)

        output.setSampleBufferDelegate(self,
                                       queue: DispatchQueue(label: "flve.camera.queue"))
        session.addOutput(output)

        previewLayer.session = session
        previewLayer.videoGravity = .resizeAspectFill
        previewLayer.frame = view.bounds

        view.layer.addSublayer(previewLayer)
        session.commitConfiguration()
    }

    func setupOverlay() {
        overlayView.frame = view.bounds
        overlayView.isUserInteractionEnabled = false
        view.addSubview(overlayView)
    }
}

extension FLVETestVC: AVCaptureVideoDataOutputSampleBufferDelegate {

    func captureOutput(_ output: AVCaptureOutput,
                       didOutput sampleBuffer: CMSampleBuffer,
                       from connection: AVCaptureConnection) {

        guard Date().timeIntervalSince(lastSent) > frameInterval else { return }
        lastSent = Date()

        guard let jpeg = sampleBuffer.toJPEGData() else { return }

        api.submitFrame(jpeg) { [weak self] result in
            DispatchQueue.main.async {
                self?.overlayView.update(result)
            }
        }
    }
}
