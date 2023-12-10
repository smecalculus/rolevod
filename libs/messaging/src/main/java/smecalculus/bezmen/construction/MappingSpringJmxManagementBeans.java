package smecalculus.bezmen.construction;

import static smecalculus.bezmen.configuration.MessagingDm.MappingMode.SPRING_MANAGEMENT;
import static smecalculus.bezmen.configuration.MessagingDm.ProtocolMode.JMX;

import org.springframework.context.annotation.Configuration;

@ConditionalOnMessagingProtocolModes(JMX)
@ConditionalOnMessagingMappingModes(SPRING_MANAGEMENT)
@Configuration(proxyBeanMethods = false)
public class MappingSpringJmxManagementBeans {}
