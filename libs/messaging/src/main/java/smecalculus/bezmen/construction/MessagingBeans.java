package smecalculus.bezmen.construction;

import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;

@Import({
    MessagingConfigBeans.class,
    MappingSpringWebMvcBeans.class,
    MappingSpringWebManagementBeans.class,
    MappingSpringJmxManagementBeans.class,
    ProtocolHttpBeans.class,
    ProtocolJmxBeans.class
})
@Configuration(proxyBeanMethods = false)
public class MessagingBeans {}
