package smecalculus.rolevod.configuration;

public interface PropsKeeper {
    <T> T read(String key, Class<T> type);
}
