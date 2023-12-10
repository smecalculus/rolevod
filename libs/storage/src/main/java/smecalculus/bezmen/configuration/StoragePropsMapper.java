package smecalculus.bezmen.configuration;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import smecalculus.bezmen.configuration.StorageDm.MappingMode;
import smecalculus.bezmen.configuration.StorageDm.ProtocolMode;

@Mapper
public interface StoragePropsMapper {
    @Mapping(source = "protocol", target = "protocolProps")
    @Mapping(source = "mapping", target = "mappingProps")
    StorageDm.StorageProps toDomain(StorageEm.StorageProps propsEdge);

    @Mapping(source = "mode", target = "protocolMode")
    @Mapping(source = "h2", target = "h2Props")
    @Mapping(source = "postgres", target = "postgresProps")
    StorageDm.ProtocolProps toDomain(StorageEm.ProtocolProps propsEdge);

    @Mapping(source = "mode", target = "mappingMode")
    StorageDm.MappingProps toDomain(StorageEm.MappingProps propsEdge);

    default ProtocolMode toProtocolMode(String mode) {
        return ProtocolMode.valueOf(mode.toUpperCase());
    }

    default MappingMode toMappingMode(String mode) {
        return MappingMode.valueOf(mode.toUpperCase());
    }
}
