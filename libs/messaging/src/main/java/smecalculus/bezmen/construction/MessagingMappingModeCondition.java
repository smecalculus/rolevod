package smecalculus.bezmen.construction;

import static org.springframework.context.annotation.ConfigurationCondition.ConfigurationPhase.REGISTER_BEAN;

import java.util.stream.Stream;
import lombok.NonNull;
import org.springframework.context.annotation.ConditionContext;
import org.springframework.context.annotation.ConfigurationCondition;
import org.springframework.core.type.AnnotatedTypeMetadata;
import smecalculus.bezmen.configuration.MessagingDm.MappingMode;
import smecalculus.bezmen.configuration.MessagingDm.MessagingProps;

class MessagingMappingModeCondition implements ConfigurationCondition {

    @Override
    public boolean matches(ConditionContext context, AnnotatedTypeMetadata metadata) {
        var attributes = metadata.getAnnotationAttributes(ConditionalOnMessagingMappingModes.class.getName());
        var expectedModes = (MappingMode[]) attributes.get("value");
        var props = context.getBeanFactory().getBean(MessagingProps.class);
        var actualModes = props.mappingProps().mappingModes();
        return Stream.of(expectedModes).anyMatch(actualModes::contains);
    }

    @Override
    public @NonNull ConfigurationPhase getConfigurationPhase() {
        return REGISTER_BEAN;
    }
}
