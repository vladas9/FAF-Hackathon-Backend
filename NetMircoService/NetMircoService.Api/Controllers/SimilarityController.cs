using Microsoft.AspNetCore.Mvc;
using NetMircoService.Api.Dtos.Models;
using NetMircoService.Domain;

namespace XMAS.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class SimilarityController : ControllerBase
{
    private readonly ISimilarityService _similarityService;

    public SimilarityController(ISimilarityService similarityService)
    {
        _similarityService = similarityService;
    }

    [HttpPost("compare")]
    public async Task<IActionResult> CompareTexts([FromBody] SimilarityRequestDto request)
    {
        if (string.IsNullOrWhiteSpace(request.Text1) || string.IsNullOrWhiteSpace(request.Text2))
        {
            return BadRequest(new { error = "Both Text1 and Text2 cannot be empty." });
        }

        try
        {
            var score = await _similarityService.CalculateSimilarityAsync(request.Text1, request.Text2);
            return Ok(new { similarityScore = score });
        }
        catch (Exception ex)
        {
            return StatusCode(500, new { error = ex.Message });
        }
    }
}