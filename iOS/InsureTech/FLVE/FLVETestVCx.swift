//
//
//import UIKit
//import AVFoundation
//
//final class FLVETestVCx: UIViewController {
//
//    private let cameraContainer = UIView()
//
//    private let session = AVCaptureSession()
//    private let photoOutput = AVCapturePhotoOutput()
//    private var previewLayer: AVCaptureVideoPreviewLayer!
//
//    
//    private var lastUploadedImageSize: CGSize?
//
//    private var captureCount = 0
//    private let maxCaptures = 8
//
//    private var confidenceValues: [Double] = []
//    private var livenessValues: [Double] = []
//
//    private var isCapturingSequence = false
//    
//    private let healthButton = UIButton(type: .system)
//    private let debugHealthButton = UIButton(type: .system)
//
//    private let resultTextView = UITextView()
//
//    private let captureButton = UIButton(type: .system)
//
//    override func viewDidLoad() {
//        super.viewDidLoad()
//        view.backgroundColor = .black
//        setupCamera()
//        setupUI()
//    }
//    
//    private let faceBoxLayer: CAShapeLayer = {
//        let layer = CAShapeLayer()
//        layer.lineWidth = 3
//        layer.fillColor = UIColor.clear.cgColor
//        layer.isHidden = true
//        return layer
//    }()
//
//}
//
//extension FLVETestVCx {
//    private func setupCamera() {
//        session.beginConfiguration()
//        session.sessionPreset = .photo
//
//        guard
//            let device = AVCaptureDevice.default(.builtInWideAngleCamera,
//                                                 for: .video,
//                                                 position: .front),
//            let input = try? AVCaptureDeviceInput(device: device)
//        else {
//            fatalError("Front camera not available")
//        }
//
//        if session.canAddInput(input) {
//            session.addInput(input)
//        }
//
//        if session.canAddOutput(photoOutput) {
//            session.addOutput(photoOutput)
//        }
//
//        session.commitConfiguration()
//
//        previewLayer = AVCaptureVideoPreviewLayer(session: session)
//        previewLayer.videoGravity = .resizeAspectFill
//        previewLayer.frame = view.bounds
//
//        session.startRunning()
//    }
//}
//
////extension FLVETestVCx {
////    private func setupUI() {
////        captureButton.setTitle("Capture & Liveness", for: .normal)
////        captureButton.backgroundColor = .systemBlue
////        captureButton.tintColor = .white
////        captureButton.layer.cornerRadius = 10
////        captureButton.translatesAutoresizingMaskIntoConstraints = false
////        captureButton.addTarget(self, action: #selector(captureTapped), for: .touchUpInside)
////
////        view.addSubview(captureButton)
////
////        NSLayoutConstraint.activate([
////            captureButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -30),
////            captureButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
////            captureButton.widthAnchor.constraint(equalToConstant: 220),
////            captureButton.heightAnchor.constraint(equalToConstant: 40)
////        ])
////        
////        // Health Button
////        healthButton.setTitle("Health Check", for: .normal)
////        healthButton.backgroundColor = .systemGreen
////        healthButton.tintColor = .white
////        healthButton.layer.cornerRadius = 10
////        healthButton.translatesAutoresizingMaskIntoConstraints = false
////        healthButton.addTarget(self, action: #selector(healthTapped), for: .touchUpInside)
////
////        view.addSubview(healthButton)
////
////        // Result TextView
////        resultTextView.backgroundColor = UIColor.black.withAlphaComponent(0.6)
////        resultTextView.textColor = .systemGreen
////        resultTextView.font = .monospacedSystemFont(ofSize: 13, weight: .regular)
////        resultTextView.isEditable = false
////        resultTextView.layer.cornerRadius = 8
////        resultTextView.translatesAutoresizingMaskIntoConstraints = false
////
////        view.addSubview(resultTextView)
////
////        NSLayoutConstraint.activate([
////            healthButton.bottomAnchor.constraint(equalTo: captureButton.topAnchor, constant: -12),
////            healthButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
////            healthButton.widthAnchor.constraint(equalToConstant: 220),
////            healthButton.heightAnchor.constraint(equalToConstant: 40),
////
////            resultTextView.leadingAnchor.constraint(equalTo: view.leadingAnchor, constant: 16),
////            resultTextView.trailingAnchor.constraint(equalTo: view.trailingAnchor, constant: -16),
////            resultTextView.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor, constant: 16),
////            resultTextView.heightAnchor.constraint(equalToConstant: 160)
////        ])
////        
////        // Debug Health Button
////        debugHealthButton.setTitle("Debug Health", for: .normal)
////        debugHealthButton.backgroundColor = .systemOrange
////        debugHealthButton.tintColor = .white
////        debugHealthButton.layer.cornerRadius = 10
////        debugHealthButton.translatesAutoresizingMaskIntoConstraints = false
////        debugHealthButton.addTarget(self, action: #selector(debugHealthTapped), for: .touchUpInside)
////
////        view.addSubview(debugHealthButton)
////
////        NSLayoutConstraint.activate([
////            debugHealthButton.bottomAnchor.constraint(equalTo: healthButton.topAnchor, constant: -10),
////            debugHealthButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
////            debugHealthButton.widthAnchor.constraint(equalToConstant: 220),
////            debugHealthButton.heightAnchor.constraint(equalToConstant: 40)
////        ])
////    }
////}
//
//extension FLVETestVCx {
//    private func setupUI() {
//        view.backgroundColor = .black
//
//        // Camera container
//        cameraContainer.translatesAutoresizingMaskIntoConstraints = false
//        cameraContainer.isUserInteractionEnabled = false
//        view.addSubview(cameraContainer)
//        
//        previewLayer.videoGravity = .resizeAspectFill
//        cameraContainer.layer.addSublayer(previewLayer)
//        previewLayer.addSublayer(faceBoxLayer)
//
//
//        // Buttons
//        let buttons = [captureButton, healthButton, debugHealthButton]
//        let titles = ["Capture & Liveness", "Health Check", "Debug Health"]
//        let colors: [UIColor] = [.systemBlue, .systemGreen, .systemOrange]
//
//        for (i, button) in buttons.enumerated() {
//            button.setTitle(titles[i], for: .normal)
//            button.backgroundColor = colors[i]
//            button.tintColor = .white
//            button.layer.cornerRadius = 10
//            button.translatesAutoresizingMaskIntoConstraints = false
//            view.addSubview(button)
//        }
//
//        // TextView
//        resultTextView.backgroundColor = UIColor.black.withAlphaComponent(0.6)
//        resultTextView.textColor = .systemGreen
//        resultTextView.font = .monospacedSystemFont(ofSize: 13, weight: .regular)
//        resultTextView.isEditable = false
//        resultTextView.layer.cornerRadius = 8
//        resultTextView.translatesAutoresizingMaskIntoConstraints = false
//        view.addSubview(resultTextView)
//
//        NSLayoutConstraint.activate([
//            cameraContainer.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor, constant: 10),
//            cameraContainer.centerXAnchor.constraint(equalTo: view.centerXAnchor),
//            cameraContainer.widthAnchor.constraint(equalTo: view.widthAnchor, multiplier: 0.8),
//            cameraContainer.heightAnchor.constraint(equalTo: cameraContainer.widthAnchor, multiplier: 4.0 / 3.0),
//
//            captureButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -12),
//            captureButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
//            captureButton.widthAnchor.constraint(equalToConstant: 220),
//            captureButton.heightAnchor.constraint(equalToConstant: 40),
//
//            healthButton.bottomAnchor.constraint(equalTo: captureButton.topAnchor, constant: -8),
//            healthButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
//            healthButton.widthAnchor.constraint(equalToConstant: 220),
//            healthButton.heightAnchor.constraint(equalToConstant: 40),
//
//            debugHealthButton.bottomAnchor.constraint(equalTo: healthButton.topAnchor, constant: -8),
//            debugHealthButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
//            debugHealthButton.widthAnchor.constraint(equalToConstant: 220),
//            debugHealthButton.heightAnchor.constraint(equalToConstant: 40),
//
//            resultTextView.topAnchor.constraint(equalTo: cameraContainer.bottomAnchor, constant: 10),
//            resultTextView.leadingAnchor.constraint(equalTo: view.leadingAnchor, constant: 16),
//            resultTextView.trailingAnchor.constraint(equalTo: view.trailingAnchor, constant: -16),
//            resultTextView.bottomAnchor.constraint(equalTo: debugHealthButton.topAnchor, constant: -10),
//        ])
//
//        captureButton.addTarget(self, action: #selector(captureTapped), for: .touchUpInside)
//        healthButton.addTarget(self, action: #selector(healthTapped), for: .touchUpInside)
//        debugHealthButton.addTarget(self, action: #selector(debugHealthTapped), for: .touchUpInside)
//        
//        // 🔑 Ensure buttons are on top
//        view.bringSubviewToFront(captureButton)
//        view.bringSubviewToFront(healthButton)
//        view.bringSubviewToFront(debugHealthButton)
//    }
//
//    override func viewDidLayoutSubviews() {
//        super.viewDidLayoutSubviews()
//        previewLayer.frame = cameraContainer.bounds
//        
//        let previewSize = previewLayer.bounds.size
//        print("📸 Preview size:", previewSize)
//    }
//}
//
//extension FLVETestVCx {
//
//    @objc func healthTapped() {
//        Task {
//            await performHealthCheck()
//        }
//    }
//
//    func performHealthCheck() async {
////        guard let url = URL(string: "https://farukhannan-flve.hf.space/health") else { return }
////
////        var request = URLRequest(url: url)
////        request.httpMethod = "GET"
////        request.setValue(
////            "Bearer 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",
////            forHTTPHeaderField: "Authorization"
////        )
//        
//        let baseURL = RuntimeAPIConfig.shared.resolvedBaseURL()
//        guard let url = URL(string: "\(baseURL)/health") else { return }
//
//        var request = URLRequest(url: url)
//        request.httpMethod = "GET"
//
//        if let token = RuntimeAPIConfig.shared.resolvedAuthToken() {
//            request.setValue(token, forHTTPHeaderField: "Authorization")
//        }
//
//        print(baseURL, RuntimeAPIConfig.shared.resolvedAuthToken())
//        
//        do {
//            let (data, _) = try await URLSession.shared.data(for: request)
//
//            let decoder = JSONDecoder()
//            let response = try decoder.decode(HealthResponse.self, from: data)
//            print("✅ HEALTH RESPONSE:", response)
//
//            DispatchQueue.main.async {
//                self.resultTextView.text = """
//                Health Status: \(response.status ?? "N/A")
//                Device: \(response.device ?? "N/A")
//                Detector: \(response.models?.detector ?? false)
//                Embedder: \(response.models?.embedder ?? false)
//                Liveness: \(response.models?.liveness ?? false)
//                """
//            }
//
//        } catch {
//            DispatchQueue.main.async {
//                self.resultTextView.text = "❌ Health check failed:\n\(error.localizedDescription)"
//            }
//            print("❌ Health error:", error)
//        }
//    }
//}
//
//extension FLVETestVCx {
//
//    @objc func debugHealthTapped() {
//        Task {
//            await performDebugHealthCheck()
//        }
//    }
//
//    func performDebugHealthCheck() async {
//        guard let url = URL(string: "https://farukhannan-flve.hf.space/debug") else { return }
//
//        var request = URLRequest(url: url)
//        request.httpMethod = "GET"
//        request.setValue(
//            "Bearer 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",
//            forHTTPHeaderField: "Authorization"
//        )
//
//        do {
//            let (data, _) = try await URLSession.shared.data(for: request)
//
//            let decoder = JSONDecoder()
//            let response = try decoder.decode(DebugHealthResponse.self, from: data)
//            print("✅ DEBUG RESPONSE:", response)
//
//            DispatchQueue.main.async {
//                self.resultTextView.text = """
//                Python: \(response.pythonVersion ?? "N/A")
//                Engine Initialized: \(response.engineInitialized ?? false)
//                Device: \(response.device ?? "N/A")
//                Detector Loaded: \(response.detectorLoaded ?? false)
//                Embedder Loaded: \(response.embedderLoaded ?? false)
//                Liveness Loaded: \(response.livenessLoaded ?? false)
//                Mediapipe Available: \(response.mediapipeAvailable ?? false)
//                Models Found: \(response.modelsFound?.joined(separator: ", ") ?? "N/A")
//                """
//            }
//
//        } catch {
//            DispatchQueue.main.async {
//                self.resultTextView.text = "❌ Debug health failed:\n\(error.localizedDescription)"
//            }
//            print("❌ Debug health error:", error)
//        }
//    }
//}
//
//extension FLVETestVCx: AVCapturePhotoCaptureDelegate {
////    @objc private func captureTapped() {
////        let settings = AVCapturePhotoSettings()
////        photoOutput.capturePhoto(with: settings, delegate: self)
////    }
//    
//    @objc private func captureTapped() {
//        guard !isCapturingSequence else { return }
//
//        isCapturingSequence = true
//        captureCount = 0
//        confidenceValues.removeAll()
//        livenessValues.removeAll()
//
//        captureNextFrame()
//    }
//    
//    private func captureNextFrame() {
//        guard captureCount < maxCaptures else {
//            finishSequence()
//            return
//        }
//
//        captureCount += 1
//
//        let settings = AVCapturePhotoSettings()
//        photoOutput.capturePhoto(with: settings, delegate: self)
//    }
//
//    private func drawFaceBox(
//        box: [Int],
//        imageSize: CGSize,
//        confidence: Double
//    ) {
//        guard box.count == 4 else { return }
//
//        let color: UIColor
//        switch confidence {
//        case let c where c > 0.7: color = .systemGreen
//        case let c where c < 0.3: color = .systemRed
//        default: color = .systemYellow
//        }
//
//        faceBoxLayer.strokeColor = color.cgColor
//
//        // ML box in image pixel space
//        let x1 = CGFloat(box[0])
//        let y1 = CGFloat(box[1])
//        let x2 = CGFloat(box[2])
//        let y2 = CGFloat(box[3])
//
//        let imageRect = CGRect(
//            x: x1,
//            y: y1,
//            width: x2 - x1,
//            height: y2 - y1
//        )
//
//        // Convert image rect → preview layer rect
//        let previewRect = convertImageRectToPreviewLayer(
//            imageRect: imageRect,
//            imageSize: imageSize
//        )
//
//        faceBoxLayer.path = UIBezierPath(rect: previewRect).cgPath
//        faceBoxLayer.isHidden = false
//    }
//    
//    private func convertImageRectToPreviewLayer(
//        imageRect: CGRect,
//        imageSize: CGSize
//    ) -> CGRect {
//
//        let previewSize = previewLayer.bounds.size
//        print("📸 Preview size:", previewSize)
//
//        let imageAspect = imageSize.width / imageSize.height
//        let previewAspect = previewSize.width / previewSize.height
//
//        var scale: CGFloat
//        var xOffset: CGFloat = 0
//        var yOffset: CGFloat = 0
//
//        if previewAspect > imageAspect {
//            // Preview is wider → image cropped vertically
//            scale = previewSize.width / imageSize.width
//            let scaledHeight = imageSize.height * scale
//            yOffset = (scaledHeight - previewSize.height) / 2
//        } else {
//            // Preview is taller → image cropped horizontally
//            scale = previewSize.height / imageSize.height
//            let scaledWidth = imageSize.width * scale
//            xOffset = (scaledWidth - previewSize.width) / 2
//        }
//
//        var rect = CGRect(
//            x: imageRect.origin.x * scale - xOffset,
//            y: imageRect.origin.y * scale - yOffset,
//            width: imageRect.size.width * scale,
//            height: imageRect.size.height * scale
//        )
//
//
//        return rect
//    }
//
//    
//    private func finishSequence() {
//        isCapturingSequence = false
//
//        let avgConfidence = confidenceValues.reduce(0, +) / Double(confidenceValues.count)
//        let avgLiveness = livenessValues.reduce(0, +) / Double(livenessValues.count)
//
//        DispatchQueue.main.async {
//            self.faceBoxLayer.isHidden = true
//
//            self.resultTextView.text = """
//            ✅ Sequence Complete
//
//            Avg Confidence: \(String(format: "%.2f%%", avgConfidence * 100))
//            Avg Liveness: \(String(format: "%.2f%%", avgLiveness * 100))
//
//            Frames: \(self.confidenceValues.count)
//            """
//        }
//    }
//
//
////    func photoOutput(_ output: AVCapturePhotoOutput,
////                     didFinishProcessingPhoto photo: AVCapturePhoto,
////                     error: Error?) {
////
////        guard
////            let data = photo.fileDataRepresentation(),
////            let image = UIImage(data: data)
////        else {
////            print("❌ Failed to get image")
////            return
////        }
////
////        Task {
////            await sendToLivenessJSON(image: image)    // JSON base64 (production)
////        }
////    }
//    
//    func photoOutput(_ output: AVCapturePhotoOutput,
//                     didFinishProcessingPhoto photo: AVCapturePhoto,
//                     error: Error?) {
//
//        guard
//            let data = photo.fileDataRepresentation(),
//            let image = UIImage(data: data)
//        else { return }
//
//        lastUploadedImageSize = image.size
//
//        Task {
//            await sendToLivenessJSON(image: image)
//        }
//    }
//
//}
//
//extension FLVETestVCx {
//
//    func sendToLivenessJSON(image: UIImage) async {
//        guard let url = URL(string: "https://farukhannan-flve.hf.space/liveness") else { return }
//        guard let jpegData = image.jpegData(maxKB: 80) else {
//            print("❌ Image compression failed")
//            return
//        }
//
//        let base64 = jpegData.base64EncodedString()
//
//        print("📦 Upload size:", jpegData.count / 1024, "KB")
//
//        let payload = LivenessRequest(
//            imageData: base64,
//            format: "IMAGE_FORMAT_JPEG"
//        )
//
//        var request = URLRequest(url: url)
//        request.httpMethod = "POST"
//        request.setValue(
//            "Bearer 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",
//            forHTTPHeaderField: "Authorization"
//        )
//        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
//
//        do {
//            let body = try JSONEncoder().encode(payload)
//            request.httpBody = body
//
//            let (data, _) = try await URLSession.shared.data(for: request)
//
//            let decoder = JSONDecoder()
//            let response = try decoder.decode(LivenessResponse.self, from: data)
//            print("✅ LIVENESS JSON RESPONSE:", response)
//
//            DispatchQueue.main.async {
//                let confidencePercent = (response.confidence ?? 0) * 100
//                let livenessPercent = (response.livenessScore ?? 0) * 100
//
//                self.resultTextView.text = """
//                Box: \(response.box ?? [])
//                Confidence: \(String(format: "%.2f%%", confidencePercent))
//                Liveness Score: \(String(format: "%.2f%%", livenessPercent))
//                Is Live: \(response.isLive == true ? "✅ Live" : "❌ Not Live")
//                Error: \(response.error ?? "N/A")
//                """
//            }
//            
//            let confidence = response.confidence ?? 0
//            let liveness = response.livenessScore ?? 0
//
//            confidenceValues.append(confidence)
//            livenessValues.append(liveness)
//
//            if let box = response.box,
//               let imgSize = self.lastUploadedImageSize {
//
//                self.drawFaceBox(
//                    box: box,
//                    imageSize: imgSize,
//                    confidence: confidence
//                )
//            }
//
//            // Capture next frame after a short delay
//            DispatchQueue.main.asyncAfter(deadline: .now() + 0.1) {
//                self.captureNextFrame()
//            }
//
//
//        } catch {
//            DispatchQueue.main.async {
//                self.resultTextView.text = "❌ Liveness JSON failed:\n\(error.localizedDescription)"
//            }
//            print("❌ JSON Liveness error:", error)
//        }
//    }
//}
//
//extension UIImage {
//
//    func jpegData(
//        maxKB: Int = 80,
//        maxDimension: CGFloat = 720
//    ) -> Data? {
//
//        let resized = self.resized(maxDimension: maxDimension)
//
//        let maxBytes = maxKB * 1024
//        var minQuality: CGFloat = 0.15
//        var maxQuality: CGFloat = 0.7
//        var bestData: Data?
//
//        for _ in 0..<8 {
//            let quality = (minQuality + maxQuality) / 2
//
//            guard let data = resized.jpegData(compressionQuality: quality) else {
//                return nil
//            }
//
//            if data.count > maxBytes {
//                maxQuality = quality
//            } else {
//                bestData = data
//                minQuality = quality
//            }
//        }
//
//        return bestData
//    }
//
//    private func resized(maxDimension: CGFloat) -> UIImage {
//        let maxSide = max(size.width, size.height)
//        guard maxSide > maxDimension else { return self }
//
//        let scale = maxDimension / maxSide
//        let newSize = CGSize(
//            width: size.width * scale,
//            height: size.height * scale
//        )
//
//        UIGraphicsBeginImageContextWithOptions(newSize, true, 1.0)
//        draw(in: CGRect(origin: .zero, size: newSize))
//        let img = UIGraphicsGetImageFromCurrentImageContext()
//        UIGraphicsEndImageContext()
//
//        return img ?? self
//    }
//}
