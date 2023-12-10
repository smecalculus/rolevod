package smecalculus.rolevod.validation;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import smecalculus.rolevod.configuration.ValidationDm;
import smecalculus.rolevod.configuration.ValidationDm.ValidationMode;
import smecalculus.rolevod.configuration.ValidationEm;
import smecalculus.rolevod.mapping.EdgeMapper;

@Mapper
public interface ValidationPropsMapper extends EdgeMapper {
    @Mapping(source = "mode", target = "validationMode")
    ValidationDm.ValidationProps toDomain(ValidationEm.ValidationProps props);

    default ValidationMode toValidationMode(String value) {
        return ValidationMode.valueOf(value.toUpperCase());
    }
}
