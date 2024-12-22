namespace NetMircoService.Domain;

public interface ILanguageService
{
    Task<string> DetectLanguageAsync(string text);
}