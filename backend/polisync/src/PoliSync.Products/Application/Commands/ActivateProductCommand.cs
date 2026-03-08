using Insuretech.Products.Services.V1;
using PoliSync.Infrastructure.Messaging;
using PoliSync.Infrastructure.Persistence;
using PoliSync.Products.Persistence;
using PoliSync.SharedKernel.Auth;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Commands;

// ── Command ────────────────────────────────────────────────────────────────
public sealed record ActivateProductCommand(ActivateProductRequest Request) : ICommand<ActivateProductResponse>;

// ── Handler ────────────────────────────────────────────────────────────────
public sealed class ActivateProductHandler : ICommandHandler<ActivateProductCommand, ActivateProductResponse>
{
    private readonly ProductRepository _repo;
    private readonly PoliSyncDbContext _db;
    private readonly IEventBus _bus;
    private readonly ICurrentUser _user;

    public ActivateProductHandler(
        ProductRepository repo, PoliSyncDbContext db,
        IEventBus bus, ICurrentUser user)
    {
        _repo = repo; _db = db; _bus = bus; _user = user;
    }

    public async Task<Result<ActivateProductResponse>> Handle(
        ActivateProductCommand cmd, CancellationToken ct)
    {
        if (!Guid.TryParse(cmd.Request.ProductId, out var productId))
            return Result<ActivateProductResponse>.Fail("INVALID", "Invalid product_id format.");

        var record = await _repo.GetByIdAsync(productId, ct);
        if (record is null)
            return Result<ActivateProductResponse>.NotFound($"Product '{productId}' not found.");

        if (record.Status == "PRODUCT_STATUS_ACTIVE")
            return Result<ActivateProductResponse>.Fail("CONFLICT", "Product is already active.");

        if (record.Status == "PRODUCT_STATUS_DISCONTINUED")
            return Result<ActivateProductResponse>.Fail("INVALID", "Discontinued products cannot be activated.");

        record.Status = "PRODUCT_STATUS_ACTIVE";
        await _repo.UpdateAsync(record, ct);
        await _db.SaveChangesAsync(ct);

        await _bus.PublishAsync(
            "insuretech.product.activated.v1",
            record.ProductId.ToString(),
            new
            {
                event_id     = Guid.NewGuid().ToString(),
                product_id   = record.ProductId.ToString(),
                product_code = record.ProductCode,
                product_name = record.ProductName,
                activated_by = _user.UserId.ToString(),
                timestamp    = DateTime.UtcNow,
            }, ct);

        return Result<ActivateProductResponse>.Ok(new ActivateProductResponse
        {
            Message = $"Product '{record.ProductName}' activated successfully.",
        });
    }
}
