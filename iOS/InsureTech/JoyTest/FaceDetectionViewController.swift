
import UIKit
import AVFoundation
import Vision

final class FaceDetectionViewController: UIViewController {

    // MARK: - Camera
    private let captureSession = AVCaptureSession()
    private let videoOutput = AVCaptureVideoDataOutput()
    private var previewLayer: AVCaptureVideoPreviewLayer!

    // MARK: - Vision
    private let visionQueue = DispatchQueue(label: "vision.queue")
    private var faceLayers: [CAShapeLayer] = []

    private lazy var faceRequest: VNDetectFaceLandmarksRequest = {
        VNDetectFaceLandmarksRequest { [weak self] request, _ in
            guard let self,
                  let results = request.results as? [VNFaceObservation]
            else { return }

            DispatchQueue.main.async {
                self.processFaces(results)
            }
        }
    }()

    // MARK: - Pose State
    enum FacePoseStep {
        case neutral, left, right, up, down, completed
    }

    private var currentStep: FacePoseStep = .neutral

    private func updatePoseStep(yaw: Double, pitch: Double) {
        switch currentStep {
        case .neutral:
            if abs(yaw) < 5 && abs(pitch) < 5 {
                currentStep = .left
                print("✅ Neutral verified, move left")
            }
        case .left:
            if yaw < -15 {
                currentStep = .right
                print("✅ Left verified, move right")
            }
        case .right:
            if yaw > 15 {
                currentStep = .up
                print("✅ Right verified, look up")
            }
        case .up:
            if pitch > 12 {
                currentStep = .down
                print("✅ Up verified, look down")
            }
        case .down:
            if pitch < -12 {
                currentStep = .completed
                print("✅ Full face movement completed!")
                // Optional: stop capture session or take snapshot
                captureSession.stopRunning()
            }
        case .completed:
            break
        }
    }

    // MARK: - Lifecycle
    override func viewDidLoad() {
        super.viewDidLoad()
        view.backgroundColor = .black
        configureCamera()
    }

    override func viewDidLayoutSubviews() {
        super.viewDidLayoutSubviews()
        previewLayer?.frame = view.bounds
    }
}

// MARK: - Camera Configuration
private extension FaceDetectionViewController {

    func configureCamera() {
        captureSession.sessionPreset = .high

        guard
            let device = AVCaptureDevice.default(.builtInWideAngleCamera,
                                                 for: .video,
                                                 position: .front),
            let input = try? AVCaptureDeviceInput(device: device)
        else {
            return
        }

        captureSession.addInput(input)

        videoOutput.videoSettings = [
            kCVPixelBufferPixelFormatTypeKey as String:
                kCVPixelFormatType_32BGRA
        ]

        videoOutput.setSampleBufferDelegate(self, queue: visionQueue)
        captureSession.addOutput(videoOutput)

        if let connection = videoOutput.connection(with: .video) {
            connection.videoOrientation = .portrait
            connection.isVideoMirrored = true
        }

        previewLayer = AVCaptureVideoPreviewLayer(session: captureSession)
        previewLayer.videoGravity = .resizeAspectFill
        view.layer.addSublayer(previewLayer)

        captureSession.startRunning()
    }
}

// MARK: - Vision Processing
extension FaceDetectionViewController: AVCaptureVideoDataOutputSampleBufferDelegate {

    func captureOutput(_ output: AVCaptureOutput,
                       didOutput sampleBuffer: CMSampleBuffer,
                       from connection: AVCaptureConnection) {

        guard let pixelBuffer = CMSampleBufferGetImageBuffer(sampleBuffer) else {
            return
        }

        let handler = VNImageRequestHandler(cvPixelBuffer: pixelBuffer,
                                            orientation: .leftMirrored,
                                            options: [:])

        try? handler.perform([faceRequest])
    }

    private func processFaces(_ faces: [VNFaceObservation]) {

        faceLayers.forEach { $0.removeFromSuperlayer() }
        faceLayers.removeAll()

        guard let face = faces.first else { return }

        let yaw = face.yaw?.doubleValue ?? 0      // left / right
        let pitch = face.pitch?.doubleValue ?? 0  // up / down

        let yawDegrees = yaw * 180 / .pi
        let pitchDegrees = pitch * 180 / .pi

        // Advance pose verification state
        updatePoseStep(yaw: yawDegrees, pitch: pitchDegrees)

        // Draw face box
        let rect = previewLayer.layerRectConverted(
            fromMetadataOutputRect: face.boundingBox
        )

        let layer = CAShapeLayer()
        layer.frame = rect
        layer.borderWidth = 2
        layer.cornerRadius = 6

        // Color by step completion
        switch currentStep {
        case .completed:
            layer.borderColor = UIColor.systemGreen.cgColor
        default:
            layer.borderColor = UIColor.systemYellow.cgColor
        }

        previewLayer.addSublayer(layer)
        faceLayers.append(layer)
    }
}
