namespace EphemralStorage.Core.Storage;

public class EphemeralService
{
    private readonly IStore _store;

    public EphemeralService(IStore store) => _store = store;

    public void InsertKeyValue(string key, string value, int ttl) =>
        _store.Set(key, value, ttl);

    public (string value, bool found) ReadValue(string key) =>
        _store.GetAndDelete(key);

    public List<KeyValue> ListEntries() => _store.List();
}
