import AVFoundation
import UIKit

extension CMSampleBuffer {

    func toJPEGData(compression: CGFloat = 0.7) -> Data? {
        guard let pixelBuffer = CMSampleBufferGetImageBuffer(self) else { return nil }

        let ciImage = CIImage(cvPixelBuffer: pixelBuffer)
        let context = CIContext(options: nil)

        guard let cgImage = context.createCGImage(ciImage,
                                                  from: ciImage.extent)
        else { return nil }

        let image = UIImage(
            cgImage: cgImage,
            scale: 1,
            orientation: .rightMirrored
        )

        return image.jpegData(compressionQuality: compression)
    }
}

