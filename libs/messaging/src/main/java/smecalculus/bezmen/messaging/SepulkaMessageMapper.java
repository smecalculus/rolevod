package smecalculus.bezmen.messaging;

import org.mapstruct.Mapper;
import smecalculus.bezmen.core.SepulkaMessageDm;
import smecalculus.bezmen.mapping.EdgeMapper;

@Mapper
public interface SepulkaMessageMapper extends EdgeMapper {
    SepulkaMessageDm.RegistrationRequest toDomain(SepulkaMessageEm.RegistrationRequest request);

    SepulkaMessageEm.RegistrationResponse toEdge(SepulkaMessageDm.RegistrationResponse response);
}
