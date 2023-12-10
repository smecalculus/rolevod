package smecalculus.rolevod.construction;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;
import org.springframework.context.annotation.Conditional;
import smecalculus.rolevod.configuration.ValidationDm.ValidationMode;

@Target({ElementType.TYPE, ElementType.METHOD})
@Retention(RetentionPolicy.RUNTIME)
@Conditional(ValidationModeCondition.class)
public @interface ConditionalOnValidationMode {
    ValidationMode value();
}
