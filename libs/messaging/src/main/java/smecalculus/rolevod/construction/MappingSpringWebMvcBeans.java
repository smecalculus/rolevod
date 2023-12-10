package smecalculus.rolevod.construction;

import static smecalculus.rolevod.configuration.MessagingDm.MappingMode.SPRING_MVC;
import static smecalculus.rolevod.configuration.MessagingDm.ProtocolMode.HTTP;

import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.EnableWebMvc;

@ConditionalOnMessagingProtocolModes(HTTP)
@ConditionalOnMessagingMappingModes(SPRING_MVC)
@EnableWebMvc
@Configuration(proxyBeanMethods = false)
public class MappingSpringWebMvcBeans {}
