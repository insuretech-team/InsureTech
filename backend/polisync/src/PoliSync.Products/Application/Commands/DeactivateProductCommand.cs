using Insuretech.Products.Services.V1;
using PoliSync.Infrastructure.Messaging;
using PoliSync.Infrastructure.Persistence;
using PoliSync.Products.Persistence;
using PoliSync.SharedKernel.Auth;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Commands;

public sealed record DeactivateProductCommand(DeactivateProductRequest Request) : ICommand<DeactivateProductResponse>;

public sealed class DeactivateProductHandler : ICommandHandler<DeactivateProductCommand, DeactivateProductResponse>
{
    private readonly ProductRepository _repo;
    private readonly PoliSyncDbContext _db;
    private readonly IEventBus _bus;
    private readonly ICurrentUser _user;

    public DeactivateProductHandler(ProductRepository repo, PoliSyncDbContext db, IEventBus bus, ICurrentUser user)
    { _repo = repo; _db = db; _bus = bus; _user = user; }

    public async Task<Result<DeactivateProductResponse>> Handle(DeactivateProductCommand cmd, CancellationToken ct)
    {
        if (!Guid.TryParse(cmd.Request.ProductId, out var productId))
            return Result<DeactivateProductResponse>.Fail("INVALID", "Invalid product_id format.");

        var record = await _repo.GetByIdAsync(productId, ct);
        if (record is null)
            return Result<DeactivateProductResponse>.NotFound($"Product '{productId}' not found.");

        if (record.Status == "PRODUCT_STATUS_DISCONTINUED")
            return Result<DeactivateProductResponse>.Fail("INVALID", "Discontinued products cannot be deactivated.");

        record.Status = "PRODUCT_STATUS_INACTIVE";
        await _repo.UpdateAsync(record, ct);
        await _db.SaveChangesAsync(ct);

        await _bus.PublishAsync("insuretech.product.deactivated.v1", record.ProductId.ToString(), new
        {
            event_id       = Guid.NewGuid().ToString(),
            product_id     = record.ProductId.ToString(),
            product_code   = record.ProductCode,
            reason         = cmd.Request.Reason,
            deactivated_by = _user.UserId.ToString(),
            timestamp      = DateTime.UtcNow,
        }, ct);

        return Result<DeactivateProductResponse>.Ok(new DeactivateProductResponse
        {
            Message = $"Product '{record.ProductName}' deactivated.",
        });
    }
}
