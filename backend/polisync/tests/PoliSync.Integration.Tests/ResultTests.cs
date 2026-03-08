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
        var result = Result<string>.NotFound("Resource not found");
        result.Error!.Kind.Should().Be(ResultErrorKind.NotFound);
    }

    [Fact]
    public void Unauthorized_HasUnauthorizedKind()
    {
        var result = Result<string>.Unauthorized("Access denied");
        result.Error!.Kind.Should().Be(ResultErrorKind.Unauthorized);
    }

    [Fact]
    public void Conflict_HasConflictKind()
    {
        var result = Result<string>.Conflict("Already exists");
        result.Error!.Kind.Should().Be(ResultErrorKind.Conflict);
    }

    [Fact]
    public void Map_OnSuccess_TransformsValue()
    {
        var result = Result<int>.Ok(42).Map(x => x.ToString());
        result.IsSuccess.Should().BeTrue();
        result.Value.Should().Be("42");
    }

    [Fact]
    public void Map_OnFailure_PreservesError()
    {
        var result = Result<int>.Fail("ERR", "oops").Map(x => x.ToString());
        result.IsFailure.Should().BeTrue();
        result.Error!.Code.Should().Be("ERR");
    }

    [Fact]
    public void GetValueOrThrow_OnSuccess_ReturnsValue()
    {
        var result = Result<string>.Ok("value");
        result.GetValueOrThrow().Should().Be("value");
    }

    [Fact]
    public void GetValueOrThrow_OnFailure_ThrowsInvalidOperationException()
    {
        var result = Result<string>.Fail("ERR", "oops");
        var act = () => result.GetValueOrThrow();
        act.Should().Throw<InvalidOperationException>().WithMessage("*ERR*oops*");
    }

    // ── ResultExtensions: Error → RpcException mapping ────────────────────

    [Theory]
    [InlineData(ResultErrorKind.NotFound,     StatusCode.NotFound)]
    [InlineData(ResultErrorKind.Unauthorized, StatusCode.PermissionDenied)]
    [InlineData(ResultErrorKind.Conflict,     StatusCode.AlreadyExists)]
    [InlineData(ResultErrorKind.Validation,   StatusCode.InvalidArgument)]
    [InlineData(ResultErrorKind.Internal,     StatusCode.Internal)]
    [InlineData(ResultErrorKind.DomainError,  StatusCode.InvalidArgument)]
    public void ToRpcException_MapsCorrectStatusCode(ResultErrorKind kind, StatusCode expected)
    {
        var error = new ResultError("CODE", "message", kind);
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
