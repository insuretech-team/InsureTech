using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;

namespace InsuranceEngine.Policy.Domain.Services;

/// <summary>
/// Detection service for duplicate policies (FR-063) and NID uniqueness (FR-033).
/// </summary>
public class PolicyDuplicateDetector
{
    private readonly IPolicyRepository _policyRepository;
    private readonly IEncryptionService _encryptionService;

    public PolicyDuplicateDetector(IPolicyRepository policyRepository, IEncryptionService encryptionService)
    {
        _policyRepository = policyRepository;
        _encryptionService = encryptionService;
    }

    /// <summary>
    /// FR-063: Block duplicate policy for same product + same insured within 30 days.
    /// FR-033: Validate NID uniqueness across policies.
    /// </summary>
    public async Task<Result> ValidateAsync(
        Guid customerId,
        Guid productId,
        string? nidNumber,
        Guid? excludePolicyId = null)
    {
        // FR-063: Same product + same customer within 30 days
        var sinceDate = DateTime.UtcNow.AddDays(-30);
        var hasDuplicate = await _policyRepository.ExistsByCustomerAndProductAsync(
            customerId, productId, sinceDate);

        if (hasDuplicate)
        {
            return Result.Fail(Error.Conflict(
                "A policy for this product was already created for this customer within the last 30 days. " +
                "Cross-product purchases are allowed."));
        }

        // FR-033: NID uniqueness across active policies
        if (!string.IsNullOrEmpty(nidNumber))
        {
            var encryptedNid = _encryptionService.Encrypt(nidNumber);
            var nidExists = await _policyRepository.ExistsByNidAsync(encryptedNid, excludePolicyId);

            if (nidExists)
            {
                return Result.Fail(Error.Conflict(
                    "This NID number is already associated with an existing policy. " +
                    "Please contact support if you believe this is an error."));
            }
        }

        return Result.Ok();
    }
}
