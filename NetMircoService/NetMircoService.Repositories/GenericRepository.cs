using System.Linq.Expressions;
using Microsoft.EntityFrameworkCore;
using NetMircoService.Domain;

namespace NetMircoService.Repositories;

public class GenericRepository<T> : IGenericRepository<T> where T : BaseEntity
{
    private readonly DbContext _context;
    protected readonly DbSet<T> Entities;

    /// <summary>
    ///     Initializes a new instance of the <see cref="GenericRepository{T}" /> class.
    /// </summary>
    /// <param name="context">The database context to use for repository operations.</param>
    public GenericRepository(DbContext context)
    {
        _context = context;
        Entities = context.Set<T>();
    }

    /// <inheritdoc />
    public async Task<List<T>> GetAllAsync(
        Expression<Func<T, bool>> filter = null,
        Func<IQueryable<T>, IQueryable<T>> include = null,
        bool track = true,
        int maxFeature = 100,
        int startIndex = 0)
    {
        IQueryable<T> query = Entities;

        // Avoid tracking entities
        if (!track) query = query.AsNoTracking();

        // Apply filter
        if (filter != null) query = query.Where(filter);

        // Apply include
        if (include != null) query = include(query);

        // Apply default ordering if none specified
        if (!query.Expression.ToString().Contains("OrderBy")) query = query.OrderBy(e => e.Id); // Order by Id

        // Apply pagination
        query = query.Skip(startIndex).Take(maxFeature);

        return await query.ToListAsync();
    }

    /// <inheritdoc />
    public async Task<T> GetByIdAsync(Guid id, Func<IQueryable<T>, IQueryable<T>> include = null,
        bool track = true)
    {
        IQueryable<T> query = Entities;

        // Avoid tracking entities
        if (!track) query = query.AsNoTracking();

        // Apply include
        if (include != null) query = include(query);

        return await query.FirstOrDefaultAsync(e => e.Id == id);
    }

    /// <inheritdoc />
    public async Task AddAsync(T entity)
    {
        await Entities.AddAsync(entity);
        await SaveAsync();
    }

    /// <inheritdoc />
    public async Task UpdateAsync(T entity)
    {
        Entities.Update(entity);
        await SaveAsync();
    }

    /// <inheritdoc />
    public async Task RemoveAsync(T entity)
    {
        Entities.Remove(entity);
        await SaveAsync();
    }

    /// <inheritdoc />
    public async Task SaveAsync()
    {
        await _context.SaveChangesAsync();
    }
}