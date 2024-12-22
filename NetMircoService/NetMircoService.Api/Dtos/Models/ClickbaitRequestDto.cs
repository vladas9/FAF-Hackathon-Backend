namespace NetMircoService.Api.Dtos.Models;

public class ClickbaitRequestDto
{
    public string Headline { get; set; }
    public string Description { get; set; }
    public string Url { get; set; }
    public List<string> Keywords { get; set; }
}