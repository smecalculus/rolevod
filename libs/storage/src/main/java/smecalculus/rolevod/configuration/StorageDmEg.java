package smecalculus.rolevod.configuration;

import static smecalculus.rolevod.configuration.StorageDm.MappingMode.SPRING_DATA;
import static smecalculus.rolevod.configuration.StorageDm.ProtocolMode.H2;

import smecalculus.rolevod.configuration.StorageDm.H2Props;
import smecalculus.rolevod.configuration.StorageDm.MappingMode;
import smecalculus.rolevod.configuration.StorageDm.MappingProps;
import smecalculus.rolevod.configuration.StorageDm.PostgresProps;
import smecalculus.rolevod.configuration.StorageDm.ProtocolMode;
import smecalculus.rolevod.configuration.StorageDm.ProtocolProps;
import smecalculus.rolevod.configuration.StorageDm.StorageProps;

public abstract class StorageDmEg {
    public static StorageProps.Builder storageProps() {
        return StorageProps.builder()
                .protocolProps(protocolProps().build())
                .mappingProps(mappingProps().build());
    }

    public static StorageProps.Builder storageProps(MappingMode mappingMode, ProtocolMode protocolMode) {
        return storageProps()
                .protocolProps(protocolProps(protocolMode).build())
                .mappingProps(mappingProps(mappingMode).build());
    }

    public static MappingProps.Builder mappingProps() {
        return MappingProps.builder().mappingMode(SPRING_DATA);
    }

    public static MappingProps.Builder mappingProps(MappingMode mode) {
        return mappingProps().mappingMode(mode);
    }

    public static ProtocolProps.Builder protocolProps() {
        return ProtocolProps.builder()
                .protocolMode(H2)
                .h2Props(H2Props.builder()
                        .url("jdbc:h2:mem:toy;DB_CLOSE_DELAY=-1")
                        .username("toy")
                        .password("toy")
                        .build())
                .postgresProps(PostgresProps.builder()
                        .url("jdbc:postgresql://localhost:5432/toy")
                        .schema("toy")
                        .username("toy")
                        .password("toy")
                        .build());
    }

    public static ProtocolProps.Builder protocolProps(ProtocolMode mode) {
        return protocolProps().protocolMode(mode);
    }
}
