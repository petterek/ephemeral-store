using System.Net;
using System.Net.Http.Json;
using System.Text.Json;
using Microsoft.AspNetCore.Mvc.Testing;

namespace EphemralStorage.Tests;

public class HandlerTests : IClassFixture<WebApplicationFactory<Program>>
{
    private readonly HttpClient _client;
    private readonly JsonSerializerOptions _jsonOpts = new()
    {
        PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower
    };

    public HandlerTests(WebApplicationFactory<Program> factory)
    {
        _client = factory.CreateClient();
    }

    [Fact]
    public async Task InsertAndRead()
    {
        var resp = await _client.PostAsJsonAsync("/keys",
            new { sender = "alice", datatype = "password", value = "secret", ttl = 60 });
        Assert.Equal(HttpStatusCode.Created, resp.StatusCode);

        var body = await resp.Content.ReadFromJsonAsync<JsonElement>();
        var key = body.GetProperty("key").GetString()!;
        Assert.NotEmpty(key);

        // First read succeeds
        resp = await _client.GetAsync($"/keys/{key}");
        Assert.Equal(HttpStatusCode.OK, resp.StatusCode);

        // Second read returns 404
        resp = await _client.GetAsync($"/keys/{key}");
        Assert.Equal(HttpStatusCode.NotFound, resp.StatusCode);
    }

    [Fact]
    public async Task ReadMissing()
    {
        var resp = await _client.GetAsync("/keys/nope");
        Assert.Equal(HttpStatusCode.NotFound, resp.StatusCode);
    }

    [Fact]
    public async Task InsertBadRequest()
    {
        var resp = await _client.PostAsJsonAsync("/keys",
            new { sender = "", datatype = "", value = "v" });
        Assert.Equal(HttpStatusCode.BadRequest, resp.StatusCode);
    }
}
