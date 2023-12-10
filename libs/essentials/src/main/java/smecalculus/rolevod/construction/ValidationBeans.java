package smecalculus.rolevod.construction;

import static java.util.Objects.requireNonNull;
import static smecalculus.rolevod.configuration.ValidationDm.ValidationMode.HIBERNATE_VALIDATOR;

import jakarta.validation.Validation;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;
import smecalculus.rolevod.configuration.PropsKeeper;
import smecalculus.rolevod.configuration.ValidationDm;
import smecalculus.rolevod.configuration.ValidationEm;
import smecalculus.rolevod.validation.EdgeValidator;
import smecalculus.rolevod.validation.EdgeValidatorHibernateValidator;
import smecalculus.rolevod.validation.ValidationPropsMapper;
import smecalculus.rolevod.validation.ValidationPropsMapperImpl;

@PropertySource("classpath:validation.properties")
@Configuration(proxyBeanMethods = false)
public class ValidationBeans {

    private static final Logger LOG = LoggerFactory.getLogger(ValidationBeans.class);

    @Bean
    ValidationPropsMapper validationPropsMapper() {
        return new ValidationPropsMapperImpl();
    }

    @Bean
    ValidationDm.ValidationProps validationProps(PropsKeeper keeper, ValidationPropsMapper mapper) {
        var validationProps = keeper.read("solution.validation", ValidationEm.ValidationProps.class);
        requireNonNull(validationProps.getMode(), "validation mode must not be null");
        LOG.info("Read {}", validationProps);
        return mapper.toDomain(validationProps);
    }

    @Bean
    @ConditionalOnValidationMode(HIBERNATE_VALIDATOR)
    EdgeValidator edgeValidatorHibernateValidator() {
        try (var factory = Validation.buildDefaultValidatorFactory()) {
            return new EdgeValidatorHibernateValidator(factory.getValidator());
        }
    }
}
