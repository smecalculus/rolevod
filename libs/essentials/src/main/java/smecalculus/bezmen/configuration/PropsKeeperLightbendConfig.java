package smecalculus.bezmen.configuration;

import com.typesafe.config.Config;
import com.typesafe.config.ConfigBeanFactory;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;

@RequiredArgsConstructor
public class PropsKeeperLightbendConfig implements PropsKeeper {

    @NonNull
    private Config config;

    @Override
    public <T> T read(String key, Class<T> type) {
        return ConfigBeanFactory.create(config.getConfig(key), type);
    }
}
