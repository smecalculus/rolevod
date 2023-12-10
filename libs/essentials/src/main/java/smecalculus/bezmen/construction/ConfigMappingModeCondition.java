package smecalculus.bezmen.construction;

import static org.springframework.context.annotation.ConfigurationCondition.ConfigurationPhase.REGISTER_BEAN;

import lombok.NonNull;
import org.springframework.context.annotation.ConditionContext;
import org.springframework.context.annotation.ConfigurationCondition;
import org.springframework.core.type.AnnotatedTypeMetadata;
import smecalculus.bezmen.configuration.ConfigMappingMode;

class ConfigMappingModeCondition implements ConfigurationCondition {

    @Override
    public boolean matches(ConditionContext context, AnnotatedTypeMetadata metadata) {
        var attributes = metadata.getAnnotationAttributes(ConditionalOnConfigMappingMode.class.getName());
        var expectedMode = (ConfigMappingMode) attributes.get("value");
        var actualMode = context.getEnvironment()
                .getProperty("bezmen.config.mapping.mode", ConfigMappingMode.LIGHTBEND_CONFIG.name());
        return expectedMode.name().equalsIgnoreCase(actualMode);
    }

    @Override
    public @NonNull ConfigurationPhase getConfigurationPhase() {
        return REGISTER_BEAN;
    }
}
