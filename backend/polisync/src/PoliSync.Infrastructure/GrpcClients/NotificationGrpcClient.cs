using Insuretech.Notification.Services.V1;
using Insuretech.Notification.Entity.V1;
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
                RecipientId = userId,
                Type = NotificationType.Unspecified,
                Channel = ParseChannel(channel),
                TemplateId = templateId,
                Priority = NotificationPriority.Normal,
            };
            req.TemplateData.Add(variables);
            await Client.SendNotificationAsync(req, cancellationToken: ct);
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Notification failed for user={UserId} channel={Channel}", userId, channel);
        }
    }

    private static NotificationChannel ParseChannel(string channel) =>
        channel.Trim().ToLowerInvariant() switch
        {
            "sms" => NotificationChannel.Sms,
            "email" => NotificationChannel.Email,
            "push" => NotificationChannel.Push,
            "whatsapp" => NotificationChannel.Whatsapp,
            "in_app" => NotificationChannel.InApp,
            "in-app" => NotificationChannel.InApp,
            _ => NotificationChannel.Unspecified
        };
}
