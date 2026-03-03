using System.Text;
using System.Text.Json;
using EphemralStorage.Core.Storage;
using EphemralStorage.Memory.Storage;
using EphemralStorage.Web.Storage;

var builder = WebApplication.CreateBuilder(args);

var store = new MemoryStore();
var svc = new EphemeralService(store);
var hub = new SseHub(svc);

builder.Services.AddSingleton(svc);
builder.Services.AddSingleton(hub);

var app = builder.Build();

app.UseStaticFiles();

var jsonOpts = new JsonSerializerOptions { PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower };

app.MapPost("/keys", (InsertRequest req) =>
{
    if (string.IsNullOrEmpty(req.Sender) || string.IsNullOrEmpty(req.Datatype))
        return Results.BadRequest(new ErrorResponse("sender and datatype are required"));

    var ttl = req.Ttl <= 0 ? 60 : req.Ttl;
    var raw = $"{req.Sender}.{req.Datatype}.{DateTimeOffset.UtcNow.ToUnixTimeMilliseconds()}";
    var key = Convert.ToBase64String(Encoding.UTF8.GetBytes(raw));

    svc.InsertKeyValue(key, req.Value, ttl);
    hub.Broadcast();

    return Results.Json(new InsertResponse(key), jsonOpts, statusCode: 201);
});

app.MapGet("/keys/{key}", (string key) =>
{
    var (value, found) = svc.ReadValue(key);
    if (!found)
        return Results.Json(new ErrorResponse("key not found"), jsonOpts, statusCode: 404);

    hub.Broadcast();
    return Results.Json(new ReadResponse(value), jsonOpts);
});

app.MapGet("/events", async (HttpContext ctx, CancellationToken ct) =>
{
    ctx.Response.ContentType = "text/event-stream";
    ctx.Response.Headers.CacheControl = "no-cache";
    ctx.Response.Headers.Connection = "keep-alive";

    var ch = hub.Subscribe();
    try
    {
        await foreach (var data in ch.Reader.ReadAllAsync(ct))
        {
            await ctx.Response.WriteAsync($"data: ", ct);
            await ctx.Response.Body.WriteAsync(data, ct);
            await ctx.Response.WriteAsync("\n\n", ct);
            await ctx.Response.Body.FlushAsync(ct);
        }
    }
    finally
    {
        hub.Unsubscribe(ch);
    }
});

app.MapGet("/isalive", () => Results.Ok());
app.MapGet("/isready", () => Results.Ok());

app.MapFallbackToFile("index.html");

app.Run();

public partial class Program { }
