using Grpc.Core;
using InsuranceEngine.Proto;
using InsuranceEngine.Products.Application.Interfaces;
using Microsoft.Extensions.Logging;
using System.Linq;
using System.Threading.Tasks;

namespace InsuranceEngine.Products.GrpcServices;

public class InsuranceGrpcService : InsuranceEngineService.InsuranceEngineServiceBase
{
    private readonly IProductRepository _productRepository;
    private readonly ILogger<InsuranceGrpcService> _logger;

    public InsuranceGrpcService(IProductRepository productRepository, ILogger<InsuranceGrpcService> logger)
    {
        _productRepository = productRepository;
        _logger = logger;
    }

    public override async Task<GetProductQuestionsResponse> GetProductQuestions(GetProductQuestionsRequest request, ServerCallContext context)
    {
        _logger.LogInformation($"Getting questions for product {request.ProductId}");
        
        var product = await _productRepository.GetByIdAsync(System.Guid.Parse(request.ProductId));
        if (product == null)
        {
            throw new RpcException(new Status(StatusCode.NotFound, "Product not found"));
        }

        var response = new GetProductQuestionsResponse();
        
        if (product.Questions != null)
        {
            response.Questions.AddRange(product.Questions.Select(q => new RiskAssessmentQuestionProto
            {
                Id = q.Id.ToString(),
                QuestionText = q.QuestionText,
                QuestionTextBn = q.QuestionTextBn ?? "",
                OptionsJson = q.OptionsJson ?? "[]",
                Weight = q.Weight
            }));
        }

        return response;
    }
}
