using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace InsuranceEngine.Api.RequestModels;

public class PaginationQuery
{
    [Range(1, int.MaxValue)]
    [JsonPropertyName("page")]
    public int Page { get; set; } = 1;

    [Range(1, 100)]
    [JsonPropertyName("page_size")]
    public int PageSize { get; set; } = 20;
}
