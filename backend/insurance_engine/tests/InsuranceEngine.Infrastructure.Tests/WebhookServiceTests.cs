using Xunit;
using Moq;
using FluentAssertions;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Logging.Abstractions;
using InsuranceEngine.Infrastructure.Webhooks;

namespace InsuranceEngine.Infrastructure.Tests;

public class WebhookServiceTests
{
    private readonly MockWebhookService _webhookService;

    public WebhookServiceTests()
    {
        _webhookService = new MockWebhookService(NullLogger<MockWebhookService>.Instance);
    }

    [Fact]
    public async Task CreateSubscriptionAsync_ShouldCreateSubscription()
    {
        var request = new WebhookSubscriptionRequest
        {
            SubscriberName = "Test Partner",
            TargetUrl = "https://partner.example.com/webhook",
            EventTypes = new List<string> { "policy.issued", "claim.submitted" },
            TimeoutSeconds = 30,
            MaxAttempts = 3
        };

        var subscriptionId = await _webhookService.CreateSubscriptionAsync(request);

        subscriptionId.Should().NotBeNullOrEmpty();
        subscriptionId.Should().StartWith("mock_whk_");
    }

    [Fact]
    public async Task CreateSubscriptionAsync_ShouldGenerateSecret_WhenNotProvided()
    {
        var request = new WebhookSubscriptionRequest
        {
            SubscriberName = "Test Partner",
            TargetUrl = "https://partner.example.com/webhook",
            EventTypes = new List<string> { "*" }
        };

        var subscriptionId = await _webhookService.CreateSubscriptionAsync(request);
        var subscription = await _webhookService.GetSubscriptionAsync(subscriptionId);

        subscription.Should().NotBeNull();
        subscription!.Secret.Should().StartWith("mock_secret_");
    }

    [Fact]
    public async Task GetSubscriptionAsync_ShouldReturnSubscription_WhenExists()
    {
        var request = new WebhookSubscriptionRequest
        {
            SubscriberName = "Test Partner",
            TargetUrl = "https://partner.example.com/webhook",
            EventTypes = new List<string> { "policy.issued" }
        };

        var subscriptionId = await _webhookService.CreateSubscriptionAsync(request);
        var subscription = await _webhookService.GetSubscriptionAsync(subscriptionId);

        subscription.Should().NotBeNull();
        subscription!.SubscriberName.Should().Be("Test Partner");
        subscription.TargetUrl.Should().Be("https://partner.example.com/webhook");
        subscription.IsActive.Should().BeTrue();
    }

    [Fact]
    public async Task GetSubscriptionAsync_ShouldReturnNull_WhenNotExists()
    {
        var subscription = await _webhookService.GetSubscriptionAsync("nonexistent_id");

        subscription.Should().BeNull();
    }

    [Fact]
    public async Task GetSubscriptionsByEventTypeAsync_ShouldReturnMatchingSubscriptions()
    {
        await _webhookService.CreateSubscriptionAsync(new WebhookSubscriptionRequest
        {
            SubscriberName = "Partner A",
            TargetUrl = "https://partner-a.com/webhook",
            EventTypes = new List<string> { "policy.issued" }
        });

        await _webhookService.CreateSubscriptionAsync(new WebhookSubscriptionRequest
        {
            SubscriberName = "Partner B",
            TargetUrl = "https://partner-b.com/webhook",
            EventTypes = new List<string> { "claim.submitted" }
        });

        var subscriptions = await _webhookService.GetSubscriptionsByEventTypeAsync("policy.issued");

        subscriptions.Should().NotBeEmpty();
        subscriptions.Should().Contain(s => s.SubscriberName == "Partner A");
    }

    [Fact]
    public async Task DeliverWebhookAsync_ShouldReturnSuccess()
    {
        var subscriptionId = await _webhookService.CreateSubscriptionAsync(new WebhookSubscriptionRequest
        {
            SubscriberName = "Test",
            TargetUrl = "https://test.com/webhook",
            EventTypes = new List<string> { "*" }
        });

        var payload = new { PolicyNumber = "POL-001", Status = "ISSUED" };
        var result = await _webhookService.DeliverWebhookAsync(
            subscriptionId,
            "notif-001",
            "policy.issued",
            "policy.lifecycle.events",
            payload);

        result.Should().NotBeNull();
        result.Success.Should().BeTrue();
        result.ResponseStatus.Should().Be(200);
    }

    [Fact]
    public async Task GenerateSignatureAsync_ShouldReturnMockSignature()
    {
        var signature = await _webhookService.GenerateSignatureAsync("{}", "test-secret");

        signature.Should().NotBeNullOrEmpty();
        signature.Should().Be("sha256=mock_signature");
    }

    [Fact]
    public async Task DeleteSubscriptionAsync_ShouldRemoveSubscription()
    {
        var subscriptionId = await _webhookService.CreateSubscriptionAsync(new WebhookSubscriptionRequest
        {
            SubscriberName = "Test",
            TargetUrl = "https://test.com/webhook",
            EventTypes = new List<string> { "*" }
        });

        await _webhookService.DeleteSubscriptionAsync(subscriptionId);
        var subscription = await _webhookService.GetSubscriptionAsync(subscriptionId);

        subscription.Should().BeNull();
    }

    [Fact]
    public async Task GetDeliveryAttemptsAsync_ShouldReturnEmptyList_ForNewSubscription()
    {
        var subscriptionId = await _webhookService.CreateSubscriptionAsync(new WebhookSubscriptionRequest
        {
            SubscriberName = "Test",
            TargetUrl = "https://test.com/webhook",
            EventTypes = new List<string> { "*" }
        });

        var attempts = await _webhookService.GetDeliveryAttemptsAsync(subscriptionId);

        attempts.Should().NotBeNull();
    }
}
