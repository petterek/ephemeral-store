using System.Diagnostics;
using EphemralStorage.Core.Storage;
using EphemralStorage.Memory.Storage;

namespace EphemralStorage.Tests;

public class PerformanceTests
{
    private const int Iterations = 1_000_000;
    private const int WarmupIterations = 10_000;
    private const int Runs = 3;

    [Fact]
    public void BenchInsert()
    {
        RunBench("Insert", Iterations, (svc, i) =>
        {
            svc.InsertKeyValue($"key-{i}", "value", 60);
        });
    }

    [Fact]
    public void BenchInsertAndRead()
    {
        RunBench("InsertAndRead", Iterations, (svc, i) =>
        {
            var key = $"key-{i}";
            svc.InsertKeyValue(key, "value", 60);
            svc.ReadValue(key);
        });
    }

    [Fact]
    public void BenchConcurrentInsertAndRead()
    {
        for (var run = 0; run < Runs; run++)
        {
            var svc = new EphemeralService(new MemoryStore());
            const int count = 100_000;

            // Warmup
            Parallel.For(0, 1000, i =>
            {
                svc.InsertKeyValue($"w-{i}", "v", 60);
                svc.ReadValue($"w-{i}");
            });

            svc = new EphemeralService(new MemoryStore());
            var sw = Stopwatch.StartNew();
            Parallel.For(0, count, i =>
            {
                var key = $"conc-{i}";
                svc.InsertKeyValue(key, "value", 60);
                svc.ReadValue(key);
            });
            sw.Stop();
            var nsPerOp = sw.Elapsed.TotalNanoseconds / count;
            Console.WriteLine($"ConcurrentInsertAndRead run {run + 1}: {nsPerOp:F1} ns/op ({count} ops in {sw.ElapsedMilliseconds} ms)");
        }
    }

    [Fact]
    public void BenchList()
    {
        const int entries = 1000;
        const int iterations = 50_000;

        for (var run = 0; run < Runs; run++)
        {
            var store = new MemoryStore();
            var svc = new EphemeralService(store);
            for (var i = 0; i < entries; i++)
                svc.InsertKeyValue($"key-{i}", "value", 300);

            // Warmup
            for (var i = 0; i < 100; i++)
                store.List();

            var sw = Stopwatch.StartNew();
            for (var i = 0; i < iterations; i++)
                store.List();
            sw.Stop();
            var nsPerOp = sw.Elapsed.TotalNanoseconds / iterations;
            Console.WriteLine($"List({entries} entries) run {run + 1}: {nsPerOp:F1} ns/op ({iterations} ops in {sw.ElapsedMilliseconds} ms)");
        }
    }

    private static void RunBench(string name, int iterations, Action<EphemeralService, int> action)
    {
        for (var run = 0; run < Runs; run++)
        {
            var svc = new EphemeralService(new MemoryStore());

            // Warmup
            for (var i = 0; i < WarmupIterations; i++)
                action(svc, i);

            svc = new EphemeralService(new MemoryStore());
            var sw = Stopwatch.StartNew();
            for (var i = 0; i < iterations; i++)
                action(svc, i);
            sw.Stop();
            var nsPerOp = sw.Elapsed.TotalNanoseconds / iterations;
            Console.WriteLine($"{name} run {run + 1}: {nsPerOp:F1} ns/op ({iterations} ops in {sw.ElapsedMilliseconds} ms)");
        }
    }
}
