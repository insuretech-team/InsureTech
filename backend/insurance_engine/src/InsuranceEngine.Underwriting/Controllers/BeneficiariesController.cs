using System;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Underwriting.Application.Interfaces;
using InsuranceEngine.Underwriting.Application.DTOs;

namespace InsuranceEngine.Underwriting.Controllers;

[ApiController]
[Route("api/[controller]")]
public class BeneficiariesController : ControllerBase
{
    private readonly IBeneficiaryRepository _repository;

    public BeneficiariesController(IBeneficiaryRepository repository)
    {
        _repository = repository;
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(Guid id)
    {
        var beneficiary = await _repository.GetByIdAsync(id);
        if (beneficiary == null) return NotFound();
        return Ok(beneficiary);
    }

    [HttpGet]
    public async Task<IActionResult> List([FromQuery] string? type, [FromQuery] string? status, [FromQuery] int page = 1, [FromQuery] int pageSize = 10)
    {
        var items = await _repository.ListAsync(type, status, page, pageSize);
        var total = await _repository.GetTotalCountAsync(type, status);
        return Ok(new { items, total, page, pageSize });
    }
}
