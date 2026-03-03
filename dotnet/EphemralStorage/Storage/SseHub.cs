using System.Text.Json;
using System.Threading.Channels;
using EphemralStorage.Core.Storage;

namespace EphemralStorage.Web.Storage;

public class SseHub
{
    private readonly object _lock = new();
    private readonly List<Channel<byte[]>> _clients = new();
    private readonly EphemeralService _svc;

    public SseHub(EphemeralService svc) => _svc = svc;

    public void Broadcast()
    {
        var entries = _svc.ListEntries();
        var data = JsonSerializer.SerializeToUtf8Bytes(entries,
            new JsonSerializerOptions { PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower });

        lock (_lock)
        {
            var dead = new List<Channel<byte[]>>();
            foreach (var ch in _clients)
            {
                if (!ch.Writer.TryWrite(data))
                    dead.Add(ch);
            }
            foreach (var ch in dead)
                _clients.Remove(ch);
        }
    }

    public Channel<byte[]> Subscribe()
    {
        var ch = Channel.CreateBounded<byte[]>(8);
        lock (_lock) { _clients.Add(ch); }

        // Send initial state
        var entries = _svc.ListEntries();
        var data = JsonSerializer.SerializeToUtf8Bytes(entries,
            new JsonSerializerOptions { PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower });
        ch.Writer.TryWrite(data);

        return ch;
    }

    public void Unsubscribe(Channel<byte[]> ch)
    {
        lock (_lock) { _clients.Remove(ch); }
    }
}
