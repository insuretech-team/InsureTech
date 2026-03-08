using System;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Application.Interfaces;

namespace InsuranceEngine.Application.Features.Products.Commands.UpdateProduct;

public class UpdateProductCommandHandler : IRequestHandler<UpdateProductCommand, bool>
{
    private readonly IProductRepository _productRepository;

    public UpdateProductCommandHandler(IProductRepository productRepository)
    {
        _productRepository = productRepository;
    }

    public async Task<bool> Handle(UpdateProductCommand request, CancellationToken cancellationToken)
    {
        var product = await _productRepository.GetByIdAsync(request.Id);
        if (product == null) return false;

        product.ProductName = request.ProductName;
        product.ProductNameBn = request.ProductNameBn;
        product.Description = request.Description;
        product.Category = request.Category;
        product.MinSumInsured = request.MinSumInsured;
        product.MaxSumInsured = request.MaxSumInsured;
        product.MinAge = request.MinAge;
        product.MaxAge = request.MaxAge;
        product.UpdatedAt = DateTime.UtcNow;

        await _productRepository.UpdateAsync(product);
        return true;
    }
}
