namespace XMAS.Domain;

public interface ILanguageService
{
    Task<string> DetectLanguageAsync(string text);
}