//
//  FaceCaptureViewController.swift
//  InsureTech
//
//  Created by LifeplusBD on 5/1/26.
//

import UIKit
import AVFoundation
import Vision

final class FaceCaptureViewController: UIViewController {

    // MARK: - Camera
    private let captureSession = AVCaptureSession()
    private let videoOutput = AVCaptureVideoDataOutput()
    private var previewLayer: AVCaptureVideoPreviewLayer!

    // MARK: - Completion
    var onFaceCaptured: ((String, UIImage) -> Void)?

    // MARK: - Pose state
    enum FacePoseStep: String {
        case neutral = "Look Straight"
        case left = "Turn Head Left"
        case right = "Turn Head Right"
        case completed = "Face Captured"
    }
    private var currentStep: FacePoseStep = .neutral
    private var faceCompletedTimestamp: Date?
    private let stableDuration: TimeInterval = 0.7
    private var hasCaptured = false

    // MARK: - Vision
    private lazy var faceRequest: VNDetectFaceLandmarksRequest = {
        VNDetectFaceLandmarksRequest { [weak self] request, _ in
            guard let self,
                  let results = request.results as? [VNFaceObservation] else { return }
            DispatchQueue.main.async {
                self.processFaces(results)
            }
        }
    }()
    private var faceLayer: CAShapeLayer?

    // MARK: - UI
    private let instructionLabel: UILabel = {
        let label = UILabel()
        label.text = "Look Straight"
        label.textColor = .white
        label.font = .boldSystemFont(ofSize: 20)
        label.textAlignment = .center
        label.translatesAutoresizingMaskIntoConstraints = false
        label.backgroundColor = UIColor.black.withAlphaComponent(0.4)
        label.layer.cornerRadius = 10
        label.clipsToBounds = true
        return label
    }()

    private let frameSize: CGFloat = 300

    // MARK: - Lifecycle
    override func viewDidLoad() {
        super.viewDidLoad()
        view.backgroundColor = .black
        configureCamera()
        setupOverlay()
        setupInstructionLabel()
    }

    override func viewDidLayoutSubviews() {
        super.viewDidLayoutSubviews()
        previewLayer?.frame = view.bounds
    }

    // MARK: - Camera Setup
    private func configureCamera() {
        captureSession.sessionPreset = .photo
        guard let device = AVCaptureDevice.default(.builtInWideAngleCamera, for: .video, position: .front),
              let input = try? AVCaptureDeviceInput(device: device) else { return }

        captureSession.addInput(input)
        videoOutput.videoSettings = [kCVPixelBufferPixelFormatTypeKey as String: kCVPixelFormatType_32BGRA]
        videoOutput.setSampleBufferDelegate(self, queue: DispatchQueue(label: "vision.queue"))
        captureSession.addOutput(videoOutput)

        if let conn = videoOutput.connection(with: .video) {
            conn.videoOrientation = .portrait
            conn.isVideoMirrored = true
        }

        previewLayer = AVCaptureVideoPreviewLayer(session: captureSession)
        previewLayer.videoGravity = .resizeAspectFill
        view.layer.addSublayer(previewLayer)
        DispatchQueue.global(qos: .userInitiated).async { [weak self] in
            self?.captureSession.startRunning()
        }    }

    // MARK: - Overlay & Instruction
    private func setupOverlay() {
        let overlay = UIView()
        overlay.layer.borderColor = UIColor.systemYellow.cgColor
        overlay.layer.borderWidth = 2
        overlay.frame = CGRect(
            x: (view.bounds.width - frameSize)/2,
            y: (view.bounds.height - frameSize)/2,
            width: frameSize,
            height: frameSize
        )
        overlay.isUserInteractionEnabled = false
        view.addSubview(overlay)
    }

    private func setupInstructionLabel() {
        view.addSubview(instructionLabel)
        NSLayoutConstraint.activate([
            instructionLabel.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -40),
            instructionLabel.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            instructionLabel.widthAnchor.constraint(equalToConstant: 220),
            instructionLabel.heightAnchor.constraint(equalToConstant: 50)
        ])
    }

    // MARK: - Face Processing
    private func processFaces(_ faces: [VNFaceObservation]) {
        guard let face = faces.first else { return }

        let yaw = face.yaw?.doubleValue ?? 0
        let pitch = face.pitch?.doubleValue ?? 0
        let yawDegrees = yaw * 180 / .pi
        let pitchDegrees = pitch * 180 / .pi

        // --- State machine ---
        switch currentStep {
        case .neutral:
            if abs(yawDegrees) < 5 && abs(pitchDegrees) < 5 { currentStep = .left }
        case .left:
            if yawDegrees < -15 { currentStep = .right }
        case .right:
            if yawDegrees > 15 { currentStep = .completed }
        case .completed:
            break
        }

        instructionLabel.text = currentStep.rawValue

        // --- Draw face rectangle ---
        faceLayer?.removeFromSuperlayer()
        let rect = previewLayer.layerRectConverted(fromMetadataOutputRect: face.boundingBox)
        let layer = CAShapeLayer()
        layer.frame = rect
        layer.borderWidth = 2
        layer.cornerRadius = 6
        layer.borderColor = currentStep == .completed ? UIColor.systemGreen.cgColor : UIColor.systemYellow.cgColor
        previewLayer.addSublayer(layer)
        faceLayer = layer

        // --- Automatic capture ---
        if currentStep == .completed {
            if faceCompletedTimestamp == nil {
                faceCompletedTimestamp = Date()
            } else if !hasCaptured, Date().timeIntervalSince(faceCompletedTimestamp!) >= 0.7 {
                hasCaptured = true
                captureCurrentFrame()
            }
        } else {
            faceCompletedTimestamp = nil
        }
    }

    private func captureCurrentFrame() {
        guard let layer = previewLayer else { return }
        captureSession.stopRunning()

        UIGraphicsBeginImageContextWithOptions(layer.bounds.size, false, UIScreen.main.scale)
        layer.render(in: UIGraphicsGetCurrentContext()!)
        let fullImage = UIGraphicsGetImageFromCurrentImageContext()
        UIGraphicsEndImageContext()
        guard let image = fullImage else { return }

        let cropSize = min(image.size.width, image.size.height)
        let cropRect = CGRect(
            x: (image.size.width - cropSize)/2,
            y: (image.size.height - cropSize)/2,
            width: cropSize, height: cropSize
        )
        if let cgImage = image.cgImage?.cropping(to: cropRect) {
            let croppedImage = UIImage(cgImage: cgImage, scale: image.scale, orientation: .up)

            let alert = UIAlertController(title: "Enter Name", message: nil, preferredStyle: .alert)
            alert.addTextField { $0.placeholder = "Name" }
            alert.addAction(UIAlertAction(title: "Save", style: .default, handler: { [weak self] _ in
                guard let name = alert.textFields?.first?.text, !name.isEmpty else { return }
                self?.onFaceCaptured?(name, croppedImage)
                self?.dismiss(animated: true)
            }))
            alert.addAction(UIAlertAction(title: "Cancel", style: .cancel))
            present(alert, animated: true)
        }
    }
}

// MARK: - AVCapture Delegate
extension FaceCaptureViewController: AVCaptureVideoDataOutputSampleBufferDelegate {
    func captureOutput(_ output: AVCaptureOutput,
                       didOutput sampleBuffer: CMSampleBuffer,
                       from connection: AVCaptureConnection) {
        guard let pixelBuffer = CMSampleBufferGetImageBuffer(sampleBuffer) else { return }
        let handler = VNImageRequestHandler(cvPixelBuffer: pixelBuffer,
                                            orientation: .leftMirrored,
                                            options: [:])
        try? handler.perform([faceRequest])
    }
}

