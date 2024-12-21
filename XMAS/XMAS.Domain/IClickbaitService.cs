namespace XMAS.Domain;

public interface IClickbaitService
{
    Task<string> DetectClickbaitAsync(string text);
}