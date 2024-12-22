using System.Text;
using System.Text.Json;
using Microsoft.Extensions.Configuration;
using NetMircoService.Domain;

namespace NetMircoService.Services;

public class SimilarityService : ISimilarityService
{
    private readonly HttpClient _httpClient;
    private readonly string _key;
    private readonly string _mlService;

    public SimilarityService(HttpClient httpClient, IConfiguration configuration)
    {
        _httpClient = httpClient;
        _key = configuration["MlServices:AuthKey"]!;
        _mlService = configuration["MlServices:SimilarityMLService"]!;
    }

    public async Task<string> CalculateSimilarityAsync(string text1, string text2)
    {
        _httpClient.DefaultRequestHeaders.Clear();
        _httpClient.DefaultRequestHeaders.Add("Authorization", $"Bearer {_key}");

        var payload = new
        {
            inputs = new
            {
                source_sentence = text1,
                sentences = new[] { text2 }
            }
        };
        
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