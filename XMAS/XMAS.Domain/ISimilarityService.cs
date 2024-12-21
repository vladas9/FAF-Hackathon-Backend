namespace XMAS.Domain;

public interface ISimilarityService
{
    Task<string> CalculateSimilarityAsync(string text1, string text2);
}