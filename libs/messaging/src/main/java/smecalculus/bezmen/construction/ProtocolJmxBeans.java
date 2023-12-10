package smecalculus.bezmen.construction;

import static smecalculus.bezmen.configuration.MessagingDm.ProtocolMode.JMX;

import org.springframework.boot.autoconfigure.ImportAutoConfiguration;
import org.springframework.boot.autoconfigure.jmx.JmxAutoConfiguration;
import org.springframework.context.annotation.Configuration;

@ConditionalOnMessagingProtocolModes(JMX)
@ImportAutoConfiguration({JmxAutoConfiguration.class})
@Configuration(proxyBeanMethods = false)
public class ProtocolJmxBeans {}
