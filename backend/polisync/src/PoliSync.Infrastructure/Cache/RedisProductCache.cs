using System.Text.Json;
using Microsoft.Extensions.Caching.Distributed;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;

namespace PoliSync.Infrastructure.Cache;

/// <summary>
/// Redis-backed cache for product catalog (5-minute TTL per SRS FR-028).
/// </summary>
public sealed class RedisProductCache
{
    private readonly IDistributedCache _cache;
    private readonly ILogger<RedisProductCache> _logger;
    private readonly TimeSpan _ttl;
    private static readonly JsonSerializerOptions JsonOpts = new() { PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower };

    public RedisProductCache(IDistributedCache cache, IConfiguration config, ILogger<RedisProductCache> logger)
    {
        _cache = cache;
        _logger = logger;
        var ttlSeconds = config.GetValue("Cache:ProductTtlSeconds", 300);
        _ttl = TimeSpan.FromSeconds(ttlSeconds);
    }

    public async Task<T?> GetAsync<T>(string key, CancellationToken ct = default) where T : class
    {
        try
        {
            var json = await _cache.GetStringAsync(key, ct);
            if (json is null) return null;
            return JsonSerializer.Deserialize<T>(json, JsonOpts);
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Cache GET failed for key {Key}", key);
            return null; // Cache miss on error — degrade gracefully
        }
    }

    public async Task SetAsync<T>(string key, T value, CancellationToken ct = default) where T : class
    {
        try
        {
            var json = JsonSerializer.Serialize(value, JsonOpts);
            var opts = new DistributedCacheEntryOptions { AbsoluteExpirationRelativeToNow = _ttl };
            await _cache.SetStringAsync(key, json, opts, ct);
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Cache SET failed for key {Key}", key);
        }
    }

    public async Task RemoveAsync(string key, CancellationToken ct = default)
    {
        try { await _cache.RemoveAsync(key, ct); }
        catch (Exception ex) { _logger.LogWarning(ex, "Cache REMOVE failed for key {Key}", key); }
    }

    public static string ProductKey(string productId) => $"polisync:product:{productId}";
    public static string PlanKey(string planId) => $"polisync:plan:{planId}";
    public static string PricingKey(string productId) => $"polisync:pricing:{productId}";
}
