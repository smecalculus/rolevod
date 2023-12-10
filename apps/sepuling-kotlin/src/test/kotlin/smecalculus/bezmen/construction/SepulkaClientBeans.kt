package smecalculus.bezmen.construction

import org.mockito.Mockito.mock
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.context.annotation.Import
import org.springframework.test.web.servlet.client.MockMvcWebTestClient
import smecalculus.bezmen.core.SepulkaService
import smecalculus.bezmen.messaging.SepulkaClient
import smecalculus.bezmen.messaging.SepulkaClientImpl
import smecalculus.bezmen.messaging.SepulkaClientSpringWebTest
import smecalculus.bezmen.messaging.SepulkaMessageMapperImpl
import smecalculus.bezmen.messaging.springmvc.SepulkaController
import smecalculus.bezmen.validation.EdgeValidator

@Import(ConfigBeans::class, ValidationBeans::class)
@Configuration(proxyBeanMethods = false)
class SepulkaClientBeans {
    @Bean
    fun sepulkaService(): SepulkaService {
        return mock()
    }

    @Bean
    fun internalClient(
        validator: EdgeValidator,
        service: SepulkaService,
    ): SepulkaClient {
        val mapper = SepulkaMessageMapperImpl()
        return SepulkaClientImpl(validator, mapper, service)
    }

    @Bean
    fun externalClient(internalClient: SepulkaClient): SepulkaClient {
        val client = MockMvcWebTestClient.bindToController(SepulkaController(internalClient)).build()
        return SepulkaClientSpringWebTest(client)
    }
}
