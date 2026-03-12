using System;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.SharedKernel.Interfaces;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Queries.HealthDeclarations;

public class GetHealthDeclarationQueryHandler : IRequestHandler<GetHealthDeclarationQuery, UnderwritingHealthDeclarationResponseDto?>
{
    private readonly IUnderwritingRepository _repository;
    private readonly IEncryptionService _encryptionService;

    public GetHealthDeclarationQueryHandler(
        IUnderwritingRepository repository,
        IEncryptionService encryptionService)
    {
        _repository = repository;
        _encryptionService = encryptionService;
    }

    public async Task<UnderwritingHealthDeclarationResponseDto?> Handle(GetHealthDeclarationQuery request, CancellationToken cancellationToken)
    {
        var declaration = await _repository.GetHealthDeclarationByIdAsync(request.DeclarationId);
        return declaration != null ? await MapToDto(declaration) : null;
    }

    private async Task<UnderwritingHealthDeclarationResponseDto> MapToDto(UnderwritingHealthDeclaration hd)
    {
        // Decrypt sensitive fields for response
        List<string>? preExisting = null;
        if (hd.PreExistingConditionsJson != null)
        {
            var decrypted = await _encryptionService.DecryptAsync(hd.PreExistingConditionsJson);
            preExisting = JsonSerializer.Deserialize<List<string>>(decrypted);
        }

        List<string>? familyHistory = null;
        if (hd.FamilyHistoryJson != null)
        {
            var decrypted = await _encryptionService.DecryptAsync(hd.FamilyHistoryJson);
            familyHistory = JsonSerializer.Deserialize<List<string>>(decrypted);
        }

        return new UnderwritingHealthDeclarationResponseDto(
            hd.Id,
            hd.QuoteId,
            hd.HeightCm,
            hd.WeightKg,
            hd.Bmi,
            hd.HasPreExistingConditions,
            preExisting,
            hd.IsCurrentlyHospitalized,
            hd.HasFamilyHistory,
            familyHistory,
            hd.IsSmoker,
            hd.IsAlcoholConsumer,
            hd.OccupationRiskLevel,
            hd.IsMedicalExamRequired,
            hd.IsMedicalExamCompleted,
            hd.MedicalExamDate,
            hd.CreatedAt,
            hd.UpdatedAt
        );
    }
}

public class GetHealthDeclarationByQuoteQueryHandler : IRequestHandler<GetHealthDeclarationByQuoteQuery, UnderwritingHealthDeclarationResponseDto?>
{
    private readonly IUnderwritingRepository _repository;
    private readonly IEncryptionService _encryptionService;

    public GetHealthDeclarationByQuoteQueryHandler(
        IUnderwritingRepository repository,
        IEncryptionService encryptionService)
    {
        _repository = repository;
        _encryptionService = encryptionService;
    }

    public async Task<UnderwritingHealthDeclarationResponseDto?> Handle(GetHealthDeclarationByQuoteQuery request, CancellationToken cancellationToken)
    {
        var declaration = await _repository.GetHealthDeclarationByQuoteIdAsync(request.QuoteId);
        if (declaration == null) return null;

        // Decrypt sensitive fields for response
        List<string>? preExisting = null;
        if (declaration.PreExistingConditionsJson != null)
        {
            var decrypted = await _encryptionService.DecryptAsync(declaration.PreExistingConditionsJson);
            preExisting = JsonSerializer.Deserialize<List<string>>(decrypted);
        }

        List<string>? familyHistory = null;
        if (declaration.FamilyHistoryJson != null)
        {
            var decrypted = await _encryptionService.DecryptAsync(declaration.FamilyHistoryJson);
            familyHistory = JsonSerializer.Deserialize<List<string>>(decrypted);
        }

        return new UnderwritingHealthDeclarationResponseDto(
            declaration.Id,
            declaration.QuoteId,
            declaration.HeightCm,
            declaration.WeightKg,
            declaration.Bmi,
            declaration.HasPreExistingConditions,
            preExisting,
            declaration.IsCurrentlyHospitalized,
            declaration.HasFamilyHistory,
            familyHistory,
            declaration.IsSmoker,
            declaration.IsAlcoholConsumer,
            declaration.OccupationRiskLevel,
            declaration.IsMedicalExamRequired,
            declaration.IsMedicalExamCompleted,
            declaration.MedicalExamDate,
            declaration.CreatedAt,
            declaration.UpdatedAt
        );
    }
}
