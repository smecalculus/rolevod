package smecalculus.rolevod.messaging;

import org.mapstruct.Mapper;
import smecalculus.rolevod.core.SepulkaMessageDm;
import smecalculus.rolevod.mapping.EdgeMapper;

@Mapper
public interface SepulkaMessageMapper extends EdgeMapper {
    SepulkaMessageDm.RegistrationRequest toDomain(SepulkaMessageEm.RegistrationRequest request);

    SepulkaMessageEm.RegistrationResponse toEdge(SepulkaMessageDm.RegistrationResponse response);
}
