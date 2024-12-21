using System.Text;
using System.Text.Json;
using Microsoft.Extensions.Configuration;
using XMAS.Domain;

namespace XMAS.Services;

public class ClickbaitService : IClickbaitDetectionService
{
    private readonly HttpClient _httpClient;
    private readonly string _key;
    private readonly string _mlService;

    public ClickbaitService(HttpClient httpClient, IConfiguration configuration)
    {
        _httpClient = httpClient;
        _key = configuration["MlServices:AuthKey"]!;
        _mlService = configuration["MlServices:ClickbaitMLService"]!;
    }

    public async Task<string> DetectClickbaitAsync(string text)
    {
        _httpClient.DefaultRequestHeaders.Clear();
        _httpClient.DefaultRequestHeaders.Add("Authorization", $"Bearer {_key}");

        var payload = new { inputs = text };
        var content = new StringContent(JsonSerializer.Serialize(payload), Encoding.UTF8, "application/json");

        var response = await _httpClient.PostAsync(_mlService, content);

        if (!response.IsSuccessStatusCode)
        {
            var errorResponse = await response.Content.ReadAsStringAsync();
            throw new Exception($"API call failed: {response.StatusCode} - {errorResponse}");
        }

        var jsonResponse = await response.Content.ReadAsStringAsync();
        return jsonResponse;
    }
}