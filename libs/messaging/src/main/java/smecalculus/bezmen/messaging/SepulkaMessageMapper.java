package smecalculus.bezmen.messaging;

import org.mapstruct.Mapper;
import smecalculus.bezmen.core.MessageDm;
import smecalculus.bezmen.mapping.EdgeMapper;

@Mapper
public interface SepulkaMessageMapper extends EdgeMapper {
    MessageDm.RegistrationRequest toDomain(MessageEm.RegistrationRequest request);

    MessageEm.RegistrationResponse toEdge(MessageDm.RegistrationResponse response);
}
