using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using PoliSync.Infrastructure.Clients;
using ProposalEntity = Insuretech.Policy.Entity.V1.InsuranceProposal;
using ProposalStatus = Insuretech.Policy.Entity.V1.ProposalStatus;
using PolicyEntity = Insuretech.Policy.Entity.V1.Policy;

namespace PoliSync.Policy.Infrastructure;

public sealed class GoPolicyDataGateway : IPolicyDataGateway
{
    private readonly InsuranceServiceClient _insuranceClient;

    public GoPolicyDataGateway(InsuranceServiceClient insuranceClient) =>
        _insuranceClient = insuranceClient;

    public async Task<PolicyEntity> CreatePolicyAsync(PolicyEntity policy, CancellationToken ct = default)
    {
        var r = await _insuranceClient.Client.CreatePolicyAsync(new CreatePolicyRequest { Policy = policy }, _insuranceClient.BuildCallOptions(ct));
        return r.Policy;
    }

    public async Task<PolicyEntity?> GetPolicyAsync(string policyId, CancellationToken ct = default)
    {
        try
        {
            var r = await _insuranceClient.Client.GetPolicyAsync(new GetPolicyRequest { PolicyId = policyId }, _insuranceClient.BuildCallOptions(ct));
            return r.Policy;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound) { return null; }
    }

    public async Task<PolicyEntity> UpdatePolicyAsync(PolicyEntity policy, CancellationToken ct = default)
    {
        var r = await _insuranceClient.Client.UpdatePolicyAsync(new UpdatePolicyRequest { Policy = policy }, _insuranceClient.BuildCallOptions(ct));
        return r.Policy;
    }

    public async Task<IReadOnlyList<PolicyEntity>> ListPoliciesAsync(string customerId, int page, int pageSize, CancellationToken ct = default)
    {
        var r = await _insuranceClient.Client.ListPoliciesAsync(
            new ListPoliciesRequest { CustomerId = customerId, Page = page, PageSize = pageSize },
            _insuranceClient.BuildCallOptions(ct));
        return r.Policies;
    }

    public Task DeletePolicyAsync(string policyId, CancellationToken ct = default) =>
        _insuranceClient.Client.DeletePolicyAsync(new DeletePolicyRequest { PolicyId = policyId }, _insuranceClient.BuildCallOptions(ct)).ResponseAsync;

    public async Task<ProposalEntity> CreateInsuranceProposalAsync(ProposalEntity proposal, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.CreateInsuranceProposalAsync(
            new CreateInsuranceProposalRequest { Proposal = proposal },
            _insuranceClient.BuildCallOptions(cancellationToken));

        return response.Proposal;
    }

    public async Task<ProposalEntity?> GetInsuranceProposalAsync(string proposalId, CancellationToken cancellationToken = default)
    {
        try
        {
            var response = await _insuranceClient.Client.GetInsuranceProposalAsync(
                new GetInsuranceProposalRequest { ProposalId = proposalId },
                _insuranceClient.BuildCallOptions(cancellationToken));

            return response.Proposal;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<ProposalEntity> UpdateInsuranceProposalAsync(ProposalEntity proposal, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.UpdateInsuranceProposalAsync(
            new UpdateInsuranceProposalRequest { Proposal = proposal },
            _insuranceClient.BuildCallOptions(cancellationToken));

        return response.Proposal;
    }

    public async Task<IReadOnlyList<ProposalEntity>> ListInsuranceProposalsAsync(
        string? orderId,
        string? insurerId,
        string? customerId,
        ProposalStatus? status,
        int page,
        int pageSize,
        CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.ListInsuranceProposalsAsync(
            new ListInsuranceProposalsRequest
            {
                OrderId = orderId ?? string.Empty,
                InsurerId = insurerId ?? string.Empty,
                CustomerId = customerId ?? string.Empty,
                Status = status ?? ProposalStatus.Unspecified,
                Page = page,
                PageSize = pageSize
            },
            _insuranceClient.BuildCallOptions(cancellationToken));

        return response.Proposals;
    }

    public Task DeleteInsuranceProposalAsync(string proposalId, CancellationToken cancellationToken = default)
    {
        return _insuranceClient.Client.DeleteInsuranceProposalAsync(
            new DeleteInsuranceProposalRequest { ProposalId = proposalId },
            _insuranceClient.BuildCallOptions(cancellationToken)).ResponseAsync;
    }
}

