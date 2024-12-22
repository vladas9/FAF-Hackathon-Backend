using Microsoft.AspNetCore.Mvc;
using NetMircoService.Api.Dtos.Models;
using NetMircoService.Domain;

namespace XMAS.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class SentimentController : ControllerBase
{
    private readonly ISentimentService _sentimentService;

    public SentimentController(ISentimentService sentimentService)
    {
        _sentimentService = sentimentService;
    }

    [HttpPost("analyze")]
    public IActionResult AnalyzeSentiment([FromBody] SentimentRequestDto request)
    {
        if (string.IsNullOrWhiteSpace(request.Text))
        {
            return BadRequest(new { error = "Text input cannot be empty." });
        }

        try
        {
            var rawJsonResponse = _sentimentService.AnalyzeSentiment(request.Text);
            return Content(rawJsonResponse, "application/json");
        }
        catch (Exception ex)
        {
            return StatusCode(500, new { error = ex.Message });
        }
    }
}