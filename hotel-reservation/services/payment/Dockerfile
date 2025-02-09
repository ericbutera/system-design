FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS base
WORKDIR /app
EXPOSE 8080

ENV ASPNETCORE_URLS=http://+:8080

USER app
FROM --platform=$BUILDPLATFORM mcr.microsoft.com/dotnet/sdk:8.0 AS build
ARG configuration=Release
WORKDIR /src
COPY ["payment.csproj", "./"]
RUN dotnet restore "payment.csproj"
COPY . .
WORKDIR "/src/."
RUN dotnet build "payment.csproj" -c $configuration -o /app/build

FROM build AS publish
ARG configuration=Release
RUN dotnet publish "payment.csproj" -c $configuration -o /app/publish /p:UseAppHost=false

FROM base AS final
CMD apt-get install curl
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "payment.dll"]
