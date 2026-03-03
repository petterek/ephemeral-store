using EphemralStorage.Core.Storage;
using EphemralStorage.Memory.Storage;

namespace EphemralStorage.Tests;

public class ServiceTests
{
    [Fact]
    public void InsertAndReadOnce()
    {
        var svc = new EphemeralService(new MemoryStore());
        svc.InsertKeyValue("k1", "secret", 60);

        var (val, found) = svc.ReadValue("k1");
        Assert.True(found);
        Assert.Equal("secret", val);

        var (_, found2) = svc.ReadValue("k1");
        Assert.False(found2);
    }

    [Fact]
    public async Task ReadExpired()
    {
        var svc = new EphemeralService(new MemoryStore());
        svc.InsertKeyValue("k2", "gone", 1);

        await Task.Delay(1100);

        var (_, found) = svc.ReadValue("k2");
        Assert.False(found);
    }

    [Fact]
    public void ReadMissing()
    {
        var svc = new EphemeralService(new MemoryStore());
        var (_, found) = svc.ReadValue("nope");
        Assert.False(found);
    }
}
