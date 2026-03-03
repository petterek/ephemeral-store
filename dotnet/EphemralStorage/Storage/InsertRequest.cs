namespace EphemralStorage.Web.Storage;

public record InsertRequest(string Sender, string Datatype, string Value, int Ttl);
