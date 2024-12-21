namespace XMAS.Domain;

public interface ISentimentService
{
    string AnalyzeSentiment(string text);
}