package smecalculus.bezmen.configuration;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import smecalculus.bezmen.configuration.MessagingDm.MappingMode;
import smecalculus.bezmen.configuration.MessagingDm.ProtocolMode;
import smecalculus.bezmen.mapping.EdgeMapper;

@Mapper
public interface MessagingPropsMapper extends EdgeMapper {
    @Mapping(source = "protocol", target = "protocolProps")
    @Mapping(source = "mapping", target = "mappingProps")
    MessagingDm.MessagingProps toDomain(MessagingEm.MessagingProps propsEdge);

    @Mapping(source = "modes", target = "protocolModes")
    MessagingDm.ProtocolProps toDomain(MessagingEm.ProtocolProps propsEdge);

    @Mapping(source = "modes", target = "mappingModes")
    MessagingDm.MappingProps toDomain(MessagingEm.MappingProps propsEdge);

    default ProtocolMode toProtocolMode(String value) {
        return ProtocolMode.valueOf(value.toUpperCase());
    }

    default MappingMode toMappingMode(String value) {
        return MappingMode.valueOf(value.toUpperCase());
    }
}
