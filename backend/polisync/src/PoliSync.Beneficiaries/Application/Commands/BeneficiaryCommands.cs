using PoliSync.Beneficiaries.Domain;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Beneficiaries.Application.Commands;

// ── CreateIndividualBeneficiary ──────────────────────────────────────

public record CreateIndividualBeneficiaryCommand(
    Guid UserId,
    string FullName,
    DateTime DateOfBirth,
    Gender Gender,
    string NidNumber,
    string MobileNumber,
    string? Email = null,
    Guid? PartnerId = null
) : ICommand<Guid>;

public class CreateIndividualBeneficiaryHandler : ICommandHandler<CreateIndividualBeneficiaryCommand, Guid>
{
    private readonly IBeneficiaryRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public CreateIndividualBeneficiaryHandler(IBeneficiaryRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result<Guid>> Handle(CreateIndividualBeneficiaryCommand cmd, CancellationToken ct)
    {
        var existing = await _repo.GetByUserIdAsync(cmd.UserId, ct);
        if (existing is not null)
            return Result<Guid>.Conflict("Beneficiary already exists for this user");

        // Generate a simple code for now
        var code = $"BEN-{Guid.NewGuid().ToString()[..8].ToUpper()}";
        
        var beneficiary = Beneficiary.Create(cmd.UserId, BeneficiaryType.Individual, code, cmd.PartnerId);
        var details = IndividualBeneficiary.Create(
            beneficiary.BeneficiaryId,
            cmd.FullName,
            cmd.DateOfBirth,
            cmd.Gender,
            cmd.NidNumber,
            contactInfo: $"{{ \"mobile\": \"{cmd.MobileNumber}\", \"email\": \"{cmd.Email}\" }}"
        );

        await _repo.AddAsync(beneficiary, ct);
        await _repo.AddIndividualDetailsAsync(details, ct);
        await _uow.SaveChangesAsync(ct);

        return Result<Guid>.Ok(beneficiary.BeneficiaryId);
    }
}

// ── CreateBusinessBeneficiary ────────────────────────────────────────

public record CreateBusinessBeneficiaryCommand(
    Guid UserId,
    string BusinessName,
    string TradeLicenseNumber,
    string TinNumber,
    string FocalPersonName,
    string FocalPersonMobile,
    Guid? PartnerId = null
) : ICommand<Guid>;

public class CreateBusinessBeneficiaryHandler : ICommandHandler<CreateBusinessBeneficiaryCommand, Guid>
{
    private readonly IBeneficiaryRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public CreateBusinessBeneficiaryHandler(IBeneficiaryRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result<Guid>> Handle(CreateBusinessBeneficiaryCommand cmd, CancellationToken ct)
    {
        var existing = await _repo.GetByUserIdAsync(cmd.UserId, ct);
        if (existing is not null)
            return Result<Guid>.Conflict("Beneficiary already exists for this user");

        var code = $"BEN-{Guid.NewGuid().ToString()[..8].ToUpper()}";
        
        var beneficiary = Beneficiary.Create(cmd.UserId, BeneficiaryType.Business, code, cmd.PartnerId);
        var details = BusinessBeneficiary.Create(
            beneficiary.BeneficiaryId,
            cmd.BusinessName,
            cmd.TradeLicenseNumber,
            cmd.TinNumber,
            cmd.FocalPersonName,
            cmd.FocalPersonMobile,
            cmd.PartnerId
        );

        await _repo.AddAsync(beneficiary, ct);
        await _repo.AddBusinessDetailsAsync(details, ct);
        await _uow.SaveChangesAsync(ct);

        return Result<Guid>.Ok(beneficiary.BeneficiaryId);
    }
}

// ── CompleteKyc ───────────────────────────────────────────────────────

public record CompleteKycCommand(Guid BeneficiaryId, KycStatus Status) : ICommand;

public class CompleteKycHandler : ICommandHandler<CompleteKycCommand>
{
    private readonly IBeneficiaryRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public CompleteKycHandler(IBeneficiaryRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result> Handle(CompleteKycCommand cmd, CancellationToken ct)
    {
        var beneficiary = await _repo.GetByIdAsync(cmd.BeneficiaryId, ct);
        if (beneficiary is null) return Result.NotFound("Beneficiary not found");

        beneficiary.CompleteKyc(cmd.Status);
        _repo.Update(beneficiary);
        await _uow.SaveChangesAsync(ct);

        return Result.Ok();
    }
}
