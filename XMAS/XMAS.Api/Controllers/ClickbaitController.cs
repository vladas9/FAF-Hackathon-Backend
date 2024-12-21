using System.Text.Json;
using Microsoft.AspNetCore.Mvc;
using XMAS.Api.Dtos.Models;
using XMAS.Domain;
using XMAS.Services;

namespace XMAS.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class ClickbaitController : ControllerBase
{
    private readonly IClickbaitService _clickbaitService;

    public ClickbaitController(IClickbaitService clickbaitService)
    {
        _clickbaitService = clickbaitService;
    }

    [HttpPost("detect")]
    public async Task<IActionResult> DetectClickbait([FromBody] ClickbaitRequestDto request)
    {
        if (string.IsNullOrWhiteSpace(request.Headline) && string.IsNullOrWhiteSpace(request.Description))
        {
            return BadRequest(new { error = "Headline and description cannot both be empty." });
        }

        try
        {
            var textToAnalyze = $"{request.Headline} {request.Description}";

            if (request.Keywords != null && request.Keywords.Any())
            {
                textToAnalyze += " " + string.Join(" ", request.Keywords);
            }

            var result = await _clickbaitService.DetectClickbaitAsync(textToAnalyze);
            return Ok(JsonDocument.Parse(result).RootElement);
        }
        catch (Exception ex)
        {
            return StatusCode(500, new { error = ex.Message });
        }
    }
}