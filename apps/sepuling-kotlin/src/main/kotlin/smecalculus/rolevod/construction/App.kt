package smecalculus.rolevod.construction

import org.springframework.boot.runApplication
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.context.annotation.Import
import org.springframework.context.annotation.PropertySource
import smecalculus.rolevod.configuration.MessagingDm.MappingMode.SPRING_MVC
import smecalculus.rolevod.configuration.StorageDm.MappingMode.MY_BATIS
import smecalculus.rolevod.configuration.StorageDm.MappingMode.SPRING_DATA
import smecalculus.rolevod.core.SepulkaMapper
import smecalculus.rolevod.core.SepulkaMapperImpl
import smecalculus.rolevod.core.SepulkaService
import smecalculus.rolevod.core.SepulkaServiceImpl
import smecalculus.rolevod.messaging.SepulkaClient
import smecalculus.rolevod.messaging.SepulkaClientImpl
import smecalculus.rolevod.messaging.SepulkaMessageMapper
import smecalculus.rolevod.messaging.SepulkaMessageMapperImpl
import smecalculus.rolevod.messaging.springmvc.SepulkaController
import smecalculus.rolevod.storage.SepulkaDao
import smecalculus.rolevod.storage.SepulkaDaoMyBatis
import smecalculus.rolevod.storage.SepulkaDaoSpringData
import smecalculus.rolevod.storage.SepulkaStateMapper
import smecalculus.rolevod.storage.SepulkaStateMapperImpl
import smecalculus.rolevod.storage.mybatis.SepulkaSqlMapper
import smecalculus.rolevod.storage.springdata.SepulkaRepository
import smecalculus.rolevod.validation.EdgeValidator

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
