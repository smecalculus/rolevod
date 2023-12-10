package smecalculus.bezmen.construction;

import com.typesafe.config.ConfigFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.env.Environment;
import smecalculus.bezmen.configuration.ConfigMappingMode;
import smecalculus.bezmen.configuration.ConfigProtocolMode;
import smecalculus.bezmen.configuration.PropsKeeper;
import smecalculus.bezmen.configuration.PropsKeeperLightbendConfig;
import smecalculus.bezmen.configuration.PropsKeeperSpringConfig;

@Configuration(proxyBeanMethods = false)
public class ConfigBeans {

    @Bean
    @ConditionalOnConfigProtocolMode(ConfigProtocolMode.FILE_SYSTEM)
    @ConditionalOnConfigMappingMode(ConfigMappingMode.LIGHTBEND_CONFIG)
    PropsKeeper propsKeeperLightbendConfig() {
        return new PropsKeeperLightbendConfig(ConfigFactory.load());
    }

    @Bean
    @ConditionalOnConfigProtocolMode(ConfigProtocolMode.FILE_SYSTEM)
    @ConditionalOnConfigMappingMode(ConfigMappingMode.SPRING_CONFIG)
    PropsKeeper propsKeeperSpringConfig(Environment environment) {
        return new PropsKeeperSpringConfig(environment);
    }
}
