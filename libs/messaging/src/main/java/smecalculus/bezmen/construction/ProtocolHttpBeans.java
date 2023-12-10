package smecalculus.bezmen.construction;

import static smecalculus.bezmen.configuration.MessagingDm.ProtocolMode.HTTP;

import org.springframework.boot.autoconfigure.ImportAutoConfiguration;
import org.springframework.boot.autoconfigure.web.servlet.DispatcherServletAutoConfiguration;
import org.springframework.boot.autoconfigure.web.servlet.ServletWebServerFactoryAutoConfiguration;
import org.springframework.context.annotation.Configuration;

@ConditionalOnMessagingProtocolModes(HTTP)
@ImportAutoConfiguration({ServletWebServerFactoryAutoConfiguration.class, DispatcherServletAutoConfiguration.class})
@Configuration(proxyBeanMethods = false)
public class ProtocolHttpBeans {}
