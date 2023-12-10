package smecalculus.rolevod.construction

import org.mockito.Mockito.mock
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.context.annotation.Import
import org.springframework.test.web.servlet.client.MockMvcWebTestClient
import smecalculus.rolevod.core.SepulkaService
import smecalculus.rolevod.messaging.SepulkaClient
import smecalculus.rolevod.messaging.SepulkaClientImpl
import smecalculus.rolevod.messaging.SepulkaClientSpringWebTest
import smecalculus.rolevod.messaging.SepulkaMessageMapperImpl
import smecalculus.rolevod.messaging.springmvc.SepulkaController
import smecalculus.rolevod.validation.EdgeValidator

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
