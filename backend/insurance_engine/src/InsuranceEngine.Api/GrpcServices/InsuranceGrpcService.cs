using System;
using System.Threading.Tasks;
using Grpc.Core;
using InsuranceEngine.Proto;

namespace InsuranceEngine.Api.GrpcServices;

public class InsuranceGrpcService : InsuranceEngineService.InsuranceEngineServiceBase
{
    public override Task<GetProductQuestionsResponse> GetProductQuestions(GetProductQuestionsRequest request, ServerCallContext context)
    {
        var response = new GetProductQuestionsResponse();
        // Since we reconstructed, we mocked this response to satisfy compilation.
        return Task.FromResult(response);
    }
}
