FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS base
WORKDIR /app
EXPOSE 8080
EXPOSE 8081
    
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src
COPY ["NetMircoService.Api/NetMircoService.Api.csproj", "NetMircoService.Api/"]
COPY ["NetMircoService.Domain/NetMircoService.Domain.csproj", "NetMircoService.Domain/"]
COPY ["NetMircoService.Services/NetMircoService.Services.csproj", "NetMircoService.Services/"]
RUN dotnet restore "NetMircoService.Api/NetMircoService.Api.csproj"
COPY . .
WORKDIR "/src/NetMircoService.Api"
RUN dotnet build "NetMircoService.Api.csproj" -c $BUILD_CONFIGURATION -o /app/build

FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "NetMircoService.Api.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "NetMircoService.Api.dll"]
