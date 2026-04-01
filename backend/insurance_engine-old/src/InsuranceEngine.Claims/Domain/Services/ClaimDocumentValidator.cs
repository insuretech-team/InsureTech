using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Claims.Domain.Services;

/// <summary>
/// Domain service to enforce document upload rules (FR-099/FR-081).
/// </summary>
public class ClaimDocumentValidator
{
    private static readonly string[] AllowedExtensions = { ".pdf", ".jpg", ".jpeg", ".png" };
    private const long MaxFileSizeBody = 5 * 1024 * 1024; // 5 MB
    private const long MaxTotalSizeBody = 25 * 1024 * 1024; // 25 MB

    public Result Validate(IEnumerable<ValidateDocumentRequest> documents)
    {
        if (documents == null || !documents.Any())
            return Result.Ok();

        long totalSize = 0;

        foreach (var doc in documents)
        {
            // 1. Validate Extension
            var extension = Path.GetExtension(doc.FileName).ToLowerInvariant();
            if (!AllowedExtensions.Contains(extension))
            {
                return Result.Fail("INVALID_FILE_TYPE", 
                    $"File '{doc.FileName}' has an invalid type. Only PDF, JPG, and PNG are allowed.");
            }

            // 2. Validate Individual Size
            if (doc.FileSize > MaxFileSizeBody)
            {
                return Result.Fail("FILE_TOO_LARGE", 
                    $"File '{doc.FileName}' exceeds the 5MB limit.");
            }

            totalSize += doc.FileSize;
        }

        // 3. Validate Total Size
        if (totalSize > MaxTotalSizeBody)
        {
            return Result.Fail("TOTAL_SIZE_EXCEEDED", 
                "The total size of attachments for this claim exceeds the 25MB limit.");
        }

        return Result.Ok();
    }
}

public record ValidateDocumentRequest(string FileName, long FileSize);
