package smecalculus.rolevod.construction;

import static smecalculus.rolevod.configuration.StorageDm.MappingMode.MY_BATIS;
import static smecalculus.rolevod.configuration.StorageDm.MappingMode.SPRING_DATA;
import static smecalculus.rolevod.configuration.StorageDm.ProtocolMode.H2;
import static smecalculus.rolevod.configuration.StorageDm.ProtocolMode.POSTGRES;

import org.springframework.context.annotation.Bean;
import smecalculus.rolevod.configuration.StorageDm.StorageProps;
import smecalculus.rolevod.configuration.StorageDmEg;

public class StoragePropsBeans {

    public static class SpringDataPostgres {
        @Bean
        public StorageProps storageProps() {
            return StorageDmEg.storageProps(SPRING_DATA, POSTGRES).build();
        }
    }

    public static class SpringDataH2 {
        @Bean
        public StorageProps storageProps() {
            return StorageDmEg.storageProps(SPRING_DATA, H2).build();
        }
    }

    public static class MyBatisPostgres {
        @Bean
        public StorageProps storageProps() {
            return StorageDmEg.storageProps(MY_BATIS, POSTGRES).build();
        }
    }

    public static class MyBatisH2 {
        @Bean
        public StorageProps storageProps() {
            return StorageDmEg.storageProps(MY_BATIS, H2).build();
        }
    }
}
