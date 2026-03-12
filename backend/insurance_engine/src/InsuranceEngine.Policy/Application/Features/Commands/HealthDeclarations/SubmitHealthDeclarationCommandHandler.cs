using System;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Commands.HealthDeclarations;

public class SubmitHealthDeclarationCommandHandler : IRequestHandler<SubmitHealthDeclarationCommand, Result<Guid>>
{
    private readonly IUnderwritingRepository _repository;
    private readonly IEncryptionService _encryptionService;

    public SubmitHealthDeclarationCommandHandler(
        IUnderwritingRepository repository,
        IEncryptionService encryptionService)
    {
        _repository = repository;
        _encryptionService = encryptionService;
    }

    public async Task<Result<Guid>> Handle(SubmitHealthDeclarationCommand request, CancellationToken cancellationToken)
    {
        // 1. Verify quote exists
        var quote = await _repository.GetQuoteByIdAsync(request.QuoteId);
        if (quote == null)
            return Result.Fail<Guid>(new Error("NOT_FOUND", $"Quote with ID {request.QuoteId} not found."));

        // 2. Check if health declaration already exists for this quote
        var existing = await _repository.GetHealthDeclarationByQuoteIdAsync(request.QuoteId);
        if (existing != null)
            return Result.Fail<Guid>(new Error("DUPLICATE", "Health declaration already exists for this quote. Use update instead."));

        // 3. Encrypt sensitive JSON fields
        var hdDto = request.HealthDeclaration;

        var preExistingJson = hdDto.PreExistingConditions != null ? JsonSerializer.Serialize(hdDto.PreExistingConditions) : null;
        var encryptedPreExisting = preExistingJson != null ? await _encryptionService.EncryptAsync(preExistingJson) : null;

        var familyHistoryJson = hdDto.FamilyHistory != null ? JsonSerializer.Serialize(hdDto.FamilyHistory) : null;
        var encryptedFamilyHistory = familyHistoryJson != null ? await _encryptionService.EncryptAsync(familyHistoryJson) : null;

        var medicalDocumentsJson = hdDto.MedicalDocuments != null ? JsonSerializer.Serialize(hdDto.MedicalDocuments) : null;
        var encryptedMedicalDocs = medicalDocumentsJson != null ? await _encryptionService.EncryptAsync(medicalDocumentsJson) : null;

        // 4. Create entity
        var declaration = new UnderwritingHealthDeclaration
        {
            Id = Guid.NewGuid(),
            QuoteId = request.QuoteId,
            HeightCm = hdDto.HeightCm,
            WeightKg = hdDto.WeightKg,
            Bmi = hdDto.Bmi,
            HasPreExistingConditions = hdDto.HasPreExistingConditions,
            PreExistingConditionsJson = encryptedPreExisting,
            IsCurrentlyHospitalized = hdDto.IsCurrentlyHospitalized,
            HasFamilyHistory = hdDto.HasFamilyHistory,
            FamilyHistoryJson = encryptedFamilyHistory,
            IsSmoker = hdDto.IsSmoker,
            IsAlcoholConsumer = hdDto.IsAlcoholConsumer,
            OccupationRiskLevel = hdDto.OccupationRiskLevel,
            IsMedicalExamRequired = hdDto.IsMedicalExamRequired,
            IsMedicalExamCompleted = hdDto.IsMedicalExamCompleted,
            MedicalExamResultsJson = null,
            MedicalExamDate = hdDto.MedicalExamDate,
            MedicalDocumentsJson = encryptedMedicalDocs,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        // 5. Persist
        await _repository.AddHealthDeclarationAsync(declaration);

        return Result.Ok(declaration.Id);
    }
}
