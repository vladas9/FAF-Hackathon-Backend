using XMAS.Domain;
using XMAS.Services;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

// Register custom services
builder.Services.AddScoped<ISentimentService, SentimentService>();
builder.Services.AddScoped<IClickbaitDetectionService, ClickbaitService>();

// Register background services
builder.Services.AddHttpClient<ISentimentService, SentimentService>();
builder.Services.AddHttpClient<IClickbaitDetectionService, ClickbaitService>();



var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();
app.UseAuthorization();
app.MapControllers();

app.Run();