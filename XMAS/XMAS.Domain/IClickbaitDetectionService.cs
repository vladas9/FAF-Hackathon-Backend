namespace XMAS.Domain;

public interface IClickbaitDetectionService
{
    Task<string> DetectClickbaitAsync(string text);
}