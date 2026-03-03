using BenchmarkDotNet.Attributes;
using EphemralStorage.Core.Storage;
using EphemralStorage.Memory.Storage;

namespace EphemralStorage.Tests;

[MemoryDiagnoser]
[ShortRunJob]
public class StorageBenchmarks
{
    private EphemeralService _svc = null!;
    private MemoryStore _store = null!;

    [GlobalSetup]
    public void Setup()
    {
        _store = new MemoryStore();
        _svc = new EphemeralService(_store);
    }

    [Benchmark]
    public void Insert()
    {
        _svc.InsertKeyValue("key", "value", 60);
    }

    [Benchmark]
    public void InsertAndRead()
    {
        _svc.InsertKeyValue("key", "value", 60);
        _svc.ReadValue("key");
    }

    [Benchmark]
    public void ReadMissing()
    {
        _svc.ReadValue("nonexistent");
    }

    [Benchmark]
    [Arguments(100)]
    [Arguments(1000)]
    public void ListEntries(int count)
    {
        for (var i = 0; i < count; i++)
            _svc.InsertKeyValue($"list-key-{i}", "value", 300);

        _store.List();
    }
}

[MemoryDiagnoser]
[ShortRunJob]
public class ConcurrencyBenchmarks
{
    private EphemeralService _svc = null!;

    [GlobalSetup]
    public void Setup()
    {
        _svc = new EphemeralService(new MemoryStore());
    }

    [Benchmark]
    [Arguments(100)]
    [Arguments(1000)]
    public void ConcurrentInsertAndRead(int count)
    {
        var tasks = new Task[count];
        for (var i = 0; i < count; i++)
        {
            var key = $"conc-{i}";
            tasks[i] = Task.Run(() =>
            {
                _svc.InsertKeyValue(key, "value", 60);
                _svc.ReadValue(key);
            });
        }
        Task.WaitAll(tasks);
    }
}
