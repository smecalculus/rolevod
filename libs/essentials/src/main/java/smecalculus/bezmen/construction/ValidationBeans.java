package smecalculus.bezmen.construction;

import static java.util.Objects.requireNonNull;
import static smecalculus.bezmen.configuration.ValidationMode.HIBERNATE_VALIDATOR;

import jakarta.validation.Validation;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;
import smecalculus.bezmen.configuration.PropsKeeper;
import smecalculus.bezmen.configuration.ValidationProps;
import smecalculus.bezmen.configuration.ValidationPropsEdge;
import smecalculus.bezmen.validation.EdgeValidator;
import smecalculus.bezmen.validation.EdgeValidatorHibernateValidator;
import smecalculus.bezmen.validation.ValidationPropsMapper;
import smecalculus.bezmen.validation.ValidationPropsMapperImpl;

@PropertySource("classpath:validation.properties")
@Configuration(proxyBeanMethods = false)
public class ValidationBeans {

    private static final Logger LOG = LoggerFactory.getLogger(ValidationBeans.class);

    @Bean
    ValidationPropsMapper validationPropsMapper() {
        return new ValidationPropsMapperImpl();
    }

    @Bean
    ValidationProps validationProps(PropsKeeper keeper, ValidationPropsMapper mapper) {
        var propsEdge = keeper.read("bezmen.validation", ValidationPropsEdge.class);
        requireNonNull(propsEdge.getMode(), "validation mode must not be null");
        LOG.info("Read {}", propsEdge);
        return mapper.toDomain(propsEdge);
    }

    @Bean
    @ConditionalOnValidationMode(HIBERNATE_VALIDATOR)
    EdgeValidator edgeValidatorHibernateValidator() {
        try (var factory = Validation.buildDefaultValidatorFactory()) {
            return new EdgeValidatorHibernateValidator(factory.getValidator());
        }
    }
}
