using FluentAssertions;
using Insuretech.Common.V1;
using Insuretech.Products.Entity.V1;
using Microsoft.Extensions.Logging.Abstractions;
using PoliSync.Products.Application.Commands;
using PoliSync.Products.Application.Queries;
using PoliSync.Products.Infrastructure;
using Xunit;

namespace PoliSync.Products.Tests;

public sealed class ProductDomainTests
{
    [Fact]
    public async Task CreateProduct_WithDuplicateCode_ReturnsFailure()
    {
        var repository = new FakeProductRepository();
        var existing = BuildProduct(productId: "prod-1", productCode: "HLTH-001");
        repository.Seed(existing);
        var handler = new CreateProductCommandHandler(repository, NullLogger<CreateProductCommandHandler>.Instance);

        var result = await handler.Handle(new CreateProductCommand
        {
            ProductCode = "HLTH-001",
            ProductName = "Health Basic",
            Category = ProductCategory.Health,
            Description = "Basic plan",
            BasePremiumAmount = 10_000,
            MinSumInsuredAmount = 100_000,
            MaxSumInsuredAmount = 500_000,
            MinTenureMonths = 1,
            MaxTenureMonths = 12,
            CreatedBy = "admin"
        }, CancellationToken.None);

        result.IsFailure.Should().BeTrue();
        result.Error!.Message.Should().Contain("already exists");
    }

    [Fact]
    public async Task CreateProduct_WithInvalidSumInsuredRange_ReturnsFailure()
    {
        var repository = new FakeProductRepository();
        var handler = new CreateProductCommandHandler(repository, NullLogger<CreateProductCommandHandler>.Instance);

        var result = await handler.Handle(new CreateProductCommand
        {
            ProductCode = "HLTH-001",
            ProductName = "Health Basic",
            Category = ProductCategory.Health,
            Description = "Basic plan",
            BasePremiumAmount = 10_000,
            MinSumInsuredAmount = 500_000,
            MaxSumInsuredAmount = 100_000,
            MinTenureMonths = 1,
            MaxTenureMonths = 12,
            CreatedBy = "admin"
        }, CancellationToken.None);

        result.IsFailure.Should().BeTrue();
        result.Error!.Message.Should().Contain("Min sum insured");
    }

    [Fact]
    public async Task CreateProduct_WithValidRequest_PersistsProduct()
    {
        var repository = new FakeProductRepository();
        var handler = new CreateProductCommandHandler(repository, NullLogger<CreateProductCommandHandler>.Instance);

        var result = await handler.Handle(new CreateProductCommand
        {
            ProductCode = "HLTH-001",
            ProductName = "Health Basic",
            Category = ProductCategory.Health,
            Description = "Basic plan",
            BasePremiumAmount = 10_000,
            MinSumInsuredAmount = 100_000,
            MaxSumInsuredAmount = 500_000,
            MinTenureMonths = 1,
            MaxTenureMonths = 12,
            Exclusions = ["Dental"],
            CreatedBy = "admin"
        }, CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
        result.Value.Should().NotBeNull();
        result.Value!.ProductCode.Should().Be("HLTH-001");
        result.Value.ProductName.Should().Be("Health Basic");
        result.Value.Status.Should().Be(ProductStatus.Draft);
        result.Value.BasePremium.Amount.Should().Be(10_000);
        result.Value.Exclusions.Should().ContainSingle().Which.Should().Be("Dental");
        repository.Stored.Should().ContainSingle(p => p.ProductCode == "HLTH-001");
    }

    [Fact]
    public async Task UpdateProduct_WhenMissing_ReturnsFailure()
    {
        var repository = new FakeProductRepository();
        var handler = new UpdateProductCommandHandler(repository, NullLogger<UpdateProductCommandHandler>.Instance);

        var result = await handler.Handle(new UpdateProductCommand
        {
            ProductId = "missing",
            ProductName = "Updated"
        }, CancellationToken.None);

        result.IsFailure.Should().BeTrue();
        result.Error!.Message.Should().Contain("Product not found");
    }

    [Fact]
    public async Task UpdateProduct_WithValidRequest_UpdatesFields()
    {
        var repository = new FakeProductRepository();
        repository.Seed(BuildProduct(productId: "prod-1", productCode: "HLTH-001"));
        var handler = new UpdateProductCommandHandler(repository, NullLogger<UpdateProductCommandHandler>.Instance);

        var result = await handler.Handle(new UpdateProductCommand
        {
            ProductId = "prod-1",
            ProductName = "Health Plus",
            Description = "Updated description",
            BasePremiumAmount = 20_000,
            Status = ProductStatus.Active,
            Exclusions = ["Maternity"]
        }, CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
        result.Value!.ProductName.Should().Be("Health Plus");
        result.Value.Description.Should().Be("Updated description");
        result.Value.BasePremium.Amount.Should().Be(20_000);
        result.Value.Status.Should().Be(ProductStatus.Active);
        result.Value.Exclusions.Should().BeEquivalentTo(["Maternity"]);
    }

    [Fact]
    public async Task DeleteProduct_WhenFound_RemovesProduct()
    {
        var repository = new FakeProductRepository();
        repository.Seed(BuildProduct(productId: "prod-1", productCode: "HLTH-001"));
        var handler = new DeleteProductCommandHandler(repository, NullLogger<DeleteProductCommandHandler>.Instance);

        var result = await handler.Handle(new DeleteProductCommand
        {
            ProductId = "prod-1"
        }, CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
        result.Value.Should().BeTrue();
        repository.Stored.Should().BeEmpty();
    }

    [Fact]
    public async Task GetProduct_WhenFound_ReturnsSuccess()
    {
        var repository = new FakeProductRepository();
        repository.Seed(BuildProduct(productId: "prod-1", productCode: "HLTH-001"));
        var handler = new GetProductQueryHandler(repository, NullLogger<GetProductQueryHandler>.Instance);

        var result = await handler.Handle(new GetProductQuery
        {
            ProductId = "prod-1"
        }, CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
        result.Value!.ProductId.Should().Be("prod-1");
    }

    [Fact]
    public async Task ListProducts_WithStatusFilter_FiltersResult()
    {
        var repository = new FakeProductRepository();
        repository.Seed(
            BuildProduct(productId: "prod-1", productCode: "HLTH-001", status: ProductStatus.Active),
            BuildProduct(productId: "prod-2", productCode: "HLTH-002", status: ProductStatus.Draft));
        var handler = new ListProductsQueryHandler(repository, NullLogger<ListProductsQueryHandler>.Instance);

        var result = await handler.Handle(new ListProductsQuery
        {
            Status = ProductStatus.Active
        }, CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
        result.Value!.Products.Should().ContainSingle();
        result.Value.Products[0].ProductId.Should().Be("prod-1");
        result.Value.TotalCount.Should().Be(1);
    }

    private static Product BuildProduct(
        string productId,
        string productCode,
        ProductStatus status = ProductStatus.Draft,
        ProductCategory category = ProductCategory.Health)
        => new()
        {
            ProductId = productId,
            ProductCode = productCode,
            ProductName = "Health Basic",
            Category = category,
            Description = "Basic plan",
            BasePremium = new Money { Amount = 10_000, Currency = "BDT" },
            MinSumInsured = new Money { Amount = 100_000, Currency = "BDT" },
            MaxSumInsured = new Money { Amount = 500_000, Currency = "BDT" },
            MinTenureMonths = 1,
            MaxTenureMonths = 12,
            Status = status,
            CreatedBy = "admin",
            BasePremiumCurrency = "BDT",
            MinSumInsuredCurrency = "BDT",
            MaxSumInsuredCurrency = "BDT"
        };

    private sealed class FakeProductRepository : IProductRepository
    {
        private readonly List<Product> _stored = [];
        public IReadOnlyList<Product> Stored => _stored;

        public void Seed(params Product[] products) => _stored.AddRange(products.Select(Clone));

        public Task<Product> CreateAsync(Product product, CancellationToken ct = default)
        {
            var created = Clone(product);
            if (string.IsNullOrWhiteSpace(created.ProductId))
            {
                created.ProductId = Guid.NewGuid().ToString();
            }

            _stored.Add(created);
            return Task.FromResult(Clone(created));
        }

        public Task DeleteAsync(string productId, CancellationToken ct = default)
        {
            _stored.RemoveAll(p => p.ProductId == productId);
            return Task.CompletedTask;
        }

        public Task<List<Product>> GetActiveProductsAsync(CancellationToken ct = default)
            => Task.FromResult(_stored.Where(p => p.Status == ProductStatus.Active).Select(Clone).ToList());

        public Task<List<Product>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken ct = default)
            => Task.FromResult(_stored.Select(Clone).ToList());

        public Task<List<Product>> GetByCategoryAsync(ProductCategory category, CancellationToken ct = default)
            => Task.FromResult(_stored.Where(p => p.Category == category).Select(Clone).ToList());

        public Task<Product?> GetByCodeAsync(string productCode, CancellationToken ct = default)
            => Task.FromResult(_stored.FirstOrDefault(p => p.ProductCode == productCode) is { } product ? Clone(product) : null);

        public Task<Product?> GetByIdAsync(string productId, CancellationToken ct = default)
            => Task.FromResult(_stored.FirstOrDefault(p => p.ProductId == productId) is { } product ? Clone(product) : null);

        public Task<Product> UpdateAsync(Product product, CancellationToken ct = default)
        {
            var index = _stored.FindIndex(p => p.ProductId == product.ProductId);
            if (index >= 0)
            {
                _stored[index] = Clone(product);
            }

            return Task.FromResult(Clone(product));
        }

        private static Product Clone(Product product)
        {
            var copy = new Product
            {
                ProductId = product.ProductId,
                ProductCode = product.ProductCode,
                ProductName = product.ProductName,
                Category = product.Category,
                Description = product.Description,
                BasePremium = product.BasePremium == null
                    ? null
                    : new Money { Amount = product.BasePremium.Amount, Currency = product.BasePremium.Currency },
                MinSumInsured = product.MinSumInsured == null
                    ? null
                    : new Money { Amount = product.MinSumInsured.Amount, Currency = product.MinSumInsured.Currency },
                MaxSumInsured = product.MaxSumInsured == null
                    ? null
                    : new Money { Amount = product.MaxSumInsured.Amount, Currency = product.MaxSumInsured.Currency },
                MinTenureMonths = product.MinTenureMonths,
                MaxTenureMonths = product.MaxTenureMonths,
                Status = product.Status,
                CreatedBy = product.CreatedBy,
                BasePremiumCurrency = product.BasePremiumCurrency,
                MinSumInsuredCurrency = product.MinSumInsuredCurrency,
                MaxSumInsuredCurrency = product.MaxSumInsuredCurrency
            };

            copy.Exclusions.AddRange(product.Exclusions);
            return copy;
        }
    }
}
