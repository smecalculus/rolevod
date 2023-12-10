package smecalculus.rolevod.construction;

import com.typesafe.config.ConfigFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.env.Environment;
import smecalculus.rolevod.configuration.ConfigMappingMode;
import smecalculus.rolevod.configuration.ConfigProtocolMode;
import smecalculus.rolevod.configuration.PropsKeeper;
import smecalculus.rolevod.configuration.PropsKeeperLightbendConfig;
import smecalculus.rolevod.configuration.PropsKeeperSpringConfig;

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
