using InsuranceEngine.Grpc.Clients;
using Insuretech.Notification.Entity.V1;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Infrastructure.Notifications;

public interface INotificationService
{
    Task<string> SendNotificationAsync(
        string recipientId,
        NotificationType type,
        NotificationChannel channel,
        string message,
        string? subject = null,
        NotificationPriority priority = NotificationPriority.Normal,
        CancellationToken ct = default);

    Task<string> SendEmailAsync(
        string recipientId,
        string subject,
        string body,
        Dictionary<string, string>? templateData = null,
        CancellationToken ct = default);

    Task<string> SendSmsAsync(
        string recipientId,
        string message,
        CancellationToken ct = default);

    Task SendBulkNotificationsAsync(
        IEnumerable<(string recipientId, NotificationType type, NotificationChannel channel, string message, string? subject)> notifications,
        CancellationToken ct = default);

    Task NotifyPolicyIssuedAsync(string userId, string policyNumber, string policyHolderName, CancellationToken ct = default);
    Task NotifyClaimSubmittedAsync(string userId, string claimId, string claimType, CancellationToken ct = default);
    Task NotifyClaimApprovedAsync(string userId, string claimId, decimal approvedAmount, CancellationToken ct = default);
    Task NotifyClaimRejectedAsync(string userId, string claimId, string reason, CancellationToken ct = default);
    Task NotifyRenewalReminderAsync(string userId, string policyNumber, DateTime expiryDate, CancellationToken ct = default);
    Task NotifyGracePeriodAsync(string userId, string policyNumber, int daysRemaining, CancellationToken ct = default);
    Task NotifyPolicyLapsedAsync(string userId, string policyNumber, CancellationToken ct = default);
    Task NotifyOtpAsync(string userId, string otpCode, CancellationToken ct = default);
}

public class NotificationService : INotificationService
{
    private readonly InsuranceServiceClient _client;
    private readonly ILogger<NotificationService> _logger;

    public NotificationService(InsuranceServiceClient client, ILogger<NotificationService> logger)
    {
        _client = client;
        _logger = logger;
    }

    public async Task<string> SendNotificationAsync(
        string recipientId,
        NotificationType type,
        NotificationChannel channel,
        string message,
        string? subject = null,
        NotificationPriority priority = NotificationPriority.Normal,
        CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Notification.Services.V1.SendNotificationRequest
            {
                RecipientId = recipientId,
                Type = type,
                Channel = channel,
                Message = message,
                Priority = priority
            };

            if (!string.IsNullOrEmpty(subject))
                request.Subject = subject;

            var response = await _client.Notifications.SendNotificationAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                _logger.LogError("Failed to send notification: {Error}", response.Error.Message);
                throw new InvalidOperationException($"Notification failed: {response.Error.Message}");
            }

            _logger.LogInformation("Notification sent: {NotificationId} to {RecipientId}", 
                response.NotificationId, recipientId);
            
            return response.NotificationId;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error sending notification to {RecipientId}", recipientId);
            throw;
        }
    }

    public async Task<string> SendEmailAsync(
        string recipientId,
        string subject,
        string body,
        Dictionary<string, string>? templateData = null,
        CancellationToken ct = default)
    {
        var request = new Insuretech.Notification.Services.V1.SendNotificationRequest
        {
            RecipientId = recipientId,
            Type = NotificationType.Unspecified,
            Channel = NotificationChannel.Email,
            Subject = subject,
            Message = body,
            Priority = NotificationPriority.Normal
        };

        if (templateData != null)
        {
            foreach (var kvp in templateData)
            {
                request.TemplateData[kvp.Key] = kvp.Value;
            }
        }

        var response = await _client.Notifications.SendNotificationAsync(request, _client.BuildCallOptions(ct));
        return response.NotificationId;
    }

    public async Task<string> SendSmsAsync(
        string recipientId,
        string message,
        CancellationToken ct = default)
    {
        var request = new Insuretech.Notification.Services.V1.SendNotificationRequest
        {
            RecipientId = recipientId,
            Type = NotificationType.Unspecified,
            Channel = NotificationChannel.Sms,
            Message = message,
            Priority = NotificationPriority.Normal
        };

        var response = await _client.Notifications.SendNotificationAsync(request, _client.BuildCallOptions(ct));
        return response.NotificationId;
    }

    public async Task SendBulkNotificationsAsync(
        IEnumerable<(string recipientId, NotificationType type, NotificationChannel channel, string message, string? subject)> notifications,
        CancellationToken ct = default)
    {
        var request = new Insuretech.Notification.Services.V1.SendBulkNotificationsRequest();

        foreach (var (recipientId, type, channel, message, subject) in notifications)
        {
            var notification = new Insuretech.Notification.Services.V1.SendNotificationRequest
            {
                RecipientId = recipientId,
                Type = type,
                Channel = channel,
                Message = message,
                Priority = NotificationPriority.Normal
            };

            if (!string.IsNullOrEmpty(subject))
                notification.Subject = subject;

            request.Notifications.Add(notification);
        }

        var response = await _client.Notifications.SendBulkNotificationsAsync(request, _client.BuildCallOptions(ct));
        
        _logger.LogInformation("Bulk notification sent: {SuccessCount} succeeded, {FailedCount} failed",
            response.SuccessCount, response.FailedCount);
    }

    public async Task NotifyPolicyIssuedAsync(string userId, string policyNumber, string policyHolderName, CancellationToken ct = default)
    {
        var message = $"Dear {policyHolderName}, your policy {policyNumber} has been issued successfully. " +
                      "You can now view your policy details in the LabAid app.";

        await SendNotificationAsync(
            userId,
            NotificationType.PolicyIssued,
            NotificationChannel.Email,
            message,
            "Policy Issued - " + policyNumber,
            NotificationPriority.High,
            ct);
    }

    public async Task NotifyClaimSubmittedAsync(string userId, string claimId, string claimType, CancellationToken ct = default)
    {
        var message = $"Your {claimType} claim (ID: {claimId}) has been submitted successfully. " +
                      "We will notify you once it's processed.";

        await SendNotificationAsync(
            userId,
            NotificationType.ClaimSubmitted,
            NotificationChannel.Email,
            message,
            $"Claim Submitted - {claimId}",
            NotificationPriority.Normal,
            ct);
    }

    public async Task NotifyClaimApprovedAsync(string userId, string claimId, decimal approvedAmount, CancellationToken ct = default)
    {
        var message = $"Great news! Your claim (ID: {claimId}) has been approved. " +
                      $"Approved amount: BDT {approvedAmount:N2}. The amount will be credited to your account soon.";

        await SendNotificationAsync(
            userId,
            NotificationType.ClaimApproved,
            NotificationChannel.Email,
            message,
            $"Claim Approved - {claimId}",
            NotificationPriority.High,
            ct);
    }

    public async Task NotifyClaimRejectedAsync(string userId, string claimId, string reason, CancellationToken ct = default)
    {
        var message = $"Your claim (ID: {claimId}) has been rejected. Reason: {reason}. " +
                      "If you have questions, please contact our support team.";

        await SendNotificationAsync(
            userId,
            NotificationType.ClaimRejected,
            NotificationChannel.Email,
            message,
            $"Claim Update - {claimId}",
            NotificationPriority.Normal,
            ct);
    }

    public async Task NotifyRenewalReminderAsync(string userId, string policyNumber, DateTime expiryDate, CancellationToken ct = default)
    {
        var daysUntilExpiry = (expiryDate - DateTime.UtcNow).Days;
        var message = $"Your policy {policyNumber} will expire on {expiryDate:dd MMM yyyy}. " +
                      $"Please renew it within {daysUntilExpiry} days to continue your coverage without interruption.";

        await SendNotificationAsync(
            userId,
            NotificationType.RenewalReminder,
            NotificationChannel.Email,
            message,
            $"Policy Renewal Reminder - {policyNumber}",
            NotificationPriority.Normal,
            ct);
    }

    public async Task NotifyGracePeriodAsync(string userId, string policyNumber, int daysRemaining, CancellationToken ct = default)
    {
        var message = $"IMPORTANT: Your policy {policyNumber} has entered the grace period. " +
                      $"Please renew within {daysRemaining} days to avoid policy lapse.";

        await SendNotificationAsync(
            userId,
            NotificationType.GracePeriod,
            NotificationChannel.Email,
            message,
            $"Action Required: Policy Grace Period - {policyNumber}",
            NotificationPriority.High,
            ct);
    }

    public async Task NotifyPolicyLapsedAsync(string userId, string policyNumber, CancellationToken ct = default)
    {
        var message = $"Your policy {policyNumber} has lapsed due to non-payment. " +
                      "Please contact our support team to discuss revival options.";

        await SendNotificationAsync(
            userId,
            NotificationType.PolicyLapsed,
            NotificationChannel.Email,
            message,
            $"Policy Lapsed - {policyNumber}",
            NotificationPriority.Urgent,
            ct);
    }

    public async Task NotifyOtpAsync(string userId, string otpCode, CancellationToken ct = default)
    {
        var message = $"Your LabAid verification code is: {otpCode}. This code will expire in 5 minutes. " +
                      "Do not share this code with anyone.";

        await SendNotificationAsync(
            userId,
            NotificationType.Otp,
            NotificationChannel.Sms,
            message,
            priority: NotificationPriority.High,
            ct: ct);
    }
}

public class MockNotificationService : INotificationService
{
    private readonly ILogger<MockNotificationService> _logger;

    public MockNotificationService(ILogger<MockNotificationService> logger)
    {
        _logger = logger;
    }

    public Task<string> SendNotificationAsync(
        string recipientId,
        NotificationType type,
        NotificationChannel channel,
        string message,
        string? subject = null,
        NotificationPriority priority = NotificationPriority.Normal,
        CancellationToken ct = default)
    {
        var notificationId = Guid.NewGuid().ToString();
        _logger.LogInformation(
            "[MOCK] Notification sent: {NotificationId} to {RecipientId} via {Channel} - Type: {Type}",
            notificationId, recipientId, channel, type);
        return Task.FromResult(notificationId);
    }

    public Task<string> SendEmailAsync(
        string recipientId,
        string subject,
        string body,
        Dictionary<string, string>? templateData = null,
        CancellationToken ct = default)
    {
        var notificationId = Guid.NewGuid().ToString();
        _logger.LogInformation(
            "[MOCK] Email sent: {NotificationId} to {RecipientId} - Subject: {Subject}",
            notificationId, recipientId, subject);
        return Task.FromResult(notificationId);
    }

    public Task<string> SendSmsAsync(
        string recipientId,
        string message,
        CancellationToken ct = default)
    {
        var notificationId = Guid.NewGuid().ToString();
        _logger.LogInformation(
            "[MOCK] SMS sent: {NotificationId} to {RecipientId}",
            notificationId, recipientId);
        return Task.FromResult(notificationId);
    }

    public Task SendBulkNotificationsAsync(
        IEnumerable<(string recipientId, NotificationType type, NotificationChannel channel, string message, string? subject)> notifications,
        CancellationToken ct = default)
    {
        var count = notifications.Count();
        _logger.LogInformation("[MOCK] Bulk notifications sent: {Count} notifications", count);
        return Task.CompletedTask;
    }

    public Task NotifyPolicyIssuedAsync(string userId, string policyNumber, string policyHolderName, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] PolicyIssued notification for {UserId}: Policy {PolicyNumber}", userId, policyNumber);
        return Task.CompletedTask;
    }

    public Task NotifyClaimSubmittedAsync(string userId, string claimId, string claimType, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] ClaimSubmitted notification for {UserId}: Claim {ClaimId}", userId, claimId);
        return Task.CompletedTask;
    }

    public Task NotifyClaimApprovedAsync(string userId, string claimId, decimal approvedAmount, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] ClaimApproved notification for {UserId}: Claim {ClaimId}, Amount {Amount}", 
            userId, claimId, approvedAmount);
        return Task.CompletedTask;
    }

    public Task NotifyClaimRejectedAsync(string userId, string claimId, string reason, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] ClaimRejected notification for {UserId}: Claim {ClaimId}", userId, claimId);
        return Task.CompletedTask;
    }

    public Task NotifyRenewalReminderAsync(string userId, string policyNumber, DateTime expiryDate, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] RenewalReminder notification for {UserId}: Policy {PolicyNumber}, Expiry {ExpiryDate}", 
            userId, policyNumber, expiryDate);
        return Task.CompletedTask;
    }

    public Task NotifyGracePeriodAsync(string userId, string policyNumber, int daysRemaining, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] GracePeriod notification for {UserId}: Policy {PolicyNumber}, Days {Days}", 
            userId, policyNumber, daysRemaining);
        return Task.CompletedTask;
    }

    public Task NotifyPolicyLapsedAsync(string userId, string policyNumber, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] PolicyLapsed notification for {UserId}: Policy {PolicyNumber}", userId, policyNumber);
        return Task.CompletedTask;
    }

    public Task NotifyOtpAsync(string userId, string otpCode, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] OTP notification for {UserId}: Code {OtpCode}", userId, otpCode);
        return Task.CompletedTask;
    }
}
