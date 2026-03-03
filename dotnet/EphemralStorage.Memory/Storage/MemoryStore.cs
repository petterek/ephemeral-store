using EphemralStorage.Core.Storage;

namespace EphemralStorage.Memory.Storage;

public class MemoryStore : IStore
{
    private readonly object _lock = new();
    private readonly Dictionary<string, Entry> _entries = new();

    private record Entry(string Value, DateTime ExpiresAt);

    public void Set(string key, string value, int ttlSeconds)
    {
        lock (_lock)
        {
            _entries[key] = new Entry(value, DateTime.UtcNow.AddSeconds(ttlSeconds));
        }
    }

    public (string value, bool found) GetAndDelete(string key)
    {
        lock (_lock)
        {
            if (!_entries.TryGetValue(key, out var entry))
                return ("", false);

            _entries.Remove(key);

            if (DateTime.UtcNow > entry.ExpiresAt)
                return ("", false);

            return (entry.Value, true);
        }
    }

    public List<KeyValue> List()
    {
        lock (_lock)
        {
            var now = DateTime.UtcNow;
            var expired = _entries.Where(e => now > e.Value.ExpiresAt).Select(e => e.Key).ToList();
            foreach (var key in expired)
                _entries.Remove(key);

            return _entries.Select(e => new KeyValue(
                e.Key,
                e.Value.Value,
                (int)(e.Value.ExpiresAt - now).TotalSeconds
            )).ToList();
        }
    }
}
