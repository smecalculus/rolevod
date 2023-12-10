package smecalculus.rolevod.construction;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;
import smecalculus.rolevod.configuration.PropsKeeper;
import smecalculus.rolevod.configuration.StorageDm.StorageProps;
import smecalculus.rolevod.configuration.StorageEm;
import smecalculus.rolevod.configuration.StoragePropsMapper;
import smecalculus.rolevod.configuration.StoragePropsMapperImpl;
import smecalculus.rolevod.validation.EdgeValidator;

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
        var storageProps = keeper.read("solution.storage", StorageEm.StorageProps.class);
        validator.validate(storageProps);
        LOG.info("Read {}", storageProps);
        return mapper.toDomain(storageProps);
    }
}
