using Insuretech.Notification.Services.V1;
using Microsoft.Extensions.Logging;

namespace PoliSync.Infrastructure.GrpcClients;

/// <summary>
/// Typed wrapper for the Go Notification service gRPC client.
/// Fire-and-forget — notification failures never block business flows.
/// </summary>
public sealed class NotificationGrpcClient
{
    private readonly GrpcClientFactory _factory;
    private readonly ILogger<NotificationGrpcClient> _logger;

    public NotificationGrpcClient(GrpcClientFactory factory, ILogger<NotificationGrpcClient> logger)
    { _factory = factory; _logger = logger; }

    private NotificationService.NotificationServiceClient Client =>
        _factory.GetClient("NotificationService", ch => new NotificationService.NotificationServiceClient(ch));

    public async Task SendAsync(
        string userId, string channel, string templateId,
        Dictionary<string, string> variables,
        CancellationToken ct = default)
    {
        try
        {
            var req = new SendNotificationRequest
            {
                UserId     = userId,
                Channel    = channel,
                TemplateId = templateId,
            };
            req.Variables.Add(variables);
            await Client.SendNotificationAsync(req, cancellationToken: ct);
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Notification failed for user={UserId} channel={Channel}", userId, channel);
        }
    }
}
