package smecalculus.bezmen.construction;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;
import org.springframework.context.annotation.Conditional;
import smecalculus.bezmen.configuration.ConfigMappingMode;

@Target({ElementType.TYPE, ElementType.METHOD})
@Retention(RetentionPolicy.RUNTIME)
@Conditional(ConfigMappingModeCondition.class)
public @interface ConditionalOnConfigMappingMode {
    ConfigMappingMode value();
}
