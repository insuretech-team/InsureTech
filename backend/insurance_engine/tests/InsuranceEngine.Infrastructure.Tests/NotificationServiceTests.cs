using Xunit;
using FluentAssertions;
using Microsoft.Extensions.Logging.Abstractions;
using InsuranceEngine.Infrastructure.Notifications;
using Insuretech.Notification.Entity.V1;

namespace InsuranceEngine.Infrastructure.Tests;

public class NotificationServiceTests
{
    private readonly MockNotificationService _notificationService;

    public NotificationServiceTests()
    {
        _notificationService = new MockNotificationService(NullLogger<MockNotificationService>.Instance);
    }

    [Fact]
    public async Task SendNotificationAsync_ShouldReturnNotificationId()
    {
        var recipientId = "user-001";
        var type = NotificationType.PolicyIssued;
        var channel = NotificationChannel.Email;
        var message = "Test notification";

        var result = await _notificationService.SendNotificationAsync(
            recipientId, type, channel, message);

        result.Should().NotBeNullOrEmpty();
    }

    [Fact]
    public async Task SendEmailAsync_ShouldReturnNotificationId()
    {
        var recipientId = "user-001";
        var subject = "Test Subject";
        var body = "Test body content";

        var result = await _notificationService.SendEmailAsync(
            recipientId, subject, body);

        result.Should().NotBeNullOrEmpty();
    }

    [Fact]
    public async Task SendSmsAsync_ShouldReturnNotificationId()
    {
        var recipientId = "user-001";
        var message = "Test SMS content";

        var result = await _notificationService.SendSmsAsync(recipientId, message);

        result.Should().NotBeNullOrEmpty();
    }

    [Fact]
    public async Task SendBulkNotificationsAsync_ShouldNotThrow()
    {
        var notifications = new List<(string recipientId, NotificationType type, NotificationChannel channel, string message, string? subject)>
        {
            ("user-001", NotificationType.PolicyIssued, NotificationChannel.Email, "Test 1", "Subject 1"),
            ("user-002", NotificationType.PolicyIssued, NotificationChannel.Sms, "Test 2", null)
        };

        var act = async () => await _notificationService.SendBulkNotificationsAsync(notifications);

        await act.Should().NotThrowAsync();
    }

    [Fact]
    public async Task NotifyPolicyIssuedAsync_ShouldNotThrow()
    {
        var act = async () => await _notificationService.NotifyPolicyIssuedAsync(
            "user-001", "POL-001", "John Doe");

        await act.Should().NotThrowAsync();
    }

    [Fact]
    public async Task NotifyClaimSubmittedAsync_ShouldNotThrow()
    {
        var act = async () => await _notificationService.NotifyClaimSubmittedAsync(
            "user-001", "CLM-001", "Health");

        await act.Should().NotThrowAsync();
    }

    [Fact]
    public async Task NotifyClaimApprovedAsync_ShouldNotThrow()
    {
        var act = async () => await _notificationService.NotifyClaimApprovedAsync(
            "user-001", "CLM-001", 5000);

        await act.Should().NotThrowAsync();
    }

    [Fact]
    public async Task NotifyClaimRejectedAsync_ShouldNotThrow()
    {
        var act = async () => await _notificationService.NotifyClaimRejectedAsync(
            "user-001", "CLM-001", "Insufficient documentation");

        await act.Should().NotThrowAsync();
    }

    [Fact]
    public async Task NotifyRenewalReminderAsync_ShouldNotThrow()
    {
        var act = async () => await _notificationService.NotifyRenewalReminderAsync(
            "user-001", "POL-001", DateTime.UtcNow.AddDays(30));

        await act.Should().NotThrowAsync();
    }

    [Fact]
    public async Task NotifyGracePeriodAsync_ShouldNotThrow()
    {
        var act = async () => await _notificationService.NotifyGracePeriodAsync(
            "user-001", "POL-001", 15);

        await act.Should().NotThrowAsync();
    }

    [Fact]
    public async Task NotifyPolicyLapsedAsync_ShouldNotThrow()
    {
        var act = async () => await _notificationService.NotifyPolicyLapsedAsync(
            "user-001", "POL-001");

        await act.Should().NotThrowAsync();
    }

    [Fact]
    public async Task NotifyOtpAsync_ShouldNotThrow()
    {
        var act = async () => await _notificationService.NotifyOtpAsync(
            "user-001", "123456");

        await act.Should().NotThrowAsync();
    }
}
