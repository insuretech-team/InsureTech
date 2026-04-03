using System.Collections.Concurrent;
using System.Net.Http.Json;
using System.Security.Cryptography;
using System.Text;
using System.Text.Json;
using Insuretech.Notification.Entity.V1;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Infrastructure.Webhooks;

public class WebhookSubscriptionRequest
{
    public string SubscriberName { get; set; } = string.Empty;
    public string TargetUrl { get; set; } = string.Empty;
    public string Secret { get; set; } = string.Empty;
    public List<string> EventTypes { get; set; } = new();
    public List<string> TopicGroups { get; set; } = new();
    public List<string> Topics { get; set; } = new();
    public List<string> Channels { get; set; } = new();
    public int TimeoutSeconds { get; set; } = 30;
    public int MaxAttempts { get; set; } = 3;
}

public class WebhookDeliveryResult
{
    public bool Success { get; set; }
    public int ResponseStatus { get; set; }
    public string? ResponseBody { get; set; }
    public string? ErrorMessage { get; set; }
}

public interface IWebhookService
{
    Task<string> CreateSubscriptionAsync(WebhookSubscriptionRequest request, CancellationToken ct = default);
    Task UpdateSubscriptionAsync(string subscriptionId, WebhookSubscriptionRequest request, CancellationToken ct = default);
    Task DeleteSubscriptionAsync(string subscriptionId, CancellationToken ct = default);
    Task<WebhookSubscription?> GetSubscriptionAsync(string subscriptionId, CancellationToken ct = default);
    Task<List<WebhookSubscription>> GetSubscriptionsByEventTypeAsync(string eventType, CancellationToken ct = default);
    Task<WebhookDeliveryResult> DeliverWebhookAsync(string subscriptionId, string notificationId, string lifecycleEvent, string sourceTopic, object payload, CancellationToken ct = default);
    Task<List<WebhookDeliveryAttempt>> GetDeliveryAttemptsAsync(string subscriptionId, CancellationToken ct = default);
    Task<string> GenerateSignatureAsync(string payload, string secret, CancellationToken ct = default);
}

public class WebhookService : IWebhookService
{
    private readonly ILogger<WebhookService> _logger;
    private readonly HttpClient _httpClient;
    private readonly ConcurrentDictionary<string, WebhookSubscription> _subscriptions = new();
    private readonly ConcurrentDictionary<string, List<WebhookDeliveryAttempt>> _deliveryAttempts = new();
    private readonly ConcurrentDictionary<string, List<Task>> _pendingDeliveries = new();

    public WebhookService(ILogger<WebhookService> logger, IHttpClientFactory? httpClientFactory = null)
    {
        _logger = logger;
        _httpClient = httpClientFactory?.CreateClient("WebhookClient") ?? new HttpClient();
        _httpClient.Timeout = TimeSpan.FromSeconds(30);
    }

    public Task<string> CreateSubscriptionAsync(WebhookSubscriptionRequest request, CancellationToken ct = default)
    {
        var subscriptionId = $"whk_{Guid.NewGuid():N}";
        var subscription = new WebhookSubscription
        {
            SubscriptionId = subscriptionId,
            SubscriberName = request.SubscriberName,
            TargetUrl = request.TargetUrl,
            Secret = string.IsNullOrEmpty(request.Secret) ? GenerateSecret() : request.Secret,
            IsActive = true,
            TimeoutSeconds = request.TimeoutSeconds > 0 ? request.TimeoutSeconds : 30,
            MaxAttempts = request.MaxAttempts > 0 ? request.MaxAttempts : 3,
            CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow),
            UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow)
        };

        subscription.EventTypes.AddRange(request.EventTypes);
        subscription.TopicGroups.AddRange(request.TopicGroups);
        subscription.Topics.AddRange(request.Topics);
        subscription.Channels.AddRange(request.Channels);

        _subscriptions[subscriptionId] = subscription;
        _deliveryAttempts[subscriptionId] = new List<WebhookDeliveryAttempt>();

        _logger.LogInformation("Created webhook subscription {SubscriptionId} for {SubscriberName}", 
            subscriptionId, request.SubscriberName);

        return Task.FromResult(subscriptionId);
    }

    public Task UpdateSubscriptionAsync(string subscriptionId, WebhookSubscriptionRequest request, CancellationToken ct = default)
    {
        if (!_subscriptions.TryGetValue(subscriptionId, out var subscription))
        {
            throw new KeyNotFoundException($"Subscription {subscriptionId} not found");
        }

        subscription.SubscriberName = request.SubscriberName;
        subscription.TargetUrl = request.TargetUrl;
        subscription.IsActive = true;
        subscription.TimeoutSeconds = request.TimeoutSeconds > 0 ? request.TimeoutSeconds : 30;
        subscription.MaxAttempts = request.MaxAttempts > 0 ? request.MaxAttempts : 3;
        subscription.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);

        subscription.EventTypes.Clear();
        subscription.EventTypes.AddRange(request.EventTypes);
        subscription.TopicGroups.Clear();
        subscription.TopicGroups.AddRange(request.TopicGroups);
        subscription.Topics.Clear();
        subscription.Topics.AddRange(request.Topics);
        subscription.Channels.Clear();
        subscription.Channels.AddRange(request.Channels);

        _logger.LogInformation("Updated webhook subscription {SubscriptionId}", subscriptionId);
        return Task.CompletedTask;
    }

    public Task DeleteSubscriptionAsync(string subscriptionId, CancellationToken ct = default)
    {
        if (_subscriptions.TryRemove(subscriptionId, out var subscription))
        {
            subscription.IsActive = false;
            subscription.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
            _logger.LogInformation("Deactivated webhook subscription {SubscriptionId}", subscriptionId);
        }
        return Task.CompletedTask;
    }

    public Task<WebhookSubscription?> GetSubscriptionAsync(string subscriptionId, CancellationToken ct = default)
    {
        _subscriptions.TryGetValue(subscriptionId, out var subscription);
        return Task.FromResult(subscription);
    }

    public Task<List<WebhookSubscription>> GetSubscriptionsByEventTypeAsync(string eventType, CancellationToken ct = default)
    {
        var matchingSubscriptions = _subscriptions.Values
            .Where(s => s.IsActive && (s.EventTypes.Contains(eventType) || s.EventTypes.Contains("*") || s.TopicGroups.Contains("*") || s.Topics.Contains(eventType)))
            .ToList();

        return Task.FromResult(matchingSubscriptions);
    }

    public async Task<WebhookDeliveryResult> DeliverWebhookAsync(
        string subscriptionId,
        string notificationId,
        string lifecycleEvent,
        string sourceTopic,
        object payload,
        CancellationToken ct = default)
    {
        if (!_subscriptions.TryGetValue(subscriptionId, out var subscription))
        {
            return new WebhookDeliveryResult { Success = false, ErrorMessage = "Subscription not found" };
        }

        var attemptId = $"att_{Guid.NewGuid():N}";
        var payloadJson = JsonSerializer.Serialize(payload);
        var signature = await GenerateSignatureAsync(payloadJson, subscription.Secret, ct);

        var deliveryAttempt = new WebhookDeliveryAttempt
        {
            AttemptId = attemptId,
            SubscriptionId = subscriptionId,
            NotificationId = notificationId,
            LifecycleEvent = lifecycleEvent,
            SourceTopic = sourceTopic,
            Payload = payloadJson,
            Status = "pending",
            RetryCount = 0,
            ScheduledAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow),
            LastAttemptedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow),
            CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow),
            UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow)
        };

        if (!_deliveryAttempts.ContainsKey(subscriptionId))
        {
            _deliveryAttempts[subscriptionId] = new List<WebhookDeliveryAttempt>();
        }
        _deliveryAttempts[subscriptionId].Add(deliveryAttempt);

        var result = await AttemptDeliveryAsync(subscription, deliveryAttempt, payloadJson, signature, ct);

        deliveryAttempt.Status = result.Success ? "delivered" : "failed";
        deliveryAttempt.ResponseStatus = result.ResponseStatus;
        deliveryAttempt.ResponseBody = result.ResponseBody?.Length > 1000 ? result.ResponseBody[..1000] : result.ResponseBody;
        deliveryAttempt.ErrorMessage = result.ErrorMessage;
        deliveryAttempt.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);

        _logger.LogInformation(
            "Webhook delivery {Status}: Subscription={SubscriptionId}, Attempt={AttemptId}, ResponseStatus={ResponseStatus}",
            result.Success ? "SUCCESS" : "FAILED",
            subscriptionId,
            attemptId,
            result.ResponseStatus);

        return result;
    }

    private async Task<WebhookDeliveryResult> AttemptDeliveryAsync(
        WebhookSubscription subscription,
        WebhookDeliveryAttempt attempt,
        string payloadJson,
        string signature,
        CancellationToken ct)
    {
        try
        {
            var webhookPayload = new
            {
                @event = attempt.LifecycleEvent,
                timestamp = DateTime.UtcNow.ToString("O"),
                data = JsonSerializer.Deserialize<JsonElement>(payloadJson)
            };

            var content = JsonContent.Create(webhookPayload);
            content.Headers.Add("X-Webhook-Signature", signature);
            content.Headers.Add("X-Webhook-Event", attempt.LifecycleEvent);
            content.Headers.Add("X-Webhook-Attempt-Id", attempt.AttemptId);

            using var cts = CancellationTokenSource.CreateLinkedTokenSource(ct);
            cts.CancelAfter(TimeSpan.FromSeconds(subscription.TimeoutSeconds));

            var response = await _httpClient.PostAsync(subscription.TargetUrl, content, cts.Token);

            var responseBody = await response.Content.ReadAsStringAsync(ct);

            return new WebhookDeliveryResult
            {
                Success = response.IsSuccessStatusCode,
                ResponseStatus = (int)response.StatusCode,
                ResponseBody = responseBody
            };
        }
        catch (TaskCanceledException)
        {
            return new WebhookDeliveryResult
            {
                Success = false,
                ErrorMessage = "Request timed out"
            };
        }
        catch (HttpRequestException ex)
        {
            return new WebhookDeliveryResult
            {
                Success = false,
                ErrorMessage = $"HTTP error: {ex.Message}"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Unexpected error delivering webhook to {Url}", subscription.TargetUrl);
            return new WebhookDeliveryResult
            {
                Success = false,
                ErrorMessage = ex.Message
            };
        }
    }

    public Task<List<WebhookDeliveryAttempt>> GetDeliveryAttemptsAsync(string subscriptionId, CancellationToken ct = default)
    {
        if (_deliveryAttempts.TryGetValue(subscriptionId, out var attempts))
        {
            return Task.FromResult(attempts.OrderByDescending(a => a.CreatedAt).Take(100).ToList());
        }
        return Task.FromResult(new List<WebhookDeliveryAttempt>());
    }

    public Task<string> GenerateSignatureAsync(string payload, string secret, CancellationToken ct = default)
    {
        using var hmac = new HMACSHA256(Encoding.UTF8.GetBytes(secret));
        var hash = hmac.ComputeHash(Encoding.UTF8.GetBytes(payload));
        var signature = Convert.ToHexString(hash).ToLowerInvariant();
        return Task.FromResult($"sha256={signature}");
    }

    private static string GenerateSecret()
    {
        var bytes = RandomNumberGenerator.GetBytes(32);
        return Convert.ToBase64String(bytes);
    }
}

public class MockWebhookService : IWebhookService
{
    private readonly ILogger<MockWebhookService> _logger;
    private readonly ConcurrentDictionary<string, WebhookSubscription> _subscriptions = new();
    private readonly ConcurrentDictionary<string, List<WebhookDeliveryAttempt>> _deliveryAttempts = new();

    public MockWebhookService(ILogger<MockWebhookService> logger)
    {
        _logger = logger;
    }

    public Task<string> CreateSubscriptionAsync(WebhookSubscriptionRequest request, CancellationToken ct = default)
    {
        var subscriptionId = $"mock_whk_{Guid.NewGuid():N}";
        _logger.LogInformation("[MOCK] Creating webhook subscription: {SubscriptionId} for {SubscriberName}", 
            subscriptionId, request.SubscriberName);

        var subscription = new WebhookSubscription
        {
            SubscriptionId = subscriptionId,
            SubscriberName = request.SubscriberName,
            TargetUrl = request.TargetUrl,
            Secret = "mock_secret_" + Guid.NewGuid().ToString("N"),
            IsActive = true,
            TimeoutSeconds = request.TimeoutSeconds > 0 ? request.TimeoutSeconds : 30,
            MaxAttempts = request.MaxAttempts > 0 ? request.MaxAttempts : 3,
            CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow),
            UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow)
        };
        subscription.EventTypes.AddRange(request.EventTypes);
        subscription.TopicGroups.AddRange(request.TopicGroups);
        subscription.Topics.AddRange(request.Topics);
        subscription.Channels.AddRange(request.Channels);

        _subscriptions[subscriptionId] = subscription;
        _deliveryAttempts[subscriptionId] = new List<WebhookDeliveryAttempt>();

        return Task.FromResult(subscriptionId);
    }

    public Task UpdateSubscriptionAsync(string subscriptionId, WebhookSubscriptionRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] Updating webhook subscription: {SubscriptionId}", subscriptionId);
        return Task.CompletedTask;
    }

    public Task DeleteSubscriptionAsync(string subscriptionId, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] Deleting webhook subscription: {SubscriptionId}", subscriptionId);
        if (_subscriptions.TryRemove(subscriptionId, out _))
        {
            _deliveryAttempts.TryRemove(subscriptionId, out _);
        }
        return Task.CompletedTask;
    }

    public Task<WebhookSubscription?> GetSubscriptionAsync(string subscriptionId, CancellationToken ct = default)
    {
        _subscriptions.TryGetValue(subscriptionId, out var subscription);
        return Task.FromResult(subscription);
    }

    public Task<List<WebhookSubscription>> GetSubscriptionsByEventTypeAsync(string eventType, CancellationToken ct = default)
    {
        var result = _subscriptions.Values.Where(s => s.IsActive).ToList();
        return Task.FromResult(result);
    }

    public Task<WebhookDeliveryResult> DeliverWebhookAsync(
        string subscriptionId,
        string notificationId,
        string lifecycleEvent,
        string sourceTopic,
        object payload,
        CancellationToken ct = default)
    {
        _logger.LogInformation(
            "[MOCK] Webhook delivery: Subscription={SubscriptionId}, Event={Event}, Topic={Topic}",
            subscriptionId, lifecycleEvent, sourceTopic);

        return Task.FromResult(new WebhookDeliveryResult
        {
            Success = true,
            ResponseStatus = 200,
            ResponseBody = "{\"status\": \"mock_success\"}"
        });
    }

    public Task<List<WebhookDeliveryAttempt>> GetDeliveryAttemptsAsync(string subscriptionId, CancellationToken ct = default)
    {
        if (_deliveryAttempts.TryGetValue(subscriptionId, out var attempts))
        {
            return Task.FromResult(attempts);
        }
        return Task.FromResult(new List<WebhookDeliveryAttempt>());
    }

    public Task<string> GenerateSignatureAsync(string payload, string secret, CancellationToken ct = default)
    {
        return Task.FromResult("sha256=mock_signature");
    }
}
