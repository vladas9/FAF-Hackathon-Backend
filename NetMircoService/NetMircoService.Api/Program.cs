using NetMircoService.Domain;
using NetMircoService.Services;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

// Register custom services
builder.Services.AddScoped<ISentimentService, SentimentService>();
builder.Services.AddScoped<IClickbaitService, ClickbaitService>();
builder.Services.AddScoped<ILanguageService, LanguageService>();
builder.Services.AddScoped<ISimilarityService, SimilarityService>();

// Register background services
builder.Services.AddHttpClient<ISentimentService, SentimentService>();
builder.Services.AddHttpClient<IClickbaitService, ClickbaitService>();
builder.Services.AddHttpClient<ILanguageService, LanguageService>();
builder.Services.AddHttpClient<ISimilarityService, SimilarityService>();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();
// app.UseAuthorization();
app.MapControllers();

app.Run();