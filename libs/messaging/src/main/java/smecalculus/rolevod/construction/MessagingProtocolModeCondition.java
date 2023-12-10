package smecalculus.rolevod.construction;

import static org.springframework.context.annotation.ConfigurationCondition.ConfigurationPhase.REGISTER_BEAN;

import java.util.stream.Stream;
import lombok.NonNull;
import org.springframework.context.annotation.ConditionContext;
import org.springframework.context.annotation.ConfigurationCondition;
import org.springframework.core.type.AnnotatedTypeMetadata;
import smecalculus.rolevod.configuration.MessagingDm.MessagingProps;
import smecalculus.rolevod.configuration.MessagingDm.ProtocolMode;

class MessagingProtocolModeCondition implements ConfigurationCondition {

    @Override
    public boolean matches(ConditionContext context, AnnotatedTypeMetadata metadata) {
        var attributes = metadata.getAnnotationAttributes(ConditionalOnMessagingProtocolModes.class.getName());
        var expectedModes = (ProtocolMode[]) attributes.get("value");
        var props = context.getBeanFactory().getBean(MessagingProps.class);
        var actualModes = props.protocolProps().protocolModes();
        return Stream.of(expectedModes).anyMatch(actualModes::contains);
    }

    @Override
    public @NonNull ConfigurationPhase getConfigurationPhase() {
        return REGISTER_BEAN;
    }
}
