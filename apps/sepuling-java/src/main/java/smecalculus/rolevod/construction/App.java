package smecalculus.rolevod.construction;

import static smecalculus.rolevod.configuration.MessagingDm.MappingMode.SPRING_MVC;
import static smecalculus.rolevod.configuration.StorageDm.MappingMode.MY_BATIS;
import static smecalculus.rolevod.configuration.StorageDm.MappingMode.SPRING_DATA;

import org.springframework.boot.SpringApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;
import org.springframework.context.annotation.PropertySource;
import smecalculus.rolevod.core.SepulkaMapper;
import smecalculus.rolevod.core.SepulkaMapperImpl;
import smecalculus.rolevod.core.SepulkaService;
import smecalculus.rolevod.core.SepulkaServiceImpl;
import smecalculus.rolevod.messaging.SepulkaClient;
import smecalculus.rolevod.messaging.SepulkaClientImpl;
import smecalculus.rolevod.messaging.SepulkaMessageMapper;
import smecalculus.rolevod.messaging.SepulkaMessageMapperImpl;
import smecalculus.rolevod.messaging.springmvc.SepulkaController;
import smecalculus.rolevod.storage.SepulkaDao;
import smecalculus.rolevod.storage.SepulkaDaoMyBatis;
import smecalculus.rolevod.storage.SepulkaDaoSpringData;
import smecalculus.rolevod.storage.SepulkaStateMapper;
import smecalculus.rolevod.storage.SepulkaStateMapperImpl;
import smecalculus.rolevod.storage.mybatis.SepulkaSqlMapper;
import smecalculus.rolevod.storage.springdata.SepulkaRepository;
import smecalculus.rolevod.validation.EdgeValidator;

@Import({ConfigBeans.class, ValidationBeans.class, MessagingBeans.class, StorageBeans.class})
@PropertySource("classpath:application.properties")
@Configuration(proxyBeanMethods = false)
public class App {

    public static void main(String[] args) {
        SpringApplication.run(App.class, args);
    }

    @Bean
    @ConditionalOnMessagingMappingModes(SPRING_MVC)
    SepulkaController sepulkaControllerSpringMvc(SepulkaClient client) {
        return new SepulkaController(client);
    }

    @Bean
    SepulkaMessageMapper sepulkaMessageMapper() {
        return new SepulkaMessageMapperImpl();
    }

    @Bean
    SepulkaClient sepulkaClient(EdgeValidator validator, SepulkaMessageMapper mapper, SepulkaService service) {
        return new SepulkaClientImpl(validator, mapper, service);
    }

    @Bean
    SepulkaMapper sepulkaMapper() {
        return new SepulkaMapperImpl();
    }

    @Bean
    SepulkaService sepulkaService(SepulkaMapper mapper, SepulkaDao dao) {
        return new SepulkaServiceImpl(mapper, dao);
    }

    @Bean
    SepulkaStateMapper sepulkaStateMapper() {
        return new SepulkaStateMapperImpl();
    }

    @Bean
    @ConditionalOnStorageMappingMode(SPRING_DATA)
    SepulkaDaoSpringData sepulkaDaoSpringData(SepulkaStateMapper mapper, SepulkaRepository repository) {
        return new SepulkaDaoSpringData(mapper, repository);
    }

    @Bean
    @ConditionalOnStorageMappingMode(MY_BATIS)
    SepulkaDaoMyBatis sepulkaDaoMyBatis(SepulkaStateMapper stateMapper, SepulkaSqlMapper sqlMapper) {
        return new SepulkaDaoMyBatis(stateMapper, sqlMapper);
    }
}
