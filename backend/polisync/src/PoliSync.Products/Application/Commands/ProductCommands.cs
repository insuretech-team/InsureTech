using PoliSync.Products.Domain;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Commands;

// ── CreateProduct ────────────────────────────────────────────────────

public record CreateProductCommand(
    string ProductCode,
    string ProductName,
    ProductCategory Category,
    long BasePremium,
    long MinSumInsured,
    long MaxSumInsured,
    int MinTenureMonths,
    int MaxTenureMonths,
    string CreatedBy,
    string? Description = null,
    List<string>? Exclusions = null,
    string? ProductAttributes = null
) : ICommand<Guid>;

public class CreateProductHandler : ICommandHandler<CreateProductCommand, Guid>
{
    private readonly IProductRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public CreateProductHandler(IProductRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result<Guid>> Handle(CreateProductCommand cmd, CancellationToken ct)
    {
        // Check for duplicate product code
        var existing = await _repo.GetByCodeAsync(cmd.ProductCode, ct);
        if (existing is not null)
            return Result<Guid>.Conflict($"Product with code '{cmd.ProductCode}' already exists");

        var product = Product.Create(
            productCode: cmd.ProductCode,
            productName: cmd.ProductName,
            category: cmd.Category,
            basePremium: cmd.BasePremium,
            minSumInsured: cmd.MinSumInsured,
            maxSumInsured: cmd.MaxSumInsured,
            minTenureMonths: cmd.MinTenureMonths,
            maxTenureMonths: cmd.MaxTenureMonths,
            createdBy: cmd.CreatedBy,
            description: cmd.Description,
            exclusions: cmd.Exclusions,
            productAttributes: cmd.ProductAttributes
        );

        await _repo.AddAsync(product, ct);
        await _uow.SaveChangesAsync(ct);

        return Result<Guid>.Ok(product.ProductId);
    }
}

// ── UpdateProduct ────────────────────────────────────────────────────

public record UpdateProductCommand(
    Guid ProductId,
    string? ProductName = null,
    string? Description = null,
    long? BasePremium = null,
    long? MinSumInsured = null,
    long? MaxSumInsured = null,
    int? MinTenureMonths = null,
    int? MaxTenureMonths = null,
    List<string>? Exclusions = null,
    string? ProductAttributes = null
) : ICommand;

public class UpdateProductHandler : ICommandHandler<UpdateProductCommand>
{
    private readonly IProductRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public UpdateProductHandler(IProductRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result> Handle(UpdateProductCommand cmd, CancellationToken ct)
    {
        var product = await _repo.GetByIdAsync(cmd.ProductId, ct);
        if (product is null)
            return Result.NotFound($"Product '{cmd.ProductId}' not found");

        product.Update(
            productName: cmd.ProductName,
            description: cmd.Description,
            basePremium: cmd.BasePremium,
            minSumInsured: cmd.MinSumInsured,
            maxSumInsured: cmd.MaxSumInsured,
            minTenureMonths: cmd.MinTenureMonths,
            maxTenureMonths: cmd.MaxTenureMonths,
            exclusions: cmd.Exclusions,
            productAttributes: cmd.ProductAttributes
        );

        _repo.Update(product);
        await _uow.SaveChangesAsync(ct);

        return Result.Ok();
    }
}

// ── ActivateProduct ──────────────────────────────────────────────────

public record ActivateProductCommand(Guid ProductId) : ICommand;

public class ActivateProductHandler : ICommandHandler<ActivateProductCommand>
{
    private readonly IProductRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public ActivateProductHandler(IProductRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result> Handle(ActivateProductCommand cmd, CancellationToken ct)
    {
        var product = await _repo.GetByIdAsync(cmd.ProductId, ct);
        if (product is null)
            return Result.NotFound($"Product '{cmd.ProductId}' not found");

        product.Activate();
        _repo.Update(product);
        await _uow.SaveChangesAsync(ct);

        return Result.Ok();
    }
}

// ── DeactivateProduct ────────────────────────────────────────────────

public record DeactivateProductCommand(Guid ProductId, string? Reason = null) : ICommand;

public class DeactivateProductHandler : ICommandHandler<DeactivateProductCommand>
{
    private readonly IProductRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public DeactivateProductHandler(IProductRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result> Handle(DeactivateProductCommand cmd, CancellationToken ct)
    {
        var product = await _repo.GetByIdAsync(cmd.ProductId, ct);
        if (product is null)
            return Result.NotFound($"Product '{cmd.ProductId}' not found");

        product.Deactivate(cmd.Reason);
        _repo.Update(product);
        await _uow.SaveChangesAsync(ct);

        return Result.Ok();
    }
}

// ── DiscontinueProduct ───────────────────────────────────────────────

public record DiscontinueProductCommand(Guid ProductId, string? Reason = null) : ICommand;

public class DiscontinueProductHandler : ICommandHandler<DiscontinueProductCommand>
{
    private readonly IProductRepository _repo;
    private readonly SharedKernel.Persistence.IUnitOfWork _uow;

    public DiscontinueProductHandler(IProductRepository repo, SharedKernel.Persistence.IUnitOfWork uow)
    {
        _repo = repo;
        _uow = uow;
    }

    public async Task<Result> Handle(DiscontinueProductCommand cmd, CancellationToken ct)
    {
        var product = await _repo.GetByIdAsync(cmd.ProductId, ct);
        if (product is null)
            return Result.NotFound($"Product '{cmd.ProductId}' not found");

        product.Discontinue(cmd.Reason);
        _repo.Update(product);
        await _uow.SaveChangesAsync(ct);

        return Result.Ok();
    }
}
