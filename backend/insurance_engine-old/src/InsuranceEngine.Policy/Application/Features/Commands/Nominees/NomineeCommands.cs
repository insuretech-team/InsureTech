using System;
using System.Collections.Generic;
using MediatR;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Policy.Application.Features.Commands.Nominees;

// --- Add Nominee ---
public record AddNomineeCommand(
    Guid PolicyId,
    Guid? BeneficiaryId,
    string FullName,
    string Relationship,
    double SharePercentage,
    DateTime? DateOfBirth = null,
    string? NidNumber = null,
    string? PhoneNumber = null,
    string? NomineeDobText = null
) : IRequest<Result<Guid>>;

// --- Update Nominee ---
public record UpdateNomineeCommand(
    Guid PolicyId,
    Guid NomineeId,
    string? FullName = null,
    string? Relationship = null,
    double? SharePercentage = null,
    DateTime? DateOfBirth = null,
    string? NidNumber = null,
    string? PhoneNumber = null,
    string? NomineeDobText = null
) : IRequest<Result>;

// --- Delete Nominee ---
public record DeleteNomineeCommand(Guid PolicyId, Guid NomineeId) : IRequest<Result>;

// --- List Nominees ---
public record ListNomineesQuery(Guid PolicyId) : IRequest<List<NomineeDto>>;
