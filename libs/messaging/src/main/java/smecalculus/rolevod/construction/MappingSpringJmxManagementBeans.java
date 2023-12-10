package smecalculus.rolevod.construction;

import static smecalculus.rolevod.configuration.MessagingDm.MappingMode.SPRING_MANAGEMENT;
import static smecalculus.rolevod.configuration.MessagingDm.ProtocolMode.JMX;

import org.springframework.context.annotation.Configuration;

@ConditionalOnMessagingProtocolModes(JMX)
@ConditionalOnMessagingMappingModes(SPRING_MANAGEMENT)
@Configuration(proxyBeanMethods = false)
public class MappingSpringJmxManagementBeans {}
