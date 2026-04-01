using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Underwriting.Application.Commands;

public sealed record SubmitHealthDeclarationCommand(
    string QuoteId,
    int ApplicantAge,
    int HeightCm,
    string WeightKg,
    bool HasPreExistingConditions,
    string PreExistingConditions,
    bool Smoker,
    bool AlcoholConsumer,
    string OccupationRiskLevel,
    bool IsCurrentlyHospitalized = false,
    bool HasFamilyHistory = false,
    string FamilyHistory = ""
) : ICommand<SubmitHealthDeclarationResult>;

public sealed record SubmitHealthDeclarationResult(string DeclarationId, bool MedicalExamRequired, bool AutoApprovalPossible);
