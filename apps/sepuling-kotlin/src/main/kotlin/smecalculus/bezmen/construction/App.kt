package smecalculus.bezmen.construction

import org.springframework.boot.runApplication
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.context.annotation.Import
import org.springframework.context.annotation.PropertySource
import smecalculus.bezmen.configuration.MessagingDm.MappingMode.SPRING_MVC
import smecalculus.bezmen.configuration.StorageDm.MappingMode.MY_BATIS
import smecalculus.bezmen.configuration.StorageDm.MappingMode.SPRING_DATA
import smecalculus.bezmen.core.SepulkaMapper
import smecalculus.bezmen.core.SepulkaMapperImpl
import smecalculus.bezmen.core.SepulkaService
import smecalculus.bezmen.core.SepulkaServiceImpl
import smecalculus.bezmen.messaging.SepulkaClient
import smecalculus.bezmen.messaging.SepulkaClientImpl
import smecalculus.bezmen.messaging.SepulkaMessageMapper
import smecalculus.bezmen.messaging.SepulkaMessageMapperImpl
import smecalculus.bezmen.messaging.springmvc.SepulkaController
import smecalculus.bezmen.storage.SepulkaDao
import smecalculus.bezmen.storage.SepulkaDaoMyBatis
import smecalculus.bezmen.storage.SepulkaDaoSpringData
import smecalculus.bezmen.storage.SepulkaStateMapper
import smecalculus.bezmen.storage.SepulkaStateMapperImpl
import smecalculus.bezmen.storage.mybatis.SepulkaSqlMapper
import smecalculus.bezmen.storage.springdata.SepulkaRepository
import smecalculus.bezmen.validation.EdgeValidator

fun main(args: Array<String>) {
    runApplication<App>(*args)
}

@Import(ConfigBeans::class, ValidationBeans::class, MessagingBeans::class, StorageBeans::class)
@PropertySource("classpath:application.properties")
@Configuration(proxyBeanMethods = false)
class App {
    @Bean
    @ConditionalOnMessagingMappingModes(SPRING_MVC)
    fun sepulkaControllerSpringMvc(client: SepulkaClient): SepulkaController {
        return SepulkaController(client)
    }

    @Bean
    fun sepulkaMessageMapper(): SepulkaMessageMapper {
        return SepulkaMessageMapperImpl()
    }

    @Bean
    fun sepulkaClient(
        validator: EdgeValidator,
        mapper: SepulkaMessageMapper,
        service: SepulkaService,
    ): SepulkaClient {
        return SepulkaClientImpl(validator, mapper, service)
    }

    @Bean
    fun sepulkaMapper(): SepulkaMapper {
        return SepulkaMapperImpl()
    }

    @Bean
    fun sepulkaService(
        mapper: SepulkaMapper,
        dao: SepulkaDao,
    ): SepulkaService {
        return SepulkaServiceImpl(mapper, dao)
    }

    @Bean
    fun sepulkaStateMapper(): SepulkaStateMapper {
        return SepulkaStateMapperImpl()
    }

    @Bean
    @ConditionalOnStorageMappingMode(SPRING_DATA)
    fun sepulkaDaoSpringData(
        mapper: SepulkaStateMapper,
        repository: SepulkaRepository,
    ): SepulkaDaoSpringData {
        return SepulkaDaoSpringData(mapper, repository)
    }

    @Bean
    @ConditionalOnStorageMappingMode(MY_BATIS)
    fun sepulkaDaoMyBatis(
        stateMapper: SepulkaStateMapper,
        sqlMapper: SepulkaSqlMapper,
    ): SepulkaDaoMyBatis {
        return SepulkaDaoMyBatis(stateMapper, sqlMapper)
    }
}
