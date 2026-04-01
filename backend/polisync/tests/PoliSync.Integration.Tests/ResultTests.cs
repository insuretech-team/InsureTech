using FluentAssertions;
using Grpc.Core;
using PoliSync.SharedKernel.CQRS;
using Xunit;

namespace PoliSync.Integration.Tests;

/// <summary>
/// Unit tests for the Result discriminated union.
/// </summary>
public sealed class ResultTests
{
    [Fact]
    public void Ok_IsSuccess_True()
    {
        var result = Result<string>.Ok("hello");
        result.IsSuccess.Should().BeTrue();
        result.IsFailure.Should().BeFalse();
        result.Value.Should().Be("hello");
        result.Error.Should().BeNull();
    }

    [Fact]
    public void Fail_IsFailure_True()
    {
        var result = Result<string>.Fail("ERR_CODE", "Something went wrong");
        result.IsSuccess.Should().BeFalse();
        result.IsFailure.Should().BeTrue();
        result.Error!.Code.Should().Be("ERR_CODE");
        result.Error.Message.Should().Be("Something went wrong");
    }

    [Fact]
    public void NotFound_HasNotFoundKind()
    {
        var error = Error.NotFound("Resource", "123");
        error.Code.Should().Be("NOT_FOUND");
    }

    [Fact]
    public void Unauthorized_HasUnauthorizedKind()
    {
        var error = Error.Unauthorized("Access denied");
        error.Code.Should().Be("UNAUTHORIZED");
    }

    [Fact]
    public void Conflict_HasConflictKind()
    {
        var error = Error.Conflict("Already exists");
        error.Code.Should().Be("CONFLICT");
    }

    [Fact]
    public void Map_OnSuccess_TransformsValue()
    {
        var result = Result<int>.Ok(42).Match(
            onSuccess: x => x.ToString(),
            onFailure: _ => "failed");

        result.Should().Be("42");
    }

    [Fact]
    public void Match_OnFailure_UsesFailureBranch()
    {
        var result = Result<int>.Fail("ERR", "oops").Match(
            onSuccess: x => x.ToString(),
            onFailure: error => error.Code);

        result.Should().Be("ERR");
    }

    // ── ResultExtensions: Error → RpcException mapping ────────────────────

    [Theory]
    [InlineData("NOT_FOUND", StatusCode.NotFound)]
    [InlineData("UNAUTHORIZED", StatusCode.Unauthenticated)]
    [InlineData("FORBIDDEN", StatusCode.PermissionDenied)]
    [InlineData("CONFLICT", StatusCode.AlreadyExists)]
    [InlineData("VALIDATION_ERROR", StatusCode.InvalidArgument)]
    [InlineData("INTERNAL_ERROR", StatusCode.Internal)]
    [InlineData("UNKNOWN", StatusCode.Unknown)]
    public void ToRpcException_MapsCorrectStatusCode(string code, StatusCode expected)
    {
        var error = new Error(code, "message");
        var ex = error.ToRpcException();
        ex.StatusCode.Should().Be(expected);
        ex.Status.Detail.Should().Be("message");
    }

    // ── Unit Result (no return value) ─────────────────────────────────────

    [Fact]
    public void UnitResult_Ok_IsSuccess()
    {
        var result = Result.Ok();
        result.IsSuccess.Should().BeTrue();
        result.Error.Should().BeNull();
    }

    [Fact]
    public void UnitResult_Fail_IsFailure()
    {
        var result = Result.Fail("ERR", "failed");
        result.IsFailure.Should().BeTrue();
        result.Error!.Code.Should().Be("ERR");
    }
}
