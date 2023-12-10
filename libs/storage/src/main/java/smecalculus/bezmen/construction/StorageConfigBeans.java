package smecalculus.bezmen.construction;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;
import smecalculus.bezmen.configuration.PropsKeeper;
import smecalculus.bezmen.configuration.StorageDm.StorageProps;
import smecalculus.bezmen.configuration.StorageEm;
import smecalculus.bezmen.configuration.StoragePropsMapper;
import smecalculus.bezmen.configuration.StoragePropsMapperImpl;
import smecalculus.bezmen.validation.EdgeValidator;

@PropertySource("classpath:storage.properties")
@Configuration(proxyBeanMethods = false)
public class StorageConfigBeans {

    private static final Logger LOG = LoggerFactory.getLogger(StorageConfigBeans.class);

    @Bean
    StoragePropsMapper storagePropsMapper() {
        return new StoragePropsMapperImpl();
    }

    @Bean
    StorageProps storageProps(PropsKeeper keeper, EdgeValidator validator, StoragePropsMapper mapper) {
        var propsEdge = keeper.read("bezmen.storage", StorageEm.StorageProps.class);
        validator.validate(propsEdge);
        LOG.info("Read {}", propsEdge);
        return mapper.toDomain(propsEdge);
    }
}
