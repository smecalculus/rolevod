package smecalculus.bezmen.construction;

import static org.springframework.context.annotation.ConfigurationCondition.ConfigurationPhase.REGISTER_BEAN;

import lombok.NonNull;
import org.springframework.context.annotation.ConditionContext;
import org.springframework.context.annotation.ConfigurationCondition;
import org.springframework.core.type.AnnotatedTypeMetadata;
import smecalculus.bezmen.configuration.ValidationMode;
import smecalculus.bezmen.configuration.ValidationProps;

class ValidationModeCondition implements ConfigurationCondition {

    @Override
    public boolean matches(ConditionContext context, AnnotatedTypeMetadata metadata) {
        var attributes = metadata.getAnnotationAttributes(ConditionalOnValidationMode.class.getName());
        var mode = (ValidationMode) attributes.get("value");
        var props = context.getBeanFactory().getBean(ValidationProps.class);
        return mode == props.validationMode();
    }

    @Override
    public @NonNull ConfigurationPhase getConfigurationPhase() {
        return REGISTER_BEAN;
    }
}
