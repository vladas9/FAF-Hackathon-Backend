namespace NetMircoService.Domain;

public interface ISentimentService
{
    string AnalyzeSentiment(string text);
}