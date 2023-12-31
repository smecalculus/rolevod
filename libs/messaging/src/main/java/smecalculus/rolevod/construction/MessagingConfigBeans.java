package smecalculus.rolevod.construction;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;
import smecalculus.rolevod.configuration.MessagingDm;
import smecalculus.rolevod.configuration.MessagingEm;
import smecalculus.rolevod.configuration.MessagingPropsMapper;
import smecalculus.rolevod.configuration.MessagingPropsMapperImpl;
import smecalculus.rolevod.configuration.PropsKeeper;
import smecalculus.rolevod.validation.EdgeValidator;

@PropertySource("classpath:messaging.properties")
@Configuration(proxyBeanMethods = false)
public class MessagingConfigBeans {

    private static final Logger LOG = LoggerFactory.getLogger(MessagingConfigBeans.class);

    @Bean
    MessagingPropsMapper messagingPropsMapper() {
        return new MessagingPropsMapperImpl();
    }

    @Bean
    MessagingDm.MessagingProps messagingProps(
            PropsKeeper keeper, EdgeValidator validator, MessagingPropsMapper mapper) {
        var propsEdge = keeper.read("solution.messaging", MessagingEm.MessagingProps.class);
        validator.validate(propsEdge);
        LOG.info("Read {}", propsEdge);
        return mapper.toDomain(propsEdge);
    }
}
