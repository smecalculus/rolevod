package smecalculus.bezmen.construction;

import static org.mockito.Mockito.mock;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;
import org.springframework.test.web.servlet.client.MockMvcWebTestClient;
import smecalculus.bezmen.core.SepulkaService;
import smecalculus.bezmen.messaging.SepulkaClient;
import smecalculus.bezmen.messaging.SepulkaClientImpl;
import smecalculus.bezmen.messaging.SepulkaClientSpringWebTest;
import smecalculus.bezmen.messaging.SepulkaMessageMapperImpl;
import smecalculus.bezmen.messaging.springmvc.SepulkaController;
import smecalculus.bezmen.validation.EdgeValidator;

@Import({ConfigBeans.class, ValidationBeans.class})
@Configuration(proxyBeanMethods = false)
public class SepulkaClientBeans {

    @Bean
    public SepulkaService sepulkaService() {
        return mock(SepulkaService.class);
    }

    @Bean
    SepulkaClient internalClient(EdgeValidator validator, SepulkaService service) {
        var mapper = new SepulkaMessageMapperImpl();
        return new SepulkaClientImpl(validator, mapper, service);
    }

    @Bean
    SepulkaClient externalClient(SepulkaClient internalClient) {
        var client = MockMvcWebTestClient.bindToController(new SepulkaController(internalClient))
                .build();
        return new SepulkaClientSpringWebTest(client);
    }
}
