namespace EphemralStorage.Core.Storage;

public record KeyValue(string Key, string Value, int ExpiresIn);
