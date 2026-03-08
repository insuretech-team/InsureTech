using Insuretech.Products.Services.V1;
using PoliSync.Infrastructure.Messaging;
using PoliSync.Infrastructure.Persistence;
using PoliSync.Products.Persistence;
using PoliSync.SharedKernel.Auth;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Commands;

public sealed record DiscontinueProductCommand(DiscontinueProductRequest Request) : ICommand<DiscontinueProductResponse>;

public sealed class DiscontinueProductHandler : ICommandHandler<DiscontinueProductCommand, DiscontinueProductResponse>
{
    private readonly ProductRepository _repo;
    private readonly PoliSyncDbContext _db;
    private readonly IEventBus _bus;
    private readonly ICurrentUser _user;

    public DiscontinueProductHandler(ProductRepository repo, PoliSyncDbContext db, IEventBus bus, ICurrentUser user)
    { _repo = repo; _db = db; _bus = bus; _user = user; }

    public async Task<Result<DiscontinueProductResponse>> Handle(DiscontinueProductCommand cmd, CancellationToken ct)
    {
        if (!Guid.TryParse(cmd.Request.ProductId, out var productId))
            return Result<DiscontinueProductResponse>.Fail("INVALID", "Invalid product_id format.");

        var record = await _repo.GetByIdAsync(productId, ct);
        if (record is null)
            return Result<DiscontinueProductResponse>.NotFound($"Product '{productId}' not found.");

        record.Status    = "PRODUCT_STATUS_DISCONTINUED";
        record.DeletedAt = DateTime.UtcNow; // soft delete
        await _repo.UpdateAsync(record, ct);
        await _db.SaveChangesAsync(ct);

        await _bus.PublishAsync("insuretech.product.discontinued.v1", record.ProductId.ToString(), new
        {
            event_id          = Guid.NewGuid().ToString(),
            product_id        = record.ProductId.ToString(),
            product_code      = record.ProductCode,
            reason            = cmd.Request.Reason,
            discontinued_by   = _user.UserId.ToString(),
            timestamp         = DateTime.UtcNow,
        }, ct);

        return Result<DiscontinueProductResponse>.Ok(new DiscontinueProductResponse
        {
            Message = $"Product '{record.ProductName}' discontinued.",
        });
    }
}
