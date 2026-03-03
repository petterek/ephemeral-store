namespace EphemralStorage.Core.Storage;

public interface IStore
{
    void Set(string key, string value, int ttlSeconds);
    (string value, bool found) GetAndDelete(string key);
    List<KeyValue> List();
}
