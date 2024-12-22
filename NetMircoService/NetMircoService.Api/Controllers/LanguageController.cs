using System.Text.Json;
using Microsoft.AspNetCore.Mvc;
using NetMircoService.Api.Dtos.Models;
using NetMircoService.Domain;

namespace XMAS.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class LanguageController : ControllerBase
{
    private readonly ILanguageService _languageService;

    public LanguageController(ILanguageService languageService)
    {
        _languageService = languageService;
    }

    [HttpPost("detect")]
    public async Task<IActionResult> DetectLanguage([FromBody] LanguageRequestDto request)
    {
        if (string.IsNullOrWhiteSpace(request.Text))
        {
            return BadRequest(new { error = "Text input cannot be empty." });
        }

        try
        {
            var result = await _languageService.DetectLanguageAsync(request.Text);
            return Ok(JsonDocument.Parse(result).RootElement);
        }
        catch (Exception ex)
        {
            return StatusCode(500, new { error = ex.Message });
        }
    }
}