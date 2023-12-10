package smecalculus.bezmen.configuration;

public interface PropsKeeper {
    <T> T read(String key, Class<T> type);
}
