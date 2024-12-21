using System.Linq.Expressions;

namespace Geoserver.Repositories;

public interface IGenericRepository<T> where T : class
{
    /// <summary>
    ///     Retrieves a paginated list of entities that match the specified filter.
    /// </summary>
    /// <param name="filter">An optional filter expression to apply to the entities.</param>
    /// <param name="include">An optional function to include related entities.</param>
    /// <param name="track">A boolean indicating whether the entities should be tracked by the context.</param>
    /// <param name="maxFeature">The maximum number of entities to retrieve. Defaults to 100 if not specified.</param>
    /// <param name="startIndex">The zero-based index of the first entity to retrieve. Used for pagination.</param>
    /// <returns>A list of entities that match the specified filter.</returns>
    Task<List<T>> GetAllAsync(
        Expression<Func<T, bool>> filter = null,
        Func<IQueryable<T>, IQueryable<T>> include = null,
        bool track = true, int maxFeature = 100, int startIndex = 0);

    /// <summary>
    ///     Retrieves an entity by its identifier.
    /// </summary>
    /// <param name="id">The unique identifier of the entity.</param>
    /// <param name="include">An optional function to include related entities.</param>
    /// <param name="track">A boolean indicating whether the entity should be tracked by the context.</param>
    /// <returns>The entity with the specified identifier.</returns>
    Task<T> GetByIdAsync(
        Guid id,
        Func<IQueryable<T>, IQueryable<T>> include = null,
        bool track = true);

    /// <summary>
    ///     Adds a new entity to the repository.
    /// </summary>
    /// <param name="entity">The entity to add.</param>
    Task AddAsync(T entity);

    /// <summary>
    ///     Updates an existing entity in the repository.
    /// </summary>
    /// <param name="entity">The entity to update.</param>
    Task UpdateAsync(T entity);

    /// <summary>
    ///     Removes an entity from the repository.
    /// </summary>
    /// <param name="entity">The entity to remove.</param>
    Task RemoveAsync(T entity);

    /// <summary>
    ///     Saves all changes made in the repository to the database.
    /// </summary>
    Task SaveAsync();
}