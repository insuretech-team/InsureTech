using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Policy.Application.Commands;

public sealed record AddNomineeCommand(
    string PolicyId,
    string FullName,
    string Relationship,
    double SharePercentage,
    DateTime? DateOfBirth,
    string? NidNumber,
    string? PhoneNumber,
    string? NomineeDobText) : ICommand<string>;

public sealed record UpdateNomineeCommand(
    string PolicyId,
    string NomineeId,
    string? FullName,
    string? Relationship,
    double? SharePercentage,
    DateTime? DateOfBirth,
    string? NidNumber,
    string? PhoneNumber) : ICommand<bool>;

public sealed record DeleteNomineeCommand(string PolicyId, string NomineeId) : ICommand<bool>;
