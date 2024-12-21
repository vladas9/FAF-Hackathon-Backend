using System.Text;
using System.Text.Json;
using Microsoft.Extensions.Configuration;
using XMAS.Domain;

namespace XMAS.Services;

public class SentimentService : ISentimentService
{
    private readonly HttpClient _httpClient;
    private readonly string _apiKey;

    public SentimentService(HttpClient httpClient, IConfiguration configuration)
    {
        _httpClient = httpClient;
        // https://huggingface.co/settings/tokens
        _apiKey = configuration["HuggingFace:ApiKey"]!;
    }

    public string AnalyzeSentiment(string text)
    {
        return CallSentimentApiAsync(text).GetAwaiter().GetResult();
    }

    private async Task<string> CallSentimentApiAsync(string text)
    {
        const int maxTokens = 512;

        // Truncate the text if it exceeds the maximum token limit
        if (text.Length > maxTokens)
        {
            text = text.Substring(0, maxTokens);
        }

        var apiUrl = "https://api-inference.huggingface.co/models/distilbert-base-uncased-finetuned-sst-2-english";

        _httpClient.DefaultRequestHeaders.Clear();
        _httpClient.DefaultRequestHeaders.Add("Authorization", $"Bearer {_apiKey}");

        var payload = new { inputs = text };
        var content = new StringContent(JsonSerializer.Serialize(payload), Encoding.UTF8, "application/json");

        var response = await _httpClient.PostAsync(apiUrl, content);

        if (!response.IsSuccessStatusCode)
        {
            var errorResponse = await response.Content.ReadAsStringAsync();
            throw new Exception($"API call failed: {response.StatusCode} - {errorResponse}");
        }

        var jsonResponse = await response.Content.ReadAsStringAsync();
        return jsonResponse;
    }

}